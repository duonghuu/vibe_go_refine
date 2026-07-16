package service

import (
	"context"
	"go_refine_dashboard_be/internal/domain/product/entity"
	"go_refine_dashboard_be/internal/infrastructure/repository"
	"go_refine_dashboard_be/internal/usecases/product/dto"
)

type ProductService interface {
	GetProducts(ctx context.Context, req dto.ListProductReq) ([]entity.Product, int64, error)
	CreateProduct(ctx context.Context, req dto.CreateProductReq) (*entity.Product, error)
	UpdateProduct(ctx context.Context, id uint, req dto.UpdateProductReq) (*entity.Product, error)
	UpdateProductStatus(ctx context.Context, id uint, req dto.UpdateProductStatusReq) error
	DeleteProduct(ctx context.Context, id uint) error
	BulkDeleteProducts(ctx context.Context, req dto.BulkDeleteReq) error
	BulkUpdateProductStatus(ctx context.Context, req dto.BulkStatusReq) error
}

type productServiceImpl struct {
	repo repository.ProductRepository
}

func NewProductService(repo repository.ProductRepository) ProductService {
	return &productServiceImpl{repo: repo}
}

func (s *productServiceImpl) GetProducts(ctx context.Context, req dto.ListProductReq) ([]entity.Product, int64, error) {
	page := req.Page
	if page < 0 {
		page = 0
	}
	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 10
	}
	return s.repo.FindAll(ctx, page, pageSize, req.Sort, req.Order, req.Search, req.CategoryID, req.Status, req.MinPrice, req.MaxPrice, req.MinStock)
}

func (s *productServiceImpl) CreateProduct(ctx context.Context, req dto.CreateProductReq) (*entity.Product, error) {
	product := &entity.Product{
		Name:       req.Name,
		SKU:        req.SKU,
		CategoryID: req.CategoryID,
		Price:      req.Price,
		SalePrice:  req.SalePrice,
		Stock:      req.Stock,
		Status:     req.Status,
		ImageURL:   req.ImageURL,
	}
	err := s.repo.Create(ctx, product)
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (s *productServiceImpl) UpdateProduct(ctx context.Context, id uint, req dto.UpdateProductReq) (*entity.Product, error) {
	product, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	product.Name = req.Name
	product.SKU = req.SKU
	product.CategoryID = req.CategoryID
	product.Price = req.Price
	product.SalePrice = req.SalePrice
	product.Stock = req.Stock
	product.Status = req.Status
	product.ImageURL = req.ImageURL

	err = s.repo.Update(ctx, product)
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (s *productServiceImpl) UpdateProductStatus(ctx context.Context, id uint, req dto.UpdateProductStatusReq) error {
	return s.repo.UpdateStatus(ctx, id, req.Status)
}

func (s *productServiceImpl) DeleteProduct(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}

func (s *productServiceImpl) BulkDeleteProducts(ctx context.Context, req dto.BulkDeleteReq) error {
	return s.repo.BulkDelete(ctx, req.IDs)
}

func (s *productServiceImpl) BulkUpdateProductStatus(ctx context.Context, req dto.BulkStatusReq) error {
	return s.repo.BulkUpdateStatus(ctx, req.IDs, req.Status)
}
