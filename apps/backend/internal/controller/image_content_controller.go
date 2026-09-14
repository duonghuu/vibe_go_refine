package controller

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"go_refine_dashboard_be/internal/infrastructure/queryservice"
	"go_refine_dashboard_be/internal/usecases/imagecontent/dto"
	"go_refine_dashboard_be/internal/usecases/imagecontent/service"
	"gorm.io/gorm"
)

type ImageContentController struct {
	command service.ImageContentService
	query   queryservice.ImageContentQueryService
}

func NewImageContentController(command service.ImageContentService, query queryservice.ImageContentQueryService) *ImageContentController {
	return &ImageContentController{command: command, query: query}
}

func imageContentIdentity(c *gin.Context) (uint, string, error) {
	value, ok := c.Get("userId")
	if !ok {
		return 0, "", errors.New("không tìm thấy thông tin người dùng")
	}
	userID, ok := value.(uint)
	if !ok || userID == 0 {
		return 0, "", errors.New("thông tin định danh không hợp lệ")
	}
	roleValue, ok := c.Get("role")
	if !ok {
		return 0, "", errors.New("không tìm thấy vai trò người dùng")
	}
	role, ok := roleValue.(string)
	if !ok || role == "" {
		return 0, "", errors.New("vai trò không hợp lệ")
	}
	return userID, role, nil
}

func writeImageContentError(c *gin.Context, err error) {
	status := http.StatusInternalServerError
	code := "INTERNAL_ERROR"
	message := "internal error"
	switch {
	case errors.Is(err, service.ErrTypeNotFound), errors.Is(err, service.ErrContentNotFound), errors.Is(err, service.ErrMediaNotFound), errors.Is(err, queryservice.ErrImageContentTypeNotFound), errors.Is(err, queryservice.ErrImageContentNotFound), errors.Is(err, gorm.ErrRecordNotFound):
		status, code, message = http.StatusNotFound, "NOT_FOUND", err.Error()
	case errors.Is(err, service.ErrForbidden), errors.Is(err, service.ErrMediaForbidden):
		status, code, message = http.StatusForbidden, "FORBIDDEN", err.Error()
	case errors.Is(err, service.ErrTypeCodeConflict):
		status, code, message = http.StatusConflict, "TYPE_CODE_CONFLICT", err.Error()
	case errors.Is(err, service.ErrTypeConfigConflict):
		status, code, message = http.StatusConflict, "TYPE_CONFIG_CONFLICT", err.Error()
	case errors.Is(err, service.ErrMaxItemsExceeded):
		status, code, message = http.StatusConflict, "MAX_ITEMS_EXCEEDED", err.Error()
	case errors.Is(err, service.ErrTypeInUse):
		status, code, message = http.StatusConflict, "TYPE_IN_USE", err.Error()
	case errors.Is(err, service.ErrOrderConflict):
		status, code, message = http.StatusConflict, "ORDER_CONFLICT", err.Error()
	case errors.Is(err, service.ErrFieldRequired):
		status, code, message = http.StatusBadRequest, "FIELD_REQUIRED", err.Error()
	case errors.Is(err, service.ErrFieldNotEnabled):
		status, code, message = http.StatusBadRequest, "FIELD_NOT_ENABLED", err.Error()
	case errors.Is(err, service.ErrMediaTypeInvalid):
		status, code, message = http.StatusBadRequest, "MEDIA_TYPE_INVALID", err.Error()
	case errors.Is(err, service.ErrURLInvalid), errors.Is(err, service.ErrTypeConfigInvalid), errors.Is(err, service.ErrValidation), errors.Is(err, service.ErrSortOrderInvalid), errors.Is(err, queryservice.ErrPublicTypeCodesInvalid):
		status, code, message = http.StatusBadRequest, "VALIDATION_ERROR", err.Error()
	}
	c.JSON(status, dto.ImageContentErrorResponse{Error: message, Code: code, Message: message})
}

func (c *ImageContentController) GetTypes(ctx *gin.Context) {
	var query dto.GetImageContentTypesQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ImageContentErrorResponse{Error: err.Error(), Code: "VALIDATION_ERROR", Message: "query không hợp lệ"})
		return
	}
	result, err := c.query.ListTypes(ctx.Request.Context(), query)
	if err != nil {
		writeImageContentError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, result)
}

func (c *ImageContentController) CreateType(ctx *gin.Context) {
	var req dto.CreateImageContentTypeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ImageContentErrorResponse{Error: err.Error(), Code: "VALIDATION_ERROR", Message: "body không hợp lệ"})
		return
	}
	id, err := c.command.CreateType(ctx.Request.Context(), &req)
	if err != nil {
		writeImageContentError(ctx, err)
		return
	}
	result, err := c.query.GetType(ctx.Request.Context(), id)
	if err != nil {
		writeImageContentError(ctx, err)
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"data": result})
}

func (c *ImageContentController) GetType(ctx *gin.Context) {
	var uri dto.ImageContentIDURI
	if err := ctx.ShouldBindUri(&uri); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ImageContentErrorResponse{Error: err.Error(), Code: "VALIDATION_ERROR", Message: "URI không hợp lệ"})
		return
	}
	result, err := c.query.GetType(ctx.Request.Context(), uri.ID)
	if err != nil {
		writeImageContentError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": result})
}

func (c *ImageContentController) UpdateType(ctx *gin.Context) {
	var uri dto.ImageContentIDURI
	var req dto.UpdateImageContentTypeRequest
	if err := ctx.ShouldBindUri(&uri); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ImageContentErrorResponse{Error: err.Error(), Code: "VALIDATION_ERROR", Message: "URI không hợp lệ"})
		return
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ImageContentErrorResponse{Error: err.Error(), Code: "VALIDATION_ERROR", Message: "body không hợp lệ"})
		return
	}
	if err := c.command.UpdateType(ctx.Request.Context(), uri.ID, &req); err != nil {
		writeImageContentError(ctx, err)
		return
	}
	result, err := c.query.GetType(ctx.Request.Context(), uri.ID)
	if err != nil {
		writeImageContentError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": result})
}

func (c *ImageContentController) DeleteType(ctx *gin.Context) {
	var uri dto.ImageContentIDURI
	if err := ctx.ShouldBindUri(&uri); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ImageContentErrorResponse{Error: err.Error(), Code: "VALIDATION_ERROR", Message: "URI không hợp lệ"})
		return
	}
	if err := c.command.DeleteType(ctx.Request.Context(), uri.ID); err != nil {
		writeImageContentError(ctx, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}

func (c *ImageContentController) GetContents(ctx *gin.Context) {
	var query dto.GetImageContentsQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ImageContentErrorResponse{Error: err.Error(), Code: "VALIDATION_ERROR", Message: "query không hợp lệ"})
		return
	}
	result, err := c.query.ListContents(ctx.Request.Context(), query)
	if err != nil {
		writeImageContentError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, result)
}

func (c *ImageContentController) CreateContent(ctx *gin.Context) {
	var req dto.CreateImageContentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ImageContentErrorResponse{Error: err.Error(), Code: "VALIDATION_ERROR", Message: "body không hợp lệ"})
		return
	}
	userID, role, err := imageContentIdentity(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, dto.ImageContentErrorResponse{Error: err.Error(), Code: "UNAUTHORIZED", Message: "unauthorized"})
		return
	}
	id, err := c.command.CreateContent(ctx.Request.Context(), &req, userID, role)
	if err != nil {
		writeImageContentError(ctx, err)
		return
	}
	result, err := c.query.GetContent(ctx.Request.Context(), id)
	if err != nil {
		writeImageContentError(ctx, err)
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"data": result})
}

func (c *ImageContentController) GetContent(ctx *gin.Context) {
	var uri dto.ImageContentIDURI
	if err := ctx.ShouldBindUri(&uri); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ImageContentErrorResponse{Error: err.Error(), Code: "VALIDATION_ERROR", Message: "URI không hợp lệ"})
		return
	}
	result, err := c.query.GetContent(ctx.Request.Context(), uri.ID)
	if err != nil {
		writeImageContentError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": result})
}

func (c *ImageContentController) UpdateContent(ctx *gin.Context) {
	var uri dto.ImageContentIDURI
	var req dto.UpdateImageContentRequest
	if err := ctx.ShouldBindUri(&uri); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ImageContentErrorResponse{Error: err.Error(), Code: "VALIDATION_ERROR", Message: "URI không hợp lệ"})
		return
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ImageContentErrorResponse{Error: err.Error(), Code: "VALIDATION_ERROR", Message: "body không hợp lệ"})
		return
	}
	userID, role, err := imageContentIdentity(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, dto.ImageContentErrorResponse{Error: err.Error(), Code: "UNAUTHORIZED", Message: "unauthorized"})
		return
	}
	if err := c.command.UpdateContent(ctx.Request.Context(), uri.ID, &req, userID, role); err != nil {
		writeImageContentError(ctx, err)
		return
	}
	result, err := c.query.GetContent(ctx.Request.Context(), uri.ID)
	if err != nil {
		writeImageContentError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": result})
}

func (c *ImageContentController) DeleteContent(ctx *gin.Context) {
	var uri dto.ImageContentIDURI
	if err := ctx.ShouldBindUri(&uri); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ImageContentErrorResponse{Error: err.Error(), Code: "VALIDATION_ERROR", Message: "URI không hợp lệ"})
		return
	}
	if err := c.command.DeleteContent(ctx.Request.Context(), uri.ID); err != nil {
		writeImageContentError(ctx, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}

func (c *ImageContentController) ReorderContents(ctx *gin.Context) {
	var req dto.ReorderImageContentsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ImageContentErrorResponse{Error: err.Error(), Code: "VALIDATION_ERROR", Message: "body không hợp lệ"})
		return
	}
	if _, _, err := imageContentIdentity(ctx); err != nil {
		ctx.JSON(http.StatusUnauthorized, dto.ImageContentErrorResponse{Error: err.Error(), Code: "UNAUTHORIZED", Message: "unauthorized"})
		return
	}
	if err := c.command.ReorderContents(ctx.Request.Context(), &req); err != nil {
		writeImageContentError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": []dto.ImageContentResponse{}, "message": "sắp xếp image content thành công"})
}

func (c *ImageContentController) GetPublicContents(ctx *gin.Context) {
	var query dto.GetPublicImageContentsQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ImageContentErrorResponse{Error: err.Error(), Code: "VALIDATION_ERROR", Message: "query không hợp lệ"})
		return
	}
	result, err := c.query.GetPublicContents(ctx.Request.Context(), query)
	if err != nil {
		writeImageContentError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, result)
}
