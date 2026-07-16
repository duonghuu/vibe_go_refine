package controller

import (
	"net/http"
	"strconv"

	"go_refine_dashboard_be/internal/usecases/media/service"

	"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
)

type MediaController struct {
	mediaService service.MediaService
}

func NewMediaController(mediaService service.MediaService) *MediaController {
	return &MediaController{
		mediaService: mediaService,
	}
}

func (c *MediaController) UploadFile(ctx *gin.Context) {
	spew.Dump("UploadFile")
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

	// Mock User ID since authentication is temporarily disabled
	var mockUserId uint = 1

	response, err := c.mediaService.UploadFile(ctx.Request.Context(), file, mockUserId)
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

	// Mock User ID since authentication is temporarily disabled
	var mockUserId uint = 1

	err = c.mediaService.DeleteFile(ctx.Request.Context(), uint(id), mockUserId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "deleted successfully"})
}
