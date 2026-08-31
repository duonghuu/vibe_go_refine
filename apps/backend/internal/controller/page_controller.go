package controller

import (
	"errors"
	"net/http"

	"go_refine_dashboard_be/internal/usecases/page/dto"
	"go_refine_dashboard_be/internal/usecases/page/service"

	"github.com/gin-gonic/gin"
)

type PageController struct {
	pageService service.PageService
}

func NewPageController(pageService service.PageService) *PageController {
	return &PageController{pageService: pageService}
}

func (c *PageController) CreatePage(ctx *gin.Context) {
	var req dto.CreatePageRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "validation_error",
			"message": "Dữ liệu không hợp lệ: " + err.Error(),
		})
		return
	}

	userIDValue, exists := ctx.Get("userId")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized", "message": "Không tìm thấy thông tin người dùng"})
		return
	}
	authorID, ok := userIDValue.(uint)
	if !ok || authorID == 0 {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized", "message": "Thông tin định danh không hợp lệ"})
		return
	}

	page, err := c.pageService.CreatePage(ctx.Request.Context(), authorID, &req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrPageInvalidInput):
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "validation_error", "message": err.Error()})
		case errors.Is(err, service.ErrPageSlugConflict):
			ctx.JSON(http.StatusConflict, gin.H{
				"error":   "slug_conflict",
				"message": err.Error(),
				"fields":  gin.H{"slug": err.Error()},
			})
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal_error", "message": "Không thể tạo trang"})
		}
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"data": page})
}
