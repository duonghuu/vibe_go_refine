//go:build wireinject
// +build wireinject

package di

import (
	"go_refine_dashboard_be/internal/controller"
	"go_refine_dashboard_be/internal/infrastructure/repository"
	mediaService "go_refine_dashboard_be/internal/usecases/media/service"
	"go_refine_dashboard_be/internal/usecases/product/service"
	"gorm.io/gorm"

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

func InitializeProductController(db *gorm.DB) *controller.ProductController {
	wire.Build(ProductSet)
	return &controller.ProductController{}
}

func InitializeMediaController(db *gorm.DB) *controller.MediaController {
	wire.Build(MediaSet)
	return &controller.MediaController{}
}
