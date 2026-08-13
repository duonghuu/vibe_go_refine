package controller

import (
	"net/http"
	"strconv"

	"go_refine_dashboard_be/internal/usecases/posttype/dto"
	"go_refine_dashboard_be/internal/usecases/posttype/service"

	"github.com/gin-gonic/gin"
)

type PostTypeController struct {
	postTypeService service.PostTypeService
}

func NewPostTypeController(postTypeService service.PostTypeService) *PostTypeController {
	return &PostTypeController{
		postTypeService: postTypeService,
	}
}

func (c *PostTypeController) GetPostTypes(ctx *gin.Context) {
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

	items, total, err := c.postTypeService.GetPostTypes(ctx.Request.Context(), skip, limit, sortField, sortOrder, query, status)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data":  items,
		"total": total,
	})
}

func (c *PostTypeController) GetPostTypeByID(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID không hợp lệ"})
		return
	}

	item, err := c.postTypeService.GetPostTypeByID(ctx.Request.Context(), uint(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": item})
}

func (c *PostTypeController) CreatePostType(ctx *gin.Context) {
	var req dto.CreatePostTypeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item, err := c.postTypeService.CreatePostType(ctx.Request.Context(), &req)
	if err != nil {
		if err.Error() == "Mã loại bài viết đã tồn tại" {
			ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"data": item})
}

func (c *PostTypeController) UpdatePostType(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID không hợp lệ"})
		return
	}

	var req dto.UpdatePostTypeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item, err := c.postTypeService.UpdatePostType(ctx.Request.Context(), uint(id), &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": item})
}

func (c *PostTypeController) DeletePostType(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID không hợp lệ"})
		return
	}

	if err := c.postTypeService.DeletePostType(ctx.Request.Context(), uint(id)); err != nil {
		if err.Error() == "Không thể xóa loại bài viết đang được sử dụng" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Đã xóa loại bài viết thành công"})
}
