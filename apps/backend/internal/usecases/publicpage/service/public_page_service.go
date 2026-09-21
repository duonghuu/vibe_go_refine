package service

import (
	"context"
	"encoding/json"
	"errors"
	"html"
	"log"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/redis/go-redis/v9"
	seoEntity "go_refine_dashboard_be/internal/domain/seometa/entity"
	"go_refine_dashboard_be/internal/infrastructure/cache"
	"go_refine_dashboard_be/internal/infrastructure/repository"
	"go_refine_dashboard_be/internal/usecases/publicpage/dto"
	"gorm.io/gorm"
)

var (
	ErrInvalidSlug        = errors.New("slug trang không hợp lệ")
	ErrPublicPageNotFound = errors.New("không tìm thấy trang")
)

var publicPageSlugPattern = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)
var htmlTagPattern = regexp.MustCompile(`<[^>]*>`)

type PublicPageService interface {
	GetBySlug(ctx context.Context, slug string) (*dto.GetPublicPageResponse, error)
}

type publicPageService struct {
	repo          repository.PublicPageRepository
	redisClient   *redis.Client
	publicSiteURL string
}

func NewPublicPageService(repo repository.PublicPageRepository, redisClient *redis.Client, publicSiteURL string) PublicPageService {
	return &publicPageService{repo: repo, redisClient: redisClient, publicSiteURL: strings.TrimRight(strings.TrimSpace(publicSiteURL), "/")}
}

func (s *publicPageService) GetBySlug(ctx context.Context, slug string) (*dto.GetPublicPageResponse, error) {
	slug = strings.TrimSpace(slug)
	if !publicPageSlugPattern.MatchString(slug) {
		return nil, ErrInvalidSlug
	}

	cacheKey := cache.PublicPageCacheKey(slug)
	if s.redisClient != nil {
		cached, err := s.redisClient.Get(ctx, cacheKey).Bytes()
		if err == nil {
			var response dto.GetPublicPageResponse
			if json.Unmarshal(cached, &response) == nil {
				return &response, nil
			}
			log.Printf("public page cache decode failed for key %q", cacheKey)
		} else if !errors.Is(err, redis.Nil) {
			log.Printf("public page cache read failed for key %q: %v", cacheKey, err)
		}
	}

	projection, err := s.repo.FindPublishedBySlug(ctx, slug)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrPublicPageNotFound
		}
		return nil, err
	}
	response := mapPublicPage(projection, s.publicSiteURL)
	if s.redisClient != nil {
		if payload, marshalErr := json.Marshal(response); marshalErr == nil {
			if err := s.redisClient.SetEx(ctx, cacheKey, payload, 15*time.Minute).Err(); err != nil {
				log.Printf("public page cache write failed for key %q: %v", cacheKey, err)
			}
		}
	}
	return response, nil
}

func mapPublicPage(projection *repository.PublicPageProjection, publicSiteURL string) *dto.GetPublicPageResponse {
	result := &dto.GetPublicPageResponse{Data: dto.PublicPageResponse{
		ID: projection.ID, Title: projection.Title, Slug: projection.Slug, Content: projection.Content,
		UpdatedAt: projection.UpdatedAt, Gallery: make([]dto.PublicPageMediaResponse, 0), Sections: make([]dto.PublicPageSectionResponse, 0),
		SEO: resolveSEO(projection, publicSiteURL),
	}}
	if projection.Thumbnail != nil {
		thumbnail := mapPublicMedia(*projection.Thumbnail)
		result.Data.Thumbnail = &thumbnail
	}
	for _, media := range projection.Gallery {
		result.Data.Gallery = append(result.Data.Gallery, mapPublicMedia(media))
	}
	for _, section := range projection.Sections {
		mapped := dto.PublicPageSectionResponse{
			ID: section.ID, Key: section.Key, Title: section.Title, Description: section.Description,
			BackgroundColor: section.BackgroundColor, SortOrder: section.SortOrder,
			Collections: make(map[string][]dto.PublicPageSectionItemResponse, len(section.Collections)),
		}
		if section.BackgroundMedia != nil {
			media := mapPublicSectionMedia(*section.BackgroundMedia)
			mapped.BackgroundMedia = &media
		}
		if section.FeatureMedia != nil {
			media := mapPublicSectionMedia(*section.FeatureMedia)
			mapped.FeatureMedia = &media
		}
		for collection, items := range section.Collections {
			mappedItems := make([]dto.PublicPageSectionItemResponse, 0, len(items))
			for _, item := range items {
				mappedItems = append(mappedItems, dto.PublicPageSectionItemResponse{ItemType: item.ItemType, ItemID: item.ItemID, SortOrder: item.SortOrder, Data: dto.PublicSectionItemDataResponse{
					ID: item.Data.ID, Name: item.Data.Name, Description: item.Data.Description, Title: item.Data.Title, Slug: item.Data.Slug, TypeCode: item.Data.TypeCode,
					ImageURL: item.Data.ImageURL, OriginalURL: item.Data.OriginalURL, ThumbnailURL: item.Data.ThumbnailURL, MediumURL: item.Data.MediumURL,
				}})
			}
			mapped.Collections[collection] = mappedItems
		}
		result.Data.Sections = append(result.Data.Sections, mapped)
	}
	return result
}

func mapPublicMedia(media repository.PublicPageMediaProjection) dto.PublicPageMediaResponse {
	return dto.PublicPageMediaResponse{ID: media.ID, OriginalURL: media.OriginalURL, ThumbnailURL: media.ThumbnailURL, MediumURL: media.MediumURL}
}

func mapPublicSectionMedia(media repository.PublicSectionMediaProjection) dto.PublicPageSectionMediaResponse {
	return dto.PublicPageSectionMediaResponse{ID: media.ID, OriginalURL: media.OriginalURL, ThumbnailURL: media.ThumbnailURL, MediumURL: media.MediumURL}
}

func resolveSEO(projection *repository.PublicPageProjection, publicSiteURL string) dto.PublicPageSEOResponse {
	seo := dto.PublicPageSEOResponse{ResolvedTitle: projection.Title, ResolvedDescription: stripHTML(projection.Content), Robots: "index,follow"}
	if publicSiteURL == "" {
		seo.ResolvedCanonicalURL = projection.Slug
	} else {
		seo.ResolvedCanonicalURL = publicSiteURL + "/" + projection.Slug
	}
	if projection.Thumbnail != nil {
		seo.OGImage = projection.Thumbnail.OriginalURL
		seo.TwitterImage = projection.Thumbnail.OriginalURL
	}
	if projection.SEO == nil {
		seo.OGTitle, seo.OGDescription, seo.TwitterTitle, seo.TwitterDescription = seo.ResolvedTitle, seo.ResolvedDescription, seo.ResolvedTitle, seo.ResolvedDescription
		return seo
	}
	applySEO(&seo, projection.SEO)
	return seo
}

func applySEO(result *dto.PublicPageSEOResponse, meta *seoEntity.SEOMeta) {
	if meta.MetaTitle != nil {
		result.ResolvedTitle = *meta.MetaTitle
	}
	if meta.MetaDescription != nil {
		result.ResolvedDescription = *meta.MetaDescription
	}
	if meta.CanonicalURL != nil {
		result.ResolvedCanonicalURL = *meta.CanonicalURL
	}
	result.OGTitle = seoValue(meta.OGTitle, result.ResolvedTitle)
	result.OGDescription = seoValue(meta.OGDescription, result.ResolvedDescription)
	result.OGImage = seoValue(meta.OGImage, result.OGImage)
	result.OGType = seoValue(meta.OGType, "website")
	result.TwitterTitle = seoValue(meta.TwitterTitle, result.OGTitle)
	result.TwitterDescription = seoValue(meta.TwitterDescription, result.OGDescription)
	result.TwitterImage = seoValue(meta.TwitterImage, result.OGImage)
	result.TwitterCard = seoValue(meta.TwitterCard, "summary_large_image")
	if meta.Robots != nil && strings.TrimSpace(*meta.Robots) != "" {
		result.Robots = strings.TrimSpace(*meta.Robots)
	}
	if len(meta.SchemaJSON) > 0 && json.Valid(meta.SchemaJSON) {
		result.SchemaJSON = json.RawMessage(meta.SchemaJSON)
	}
}

func seoValue(value *string, fallback string) string {
	if value == nil {
		return fallback
	}
	return strings.TrimSpace(*value)
}

func stripHTML(value string) string {
	clean := html.UnescapeString(htmlTagPattern.ReplaceAllString(value, " "))
	clean = strings.Join(strings.Fields(clean), " ")
	if utf8.RuneCountInString(clean) > 500 {
		runes := []rune(clean)
		return string(runes[:500])
	}
	return clean
}
