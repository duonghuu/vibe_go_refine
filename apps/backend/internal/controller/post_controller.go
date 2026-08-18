package controller

import (
	"fmt"
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

func (c *PostController) GetPosts(ctx *gin.Context) {
	var query dto.GetPostsQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Tham số truy vấn không hợp lệ: " + err.Error()})
		return
	}

	resp, err := c.postService.GetPosts(ctx.Request.Context(), query)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

func (c *PostController) DeletePost(ctx *gin.Context) {
	idStr := ctx.Param("post_id")
	var id uint
	if _, err := fmt.Sscanf(idStr, "%d", &id); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID không hợp lệ"})
		return
	}

	_, exists := ctx.Get("userId")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Không tìm thấy thông tin người dùng"})
		return
	}

	if err := c.postService.DeletePost(ctx.Request.Context(), id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Xóa bài viết thành công"})
}
