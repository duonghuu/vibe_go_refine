package controller

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	publicPageDTO "go_refine_dashboard_be/internal/usecases/publicpage/dto"
	publicPageService "go_refine_dashboard_be/internal/usecases/publicpage/service"
)

type PublicPageController struct {
	service publicPageService.PublicPageService
}

func NewPublicPageController(service publicPageService.PublicPageService) *PublicPageController {
	return &PublicPageController{service: service}
}

func (c *PublicPageController) GetBySlug(ctx *gin.Context) {
	var uri publicPageDTO.GetPublicPageURI
	if err := ctx.ShouldBindUri(&uri); err != nil {
		ctx.JSON(http.StatusBadRequest, publicPageDTO.PublicErrorResponse{Error: "slug không hợp lệ", Code: "INVALID_SLUG"})
		return
	}

	result, err := c.service.GetBySlug(ctx.Request.Context(), uri.Slug)
	if err != nil {
		switch {
		case errors.Is(err, publicPageService.ErrInvalidSlug):
			ctx.JSON(http.StatusBadRequest, publicPageDTO.PublicErrorResponse{Error: "slug không hợp lệ", Code: "INVALID_SLUG"})
		case errors.Is(err, publicPageService.ErrPublicPageNotFound):
			ctx.JSON(http.StatusNotFound, publicPageDTO.PublicErrorResponse{Error: "không tìm thấy trang", Code: "PAGE_NOT_FOUND"})
		default:
			ctx.JSON(http.StatusInternalServerError, publicPageDTO.PublicErrorResponse{Error: "internal error", Code: "INTERNAL_ERROR"})
		}
		return
	}

	ctx.JSON(http.StatusOK, result)
}
