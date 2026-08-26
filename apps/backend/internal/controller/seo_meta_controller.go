package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go_refine_dashboard_be/internal/infrastructure/repository"
	"go_refine_dashboard_be/internal/usecases/seometa/dto"
	"go_refine_dashboard_be/internal/usecases/seometa/service"
	"net/http"
)

type SEOMetaController struct{ service service.SEOService }

func NewSEOMetaController(s service.SEOService) *SEOMetaController {
	return &SEOMetaController{service: s}
}
func (c *SEOMetaController) Get(ctx *gin.Context) {
	var uri dto.EntityURI
	if err := ctx.ShouldBindUri(&uri); err != nil {
		writeSEOError(ctx, http.StatusBadRequest, err)
		return
	}
	out, err := c.service.Get(ctx.Request.Context(), uri.EntityType, uri.EntityID)
	if err != nil {
		writeSEOError(ctx, status(err), err)
		return
	}
	ctx.JSON(http.StatusOK, out)
}
func (c *SEOMetaController) Put(ctx *gin.Context) {
	var uri dto.EntityURI
	if err := ctx.ShouldBindUri(&uri); err != nil {
		writeSEOError(ctx, http.StatusBadRequest, err)
		return
	}
	var req dto.UpsertRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		writeSEOError(ctx, http.StatusBadRequest, err)
		return
	}
	out, err := c.service.Upsert(ctx.Request.Context(), uri.EntityType, uri.EntityID, &req)
	if err != nil {
		writeSEOError(ctx, status(err), err)
		return
	}
	ctx.JSON(http.StatusOK, out)
}
func (c *SEOMetaController) Delete(ctx *gin.Context) {
	var uri dto.EntityURI
	if err := ctx.ShouldBindUri(&uri); err != nil {
		writeSEOError(ctx, http.StatusBadRequest, err)
		return
	}
	if err := c.service.Delete(ctx.Request.Context(), uri.EntityType, uri.EntityID); err != nil {
		writeSEOError(ctx, status(err), err)
		return
	}
	ctx.JSON(http.StatusOK, dto.DeleteResponse{Message: "Xóa SEO thành công"})
}
func status(err error) int {
	switch {
	case errors.Is(err, service.ErrValidation):
		return http.StatusUnprocessableEntity
	case errors.Is(err, repository.ErrEntityNotFound):
		return http.StatusNotFound
	case errors.Is(err, repository.ErrEntityTypeNotSupported), errors.Is(err, repository.ErrEntityTypeNotConfigured):
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
func writeSEOError(c *gin.Context, code int, err error) {
	c.JSON(code, dto.ErrorResponse{Error: "seo_meta_error", Message: err.Error()})
}
