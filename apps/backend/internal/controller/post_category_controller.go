package controller

import (
	"fmt"
	"net/http"

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
