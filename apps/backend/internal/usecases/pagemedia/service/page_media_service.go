package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	mediaEntity "go_refine_dashboard_be/internal/domain/media/entity"
	pageEntity "go_refine_dashboard_be/internal/domain/page/entity"
	"go_refine_dashboard_be/internal/infrastructure/cache"
	"go_refine_dashboard_be/internal/infrastructure/repository"
	"go_refine_dashboard_be/internal/usecases/pagemedia/dto"
	"gorm.io/gorm"
)

var (
	ErrPageNotFound       = errors.New("trang không tồn tại")
	ErrMediaNotFound      = errors.New("media không tồn tại")
	ErrMediaForbidden     = errors.New("không có quyền sử dụng media này")
	ErrCollectionInvalid  = errors.New("collection không hợp lệ")
	ErrThumbnailLimit     = errors.New("thumbnail chỉ được có một ảnh")
	ErrMediaDuplicate     = errors.New("media bị trùng trong collection")
	ErrSortOrderInvalid   = errors.New("sort order không hợp lệ")
	ErrMediaTypeInvalid   = errors.New("media phải là file hình ảnh")
	ErrPageMediaNotFound  = errors.New("liên kết page media không tồn tại")
	ErrCollectionRequired = errors.New("collection là bắt buộc")
)

type PageMediaService interface {
	GetPageMedia(ctx context.Context, pageID uint, collection string) ([]dto.PageMediaResponse, int64, error)
	CreatePageMedia(ctx context.Context, pageID uint, req *dto.CreatePageMediaRequest, userID uint, role string) (*dto.PageMediaResponse, error)
	SyncPageMedia(ctx context.Context, pageID uint, collection string, req *dto.SyncPageMediaRequest, userID uint, role string) ([]dto.PageMediaResponse, error)
	DeletePageMedia(ctx context.Context, pageID, mediaID uint, collection string, userID uint, role string) error
}

type pageMediaService struct {
	repo            repository.PageMediaRepository
	db              *gorm.DB
	redisClient     *redis.Client
	publicPageCache cache.PublicPageCacheInvalidator
}

func NewPageMediaService(repo repository.PageMediaRepository, db *gorm.DB, redisClient *redis.Client, publicPageCache cache.PublicPageCacheInvalidator) PageMediaService {
	return &pageMediaService{repo: repo, db: db, redisClient: redisClient, publicPageCache: publicPageCache}
}

func validCollection(collection string) bool {
	return collection == "thumbnail" || collection == "gallery"
}

func canManage(media mediaEntity.Media, userID uint, role string) bool {
	return role == "ADMIN" || (role == "STAFF" && media.OwnerID == userID)
}

func validateMedia(media []mediaEntity.Media, ids []uint, userID uint, role string) error {
	if len(media) != len(ids) {
		return ErrMediaNotFound
	}
	seen := make(map[uint]struct{}, len(ids))
	for _, id := range ids {
		if _, ok := seen[id]; ok {
			return ErrMediaDuplicate
		}
		seen[id] = struct{}{}
	}
	for _, item := range media {
		if !canManage(item, userID, role) {
			return ErrMediaForbidden
		}
		if !strings.HasPrefix(strings.ToLower(item.MimeType), "image/") {
			return ErrMediaTypeInvalid
		}
	}
	return nil
}

func mapPageMedia(item pageEntity.PageMedia) dto.PageMediaResponse {
	return dto.PageMediaResponse{ID: item.ID, PageID: item.PageID, MediaID: item.MediaID, Collection: string(item.Collection), SortOrder: item.SortOrder, Media: dto.MapMediaResponse(item.Media)}
}

func (s *pageMediaService) cacheKey(pageID uint, collection string) string {
	if collection == "" {
		return fmt.Sprintf("cache:page:%d:media:all", pageID)
	}
	return fmt.Sprintf("cache:page:%d:media:%s", pageID, collection)
}

func (s *pageMediaService) invalidateCache(ctx context.Context, pageID uint) {
	if s.redisClient == nil {
		return
	}
	keys := []string{s.cacheKey(pageID, ""), s.cacheKey(pageID, "thumbnail"), s.cacheKey(pageID, "gallery"), fmt.Sprintf("seo:page:%d", pageID)}
	for _, key := range keys {
		if err := s.redisClient.Del(ctx, key).Err(); err != nil {
			log.Printf("page media cache invalidation failed for %q: %v", key, err)
		}
	}
}

func (s *pageMediaService) ensurePage(ctx context.Context, pageID uint) error {
	var page pageEntity.Page
	if err := s.db.WithContext(ctx).First(&page, pageID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrPageNotFound
		}
		return fmt.Errorf("kiểm tra trang: %w", err)
	}
	return nil
}

func (s *pageMediaService) GetPageMedia(ctx context.Context, pageID uint, collection string) ([]dto.PageMediaResponse, int64, error) {
	if !validCollection(collection) && collection != "" {
		return nil, 0, ErrCollectionInvalid
	}
	if err := s.ensurePage(ctx, pageID); err != nil {
		return nil, 0, err
	}
	key := s.cacheKey(pageID, collection)
	if s.redisClient != nil {
		if payload, err := s.redisClient.Get(ctx, key).Result(); err == nil {
			var result []dto.PageMediaResponse
			if json.Unmarshal([]byte(payload), &result) == nil {
				return result, int64(len(result)), nil
			}
		}
	}
	items, total, err := s.repo.FindByPageID(ctx, pageID, collection)
	if err != nil {
		return nil, 0, fmt.Errorf("lấy media trang: %w", err)
	}
	result := make([]dto.PageMediaResponse, 0, len(items))
	for _, item := range items {
		result = append(result, mapPageMedia(item))
	}
	if s.redisClient != nil {
		if payload, marshalErr := json.Marshal(result); marshalErr == nil {
			_ = s.redisClient.SetEx(ctx, key, payload, 15*time.Minute).Err()
		}
	}
	return result, total, nil
}

func (s *pageMediaService) CreatePageMedia(ctx context.Context, pageID uint, req *dto.CreatePageMediaRequest, userID uint, role string) (*dto.PageMediaResponse, error) {
	if !validCollection(req.Collection) {
		return nil, ErrCollectionInvalid
	}
	var item pageEntity.PageMedia
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if _, err := s.repo.FindPageForUpdate(ctx, tx, pageID); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrPageNotFound
			}
			return err
		}
		media, err := s.repo.FindMediaByIDsForUpdate(ctx, tx, []uint{req.MediaID})
		if err != nil {
			return err
		}
		if err := validateMedia(media, []uint{req.MediaID}, userID, role); err != nil {
			return err
		}
		if req.Collection == "thumbnail" {
			if err := s.repo.SoftDeleteNotIn(ctx, tx, pageID, req.Collection, []uint{req.MediaID}); err != nil {
				return err
			}
		}
		sortOrder := 0
		if req.SortOrder != nil {
			sortOrder = *req.SortOrder
		} else if req.Collection == "gallery" {
			sortOrder, err = s.repo.NextSortOrder(ctx, tx, pageID, req.Collection)
			if err != nil {
				return err
			}
		}
		item = pageEntity.PageMedia{PageID: pageID, MediaID: req.MediaID, Collection: pageEntity.PageMediaCollection(req.Collection), SortOrder: sortOrder}
		if err := s.repo.CreateOrRestore(ctx, tx, &item); err != nil {
			return err
		}
		return s.repo.MarkAttached(ctx, tx, []uint{req.MediaID})
	})
	if err != nil {
		return nil, fmt.Errorf("liên kết media: %w", err)
	}
	s.invalidateCache(ctx, pageID)
	s.publicPageCache.InvalidatePage(ctx, pageID)
	items, _, err := s.GetPageMedia(ctx, pageID, req.Collection)
	if err != nil {
		return nil, err
	}
	for _, result := range items {
		if result.MediaID == req.MediaID {
			return &result, nil
		}
	}
	return nil, ErrPageMediaNotFound
}

func (s *pageMediaService) SyncPageMedia(ctx context.Context, pageID uint, collection string, req *dto.SyncPageMediaRequest, userID uint, role string) ([]dto.PageMediaResponse, error) {
	if !validCollection(collection) {
		return nil, ErrCollectionInvalid
	}
	if collection == "thumbnail" && len(req.Media) > 1 {
		return nil, ErrThumbnailLimit
	}
	ids := make([]uint, 0, len(req.Media))
	seenOrder := make(map[int]struct{}, len(req.Media))
	for index, item := range req.Media {
		if item.SortOrder != index {
			return nil, ErrSortOrderInvalid
		}
		if _, ok := seenOrder[item.SortOrder]; ok {
			return nil, ErrMediaDuplicate
		}
		seenOrder[item.SortOrder] = struct{}{}
		ids = append(ids, item.ID)
	}
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if _, err := s.repo.FindPageForUpdate(ctx, tx, pageID); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrPageNotFound
			}
			return err
		}
		media, err := s.repo.FindMediaByIDsForUpdate(ctx, tx, ids)
		if err != nil {
			return err
		}
		if err := validateMedia(media, ids, userID, role); err != nil {
			return err
		}
		if err := s.repo.SoftDeleteNotIn(ctx, tx, pageID, collection, ids); err != nil {
			return err
		}
		for index, id := range ids {
			if err := s.repo.CreateOrRestore(ctx, tx, &pageEntity.PageMedia{PageID: pageID, MediaID: id, Collection: pageEntity.PageMediaCollection(collection), SortOrder: index}); err != nil {
				return err
			}
		}
		return s.repo.MarkAttached(ctx, tx, ids)
	})
	if err != nil {
		return nil, fmt.Errorf("đồng bộ media: %w", err)
	}
	s.invalidateCache(ctx, pageID)
	s.publicPageCache.InvalidatePage(ctx, pageID)
	result, _, err := s.GetPageMedia(ctx, pageID, collection)
	return result, err
}

func (s *pageMediaService) DeletePageMedia(ctx context.Context, pageID, mediaID uint, collection string, userID uint, role string) error {
	if collection != "" && !validCollection(collection) {
		return ErrCollectionInvalid
	}
	var deleted bool
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if _, err := s.repo.FindPageForUpdate(ctx, tx, pageID); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrPageNotFound
			}
			return err
		}
		media, err := s.repo.FindMediaByIDsForUpdate(ctx, tx, []uint{mediaID})
		if err != nil {
			return err
		}
		if len(media) != 1 {
			return ErrMediaNotFound
		}
		if !canManage(media[0], userID, role) {
			return ErrMediaForbidden
		}
		links, err := s.repo.FindActiveLinks(ctx, tx, pageID, mediaID, collection)
		if err != nil {
			return err
		}
		if len(links) == 0 {
			if collection == "" {
				return ErrPageMediaNotFound
			}
			return ErrPageMediaNotFound
		}
		if collection == "" && len(links) > 1 {
			return ErrCollectionRequired
		}
		selected := string(links[0].Collection)
		if collection != "" {
			selected = collection
		}
		deleted, err = s.repo.SoftDeleteLink(ctx, tx, pageID, mediaID, selected)
		return err
	})
	if err != nil {
		return err
	}
	if !deleted {
		return ErrPageMediaNotFound
	}
	s.invalidateCache(ctx, pageID)
	s.publicPageCache.InvalidatePage(ctx, pageID)
	return nil
}
