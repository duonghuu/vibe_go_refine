package controller

import (
	"fmt"
	"net/http"
	"strconv"

	"go_refine_dashboard_be/internal/usecases/postcategory/dto"
	"go_refine_dashboard_be/internal/usecases/postcategory/service"

	"github.com/gin-gonic/gin"
)

type PostCategoryController struct {
	postCategoryService service.PostCategoryService
}

func NewPostCategoryController(postCategoryService service.PostCategoryService) *PostCategoryController {
	return &PostCategoryController{
		postCategoryService: postCategoryService,
	}
}

func (c *PostCategoryController) GetPostCategoryTree(ctx *gin.Context) {
	fmt.Printf("=====TRACE: [2026/08/18 11:21:43] [apps/backend/internal/controller/post_category_controller.go:23] - message\n")
	typeCode := ctx.Query("typeCode")
	if typeCode == "" {
		typeCode = ctx.Query("type_code")
	}

	if typeCode == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Thiếu tham số typeCode"})
		return
	}

	items, err := c.postCategoryService.GetPostCategoryTree(ctx.Request.Context(), typeCode)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": items,
	})
}

func (c *PostCategoryController) CreatePostCategory(ctx *gin.Context) {
	var req dto.CreatePostCategoryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item, err := c.postCategoryService.CreatePostCategory(ctx.Request.Context(), &req)
	if err != nil {
		if err.Error() == "Slug đã tồn tại" {
			ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		if err.Error() == "Danh mục cha không tồn tại" || err.Error() == "Danh mục cha không thuộc cùng loại bài viết" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"data": item})
}

func (c *PostCategoryController) GetPostCategory(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID không hợp lệ"})
		return
	}

	item, err := c.postCategoryService.GetPostCategoryByID(ctx.Request.Context(), uint(id))
	if err != nil {
		if err.Error() == "Không tìm thấy danh mục bài viết" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": item})
}

func (c *PostCategoryController) UpdatePostCategory(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID không hợp lệ"})
		return
	}

	var req dto.UpdatePostCategoryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item, err := c.postCategoryService.UpdatePostCategory(ctx.Request.Context(), uint(id), &req)
	if err != nil {
		if err.Error() == "Slug đã tồn tại" {
			ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		if err.Error() == "Không tìm thấy danh mục bài viết" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if err.Error() == "Danh mục cha không tồn tại" || err.Error() == "Danh mục cha không thuộc cùng loại bài viết" || err.Error() == "Không thể chọn danh mục hiện tại làm danh mục cha" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": item})
}

func (c *PostCategoryController) GetPostCategories(ctx *gin.Context) {
	skip, _ := strconv.Atoi(ctx.DefaultQuery("_start", "0"))
	end, _ := strconv.Atoi(ctx.DefaultQuery("_end", "10"))
	limit := end - skip
	if limit < 0 {
		limit = 10
	}
	sortField := ctx.DefaultQuery("_sort", "sort_order")
	sortOrder := ctx.DefaultQuery("_order", "ASC")
	query := ctx.Query("q")
	status := ctx.Query("status")
	
	typeCode := ctx.Query("typeCode")
	if typeCode == "" {
		typeCode = ctx.Query("type_code")
	}

	result, err := c.postCategoryService.GetPostCategories(ctx.Request.Context(), skip, limit, sortField, sortOrder, query, status, typeCode)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data":  result.Data,
		"total": result.Total,
	})
}

func (c *PostCategoryController) DeletePostCategory(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID không hợp lệ"})
		return
	}

	if err := c.postCategoryService.DeletePostCategory(ctx.Request.Context(), uint(id)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Đã xóa danh mục thành công"})
}
