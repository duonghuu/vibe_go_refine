package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go_refine_dashboard_be/internal/usecases/pagesection/dto"
	"go_refine_dashboard_be/internal/usecases/pagesection/service"
	"net/http"
)

type PageSectionController struct{ sectionService service.PageSectionService }

func NewPageSectionController(s service.PageSectionService) *PageSectionController {
	return &PageSectionController{sectionService: s}
}
func sectionIdentity(c *gin.Context) (uint, string, error) {
	v, ok := c.Get("userId")
	if !ok {
		return 0, "", errors.New("unauthorized")
	}
	id, ok := v.(uint)
	if !ok || id == 0 {
		return 0, "", errors.New("unauthorized")
	}
	r, ok := c.Get("role")
	if !ok {
		return 0, "", errors.New("unauthorized")
	}
	role, ok := r.(string)
	if !ok || role == "" {
		return 0, "", errors.New("unauthorized")
	}
	return id, role, nil
}
func writeSectionError(c *gin.Context, err error) {
	code := "INTERNAL_ERROR"
	status := http.StatusInternalServerError
	switch {
	case errors.Is(err, service.ErrInvalidInput):
		code = "VALIDATION_ERROR"
		status = http.StatusBadRequest
	case errors.Is(err, service.ErrPageNotFound):
		code = "PAGE_NOT_FOUND"
		status = http.StatusNotFound
	case errors.Is(err, service.ErrSectionNotFound):
		code = "SECTION_NOT_FOUND"
		status = http.StatusNotFound
	case errors.Is(err, service.ErrSourceNotFound):
		code = "SOURCE_NOT_FOUND"
		status = http.StatusNotFound
	case errors.Is(err, service.ErrMediaForbidden):
		code = "FORBIDDEN"
		status = http.StatusForbidden
	case errors.Is(err, service.ErrMediaTypeInvalid):
		code = "VALIDATION_ERROR"
		status = http.StatusBadRequest
	case errors.Is(err, service.ErrSectionKeyConflict):
		code = "SECTION_KEY_CONFLICT"
		status = http.StatusConflict
	case errors.Is(err, service.ErrSectionOrderConflict):
		code = "SECTION_ORDER_CONFLICT"
		status = http.StatusConflict
	case errors.Is(err, service.ErrCollectionTypeConflict):
		code = "COLLECTION_TYPE_CONFLICT"
		status = http.StatusConflict
	case errors.Is(err, service.ErrSectionItemConflict):
		code = "SECTION_ITEM_CONFLICT"
		status = http.StatusConflict
	}
	c.JSON(status, gin.H{"error": code, "code": code, "message": err.Error()})
}
func bindURI(c *gin.Context, v interface{}) bool {
	if err := c.ShouldBindUri(v); err != nil {
		writeSectionError(c, service.ErrInvalidInput)
		return false
	}
	return true
}
func (c *PageSectionController) GetSections(ctx *gin.Context) {
	var u dto.PageSectionsURI
	var q dto.GetPageSectionsQuery
	if !bindURI(ctx, &u) {
		return
	}
	if err := ctx.ShouldBindQuery(&q); err != nil {
		writeSectionError(ctx, service.ErrInvalidInput)
		return
	}
	v, err := c.sectionService.GetSections(ctx.Request.Context(), u.PageID, q)
	if err != nil {
		writeSectionError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, v)
}
func (c *PageSectionController) CreateSection(ctx *gin.Context) {
	var u dto.PageSectionsURI
	var req dto.CreatePageSectionRequest
	if !bindURI(ctx, &u) {
		return
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		writeSectionError(ctx, service.ErrInvalidInput)
		return
	}
	id, role, err := sectionIdentity(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "UNAUTHORIZED", "code": "UNAUTHORIZED"})
		return
	}
	v, err := c.sectionService.CreateSection(ctx.Request.Context(), u.PageID, &req, id, role)
	if err != nil {
		writeSectionError(ctx, err)
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"data": v})
}
func (c *PageSectionController) GetSection(ctx *gin.Context) {
	var u dto.PageSectionURI
	if !bindURI(ctx, &u) {
		return
	}
	v, err := c.sectionService.GetSection(ctx.Request.Context(), u.PageID, u.SectionID)
	if err != nil {
		writeSectionError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": v})
}
func (c *PageSectionController) UpdateSection(ctx *gin.Context) {
	var u dto.PageSectionURI
	var req dto.UpdatePageSectionRequest
	if !bindURI(ctx, &u) {
		return
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		writeSectionError(ctx, service.ErrInvalidInput)
		return
	}
	id, role, err := sectionIdentity(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "UNAUTHORIZED", "code": "UNAUTHORIZED"})
		return
	}
	v, err := c.sectionService.UpdateSection(ctx.Request.Context(), u.PageID, u.SectionID, &req, id, role)
	if err != nil {
		writeSectionError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": v, "message": "Cập nhật Section thành công"})
}
func (c *PageSectionController) DeleteSection(ctx *gin.Context) {
	var u dto.PageSectionURI
	if !bindURI(ctx, &u) {
		return
	}
	if err := c.sectionService.DeleteSection(ctx.Request.Context(), u.PageID, u.SectionID); err != nil {
		writeSectionError(ctx, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}
func (c *PageSectionController) ReorderSections(ctx *gin.Context) {
	var u dto.PageSectionsURI
	var req dto.ReorderPageSectionsRequest
	if !bindURI(ctx, &u) {
		return
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		writeSectionError(ctx, service.ErrInvalidInput)
		return
	}
	if err := c.sectionService.ReorderSections(ctx.Request.Context(), u.PageID, &req); err != nil {
		writeSectionError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Sắp xếp Section thành công", "data": []string{}})
}
func (c *PageSectionController) GetItems(ctx *gin.Context) {
	var u dto.PageSectionURI
	var q dto.GetSectionItemsQuery
	if !bindURI(ctx, &u) {
		return
	}
	if err := ctx.ShouldBindQuery(&q); err != nil {
		writeSectionError(ctx, service.ErrInvalidInput)
		return
	}
	v, err := c.sectionService.GetItems(ctx.Request.Context(), u.PageID, u.SectionID, q)
	if err != nil {
		writeSectionError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, v)
}
func (c *PageSectionController) SyncItems(ctx *gin.Context) {
	var u dto.SectionCollectionURI
	var req dto.SyncSectionCollectionRequest
	if !bindURI(ctx, &u) {
		return
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		writeSectionError(ctx, service.ErrInvalidInput)
		return
	}
	id, role, err := sectionIdentity(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "UNAUTHORIZED", "code": "UNAUTHORIZED"})
		return
	}
	v, err := c.sectionService.SyncItems(ctx.Request.Context(), u.PageID, u.SectionID, u.Collection, &req, id, role)
	if err != nil {
		writeSectionError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Đồng bộ collection thành công", "data": v.Data})
}
