//go:build wireinject
// +build wireinject

package di

import (
	"go_refine_dashboard_be/internal/controller"
	"go_refine_dashboard_be/internal/infrastructure/repository"
	mediaService "go_refine_dashboard_be/internal/usecases/media/service"
	"go_refine_dashboard_be/internal/usecases/product/service"
	categoryService "go_refine_dashboard_be/internal/usecases/category/service"
	userService "go_refine_dashboard_be/internal/usecases/user/service"
	authService "go_refine_dashboard_be/internal/usecases/auth/service"
	"gorm.io/gorm"
	"github.com/redis/go-redis/v9"

	"github.com/google/wire"
)

var ProductSet = wire.NewSet(
	repository.NewProductRepository,
	service.NewProductService,
	controller.NewProductController,
)

var MediaSet = wire.NewSet(
	repository.NewMediaRepository,
	mediaService.NewMediaService,
	controller.NewMediaController,
)

var CategorySet = wire.NewSet(
	repository.NewCategoryRepository,
	categoryService.NewCategoryService,
	controller.NewCategoryController,
)

var UserSet = wire.NewSet(
	repository.NewUserRepository,
	userService.NewUserService,
	controller.NewUserController,
)

var AuthSet = wire.NewSet(
	repository.NewUserRepository,
	authService.NewAuthService,
	controller.NewAuthController,
)

func InitializeProductController(db *gorm.DB) *controller.ProductController {
	wire.Build(ProductSet)
	return &controller.ProductController{}
}

func InitializeMediaController(db *gorm.DB) *controller.MediaController {
	wire.Build(MediaSet)
	return &controller.MediaController{}
}

func InitializeCategoryController(db *gorm.DB) *controller.CategoryController {
	wire.Build(CategorySet)
	return &controller.CategoryController{}
}

func InitializeUserController(db *gorm.DB) *controller.UserController {
	wire.Build(UserSet)
	return &controller.UserController{}
}

func InitializeAuthController(db *gorm.DB, redisClient *redis.Client) *controller.AuthController {
	wire.Build(AuthSet)
	return &controller.AuthController{}
}
