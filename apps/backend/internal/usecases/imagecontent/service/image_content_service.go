package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/url"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/redis/go-redis/v9"
	imageContentEntity "go_refine_dashboard_be/internal/domain/imagecontent/entity"
	mediaEntity "go_refine_dashboard_be/internal/domain/media/entity"
	"go_refine_dashboard_be/internal/infrastructure/repository"
	"go_refine_dashboard_be/internal/usecases/imagecontent/dto"
	"gorm.io/gorm"
)

var (
	ErrValidation         = errors.New("dữ liệu không hợp lệ")
	ErrTypeConfigInvalid  = errors.New("cấu hình image content type không hợp lệ")
	ErrTypeCodeConflict   = errors.New("code image content type đã tồn tại")
	ErrTypeConfigConflict = errors.New("cấu hình image content type không tương thích")
	ErrTypeInUse          = errors.New("image content type đang được sử dụng")
	ErrMaxItemsExceeded   = errors.New("image content type đã đạt giới hạn item")
	ErrFieldRequired      = errors.New("field bắt buộc chưa được điền")
	ErrFieldNotEnabled    = errors.New("field chưa được bật trong type")
	ErrMediaNotFound      = errors.New("media không tồn tại")
	ErrMediaTypeInvalid   = errors.New("media phải là file hình ảnh")
	ErrMediaForbidden     = errors.New("không có quyền sử dụng media này")
	ErrSortOrderInvalid   = errors.New("sort order không hợp lệ")
	ErrOrderConflict      = errors.New("thứ tự image content đã thay đổi")
	ErrTypeNotFound       = errors.New("image content type không tồn tại")
	ErrContentNotFound    = errors.New("image content không tồn tại")
	ErrURLInvalid         = errors.New("url không hợp lệ")
	ErrForbidden          = errors.New("không có quyền thực hiện thao tác")
)

type ImageContentService interface {
	CreateType(ctx context.Context, req *dto.CreateImageContentTypeRequest) (uint, error)
	UpdateType(ctx context.Context, id uint, req *dto.UpdateImageContentTypeRequest) error
	DeleteType(ctx context.Context, id uint) error
	CreateContent(ctx context.Context, req *dto.CreateImageContentRequest, userID uint, role string) (uint, error)
	UpdateContent(ctx context.Context, id uint, req *dto.UpdateImageContentRequest, userID uint, role string) error
	DeleteContent(ctx context.Context, id uint) error
	ReorderContents(ctx context.Context, req *dto.ReorderImageContentsRequest) error
}

type imageContentService struct {
	repo        repository.ImageContentRepository
	db          *gorm.DB
	redisClient *redis.Client
}

func NewImageContentService(repo repository.ImageContentRepository, db *gorm.DB, redisClient *redis.Client) ImageContentService {
	return &imageContentService{repo: repo, db: db, redisClient: redisClient}
}

var typeCodePattern = regexp.MustCompile(`^[A-Z][A-Z0-9_]{0,49}$`)

func (s *imageContentService) CreateType(ctx context.Context, req *dto.CreateImageContentTypeRequest) (uint, error) {
	value, err := buildTypeFromCreateRequest(req)
	if err != nil {
		return 0, err
	}
	var id uint
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		existing, findErr := s.repo.FindTypeByCodeForUpdate(ctx, tx, value.Code)
		if findErr == nil {
			if !existing.DeletedAt.Valid {
				return ErrTypeCodeConflict
			}
			value.ID = existing.ID
			value.DeletedAt = existing.DeletedAt
		}
		if findErr != nil && !errors.Is(findErr, gorm.ErrRecordNotFound) {
			return findErr
		}
		if err := s.repo.CreateOrRestoreType(ctx, tx, &value); err != nil {
			return err
		}
		id = value.ID
		return nil
	})
	if err != nil {
		return 0, fmt.Errorf("tạo image content type: %w", err)
	}
	s.invalidatePublicCache(ctx)
	return id, nil
}

func (s *imageContentService) UpdateType(ctx context.Context, id uint, req *dto.UpdateImageContentTypeRequest) error {
	input, err := buildTypeUpdateInput(req)
	if err != nil {
		return err
	}
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		value, findErr := s.repo.FindTypeForUpdate(ctx, tx, id)
		if errors.Is(findErr, gorm.ErrRecordNotFound) {
			return ErrTypeNotFound
		}
		if findErr != nil {
			return findErr
		}
		items, findErr := s.repo.FindContentsByTypeForUpdate(ctx, tx, value.Code)
		if findErr != nil {
			return findErr
		}
		if input.maxItems != nil && uint(len(items)) > *input.maxItems {
			return ErrTypeConfigConflict
		}
		if err := validateRequiredFields(input.fieldConfig, items); err != nil {
			return err
		}
		value.Name = input.name
		value.Status = input.status
		value.SortOrder = input.sortOrder
		value.MaxItems = input.maxItems
		value.FieldConfig = input.fieldConfig
		return s.repo.UpdateType(ctx, tx, value)
	})
	if err != nil {
		return fmt.Errorf("cập nhật image content type: %w", err)
	}
	s.invalidatePublicCache(ctx)
	return nil
}

func (s *imageContentService) DeleteType(ctx context.Context, id uint) error {
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		value, findErr := s.repo.FindTypeForUpdate(ctx, tx, id)
		if errors.Is(findErr, gorm.ErrRecordNotFound) {
			return ErrTypeNotFound
		}
		if findErr != nil {
			return findErr
		}
		items, findErr := s.repo.FindContentsByTypeForUpdate(ctx, tx, value.Code)
		if findErr != nil {
			return findErr
		}
		if len(items) > 0 {
			return ErrTypeInUse
		}
		return s.repo.SoftDeleteType(ctx, tx, id)
	})
	if err != nil {
		return fmt.Errorf("xóa image content type: %w", err)
	}
	s.invalidatePublicCache(ctx)
	return nil
}

func (s *imageContentService) CreateContent(ctx context.Context, req *dto.CreateImageContentRequest, userID uint, role string) (uint, error) {
	input, err := buildContentInput(req)
	if err != nil {
		return 0, err
	}
	if err := validateRole(role); err != nil {
		return 0, err
	}
	var id uint
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		typeValue, findErr := s.repo.FindTypeByCodeForUpdate(ctx, tx, input.typeCode)
		if errors.Is(findErr, gorm.ErrRecordNotFound) || (findErr == nil && typeValue.DeletedAt.Valid) {
			return ErrTypeNotFound
		}
		if findErr != nil {
			return findErr
		}
		items, findErr := s.repo.FindContentsByTypeForUpdate(ctx, tx, input.typeCode)
		if findErr != nil {
			return findErr
		}
		if typeValue.MaxItems != nil && uint(len(items)) >= *typeValue.MaxItems {
			return ErrMaxItemsExceeded
		}
		if err := validateContentMetadata(typeValue.FieldConfig, input.name, input.description, input.secondaryDescription, input.targetURL); err != nil {
			return err
		}
		media, findErr := s.repo.FindMediaForUpdate(ctx, tx, input.mediaID)
		if errors.Is(findErr, gorm.ErrRecordNotFound) {
			return ErrMediaNotFound
		}
		if findErr != nil {
			return findErr
		}
		if err := validateMedia(*media, userID, role); err != nil {
			return err
		}
		sortOrder := len(items)
		if input.sortOrder != nil {
			if *input.sortOrder > len(items) {
				return ErrSortOrderInvalid
			}
			sortOrder = *input.sortOrder
			if err := s.repo.ShiftContentOrders(ctx, tx, input.typeCode, sortOrder); err != nil {
				return err
			}
		}
		value := imageContentEntity.ImageContent{
			TypeCode:             input.typeCode,
			MediaID:              input.mediaID,
			Name:                 input.name,
			Description:          input.description,
			SecondaryDescription: input.secondaryDescription,
			TargetURL:            input.targetURL,
			SortOrder:            sortOrder,
			Status:               input.status,
			CreatedBy:            userID,
		}
		if err := s.repo.CreateContent(ctx, tx, &value); err != nil {
			return err
		}
		id = value.ID
		return s.repo.MarkMediaAttached(ctx, tx, []uint{input.mediaID})
	})
	if err != nil {
		return 0, fmt.Errorf("tạo image content: %w", err)
	}
	s.invalidatePublicCache(ctx)
	return id, nil
}

func (s *imageContentService) UpdateContent(ctx context.Context, id uint, req *dto.UpdateImageContentRequest, userID uint, role string) error {
	input, err := buildContentUpdateInput(req)
	if err != nil {
		return err
	}
	if err := validateRole(role); err != nil {
		return err
	}
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		typeCode, findErr := s.repo.FindContentTypeCode(ctx, tx, id)
		if errors.Is(findErr, gorm.ErrRecordNotFound) {
			return ErrContentNotFound
		}
		if findErr != nil {
			return findErr
		}
		typeValue, findErr := s.repo.FindTypeByCodeForUpdate(ctx, tx, typeCode)
		if errors.Is(findErr, gorm.ErrRecordNotFound) || (findErr == nil && typeValue.DeletedAt.Valid) {
			return ErrTypeNotFound
		}
		if findErr != nil {
			return findErr
		}
		value, findErr := s.repo.FindContentForUpdate(ctx, tx, id)
		if errors.Is(findErr, gorm.ErrRecordNotFound) {
			return ErrContentNotFound
		}
		if findErr != nil {
			return findErr
		}
		if err := validateContentMetadata(typeValue.FieldConfig, input.name, input.description, input.secondaryDescription, input.targetURL); err != nil {
			return err
		}
		if input.mediaID != value.MediaID {
			media, mediaErr := s.repo.FindMediaForUpdate(ctx, tx, input.mediaID)
			if errors.Is(mediaErr, gorm.ErrRecordNotFound) {
				return ErrMediaNotFound
			}
			if mediaErr != nil {
				return mediaErr
			}
			if err := validateMedia(*media, userID, role); err != nil {
				return err
			}
		}
		mediaChanged := input.mediaID != value.MediaID
		value.MediaID = input.mediaID
		value.Name = input.name
		value.Description = input.description
		value.SecondaryDescription = input.secondaryDescription
		value.TargetURL = input.targetURL
		value.Status = input.status
		if err := s.repo.UpdateContent(ctx, tx, value); err != nil {
			return err
		}
		if mediaChanged {
			return s.repo.MarkMediaAttached(ctx, tx, []uint{input.mediaID})
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("cập nhật image content: %w", err)
	}
	s.invalidatePublicCache(ctx)
	return nil
}

func (s *imageContentService) DeleteContent(ctx context.Context, id uint) error {
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		typeCode, findErr := s.repo.FindContentTypeCode(ctx, tx, id)
		if errors.Is(findErr, gorm.ErrRecordNotFound) {
			return ErrContentNotFound
		}
		if findErr != nil {
			return findErr
		}
		typeValue, findErr := s.repo.FindTypeByCodeForUpdate(ctx, tx, typeCode)
		if errors.Is(findErr, gorm.ErrRecordNotFound) || (findErr == nil && typeValue.DeletedAt.Valid) {
			return ErrTypeNotFound
		} else if findErr != nil {
			return findErr
		}
		value, findErr := s.repo.FindContentForUpdate(ctx, tx, id)
		if errors.Is(findErr, gorm.ErrRecordNotFound) {
			return ErrContentNotFound
		}
		if findErr != nil {
			return findErr
		}
		items, findErr := s.repo.FindContentsByTypeForUpdate(ctx, tx, typeCode)
		if findErr != nil {
			return findErr
		}
		if err := s.repo.SoftDeleteContent(ctx, tx, value.ID); err != nil {
			return err
		}
		remaining := make([]imageContentEntity.ImageContent, 0, len(items)-1)
		for _, item := range items {
			if item.ID != value.ID {
				remaining = append(remaining, item)
			}
		}
		return s.repo.NormalizeContentOrders(ctx, tx, remaining)
	})
	if err != nil {
		return fmt.Errorf("xóa image content: %w", err)
	}
	s.invalidatePublicCache(ctx)
	return nil
}

func (s *imageContentService) ReorderContents(ctx context.Context, req *dto.ReorderImageContentsRequest) error {
	if req == nil {
		return ErrValidation
	}
	typeCode := strings.ToUpper(strings.TrimSpace(req.TypeCode))
	if !typeCodePattern.MatchString(typeCode) {
		return ErrValidation
	}
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		typeValue, findErr := s.repo.FindTypeByCodeForUpdate(ctx, tx, typeCode)
		if errors.Is(findErr, gorm.ErrRecordNotFound) || (findErr == nil && typeValue.DeletedAt.Valid) {
			return ErrTypeNotFound
		} else if findErr != nil {
			return findErr
		}
		items, findErr := s.repo.FindContentsByTypeForUpdate(ctx, tx, typeCode)
		if findErr != nil {
			return findErr
		}
		if len(items) != len(req.Items) {
			return ErrOrderConflict
		}
		known := make(map[uint]struct{}, len(items))
		for _, item := range items {
			known[item.ID] = struct{}{}
		}
		seen := make(map[uint]struct{}, len(req.Items))
		updates := make([]repository.ImageContentOrderUpdate, 0, len(req.Items))
		for _, requested := range req.Items {
			if requested.SortOrder < 0 || requested.SortOrder >= len(items) {
				return ErrOrderConflict
			}
			if _, ok := known[requested.ID]; !ok {
				return ErrOrderConflict
			}
			if _, ok := seen[requested.ID]; ok {
				return ErrOrderConflict
			}
			seen[requested.ID] = struct{}{}
			updates = append(updates, repository.ImageContentOrderUpdate{ID: requested.ID, SortOrder: requested.SortOrder})
		}
		if len(seen) != len(known) {
			return ErrOrderConflict
		}
		orders := make(map[int]struct{}, len(updates))
		for _, update := range updates {
			orders[update.SortOrder] = struct{}{}
		}
		if len(orders) != len(items) {
			return ErrOrderConflict
		}
		return s.repo.BulkUpdateContentOrder(ctx, tx, typeCode, updates)
	})
	if err != nil {
		return fmt.Errorf("sắp xếp image content: %w", err)
	}
	s.invalidatePublicCache(ctx)
	return nil
}

type typeInput struct {
	code        string
	name        string
	status      imageContentEntity.ImageContentStatus
	sortOrder   int
	maxItems    *uint
	fieldConfig imageContentEntity.ImageContentFieldConfig
}

func buildTypeFromCreateRequest(req *dto.CreateImageContentTypeRequest) (imageContentEntity.ImageContentType, error) {
	if req == nil || req.FieldConfig == nil {
		return imageContentEntity.ImageContentType{}, ErrValidation
	}
	input, err := buildTypeInput(req.Code, req.Name, req.Status, req.SortOrder, req.MaxItems, req.FieldConfig)
	if err != nil {
		return imageContentEntity.ImageContentType{}, err
	}
	return imageContentEntity.ImageContentType{Code: input.code, Name: input.name, Status: input.status, SortOrder: input.sortOrder, MaxItems: input.maxItems, FieldConfig: input.fieldConfig}, nil
}

func buildTypeUpdateInput(req *dto.UpdateImageContentTypeRequest) (typeInput, error) {
	if req == nil || req.FieldConfig == nil {
		return typeInput{}, ErrValidation
	}
	return buildTypeInput("UNCHANGED", req.Name, req.Status, req.SortOrder, req.MaxItems, req.FieldConfig)
}

func buildTypeInput(code, name, status string, sortOrder int, maxItems *uint, config *dto.ImageContentFieldConfigRequest) (typeInput, error) {
	normalizedCode := strings.ToUpper(strings.TrimSpace(code))
	if normalizedCode != "UNCHANGED" && !typeCodePattern.MatchString(normalizedCode) {
		return typeInput{}, ErrTypeConfigInvalid
	}
	normalizedName := strings.TrimSpace(name)
	if normalizedName == "" || utf8.RuneCountInString(normalizedName) > 100 || sortOrder < 0 {
		return typeInput{}, ErrTypeConfigInvalid
	}
	if status != string(imageContentEntity.ImageContentStatusActive) && status != string(imageContentEntity.ImageContentStatusInactive) {
		return typeInput{}, ErrTypeConfigInvalid
	}
	if maxItems != nil && (*maxItems < 1 || *maxItems > 1000) {
		return typeInput{}, ErrTypeConfigInvalid
	}
	fieldConfig := imageContentEntity.ImageContentFieldConfig{
		Name:                 imageContentEntity.ImageContentFieldRule{Enabled: config.Name.Enabled, Required: config.Name.Required},
		Description:          imageContentEntity.ImageContentFieldRule{Enabled: config.Description.Enabled, Required: config.Description.Required},
		SecondaryDescription: imageContentEntity.ImageContentFieldRule{Enabled: config.SecondaryDescription.Enabled, Required: config.SecondaryDescription.Required},
		URL:                  imageContentEntity.ImageContentFieldRule{Enabled: config.URL.Enabled, Required: config.URL.Required},
	}
	if err := validateFieldConfig(fieldConfig); err != nil {
		return typeInput{}, err
	}
	return typeInput{code: normalizedCode, name: normalizedName, status: imageContentEntity.ImageContentStatus(status), sortOrder: sortOrder, maxItems: maxItems, fieldConfig: fieldConfig}, nil
}

type contentInput struct {
	typeCode             string
	mediaID              uint
	name                 *string
	description          *string
	secondaryDescription *string
	targetURL            *string
	sortOrder            *int
	status               imageContentEntity.ImageContentStatus
}

func buildContentInput(req *dto.CreateImageContentRequest) (contentInput, error) {
	if req == nil {
		return contentInput{}, ErrValidation
	}
	input := contentInput{
		typeCode:             strings.ToUpper(strings.TrimSpace(req.TypeCode)),
		mediaID:              req.MediaID,
		name:                 normalizeString(req.Name),
		description:          normalizeString(req.Description),
		secondaryDescription: normalizeString(req.SecondaryDescription),
		targetURL:            normalizeString(req.URL),
		sortOrder:            req.SortOrder,
		status:               imageContentEntity.ImageContentStatus(req.Status),
	}
	if !typeCodePattern.MatchString(input.typeCode) || input.mediaID == 0 || (input.status != imageContentEntity.ImageContentStatusActive && input.status != imageContentEntity.ImageContentStatusInactive) {
		return contentInput{}, ErrValidation
	}
	if err := validateContentLengths(input); err != nil {
		return contentInput{}, err
	}
	return input, nil
}

func buildContentUpdateInput(req *dto.UpdateImageContentRequest) (contentInput, error) {
	if req == nil {
		return contentInput{}, ErrValidation
	}
	input := contentInput{
		mediaID:              req.MediaID,
		name:                 normalizeString(req.Name),
		description:          normalizeString(req.Description),
		secondaryDescription: normalizeString(req.SecondaryDescription),
		targetURL:            normalizeString(req.URL),
		status:               imageContentEntity.ImageContentStatus(req.Status),
	}
	if input.mediaID == 0 || (input.status != imageContentEntity.ImageContentStatusActive && input.status != imageContentEntity.ImageContentStatusInactive) {
		return contentInput{}, ErrValidation
	}
	if err := validateContentLengths(input); err != nil {
		return contentInput{}, err
	}
	return input, nil
}

func normalizeString(value *string) *string {
	if value == nil {
		return nil
	}
	normalized := strings.TrimSpace(*value)
	if normalized == "" {
		return nil
	}
	return &normalized
}

func validateContentLengths(input contentInput) error {
	if input.name != nil && utf8.RuneCountInString(*input.name) > 255 {
		return ErrValidation
	}
	if input.description != nil && utf8.RuneCountInString(*input.description) > 5000 {
		return ErrValidation
	}
	if input.secondaryDescription != nil && utf8.RuneCountInString(*input.secondaryDescription) > 5000 {
		return ErrValidation
	}
	if input.targetURL != nil {
		if utf8.RuneCountInString(*input.targetURL) > 500 {
			return ErrValidation
		}
		if err := validateURL(*input.targetURL); err != nil {
			return err
		}
	}
	return nil
}

func validateFieldConfig(config imageContentEntity.ImageContentFieldConfig) error {
	rules := []imageContentEntity.ImageContentFieldRule{config.Name, config.Description, config.SecondaryDescription, config.URL}
	for _, rule := range rules {
		if rule.Required && !rule.Enabled {
			return ErrTypeConfigInvalid
		}
	}
	return nil
}

func validateRequiredFields(config imageContentEntity.ImageContentFieldConfig, items []imageContentEntity.ImageContent) error {
	for _, item := range items {
		if config.Name.Required && isEmpty(item.Name) || config.Description.Required && isEmpty(item.Description) || config.SecondaryDescription.Required && isEmpty(item.SecondaryDescription) || config.URL.Required && isEmpty(item.TargetURL) {
			return ErrTypeConfigConflict
		}
	}
	return nil
}

func validateContentMetadata(config imageContentEntity.ImageContentFieldConfig, name, description, secondaryDescription, targetURL *string) error {
	if err := validateFieldConfig(config); err != nil {
		return err
	}
	values := []struct {
		name  string
		rule  imageContentEntity.ImageContentFieldRule
		value *string
	}{
		{name: "name", rule: config.Name, value: name},
		{name: "description", rule: config.Description, value: description},
		{name: "secondaryDescription", rule: config.SecondaryDescription, value: secondaryDescription},
		{name: "url", rule: config.URL, value: targetURL},
	}
	for _, field := range values {
		if !field.rule.Enabled && !isEmpty(field.value) {
			return fmt.Errorf("%w: %s", ErrFieldNotEnabled, field.name)
		}
		if field.rule.Required && isEmpty(field.value) {
			return fmt.Errorf("%w: %s", ErrFieldRequired, field.name)
		}
	}
	return nil
}

func isEmpty(value *string) bool {
	return value == nil || strings.TrimSpace(*value) == ""
}

func validateURL(value string) error {
	for _, character := range value {
		if character < 0x20 || character == 0x7f {
			return ErrURLInvalid
		}
	}
	if strings.HasPrefix(value, "//") {
		return ErrURLInvalid
	}
	parsed, err := url.Parse(value)
	if err != nil {
		return ErrURLInvalid
	}
	if parsed.IsAbs() {
		scheme := strings.ToLower(parsed.Scheme)
		if (scheme != "http" && scheme != "https") || parsed.Host == "" {
			return ErrURLInvalid
		}
		return nil
	}
	if !strings.HasPrefix(value, "/") {
		return ErrURLInvalid
	}
	return nil
}

func validateMedia(media mediaEntity.Media, userID uint, role string) error {
	if !strings.HasPrefix(strings.ToLower(media.MimeType), "image/") {
		return ErrMediaTypeInvalid
	}
	if role == "STAFF" && media.OwnerID != userID {
		return ErrMediaForbidden
	}
	return nil
}

func validateRole(role string) error {
	if role != "ADMIN" && role != "STAFF" {
		return ErrForbidden
	}
	return nil
}

func (s *imageContentService) invalidatePublicCache(ctx context.Context) {
	if s.redisClient == nil {
		return
	}
	if err := s.redisClient.Del(ctx, "public:image-contents:v1").Err(); err != nil {
		log.Printf("image content public cache invalidation failed: %v", err)
	}
}
