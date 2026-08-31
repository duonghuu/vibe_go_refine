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

func (c *PageController) GetPages(ctx *gin.Context) {
	var query dto.GetPagesQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "validation_error",
			"message": "Tham số truy vấn không hợp lệ: " + err.Error(),
		})
		return
	}

	pages, err := c.pageService.GetPages(ctx.Request.Context(), query)
	if err != nil {
		if errors.Is(err, service.ErrPageInvalidInput) {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "validation_error", "message": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "internal_error",
			"message": "Không thể lấy danh sách trang",
		})
		return
	}

	ctx.JSON(http.StatusOK, pages)
}

func (c *PageController) GetPageByID(ctx *gin.Context) {
	var uri dto.GetPageByIDURI
	if err := ctx.ShouldBindUri(&uri); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "validation_error", "message": "ID trang không hợp lệ"})
		return
	}

	page, err := c.pageService.GetPageByID(ctx.Request.Context(), uri.PageID)
	if err != nil {
		if errors.Is(err, service.ErrPageNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "page_not_found", "message": "Không tìm thấy trang"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal_error", "message": "Không thể lấy chi tiết trang"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": page})
}

func (c *PageController) UpdatePage(ctx *gin.Context) {
	var uri dto.UpdatePageURI
	if err := ctx.ShouldBindUri(&uri); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "validation_error", "message": "ID trang không hợp lệ"})
		return
	}

	var req dto.UpdatePageRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "validation_error", "message": "Dữ liệu không hợp lệ: " + err.Error()})
		return
	}

	page, err := c.pageService.UpdatePage(ctx.Request.Context(), uri.PageID, &req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrPageNotFound):
			ctx.JSON(http.StatusNotFound, gin.H{"error": "page_not_found", "message": "Không tìm thấy trang"})
		case errors.Is(err, service.ErrPageInvalidInput):
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "validation_error", "message": err.Error()})
		case errors.Is(err, service.ErrPageSlugConflict):
			ctx.JSON(http.StatusConflict, gin.H{
				"error":   "slug_conflict",
				"message": err.Error(),
				"fields":  gin.H{"slug": err.Error()},
			})
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal_error", "message": "Không thể cập nhật trang"})
		}
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": page, "message": "Cập nhật trang thành công"})
}

func (c *PageController) DeletePage(ctx *gin.Context) {
	var uri dto.DeletePageURI
	if err := ctx.ShouldBindUri(&uri); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "validation_error",
			"message": "ID trang không hợp lệ",
		})
		return
	}

	if err := c.pageService.DeletePage(ctx.Request.Context(), uri.PageID); err != nil {
		if errors.Is(err, service.ErrPageNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error":   "page_not_found",
				"message": "Không tìm thấy trang",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "internal_error",
			"message": "Không thể xóa trang",
		})
		return
	}

	ctx.JSON(http.StatusOK, dto.DeletePageResponse{Message: "Xóa trang thành công"})
}
