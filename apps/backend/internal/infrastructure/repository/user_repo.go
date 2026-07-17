package repository

import (
	"context"

	"go_refine_dashboard_be/internal/domain/user/entity"

	"gorm.io/gorm"
)

type UserRepository interface {
	FindAll(ctx context.Context, skip, limit int, sortField, sortOrder, query, role, status string) ([]entity.User, int64, error)
	FindByID(ctx context.Context, id uint) (*entity.User, error)
	FindByEmail(ctx context.Context, email string) (*entity.User, error)
	Create(ctx context.Context, user *entity.User) error
	Update(ctx context.Context, user *entity.User) error
	BulkUpdateStatus(ctx context.Context, ids []uint, status string) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) FindAll(ctx context.Context, skip, limit int, sortField, sortOrder, query, role, status string) ([]entity.User, int64, error) {
	var users []entity.User
	var total int64

	q := r.db.WithContext(ctx).Model(&entity.User{})

	if query != "" {
		q = q.Where("name LIKE ? OR email LIKE ?", "%"+query+"%", "%"+query+"%")
	}
	if role != "" {
		q = q.Where("role = ?", role)
	}
	if status != "" {
		q = q.Where("status = ?", status)
	}

	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if sortField != "" {
		q = q.Order(sortField + " " + sortOrder)
	} else {
		q = q.Order("id DESC")
	}

	if limit > 0 {
		q = q.Offset(skip).Limit(limit)
	}

	if err := q.Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (r *userRepository) FindByID(ctx context.Context, id uint) (*entity.User, error) {
	var user entity.User
	if err := r.db.WithContext(ctx).First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	var user entity.User
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Create(ctx context.Context, user *entity.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *userRepository) Update(ctx context.Context, user *entity.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

func (r *userRepository) BulkUpdateStatus(ctx context.Context, ids []uint, status string) error {
	return r.db.WithContext(ctx).Model(&entity.User{}).Where("id IN ?", ids).Update("status", status).Error
}
