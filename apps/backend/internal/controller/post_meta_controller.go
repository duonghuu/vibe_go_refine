package controller

import (
	"net/http"

	"go_refine_dashboard_be/internal/usecases/postmeta/dto"
	"go_refine_dashboard_be/internal/usecases/postmeta/service"

	"github.com/gin-gonic/gin"
)

type PostMetaController struct {
	postMetaService service.PostMetaService
}

func NewPostMetaController(postMetaService service.PostMetaService) *PostMetaController {
	return &PostMetaController{postMetaService: postMetaService}
}

func (c *PostMetaController) GetPostMeta(ctx *gin.Context) {
	var req dto.GetPostMetaRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := c.postMetaService.GetByPostID(ctx.Request.Context(), &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, res)
}

func (c *PostMetaController) SyncPostMeta(ctx *gin.Context) {
	var req dto.SyncPostMetaRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := c.postMetaService.Sync(ctx.Request.Context(), &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, res)
}

func (c *PostMetaController) DeletePostMeta(ctx *gin.Context) {
	var req dto.DeletePostMetaRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := c.postMetaService.Delete(ctx.Request.Context(), &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, res)
}
