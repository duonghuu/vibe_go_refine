package controller

import (
	"net/http"
	"strconv"

	"go_refine_dashboard_be/internal/usecases/category/dto"
	"go_refine_dashboard_be/internal/usecases/category/service"

	"github.com/gin-gonic/gin"
)

type CategoryController struct {
	categoryService service.CategoryService
}

func NewCategoryController(categoryService service.CategoryService) *CategoryController {
	return &CategoryController{
		categoryService: categoryService,
	}
}

func (c *CategoryController) GetCategories(ctx *gin.Context) {
	skip, _ := strconv.Atoi(ctx.DefaultQuery("_start", "0"))
	end, _ := strconv.Atoi(ctx.DefaultQuery("_end", "10"))
	limit := end - skip
	if limit < 0 {
		limit = 10
	}
	sortField := ctx.DefaultQuery("_sort", "id")
	sortOrder := ctx.DefaultQuery("_order", "DESC")
	query := ctx.Query("q")
	status := ctx.Query("status")

	var parentID *uint
	if pid := ctx.Query("parentId"); pid != "" {
		if pidInt, err := strconv.ParseUint(pid, 10, 32); err == nil {
			pidUint := uint(pidInt)
			parentID = &pidUint
		}
	}

	categories, total, err := c.categoryService.GetCategories(ctx.Request.Context(), skip, limit, sortField, sortOrder, query, status, parentID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Refine expects an array if it's data property or direct array. But we'll follow Refine's provider { data, total }
	// We'll add custom headers for total count if refine data provider uses headers, but returning JSON with data/total is also common.
	// We will respond with JSON
	ctx.JSON(http.StatusOK, gin.H{
		"data":  categories,
		"total": total,
	})
}

func (c *CategoryController) GetCategoryByID(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID không hợp lệ"})
		return
	}

	category, err := c.categoryService.GetCategoryByID(ctx.Request.Context(), uint(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// Refine typically expects the data property or direct object. Let's wrap in "data"
	ctx.JSON(http.StatusOK, gin.H{"data": category})
}

func (c *CategoryController) CreateCategory(ctx *gin.Context) {
	var req dto.CreateCategoryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	category, err := c.categoryService.CreateCategory(ctx.Request.Context(), &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"data": category})
}

func (c *CategoryController) UpdateCategory(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID không hợp lệ"})
		return
	}

	var req dto.UpdateCategoryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	category, err := c.categoryService.UpdateCategory(ctx.Request.Context(), uint(id), &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": category})
}

func (c *CategoryController) DeleteCategory(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID không hợp lệ"})
		return
	}

	if err := c.categoryService.DeleteCategory(ctx.Request.Context(), uint(id)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Đã xóa danh mục thành công"})
}

func (c *CategoryController) BulkUpdateStatus(ctx *gin.Context) {
	var req dto.BulkUpdateCategoryStatusRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.categoryService.BulkUpdateStatus(ctx.Request.Context(), &req); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Cập nhật trạng thái hàng loạt thành công"})
}

func (c *CategoryController) ReorderCategories(ctx *gin.Context) {
	var req dto.ReorderCategoriesRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.categoryService.ReorderCategories(ctx.Request.Context(), &req); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Sắp xếp danh mục thành công"})
}
