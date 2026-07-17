package service

import (
	"context"
	"errors"

	"go_refine_dashboard_be/internal/domain/user/entity"
	"go_refine_dashboard_be/internal/infrastructure/repository"
	"go_refine_dashboard_be/internal/usecases/user/dto"

	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	GetUsers(ctx context.Context, skip, limit int, sortField, sortOrder, query, role, status string) ([]entity.User, int64, error)
	GetUserByID(ctx context.Context, id uint) (*entity.User, error)
	CreateUser(ctx context.Context, req *dto.CreateUserRequest) (*entity.User, error)
	UpdateUser(ctx context.Context, id uint, req *dto.UpdateUserRequest) (*entity.User, error)
	UpdateStatus(ctx context.Context, id uint, req *dto.UpdateUserStatusRequest) error
	BulkUpdateStatus(ctx context.Context, req *dto.BulkUpdateUserStatusRequest) error
	ResetPassword(ctx context.Context, id uint, req *dto.ResetPasswordRequest) error
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (s *userService) GetUsers(ctx context.Context, skip, limit int, sortField, sortOrder, query, role, status string) ([]entity.User, int64, error) {
	return s.userRepo.FindAll(ctx, skip, limit, sortField, sortOrder, query, role, status)
}

func (s *userService) GetUserByID(ctx context.Context, id uint) (*entity.User, error) {
	return s.userRepo.FindByID(ctx, id)
}

func (s *userService) CreateUser(ctx context.Context, req *dto.CreateUserRequest) (*entity.User, error) {
	existing, _ := s.userRepo.FindByEmail(ctx, req.Email)
	if existing != nil {
		return nil, errors.New("email đã tồn tại")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), 12)
	if err != nil {
		return nil, err
	}

	user := &entity.User{
		Email:    req.Email,
		Password: string(hashedPassword),
		Name:     req.Name,
		Role:     req.Role,
		Status:   "ACTIVE",
	}
	if req.Status != "" {
		user.Status = req.Status
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) UpdateUser(ctx context.Context, id uint, req *dto.UpdateUserRequest) (*entity.User, error) {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.New("người dùng không tồn tại")
	}

	if req.Name != nil {
		user.Name = *req.Name
	}
	if req.Role != nil {
		user.Role = *req.Role
	}
	if req.Status != nil {
		user.Status = *req.Status
	}

	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) UpdateStatus(ctx context.Context, id uint, req *dto.UpdateUserStatusRequest) error {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return errors.New("người dùng không tồn tại")
	}

	user.Status = req.Status
	return s.userRepo.Update(ctx, user)
}

func (s *userService) BulkUpdateStatus(ctx context.Context, req *dto.BulkUpdateUserStatusRequest) error {
	return s.userRepo.BulkUpdateStatus(ctx, req.IDs, req.Status)
}

func (s *userService) ResetPassword(ctx context.Context, id uint, req *dto.ResetPasswordRequest) error {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return errors.New("người dùng không tồn tại")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), 12)
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)
	return s.userRepo.Update(ctx, user)
}
