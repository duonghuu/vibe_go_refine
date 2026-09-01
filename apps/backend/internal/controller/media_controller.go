package controller

import (
	"errors"
	"net/http"
	"strconv"

	"go_refine_dashboard_be/internal/usecases/media/dto"
	"go_refine_dashboard_be/internal/usecases/media/service"

	"github.com/gin-gonic/gin"
)

type MediaController struct {
	mediaService service.MediaService
}

func (c *MediaController) ListMedia(ctx *gin.Context) {
	var query dto.GetMediaListQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "VALIDATION_ERROR", "code": "VALIDATION_ERROR"})
		return
	}
	userID, role, err := mediaIdentity(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "UNAUTHORIZED", "code": "UNAUTHORIZED"})
		return
	}
	result, err := c.mediaService.ListMedia(ctx.Request.Context(), query, userID, role)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "INTERNAL_ERROR", "code": "INTERNAL_ERROR"})
		return
	}
	ctx.JSON(http.StatusOK, result)
}

func NewMediaController(mediaService service.MediaService) *MediaController {
	return &MediaController{
		mediaService: mediaService,
	}
}

func (c *MediaController) UploadFile(ctx *gin.Context) {
	// Parse Multipart Form with 5MB limit
	err := ctx.Request.ParseMultipartForm(5 << 20)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "File size exceeds 5MB limit or invalid request"})
		return
	}

	file, err := ctx.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}

	userID, role, err := mediaIdentity(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	_ = role
	response, err := c.mediaService.UploadFile(ctx.Request.Context(), file, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *MediaController) DeleteFile(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	userID, role, err := mediaIdentity(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	err = c.mediaService.DeleteFile(ctx.Request.Context(), uint(id), userID, role)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrMediaNotFound):
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error(), "code": "NOT_FOUND"})
		case errors.Is(err, service.ErrMediaForbidden):
			ctx.JSON(http.StatusForbidden, gin.H{"error": err.Error(), "code": "FORBIDDEN"})
		case errors.Is(err, service.ErrMediaAttached):
			ctx.JSON(http.StatusConflict, gin.H{"error": err.Error(), "code": "MEDIA_ATTACHED"})
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal_error", "code": "INTERNAL_ERROR"})
		}
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "deleted successfully"})
}

func mediaIdentity(ctx *gin.Context) (uint, string, error) {
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
		return 0, "", errors.New("vai trò không hợp lệ")
	}
	return userID, role, nil
}
