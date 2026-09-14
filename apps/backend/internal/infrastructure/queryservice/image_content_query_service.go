package queryservice

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	imageContentEntity "go_refine_dashboard_be/internal/domain/imagecontent/entity"
	"go_refine_dashboard_be/internal/infrastructure/repository"
	"go_refine_dashboard_be/internal/usecases/imagecontent/dto"
	"gorm.io/gorm"
)

var (
	ErrImageContentTypeNotFound = errors.New("image content type không tồn tại")
	ErrImageContentNotFound     = errors.New("image content không tồn tại")
	ErrPublicTypeCodesInvalid   = errors.New("typeCodes không hợp lệ")
)

const publicImageContentCacheKey = "public:image-contents:v1"

type ImageContentQueryService interface {
	ListTypes(ctx context.Context, query dto.GetImageContentTypesQuery) (*dto.ImageContentTypeListResponse, error)
	GetType(ctx context.Context, id uint) (*dto.ImageContentTypeResponse, error)
	ListContents(ctx context.Context, query dto.GetImageContentsQuery) (*dto.ImageContentListResponse, error)
	GetContent(ctx context.Context, id uint) (*dto.ImageContentResponse, error)
	GetPublicContents(ctx context.Context, query dto.GetPublicImageContentsQuery) (*dto.PublicImageContentResponse, error)
}

type imageContentQueryService struct {
	repo        repository.ImageContentRepository
	redisClient *redis.Client
}

func NewImageContentQueryService(repo repository.ImageContentRepository, redisClient *redis.Client) ImageContentQueryService {
	return &imageContentQueryService{repo: repo, redisClient: redisClient}
}

func (s *imageContentQueryService) ListTypes(ctx context.Context, query dto.GetImageContentTypesQuery) (*dto.ImageContentTypeListResponse, error) {
	start, end := normalizePagination(query.Start, query.End)
	values, counts, total, err := s.repo.FindTypes(ctx, repository.TypeListFilter{
		Offset: start,
		Limit:  end - start,
		Search: strings.TrimSpace(query.Search),
		Status: strings.TrimSpace(query.Status),
		SortBy: query.SortBy,
		Order:  query.Order,
	})
	if err != nil {
		return nil, fmt.Errorf("lấy danh sách image content type: %w", err)
	}
	data := make([]dto.ImageContentTypeResponse, 0, len(values))
	for _, value := range values {
		data = append(data, mapTypeResponse(value, counts[value.Code]))
	}
	return &dto.ImageContentTypeListResponse{Data: data, Total: total}, nil
}

func (s *imageContentQueryService) GetType(ctx context.Context, id uint) (*dto.ImageContentTypeResponse, error) {
	value, count, err := s.repo.FindTypeByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrImageContentTypeNotFound
		}
		return nil, fmt.Errorf("lấy image content type: %w", err)
	}
	response := mapTypeResponse(*value, count)
	return &response, nil
}

func (s *imageContentQueryService) ListContents(ctx context.Context, query dto.GetImageContentsQuery) (*dto.ImageContentListResponse, error) {
	start, end := normalizePagination(query.Start, query.End)
	values, total, err := s.repo.FindContents(ctx, repository.ContentListFilter{
		Offset:   start,
		Limit:    end - start,
		Search:   strings.TrimSpace(query.Search),
		TypeCode: strings.ToUpper(strings.TrimSpace(query.TypeCode)),
		Status:   strings.TrimSpace(query.Status),
		SortBy:   query.SortBy,
		Order:    query.Order,
	})
	if err != nil {
		return nil, fmt.Errorf("lấy danh sách image content: %w", err)
	}
	data := make([]dto.ImageContentResponse, 0, len(values))
	for _, value := range values {
		data = append(data, mapContentResponse(value))
	}
	return &dto.ImageContentListResponse{Data: data, Total: total}, nil
}

func (s *imageContentQueryService) GetContent(ctx context.Context, id uint) (*dto.ImageContentResponse, error) {
	value, err := s.repo.FindContentByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrImageContentNotFound
		}
		return nil, fmt.Errorf("lấy image content: %w", err)
	}
	response := mapContentResponse(*value)
	return &response, nil
}

func (s *imageContentQueryService) GetPublicContents(ctx context.Context, query dto.GetPublicImageContentsQuery) (*dto.PublicImageContentResponse, error) {
	codes, err := normalizePublicTypeCodes(query.TypeCodes)
	if err != nil {
		return nil, err
	}

	aggregate, err := s.getPublicAggregate(ctx)
	if err != nil {
		return nil, err
	}
	if len(codes) == 0 {
		return aggregate, nil
	}
	allowed := make(map[string]struct{}, len(codes))
	for _, code := range codes {
		allowed[code] = struct{}{}
	}
	filtered := make([]dto.PublicImageContentGroupResponse, 0, len(codes))
	for _, group := range aggregate.Data {
		if _, ok := allowed[group.TypeCode]; ok {
			filtered = append(filtered, group)
		}
	}
	return &dto.PublicImageContentResponse{Data: filtered}, nil
}

func (s *imageContentQueryService) getPublicAggregate(ctx context.Context) (*dto.PublicImageContentResponse, error) {
	if s.redisClient != nil {
		if payload, err := s.redisClient.Get(ctx, publicImageContentCacheKey).Result(); err == nil {
			var cached dto.PublicImageContentResponse
			if json.Unmarshal([]byte(payload), &cached) == nil && cached.Data != nil {
				return &cached, nil
			}
			log.Printf("image content public cache payload is invalid")
		} else if !errors.Is(err, redis.Nil) {
			log.Printf("image content public cache read failed: %v", err)
		}
	}

	rows, err := s.repo.FindPublicProjection(ctx)
	if err != nil {
		return nil, fmt.Errorf("lấy public image content projection: %w", err)
	}
	response := &dto.PublicImageContentResponse{Data: make([]dto.PublicImageContentGroupResponse, 0, len(rows))}
	for _, row := range rows {
		items := make([]dto.PublicImageContentItemResponse, 0, len(row.Items))
		for _, item := range row.Items {
			items = append(items, mapPublicItem(item, row.Type.FieldConfig))
		}
		response.Data = append(response.Data, dto.PublicImageContentGroupResponse{
			TypeCode: row.Type.Code,
			Name:     row.Type.Name,
			Items:    items,
		})
	}
	if s.redisClient != nil {
		if payload, marshalErr := json.Marshal(response); marshalErr == nil {
			if setErr := s.redisClient.SetEx(ctx, publicImageContentCacheKey, payload, 15*time.Minute).Err(); setErr != nil {
				log.Printf("image content public cache write failed: %v", setErr)
			}
		}
	}
	return response, nil
}

func normalizePagination(start, end int) (int, int) {
	if start < 0 {
		start = 0
	}
	if end <= 0 {
		end = start + 20
	}
	if end <= start {
		end = start + 20
	}
	if end-start > 100 {
		end = start + 100
	}
	return start, end
}

var publicTypeCodePattern = regexp.MustCompile(`^[A-Z][A-Z0-9_]{0,49}$`)

func normalizePublicTypeCodes(raw string) ([]string, error) {
	if strings.TrimSpace(raw) == "" {
		return nil, nil
	}
	seen := make(map[string]struct{})
	codes := make([]string, 0, 20)
	for _, part := range strings.Split(raw, ",") {
		code := strings.ToUpper(strings.TrimSpace(part))
		if code == "" || !publicTypeCodePattern.MatchString(code) {
			return nil, ErrPublicTypeCodesInvalid
		}
		if _, exists := seen[code]; exists {
			continue
		}
		seen[code] = struct{}{}
		codes = append(codes, code)
		if len(codes) > 20 {
			return nil, ErrPublicTypeCodesInvalid
		}
	}
	return codes, nil
}

func mapTypeResponse(value imageContentEntity.ImageContentType, itemCount int64) dto.ImageContentTypeResponse {
	return dto.ImageContentTypeResponse{
		ID:          value.ID,
		Code:        value.Code,
		Name:        value.Name,
		Status:      string(value.Status),
		SortOrder:   value.SortOrder,
		MaxItems:    value.MaxItems,
		FieldConfig: mapFieldConfigResponse(value.FieldConfig),
		ItemCount:   itemCount,
		CreatedAt:   value.CreatedAt,
		UpdatedAt:   value.UpdatedAt,
	}
}

func mapContentResponse(value imageContentEntity.ImageContent) dto.ImageContentResponse {
	response := dto.ImageContentResponse{
		ID:                   value.ID,
		TypeCode:             value.TypeCode,
		MediaID:              value.MediaID,
		Name:                 value.Name,
		Description:          value.Description,
		SecondaryDescription: value.SecondaryDescription,
		URL:                  value.TargetURL,
		SortOrder:            value.SortOrder,
		Status:               string(value.Status),
		CreatedBy:            value.CreatedBy,
		CreatedAt:            value.CreatedAt,
		UpdatedAt:            value.UpdatedAt,
	}
	if value.Type != nil {
		response.TypeName = value.Type.Name
		response.FieldConfig = mapFieldConfigResponse(value.Type.FieldConfig)
	}
	if value.Media != nil {
		response.Media = dto.AdminImageContentMediaResponse{
			ID:           value.Media.ID,
			FileName:     value.Media.FileName,
			OriginalURL:  value.Media.OriginalUrl,
			ThumbnailURL: value.Media.ThumbnailUrl,
			MediumURL:    value.Media.MediumUrl,
			MimeType:     value.Media.MimeType,
			Status:       string(value.Media.Status),
		}
	}
	return response
}

func mapPublicItem(value imageContentEntity.ImageContent, config imageContentEntity.ImageContentFieldConfig) dto.PublicImageContentItemResponse {
	item := dto.PublicImageContentItemResponse{
		ID:        value.ID,
		SortOrder: value.SortOrder,
	}
	if config.Name.Enabled {
		item.Name = value.Name
	}
	if config.Description.Enabled {
		item.Description = value.Description
	}
	if config.SecondaryDescription.Enabled {
		item.SecondaryDescription = value.SecondaryDescription
	}
	if config.URL.Enabled {
		item.URL = value.TargetURL
	}
	if value.Media != nil {
		item.Media = dto.PublicImageContentMediaResponse{
			ID:           value.Media.ID,
			OriginalURL:  value.Media.OriginalUrl,
			ThumbnailURL: value.Media.ThumbnailUrl,
			MediumURL:    value.Media.MediumUrl,
			MimeType:     value.Media.MimeType,
		}
	}
	return item
}

func mapFieldConfigResponse(config imageContentEntity.ImageContentFieldConfig) dto.ImageContentFieldConfigResponse {
	return dto.ImageContentFieldConfigResponse{
		Name:                 dto.ImageContentFieldRuleResponse{Enabled: config.Name.Enabled, Required: config.Name.Required},
		Description:          dto.ImageContentFieldRuleResponse{Enabled: config.Description.Enabled, Required: config.Description.Required},
		SecondaryDescription: dto.ImageContentFieldRuleResponse{Enabled: config.SecondaryDescription.Enabled, Required: config.SecondaryDescription.Required},
		URL:                  dto.ImageContentFieldRuleResponse{Enabled: config.URL.Enabled, Required: config.URL.Required},
	}
}
