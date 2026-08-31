package repository

import (
	"context"
	"errors"

	"go_refine_dashboard_be/internal/domain/seometa/entity"
	"gorm.io/gorm"
)

var ErrEntityNotFound = errors.New("entity not found")
var ErrEntityTypeNotSupported = errors.New("entity type not supported")
var ErrEntityTypeNotConfigured = errors.New("entity type not configured")

type EntityType string

const (
	EntityPost            EntityType = "post"
	EntityProduct         EntityType = "product"
	EntityPostCategory    EntityType = "post_category"
	EntityProductCategory EntityType = "product_category"
	EntityPage            EntityType = "page"
)

type EntityReference struct {
	Type         EntityType
	ID           uint
	Title        string
	Slug         string
	Description  string
	Excerpt      string
	Content      string
	ThumbnailURL string
}

type EntityRegistry interface {
	Resolve(ctx context.Context, entityType EntityType, entityID uint) (*EntityReference, error)
}

type entityRegistry struct{ db *gorm.DB }

func NewEntityRegistry(db *gorm.DB) EntityRegistry { return &entityRegistry{db: db} }

func (r *entityRegistry) Resolve(ctx context.Context, entityType EntityType, entityID uint) (*EntityReference, error) {
	if entityID == 0 {
		return nil, ErrEntityNotFound
	}
	ref := &EntityReference{Type: entityType, ID: entityID}
	var err error
	switch entityType {
	case EntityPost:
		var row struct {
			ID                   uint
			Title, Slug, Content string
		}
		err = r.db.WithContext(ctx).Table("posts").Select("id, title, slug, content").Where("id = ? AND deleted_at IS NULL", entityID).First(&row).Error
		ref.Title, ref.Slug, ref.Content = row.Title, row.Slug, row.Content
	case EntityProduct:
		var row struct {
			ID                                uint
			Name, Slug, Description, ImageURL string
		}
		err = r.db.WithContext(ctx).Table("products").Select("id, name, slug, description, image_url").Where("id = ? AND deleted_at IS NULL", entityID).First(&row).Error
		ref.Title, ref.Slug, ref.Description, ref.ThumbnailURL = row.Name, row.Slug, row.Description, row.ImageURL
	case EntityPostCategory:
		var row struct {
			ID                                uint
			Name, Slug, Description, ImageURL string
		}
		err = r.db.WithContext(ctx).Table("post_categories").Select("id, name, slug, description, image_url").Where("id = ? AND deleted_at IS NULL", entityID).First(&row).Error
		ref.Title, ref.Slug, ref.Description, ref.ThumbnailURL = row.Name, row.Slug, row.Description, row.ImageURL
	case EntityProductCategory:
		var row struct {
			ID                                uint
			Name, Slug, Description, ImageURL string
		}
		err = r.db.WithContext(ctx).Table("categories").Select("id, name, slug, description, image_url").Where("id = ? AND deleted_at IS NULL", entityID).First(&row).Error
		ref.Title, ref.Slug, ref.Description, ref.ThumbnailURL = row.Name, row.Slug, row.Description, row.ImageURL
	case EntityPage:
		var row struct {
			ID                   uint
			Title, Slug, Content string
		}
		err = r.db.WithContext(ctx).Table("pages").Select("id, title, slug, content").Where("id = ? AND deleted_at IS NULL", entityID).First(&row).Error
		ref.Title, ref.Slug, ref.Content = row.Title, row.Slug, row.Content
	default:
		return nil, ErrEntityTypeNotSupported
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrEntityNotFound
	}
	if err != nil {
		return nil, err
	}
	return ref, nil
}

type SEOMetaRepository interface {
	FindByEntity(ctx context.Context, entityType string, entityID uint) (*entity.SEOMeta, error)
	Upsert(ctx context.Context, seo *entity.SEOMeta) (*entity.SEOMeta, error)
	DeleteByEntity(ctx context.Context, entityType string, entityID uint) error
}

type seoMetaRepository struct{ db *gorm.DB }

func NewSEOMetaRepository(db *gorm.DB) SEOMetaRepository { return &seoMetaRepository{db: db} }

func (r *seoMetaRepository) FindByEntity(ctx context.Context, entityType string, entityID uint) (*entity.SEOMeta, error) {
	var item entity.SEOMeta
	err := r.db.WithContext(ctx).Where("entity_type = ? AND entity_id = ?", entityType, entityID).First(&item).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *seoMetaRepository) Upsert(ctx context.Context, seo *entity.SEOMeta) (*entity.SEOMeta, error) {
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var current entity.SEOMeta
		err := tx.Where("entity_type = ? AND entity_id = ?", seo.EntityType, seo.EntityID).First(&current).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return tx.Create(seo).Error
		}
		if err != nil {
			return err
		}
		seo.ID, seo.CreatedAt = current.ID, current.CreatedAt
		return tx.Model(&current).Updates(map[string]interface{}{
			"meta_title": seo.MetaTitle, "meta_description": seo.MetaDescription, "meta_keywords": seo.MetaKeywords,
			"canonical_url": seo.CanonicalURL, "og_title": seo.OGTitle, "og_description": seo.OGDescription,
			"og_image": seo.OGImage, "og_type": seo.OGType, "twitter_title": seo.TwitterTitle,
			"twitter_description": seo.TwitterDescription, "twitter_image": seo.TwitterImage,
			"twitter_card": seo.TwitterCard, "robots": seo.Robots, "schema_json": seo.SchemaJSON,
		}).Error
	})
	if err != nil {
		return nil, err
	}
	return r.FindByEntity(ctx, seo.EntityType, seo.EntityID)
}

func (r *seoMetaRepository) DeleteByEntity(ctx context.Context, entityType string, entityID uint) error {
	return r.db.WithContext(ctx).Where("entity_type = ? AND entity_id = ?", entityType, entityID).Delete(&entity.SEOMeta{}).Error
}
