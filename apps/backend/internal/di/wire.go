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
	postTypeService "go_refine_dashboard_be/internal/usecases/posttype/service"
	postMediaService "go_refine_dashboard_be/internal/usecases/postmedia/service"
	postMetaService "go_refine_dashboard_be/internal/usecases/postmeta/service"
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

var PostTypeSet = wire.NewSet(
	repository.NewPostTypeRepository,
	postTypeService.NewPostTypeService,
	controller.NewPostTypeController,
)

var PostMediaSet = wire.NewSet(
	repository.NewPostMediaRepository,
	postMediaService.NewPostMediaService,
	controller.NewPostMediaController,
)

var PostMetaSet = wire.NewSet(
	repository.NewPostMetaRepository,
	postMetaService.NewPostMetaService,
	controller.NewPostMetaController,
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

func InitializePostTypeController(db *gorm.DB, redisClient *redis.Client) *controller.PostTypeController {
	wire.Build(PostTypeSet)
	return &controller.PostTypeController{}
}

func InitializePostMediaController(db *gorm.DB, redisClient *redis.Client) *controller.PostMediaController {
	wire.Build(PostMediaSet)
	return &controller.PostMediaController{}
}

func InitializePostMetaController(db *gorm.DB, redisClient *redis.Client) *controller.PostMetaController {
	wire.Build(PostMetaSet)
	return &controller.PostMetaController{}
}
