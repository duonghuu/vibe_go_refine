package controller

import (
	"net/http"

	"go_refine_dashboard_be/internal/usecases/post/dto"
	"go_refine_dashboard_be/internal/usecases/post/service"

	"github.com/gin-gonic/gin"
)

type PostController struct {
	postService service.PostService
}

func NewPostController(postService service.PostService) *PostController {
	return &PostController{
		postService: postService,
	}
}

func (c *PostController) CreatePost(ctx *gin.Context) {
	var req dto.CreatePostRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Dữ liệu không hợp lệ: " + err.Error()})
		return
	}

	userIDVal, exists := ctx.Get("userId")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Không tìm thấy thông tin người dùng"})
		return
	}

	authorID, ok := userIDVal.(uint)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Thông tin định danh không hợp lệ"})
		return
	}

	post, err := c.postService.CreatePost(ctx.Request.Context(), authorID, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"data": post})
}
