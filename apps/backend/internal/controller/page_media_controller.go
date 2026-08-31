package controller

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"go_refine_dashboard_be/internal/usecases/pagemedia/dto"
	"go_refine_dashboard_be/internal/usecases/pagemedia/service"
	"gorm.io/gorm"
)

type PageMediaController struct{ pageMediaService service.PageMediaService }

func NewPageMediaController(pageMediaService service.PageMediaService) *PageMediaController {
	return &PageMediaController{pageMediaService: pageMediaService}
}

func pageMediaIdentity(c *gin.Context) (uint, string, error) {
	v, ok := c.Get("userId")
	if !ok {
		return 0, "", errors.New("không tìm thấy thông tin người dùng")
	}
	id, ok := v.(uint)
	if !ok || id == 0 {
		return 0, "", errors.New("thông tin định danh không hợp lệ")
	}
	role, ok := c.Get("role")
	if !ok {
		return 0, "", errors.New("không tìm thấy vai trò người dùng")
	}
	roleValue, ok := role.(string)
	if !ok || roleValue == "" {
		return 0, "", errors.New("vai trò không hợp lệ")
	}
	return id, roleValue, nil
}

func writePageMediaError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, service.ErrPageNotFound), errors.Is(err, service.ErrMediaNotFound), errors.Is(err, service.ErrPageMediaNotFound), errors.Is(err, gorm.ErrRecordNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error(), "code": "NOT_FOUND"})
	case errors.Is(err, service.ErrMediaForbidden):
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error(), "code": "FORBIDDEN"})
	case errors.Is(err, service.ErrCollectionRequired), errors.Is(err, service.ErrCollectionInvalid), errors.Is(err, service.ErrThumbnailLimit), errors.Is(err, service.ErrMediaDuplicate), errors.Is(err, service.ErrSortOrderInvalid), errors.Is(err, service.ErrMediaTypeInvalid):
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "code": "VALIDATION_ERROR"})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal_error", "code": "INTERNAL_ERROR"})
	}
}

func (c *PageMediaController) GetPageMedia(ctx *gin.Context) {
	var uri dto.GetPageMediaURI
	var query dto.GetPageMediaQuery
	if err := ctx.ShouldBindUri(&uri); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "URI không hợp lệ", "code": "VALIDATION_ERROR"})
		return
	}
	if err := ctx.ShouldBindQuery(&query); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "query không hợp lệ", "code": "VALIDATION_ERROR"})
		return
	}
	items, total, err := c.pageMediaService.GetPageMedia(ctx.Request.Context(), uri.PageID, query.Collection)
	if err != nil {
		writePageMediaError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": items, "total": total})
}

func (c *PageMediaController) CreatePageMedia(ctx *gin.Context) {
	var uri dto.GetPageMediaURI
	var req dto.CreatePageMediaRequest
	if err := ctx.ShouldBindUri(&uri); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "URI không hợp lệ", "code": "VALIDATION_ERROR"})
		return
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "code": "VALIDATION_ERROR"})
		return
	}
	id, role, err := pageMediaIdentity(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error(), "code": "UNAUTHORIZED"})
		return
	}
	item, err := c.pageMediaService.CreatePageMedia(ctx.Request.Context(), uri.PageID, &req, id, role)
	if err != nil {
		writePageMediaError(ctx, err)
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"data": item})
}

func (c *PageMediaController) SyncPageMedia(ctx *gin.Context) {
	var uri dto.SyncPageMediaURI
	var req dto.SyncPageMediaRequest
	if err := ctx.ShouldBindUri(&uri); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "URI không hợp lệ", "code": "VALIDATION_ERROR"})
		return
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "code": "VALIDATION_ERROR"})
		return
	}
	id, role, err := pageMediaIdentity(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error(), "code": "UNAUTHORIZED"})
		return
	}
	items, err := c.pageMediaService.SyncPageMedia(ctx.Request.Context(), uri.PageID, uri.Collection, &req, id, role)
	if err != nil {
		writePageMediaError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Đồng bộ hình ảnh trang thành công", "data": items})
}

func (c *PageMediaController) DeletePageMedia(ctx *gin.Context) {
	var uri dto.DeletePageMediaURI
	var query dto.DeletePageMediaQuery
	if err := ctx.ShouldBindUri(&uri); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "URI không hợp lệ", "code": "VALIDATION_ERROR"})
		return
	}
	if err := ctx.ShouldBindQuery(&query); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "query không hợp lệ", "code": "VALIDATION_ERROR"})
		return
	}
	id, role, err := pageMediaIdentity(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error(), "code": "UNAUTHORIZED"})
		return
	}
	if err := c.pageMediaService.DeletePageMedia(ctx.Request.Context(), uri.PageID, uri.MediaID, query.Collection, id, role); err != nil {
		writePageMediaError(ctx, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}
