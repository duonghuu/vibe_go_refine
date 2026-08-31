package repository

import (
	"context"
	"strings"

	"go_refine_dashboard_be/internal/domain/page/entity"

	"gorm.io/gorm"
)

type PageRepository interface {
	Create(ctx context.Context, page *entity.Page) error
	CheckSlugExists(ctx context.Context, slug string) (bool, error)
	CheckSlugExistsExceptID(ctx context.Context, slug string, id uint) (bool, error)
	FindAndCount(ctx context.Context, filter PageListFilter) ([]entity.Page, int64, error)
	FindByID(ctx context.Context, id uint) (*entity.Page, error)
	Delete(ctx context.Context, id uint) error
	Update(ctx context.Context, page *entity.Page) error
}

type PageListFilter struct {
	TitleLike string
	Status    entity.PageStatus
	Offset    int
	Limit     int
	Sort      string
}

type pageRepository struct {
	db *gorm.DB
}

func NewPageRepository(db *gorm.DB) PageRepository {
	return &pageRepository{db: db}
}

func (r *pageRepository) Create(ctx context.Context, page *entity.Page) error {
	return r.db.WithContext(ctx).Create(page).Error
}

func (r *pageRepository) CheckSlugExists(ctx context.Context, slug string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).
		Model(&entity.Page{}).
		Where("slug = ?", slug).
		Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *pageRepository) CheckSlugExistsExceptID(ctx context.Context, slug string, id uint) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).
		Model(&entity.Page{}).
		Where("slug = ? AND id <> ?", slug, id).
		Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *pageRepository) FindAndCount(ctx context.Context, filter PageListFilter) ([]entity.Page, int64, error) {
	pages := make([]entity.Page, 0)
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.Page{})
	if filter.TitleLike != "" {
		escaped := strings.NewReplacer(`\\`, `\\\\`, `%`, `\\%`, `_`, `\\_`).Replace(filter.TitleLike)
		pattern := "%" + escaped + "%"
		query = query.Where("(title LIKE ? OR slug LIKE ?)", pattern, pattern)
	}
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if filter.Sort != "" {
		query = query.Order(filter.Sort)
	}
	if filter.Limit > 0 {
		query = query.Limit(filter.Limit).Offset(filter.Offset)
	}
	if err := query.Find(&pages).Error; err != nil {
		return nil, 0, err
	}

	return pages, total, nil
}

func (r *pageRepository) FindByID(ctx context.Context, id uint) (*entity.Page, error) {
	var page entity.Page
	if err := r.db.WithContext(ctx).First(&page, id).Error; err != nil {
		return nil, err
	}
	return &page, nil
}

func (r *pageRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&entity.Page{}, id).Error
}

func (r *pageRepository) Update(ctx context.Context, page *entity.Page) error {
	if err := r.db.WithContext(ctx).
		Model(&entity.Page{}).
		Where("id = ?", page.ID).
		Updates(map[string]interface{}{
			"title":   page.Title,
			"slug":    page.Slug,
			"content": page.Content,
			"status":  page.Status,
		}).Error; err != nil {
		return err
	}

	return r.db.WithContext(ctx).First(page, page.ID).Error
}
