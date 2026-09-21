//go:build wireinject
// +build wireinject

package di

import (
	"github.com/redis/go-redis/v9"
	"go_refine_dashboard_be/internal/controller"
	publicPageCache "go_refine_dashboard_be/internal/infrastructure/cache"
	imageContentQuery "go_refine_dashboard_be/internal/infrastructure/queryservice"
	"go_refine_dashboard_be/internal/infrastructure/repository"
	authService "go_refine_dashboard_be/internal/usecases/auth/service"
	categoryService "go_refine_dashboard_be/internal/usecases/category/service"
	imageContentService "go_refine_dashboard_be/internal/usecases/imagecontent/service"
	mediaService "go_refine_dashboard_be/internal/usecases/media/service"
	pageService "go_refine_dashboard_be/internal/usecases/page/service"
	pageMediaService "go_refine_dashboard_be/internal/usecases/pagemedia/service"
	pageSectionService "go_refine_dashboard_be/internal/usecases/pagesection/service"
	postService "go_refine_dashboard_be/internal/usecases/post/service"
	postCategoryService "go_refine_dashboard_be/internal/usecases/postcategory/service"
	postMediaService "go_refine_dashboard_be/internal/usecases/postmedia/service"
	postMetaService "go_refine_dashboard_be/internal/usecases/postmeta/service"
	postTypeService "go_refine_dashboard_be/internal/usecases/posttype/service"
	"go_refine_dashboard_be/internal/usecases/product/service"
	publicPageService "go_refine_dashboard_be/internal/usecases/publicpage/service"
	seoMetaService "go_refine_dashboard_be/internal/usecases/seometa/service"
	userService "go_refine_dashboard_be/internal/usecases/user/service"
	"gorm.io/gorm"

	"github.com/google/wire"
)

var ProductSet = wire.NewSet(
	repository.NewProductRepository,
	service.NewProductService,
	controller.NewProductController,
)

var PublicPageCacheSet = wire.NewSet(
	repository.NewPublicPageCacheRepository,
	publicPageCache.NewPublicPageCacheInvalidator,
)

var MediaSet = wire.NewSet(
	repository.NewMediaRepository,
	PublicPageCacheSet,
	mediaService.NewMediaService,
	controller.NewMediaController,
)

var CategorySet = wire.NewSet(
	repository.NewCategoryRepository,
	PublicPageCacheSet,
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
	PublicPageCacheSet,
	postMediaService.NewPostMediaService,
	controller.NewPostMediaController,
)

var PostMetaSet = wire.NewSet(
	repository.NewPostMetaRepository,
	postMetaService.NewPostMetaService,
	controller.NewPostMetaController,
)

var PostSet = wire.NewSet(
	repository.NewPostRepository,
	repository.NewPostCategoryRepository,
	PublicPageCacheSet,
	postService.NewPostService,
	controller.NewPostController,
)

var PostCategorySet = wire.NewSet(
	repository.NewPostCategoryRepository,
	postCategoryService.NewPostCategoryService,
	controller.NewPostCategoryController,
)

var SEOMetaSet = wire.NewSet(
	repository.NewSEOMetaRepository,
	repository.NewEntityRegistry,
	PublicPageCacheSet,
	seoMetaService.NewSEOService,
	controller.NewSEOMetaController,
)

var PageSet = wire.NewSet(
	repository.NewPageRepository,
	PublicPageCacheSet,
	pageService.NewPageService,
	controller.NewPageController,
)

var PageMediaSet = wire.NewSet(
	repository.NewPageMediaRepository,
	PublicPageCacheSet,
	pageMediaService.NewPageMediaService,
	controller.NewPageMediaController,
)

var PageSectionSet = wire.NewSet(
	repository.NewPageSectionRepository,
	PublicPageCacheSet,
	pageSectionService.NewPageSectionService,
	controller.NewPageSectionController,
)

var ImageContentSet = wire.NewSet(
	repository.NewImageContentRepository,
	imageContentQuery.NewImageContentQueryService,
	imageContentService.NewImageContentService,
	controller.NewImageContentController,
)

var PublicPageSet = wire.NewSet(
	repository.NewPublicPageRepository,
	publicPageService.NewPublicPageService,
	controller.NewPublicPageController,
)

func InitializeProductController(db *gorm.DB) *controller.ProductController {
	wire.Build(ProductSet)
	return &controller.ProductController{}
}

func InitializeMediaController(db *gorm.DB, redisClient *redis.Client) *controller.MediaController {
	wire.Build(MediaSet)
	return &controller.MediaController{}
}

func InitializeCategoryController(db *gorm.DB, redisClient *redis.Client) *controller.CategoryController {
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

func InitializePostController(db *gorm.DB, redisClient *redis.Client) *controller.PostController {
	wire.Build(PostSet)
	return &controller.PostController{}
}

func InitializePostCategoryController(db *gorm.DB, redisClient *redis.Client) *controller.PostCategoryController {
	wire.Build(PostCategorySet)
	return &controller.PostCategoryController{}
}

func InitializeSEOMetaController(db *gorm.DB, redisClient *redis.Client) *controller.SEOMetaController {
	wire.Build(SEOMetaSet)
	return &controller.SEOMetaController{}
}

func InitializePageController(db *gorm.DB, redisClient *redis.Client) *controller.PageController {
	wire.Build(PageSet)
	return &controller.PageController{}
}

func InitializePageMediaController(db *gorm.DB, redisClient *redis.Client) *controller.PageMediaController {
	wire.Build(PageMediaSet)
	return &controller.PageMediaController{}
}

func InitializePageSectionController(db *gorm.DB, redisClient *redis.Client) *controller.PageSectionController {
	wire.Build(PageSectionSet)
	return &controller.PageSectionController{}
}

func InitializeImageContentController(db *gorm.DB, redisClient *redis.Client) *controller.ImageContentController {
	wire.Build(ImageContentSet)
	return &controller.ImageContentController{}
}

func InitializePublicPageController(db *gorm.DB, redisClient *redis.Client, publicSiteURL string) *controller.PublicPageController {
	wire.Build(PublicPageSet)
	return &controller.PublicPageController{}
}
