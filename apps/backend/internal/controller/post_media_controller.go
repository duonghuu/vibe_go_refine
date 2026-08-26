package controller

import (
	"errors"
	"net/http"

	"go_refine_dashboard_be/internal/usecases/postmedia/dto"
	"go_refine_dashboard_be/internal/usecases/postmedia/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PostMediaController struct {
	postMediaService service.PostMediaService
}

func NewPostMediaController(postMediaService service.PostMediaService) *PostMediaController {
	return &PostMediaController{postMediaService: postMediaService}
}

func postMediaIdentity(ctx *gin.Context) (uint, string, error) {
	value, exists := ctx.Get("userId")
	if !exists {
		return 0, "", errors.New("không tìm thấy thông tin người dùng")
	}
	userID, ok := value.(uint)
	if !ok || userID == 0 {
		return 0, "", errors.New("thông tin định danh không hợp lệ")
	}
	roleValue, exists := ctx.Get("role")
	if !exists {
		return 0, "", errors.New("không tìm thấy vai trò người dùng")
	}
	role, ok := roleValue.(string)
	if !ok || role == "" {
		return 0, "", errors.New("vai trò người dùng không hợp lệ")
	}
	return userID, role, nil
}

func writePostMediaError(ctx *gin.Context, err error) {
	switch {
	case errors.Is(err, service.ErrPostNotFound), errors.Is(err, service.ErrMediaNotFound), errors.Is(err, gorm.ErrRecordNotFound):
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error(), "code": "NOT_FOUND"})
	case errors.Is(err, service.ErrMediaForbidden):
		ctx.JSON(http.StatusForbidden, gin.H{"error": err.Error(), "code": "FORBIDDEN"})
	case errors.Is(err, service.ErrCollectionInvalid), errors.Is(err, service.ErrThumbnailLimit), errors.Is(err, service.ErrMediaDuplicate), errors.Is(err, service.ErrMediaTypeInvalid):
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "code": "VALIDATION_ERROR"})
	default:
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "code": "INTERNAL_ERROR"})
	}
}

func (c *PostMediaController) GetPostMedia(ctx *gin.Context) {
	var uri struct {
		PostID uint `uri:"post_id" binding:"required,min=1"`
	}
	if err := ctx.ShouldBindUri(&uri); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Post ID không hợp lệ", "code": "VALIDATION_ERROR"})
		return
	}

	collection := ctx.Query("collection")
	items, total, err := c.postMediaService.GetPostMedia(ctx.Request.Context(), uri.PostID, collection)
	if err != nil {
		writePostMediaError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": items, "total": total})
}

func (c *PostMediaController) CreatePostMedia(ctx *gin.Context) {
	var uri struct {
		PostID uint `uri:"post_id" binding:"required,min=1"`
	}
	if err := ctx.ShouldBindUri(&uri); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Post ID không hợp lệ", "code": "VALIDATION_ERROR"})
		return
	}
	var req dto.CreatePostMediaRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "code": "VALIDATION_ERROR"})
		return
	}
	userID, role, err := postMediaIdentity(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error(), "code": "UNAUTHORIZED"})
		return
	}
	item, err := c.postMediaService.CreatePostMedia(ctx.Request.Context(), uri.PostID, &req, userID, role)
	if err != nil {
		writePostMediaError(ctx, err)
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"data": item})
}

func (c *PostMediaController) SyncPostMedia(ctx *gin.Context) {
	var uri struct {
		PostID     uint   `uri:"post_id" binding:"required,min=1"`
		Collection string `uri:"collection" binding:"required"`
	}
	if err := ctx.ShouldBindUri(&uri); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "URI không hợp lệ", "code": "VALIDATION_ERROR"})
		return
	}
	var req dto.SyncPostMediaRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "code": "VALIDATION_ERROR"})
		return
	}
	userID, role, err := postMediaIdentity(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error(), "code": "UNAUTHORIZED"})
		return
	}
	items, err := c.postMediaService.SyncPostMedia(ctx.Request.Context(), uri.PostID, uri.Collection, &req, userID, role)
	if err != nil {
		writePostMediaError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Đồng bộ thành công", "data": items})
}

func (c *PostMediaController) DeletePostMedia(ctx *gin.Context) {
	var uri struct {
		PostID  uint `uri:"post_id" binding:"required,min=1"`
		MediaID uint `uri:"media_id" binding:"required,min=1"`
	}
	if err := ctx.ShouldBindUri(&uri); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "URI không hợp lệ", "code": "VALIDATION_ERROR"})
		return
	}
	userID, role, err := postMediaIdentity(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error(), "code": "UNAUTHORIZED"})
		return
	}
	if err := c.postMediaService.DeletePostMedia(ctx.Request.Context(), uri.PostID, uri.MediaID, ctx.Query("collection"), userID, role); err != nil {
		writePostMediaError(ctx, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}
