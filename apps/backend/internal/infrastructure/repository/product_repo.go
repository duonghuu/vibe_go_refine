package repository

import (
	"context"
	"go_refine_dashboard_be/internal/domain/product/entity"
	"gorm.io/gorm"
)

type ProductRepository interface {
	FindAll(ctx context.Context, page, pageSize int, sort, order, search string, categoryID *uint, status string, minPrice, maxPrice *float64, minStock *int) ([]entity.Product, int64, error)
	Create(ctx context.Context, product *entity.Product) error
	Update(ctx context.Context, product *entity.Product) error
	FindByID(ctx context.Context, id uint) (*entity.Product, error)
	UpdateStatus(ctx context.Context, id uint, status string) error
	Delete(ctx context.Context, id uint) error
	BulkDelete(ctx context.Context, ids []uint) error
	BulkUpdateStatus(ctx context.Context, ids []uint, status string) error
	CountBySKU(ctx context.Context, sku string, excludeID uint) (int64, error)
	CountBySlug(ctx context.Context, slug string, excludeID uint) (int64, error)
	CategoryExists(ctx context.Context, categoryID uint) (bool, error)
}

type productRepositoryImpl struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepositoryImpl{db: db}
}

func (r *productRepositoryImpl) FindAll(ctx context.Context, page, pageSize int, sort, order, search string, categoryID *uint, status string, minPrice, maxPrice *float64, minStock *int) ([]entity.Product, int64, error) {
	var products []entity.Product
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.Product{}).Preload("Category")

	if search != "" {
		query = query.Where("name LIKE ? OR sku LIKE ?", "%"+search+"%", "%"+search+"%")
	}
	if categoryID != nil {
		query = query.Where("category_id = ?", *categoryID)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if minPrice != nil {
		query = query.Where("price >= ?", *minPrice)
	}
	if maxPrice != nil {
		query = query.Where("price <= ?", *maxPrice)
	}
	if minStock != nil {
		query = query.Where("stock >= ?", *minStock)
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	if sort != "" {
		if order == "desc" || order == "DESC" {
			query = query.Order(sort + " DESC")
		} else {
			query = query.Order(sort + " ASC")
		}
	} else {
		query = query.Order("created_at DESC") // Default sort
	}

	offset := page * pageSize
	err = query.Offset(offset).Limit(pageSize).Find(&products).Error
	return products, total, err
}

func (r *productRepositoryImpl) Create(ctx context.Context, product *entity.Product) error {
	return r.db.WithContext(ctx).Create(product).Error
}

func (r *productRepositoryImpl) Update(ctx context.Context, product *entity.Product) error {
	return r.db.WithContext(ctx).Save(product).Error
}

func (r *productRepositoryImpl) FindByID(ctx context.Context, id uint) (*entity.Product, error) {
	var product entity.Product
	err := r.db.WithContext(ctx).First(&product, id).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *productRepositoryImpl) UpdateStatus(ctx context.Context, id uint, status string) error {
	return r.db.WithContext(ctx).Model(&entity.Product{}).Where("id = ?", id).Update("status", status).Error
}

func (r *productRepositoryImpl) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&entity.Product{}, id).Error
}

func (r *productRepositoryImpl) BulkDelete(ctx context.Context, ids []uint) error {
	return r.db.WithContext(ctx).Where("id IN ?", ids).Delete(&entity.Product{}).Error
}

func (r *productRepositoryImpl) BulkUpdateStatus(ctx context.Context, ids []uint, status string) error {
	return r.db.WithContext(ctx).Model(&entity.Product{}).Where("id IN ?", ids).Update("status", status).Error
}

func (r *productRepositoryImpl) CountBySKU(ctx context.Context, sku string, excludeID uint) (int64, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&entity.Product{}).Where("sku = ?", sku)
	if excludeID > 0 {
		query = query.Where("id != ?", excludeID)
	}
	err := query.Count(&count).Error
	return count, err
}

func (r *productRepositoryImpl) CountBySlug(ctx context.Context, slug string, excludeID uint) (int64, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&entity.Product{}).Where("slug = ?", slug)
	if excludeID > 0 {
		query = query.Where("id != ?", excludeID)
	}
	err := query.Count(&count).Error
	return count, err
}

func (r *productRepositoryImpl) CategoryExists(ctx context.Context, categoryID uint) (bool, error) {
	var count int64
	// Because ProductRepo has db, we can check Category table directly to avoid wiring CategoryRepo
	err := r.db.WithContext(ctx).Table("categories").Where("id = ? AND deleted_at IS NULL", categoryID).Count(&count).Error
	return count > 0, err
}
