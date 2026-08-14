package controller

import (
	"net/http"
	"strconv"

	"go_refine_dashboard_be/internal/usecases/postmedia/dto"
	"go_refine_dashboard_be/internal/usecases/postmedia/service"

	"github.com/gin-gonic/gin"
)

type PostMediaController struct {
	postMediaService service.PostMediaService
}

func NewPostMediaController(postMediaService service.PostMediaService) *PostMediaController {
	return &PostMediaController{
		postMediaService: postMediaService,
	}
}

func (c *PostMediaController) GetPostMedia(ctx *gin.Context) {
	postID, err := strconv.ParseUint(ctx.Param("post_id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Post ID không hợp lệ"})
		return
	}

	collection := ctx.Query("collection")

	items, total, err := c.postMediaService.GetPostMedia(ctx.Request.Context(), uint(postID), collection)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data":  items,
		"total": total,
	})
}

func (c *PostMediaController) CreatePostMedia(ctx *gin.Context) {
	postID, err := strconv.ParseUint(ctx.Param("post_id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Post ID không hợp lệ"})
		return
	}

	var req dto.CreatePostMediaRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item, err := c.postMediaService.CreatePostMedia(ctx.Request.Context(), uint(postID), &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"data": item})
}

func (c *PostMediaController) SyncPostMedia(ctx *gin.Context) {
	postID, err := strconv.ParseUint(ctx.Param("post_id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Post ID không hợp lệ"})
		return
	}

	collection := ctx.Param("collection")
	if collection == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Collection không được trống"})
		return
	}

	var req dto.SyncPostMediaRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.postMediaService.SyncPostMedia(ctx.Request.Context(), uint(postID), collection, &req); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Đồng bộ thành công"})
}

func (c *PostMediaController) DeletePostMedia(ctx *gin.Context) {
	postID, err := strconv.ParseUint(ctx.Param("post_id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Post ID không hợp lệ"})
		return
	}

	mediaID, err := strconv.ParseUint(ctx.Param("media_id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Media ID không hợp lệ"})
		return
	}

	collection := ctx.Query("collection")

	if err := c.postMediaService.DeletePostMedia(ctx.Request.Context(), uint(postID), uint(mediaID), collection); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}
