package service

import (
	"context"
	"errors"
	"log"
	"regexp"
	"strings"

	"github.com/redis/go-redis/v9"
	mediaEntity "go_refine_dashboard_be/internal/domain/media/entity"
	pageEntity "go_refine_dashboard_be/internal/domain/page/entity"
	"go_refine_dashboard_be/internal/infrastructure/repository"
	"go_refine_dashboard_be/internal/usecases/pagesection/dto"
	"gorm.io/gorm"
)

var (
	ErrPageNotFound           = errors.New("trang không tồn tại")
	ErrSectionNotFound        = errors.New("section không tồn tại")
	ErrSourceNotFound         = errors.New("nguồn dữ liệu không tồn tại")
	ErrInvalidInput           = errors.New("dữ liệu section không hợp lệ")
	ErrSectionKeyConflict     = errors.New("key section đã tồn tại")
	ErrSectionOrderConflict   = errors.New("thứ tự section đã thay đổi, vui lòng tải lại")
	ErrCollectionTypeConflict = errors.New("collection đang thuộc item type khác")
	ErrSectionItemConflict    = errors.New("item section đã tồn tại")
	ErrMediaForbidden         = errors.New("không có quyền sử dụng media này")
	ErrMediaTypeInvalid       = errors.New("media phải là file hình ảnh")
)
var sectionKeyPattern = regexp.MustCompile(`^[a-z0-9]+(?:_[a-z0-9]+)*$`)
var collectionPattern = sectionKeyPattern
var colorPattern = regexp.MustCompile(`^#[0-9A-Fa-f]{6}([0-9A-Fa-f]{2})?$`)

type PageSectionService interface {
	GetSections(context.Context, uint, dto.GetPageSectionsQuery) (*dto.PageSectionListResponse, error)
	GetSection(context.Context, uint, uint) (*dto.PageSectionResponse, error)
	CreateSection(context.Context, uint, *dto.CreatePageSectionRequest, uint, string) (*dto.PageSectionResponse, error)
	UpdateSection(context.Context, uint, uint, *dto.UpdatePageSectionRequest, uint, string) (*dto.PageSectionResponse, error)
	DeleteSection(context.Context, uint, uint) error
	ReorderSections(context.Context, uint, *dto.ReorderPageSectionsRequest) error
	GetItems(context.Context, uint, uint, dto.GetSectionItemsQuery) (*dto.PageSectionItemListResponse, error)
	SyncItems(context.Context, uint, uint, string, *dto.SyncSectionCollectionRequest, uint, string) (*dto.PageSectionItemListResponse, error)
}
type pageSectionService struct {
	repo        repository.PageSectionRepository
	db          *gorm.DB
	redisClient *redis.Client
}

func NewPageSectionService(repo repository.PageSectionRepository, db *gorm.DB, redisClient *redis.Client) PageSectionService {
	return &pageSectionService{repo: repo, db: db, redisClient: redisClient}
}
func (s *pageSectionService) invalidate(ctx context.Context, pageID uint) {
	if s.redisClient == nil {
		return
	}
	var p struct{ Slug string }
	if s.db.WithContext(ctx).Table("pages").Select("slug").Where("id = ?", pageID).Scan(&p).Error == nil && p.Slug != "" {
		if err := s.redisClient.Del(ctx, "public:pages:slug:"+p.Slug).Err(); err != nil {
			log.Printf("section cache invalidation failed: %v", err)
		}
	}
}
func ensurePage(repo repository.PageSectionRepository, ctx context.Context, tx *gorm.DB, id uint) error {
	_, err := repo.FindPageForUpdate(ctx, tx, id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrPageNotFound
	}
	return err
}
func validateMeta(key, name, title, description, status string, color *string) error {
	if !sectionKeyPattern.MatchString(strings.TrimSpace(key)) || strings.TrimSpace(name) == "" || len(description) > 5000 {
		return ErrInvalidInput
	}
	if status != "ACTIVE" && status != "INACTIVE" {
		return ErrInvalidInput
	}
	if color != nil && *color != "" && !colorPattern.MatchString(strings.TrimSpace(*color)) {
		return ErrInvalidInput
	}
	_ = title
	return nil
}
func validateMediaIDs(ctx context.Context, repo repository.PageSectionRepository, tx *gorm.DB, ids []uint, userID uint, role string) error {
	if len(ids) == 0 {
		return nil
	}
	rows, err := repo.FindMediaByIDsForUpdate(ctx, tx, ids)
	if err != nil {
		return err
	}
	if len(rows) != len(ids) {
		return ErrSourceNotFound
	}
	seen := map[uint]bool{}
	for _, m := range rows {
		if seen[m.ID] {
			return ErrInvalidInput
		}
		seen[m.ID] = true
		if !strings.HasPrefix(strings.ToLower(m.MimeType), "image/") {
			return ErrMediaTypeInvalid
		}
		if role != "ADMIN" && m.OwnerID != userID {
			return ErrMediaForbidden
		}
	}
	return nil
}
func (s *pageSectionService) GetSections(ctx context.Context, pageID uint, q dto.GetPageSectionsQuery) (*dto.PageSectionListResponse, error) {
	cur := q.Current
	if cur < 1 {
		cur = 1
	}
	size := q.PageSize
	if size < 1 {
		size = 100
	}
	if size > 100 {
		size = 100
	}
	var page pageEntity.Page
	if err := s.db.WithContext(ctx).First(&page, pageID).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrPageNotFound
	} else if err != nil {
		return nil, err
	}
	rows, total, err := s.repo.FindSections(ctx, pageID, q.Status, (cur-1)*size, size)
	if err != nil {
		return nil, err
	}
	ids := make([]uint, 0, len(rows))
	for _, v := range rows {
		ids = append(ids, v.ID)
	}
	summary, err := s.repo.FindCollectionSummaries(ctx, ids)
	if err != nil {
		return nil, err
	}
	out := make([]dto.PageSectionResponse, 0, len(rows))
	for _, v := range rows {
		out = append(out, mapSection(v, summary[v.ID]))
	}
	return &dto.PageSectionListResponse{Data: out, Total: total}, nil
}
func (s *pageSectionService) GetSection(ctx context.Context, pageID, sectionID uint) (*dto.PageSectionResponse, error) {
	rows, _, err := s.repo.FindSections(ctx, pageID, "", 0, 100)
	if err != nil {
		return nil, err
	}
	for _, v := range rows {
		if v.ID == sectionID {
			sum, err := s.repo.FindCollectionSummaries(ctx, []uint{v.ID})
			if err != nil {
				return nil, err
			}
			r := mapSection(v, sum[v.ID])
			return &r, nil
		}
	}
	return nil, ErrSectionNotFound
}
func (s *pageSectionService) CreateSection(ctx context.Context, pageID uint, req *dto.CreatePageSectionRequest, userID uint, role string) (*dto.PageSectionResponse, error) {
	if err := validateMeta(req.Key, req.Name, req.Title, req.Description, req.Status, req.BackgroundColor); err != nil {
		return nil, err
	}
	if req.SortOrder != nil && *req.SortOrder < 0 {
		return nil, ErrInvalidInput
	}
	var created pageEntity.PageSection
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := ensurePage(s.repo, ctx, tx, pageID); err != nil {
			return err
		}
		active, err := s.repo.FindActiveSectionsForUpdate(ctx, tx, pageID)
		if err != nil {
			return err
		}
		if len(active) >= 100 {
			return ErrInvalidInput
		}
		order := len(active)
		if req.SortOrder != nil {
			order = *req.SortOrder
			if order > len(active) {
				return ErrInvalidInput
			}
		}
		ids := []uint{}
		if req.BackgroundMediaID != nil {
			ids = append(ids, *req.BackgroundMediaID)
		}
		if req.FeatureMediaID != nil {
			ids = append(ids, *req.FeatureMediaID)
		}
		if err := validateMediaIDs(ctx, s.repo, tx, ids, userID, role); err != nil {
			return err
		}
		created = pageEntity.PageSection{PageID: pageID, Key: strings.TrimSpace(req.Key), Name: strings.TrimSpace(req.Name), Title: strings.TrimSpace(req.Title), Description: strings.TrimSpace(req.Description), BackgroundColor: req.BackgroundColor, BackgroundMediaID: req.BackgroundMediaID, FeatureMediaID: req.FeatureMediaID, SortOrder: order, Status: pageEntity.PageSectionStatus(req.Status)}
		if err := s.repo.CreateOrRestoreSection(ctx, tx, &created); err != nil {
			if errors.Is(err, gorm.ErrDuplicatedKey) {
				return ErrSectionKeyConflict
			}
			return err
		}
		return s.repo.MarkMediaAttached(ctx, tx, ids)
	})
	if err != nil {
		return nil, err
	}
	s.invalidate(ctx, pageID)
	return s.GetSection(ctx, pageID, created.ID)
}
func (s *pageSectionService) UpdateSection(ctx context.Context, pageID, sectionID uint, req *dto.UpdatePageSectionRequest, userID uint, role string) (*dto.PageSectionResponse, error) {
	if err := validateMeta(req.Key, req.Name, req.Title, req.Description, req.Status, req.BackgroundColor); err != nil {
		return nil, err
	}
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := ensurePage(s.repo, ctx, tx, pageID); err != nil {
			return err
		}
		old, err := s.repo.FindSectionForUpdate(ctx, tx, pageID, sectionID)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrSectionNotFound
		}
		if err != nil {
			return err
		}
		ids := []uint{}
		if req.BackgroundMediaID != nil {
			ids = append(ids, *req.BackgroundMediaID)
		}
		if req.FeatureMediaID != nil {
			ids = append(ids, *req.FeatureMediaID)
		}
		if err := validateMediaIDs(ctx, s.repo, tx, ids, userID, role); err != nil {
			return err
		}
		old.Key = strings.TrimSpace(req.Key)
		old.Name = strings.TrimSpace(req.Name)
		old.Title = strings.TrimSpace(req.Title)
		old.Description = strings.TrimSpace(req.Description)
		old.BackgroundColor = req.BackgroundColor
		old.BackgroundMediaID = req.BackgroundMediaID
		old.FeatureMediaID = req.FeatureMediaID
		old.Status = pageEntity.PageSectionStatus(req.Status)
		if err := s.repo.UpdateSection(ctx, tx, old); err != nil {
			return err
		}
		return s.repo.MarkMediaAttached(ctx, tx, ids)
	})
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, ErrSectionKeyConflict
		}
		return nil, err
	}
	s.invalidate(ctx, pageID)
	return s.GetSection(ctx, pageID, sectionID)
}
func (s *pageSectionService) DeleteSection(ctx context.Context, pageID, sectionID uint) error {
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := ensurePage(s.repo, ctx, tx, pageID); err != nil {
			return err
		}
		if _, err := s.repo.FindSectionForUpdate(ctx, tx, pageID, sectionID); errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrSectionNotFound
		} else if err != nil {
			return err
		}
		if err := s.repo.SoftDeleteSectionAndItems(ctx, tx, sectionID); err != nil {
			return err
		}
		rows, err := s.repo.FindActiveSectionsForUpdate(ctx, tx, pageID)
		if err != nil {
			return err
		}
		updates := make([]repository.SectionOrderUpdate, 0, len(rows))
		for i, v := range rows {
			updates = append(updates, repository.SectionOrderUpdate{ID: v.ID, SortOrder: i})
		}
		return s.repo.BulkUpdateSectionOrder(ctx, tx, pageID, updates)
	})
	if err == nil {
		s.invalidate(ctx, pageID)
	}
	return err
}
func (s *pageSectionService) ReorderSections(ctx context.Context, pageID uint, req *dto.ReorderPageSectionsRequest) error {
	seen := map[uint]bool{}
	for i, v := range req.Sections {
		if seen[v.ID] || v.SortOrder != i {
			return ErrSectionOrderConflict
		}
		seen[v.ID] = true
	}
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := ensurePage(s.repo, ctx, tx, pageID); err != nil {
			return err
		}
		rows, err := s.repo.FindActiveSectionsForUpdate(ctx, tx, pageID)
		if err != nil {
			return err
		}
		if len(rows) != len(req.Sections) {
			return ErrSectionOrderConflict
		}
		for _, v := range rows {
			if !seen[v.ID] {
				return ErrSectionOrderConflict
			}
		}
		updates := make([]repository.SectionOrderUpdate, 0, len(req.Sections))
		for _, v := range req.Sections {
			updates = append(updates, repository.SectionOrderUpdate{ID: v.ID, SortOrder: v.SortOrder})
		}
		if err := s.repo.BulkUpdateSectionOrder(ctx, tx, pageID, updates); err != nil {
			return err
		}
		return nil
	})
	if err == nil {
		s.invalidate(ctx, pageID)
	}
	return err
}
func (s *pageSectionService) GetItems(ctx context.Context, pageID, sectionID uint, q dto.GetSectionItemsQuery) (*dto.PageSectionItemListResponse, error) {
	if !collectionPattern.MatchString(q.Collection) {
		return nil, ErrInvalidInput
	}
	var section pageEntity.PageSection
	if err := s.db.WithContext(ctx).Where("page_id = ? AND id = ?", pageID, sectionID).First(&section).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		var page pageEntity.Page
		if pageErr := s.db.WithContext(ctx).First(&page, pageID).Error; errors.Is(pageErr, gorm.ErrRecordNotFound) {
			return nil, ErrPageNotFound
		}
		return nil, ErrSectionNotFound
	} else if err != nil {
		return nil, err
	}
	cur := q.Current
	if cur < 1 {
		cur = 1
	}
	size := q.PageSize
	if size < 1 {
		size = 100
	}
	if size > 100 {
		size = 100
	}
	rows, total, err := s.repo.FindSectionItems(ctx, pageID, sectionID, q.Collection, (cur-1)*size, size)
	if err != nil {
		return nil, err
	}
	out, err := s.mapItems(ctx, rows)
	if err != nil {
		return nil, err
	}
	return &dto.PageSectionItemListResponse{Data: out, Total: total}, nil
}
func (s *pageSectionService) SyncItems(ctx context.Context, pageID, sectionID uint, collection string, req *dto.SyncSectionCollectionRequest, userID uint, role string) (*dto.PageSectionItemListResponse, error) {
	if !collectionPattern.MatchString(collection) || len(req.Items) > 100 {
		return nil, ErrInvalidInput
	}
	typ := pageEntity.PageSectionItemType(req.ItemType)
	seen := map[uint]bool{}
	for i, v := range req.Items {
		if seen[v.ItemID] || v.SortOrder != i {
			return nil, ErrInvalidInput
		}
		seen[v.ItemID] = true
	}
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := ensurePage(s.repo, ctx, tx, pageID); err != nil {
			return err
		}
		if _, err := s.repo.FindSectionForUpdate(ctx, tx, pageID, sectionID); errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrSectionNotFound
		} else if err != nil {
			return err
		}
		old, err := s.repo.FindCollectionItemsForUpdate(ctx, tx, sectionID, collection)
		if err != nil {
			return err
		}
		for _, v := range old {
			if v.ItemType != typ && v.DeletedAt.Valid == false {
				return ErrCollectionTypeConflict
			}
		}
		ids := make([]uint, 0, len(req.Items))
		for _, v := range req.Items {
			ids = append(ids, v.ItemID)
		}
		if typ == pageEntity.PageSectionItemTypeMedia {
			if err := validateMediaIDs(ctx, s.repo, tx, ids, userID, role); err != nil {
				return err
			}
		} else {
			valid, err := s.repo.FindSourceIDsForUpdate(ctx, tx, typ, ids)
			if err != nil {
				return err
			}
			if len(valid) != len(ids) {
				return ErrSourceNotFound
			}
		}
		if err := s.repo.SoftDeleteItemsNotIn(ctx, tx, sectionID, collection, typ, ids); err != nil {
			return err
		}
		for _, v := range req.Items {
			if err := s.repo.CreateOrRestoreItem(ctx, tx, &pageEntity.PageSectionItem{SectionID: sectionID, ItemType: typ, ItemID: v.ItemID, Collection: collection, SortOrder: v.SortOrder}); err != nil {
				if errors.Is(err, gorm.ErrDuplicatedKey) {
					return ErrSectionItemConflict
				}
				return err
			}
		}
		if typ == pageEntity.PageSectionItemTypeMedia {
			return s.repo.MarkMediaAttached(ctx, tx, ids)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	s.invalidate(ctx, pageID)
	return s.GetItems(ctx, pageID, sectionID, dto.GetSectionItemsQuery{Collection: collection, PageSize: 100})
}
func (s *pageSectionService) mapItems(ctx context.Context, rows []pageEntity.PageSectionItem) ([]dto.PageSectionItemResponse, error) {
	byType := map[pageEntity.PageSectionItemType][]uint{}
	for _, v := range rows {
		byType[v.ItemType] = append(byType[v.ItemType], v.ItemID)
	}
	data := map[pageEntity.PageSectionItemType]map[uint]repository.SourceData{}
	for typ, ids := range byType {
		v, err := s.repo.FindSourceData(ctx, typ, ids)
		if err != nil {
			return nil, err
		}
		data[typ] = v
	}
	out := make([]dto.PageSectionItemResponse, 0, len(rows))
	for _, v := range rows {
		out = append(out, dto.PageSectionItemResponse{ID: v.ID, SectionID: v.SectionID, ItemType: string(v.ItemType), ItemID: v.ItemID, Collection: v.Collection, SortOrder: v.SortOrder, Data: mapSource(data[v.ItemType][v.ItemID])})
	}
	return out, nil
}
func mapSection(v pageEntity.PageSection, sum []repository.SectionSummary) dto.PageSectionResponse {
	r := dto.PageSectionResponse{ID: v.ID, PageID: v.PageID, Key: v.Key, Name: v.Name, Title: v.Title, Description: v.Description, BackgroundColor: v.BackgroundColor, BackgroundMediaID: v.BackgroundMediaID, FeatureMedia: mapMedia(v.FeatureMedia), FeatureMediaID: v.FeatureMediaID, BackgroundMedia: mapMedia(v.BackgroundMedia), SortOrder: v.SortOrder, Status: string(v.Status), Collections: make([]dto.SectionCollectionSummary, 0, len(sum)), CreatedAt: v.CreatedAt, UpdatedAt: v.UpdatedAt}
	for _, x := range sum {
		r.Collections = append(r.Collections, dto.SectionCollectionSummary{Collection: x.Collection, ItemType: x.ItemType, Total: x.Total})
	}
	return r
}
func mapMedia(v *mediaEntity.Media) *dto.SectionMediaResponse {
	if v == nil {
		return nil
	}
	return &dto.SectionMediaResponse{ID: v.ID, FileName: v.FileName, OriginalURL: v.OriginalUrl, ThumbnailURL: v.ThumbnailUrl, MediumURL: v.MediumUrl, MimeType: v.MimeType}
}
func mapSource(v repository.SourceData) dto.SectionItemDataResponse {
	return dto.SectionItemDataResponse{ID: v.ID, Name: v.Name, Title: v.Title, Slug: v.Slug, TypeCode: v.TypeCode, ImageURL: v.ImageURL, Status: v.Status, FileName: v.FileName, OriginalURL: v.OriginalURL, ThumbnailURL: v.ThumbnailURL, MediumURL: v.MediumURL, MimeType: v.MimeType}
}
