package controller

import (
	"net/http"

	"go_refine_dashboard_be/internal/usecases/auth/dto"
	"go_refine_dashboard_be/internal/usecases/auth/service"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	authService service.AuthService
}

func NewAuthController(authService service.AuthService) *AuthController {
	return &AuthController{
		authService: authService,
	}
}

func (c *AuthController) Login(ctx *gin.Context) {
	var req dto.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, refreshToken, err := c.authService.Login(ctx.Request.Context(), req)
	if err != nil {
		status := http.StatusInternalServerError
		if err == service.ErrInvalidCredentials || err == service.ErrUserInactive {
			status = http.StatusUnauthorized
		}
		ctx.JSON(status, gin.H{"error": err.Error()})
		return
	}

	// Set refresh token in HttpOnly cookie
	ctx.SetCookie("refreshToken", refreshToken, int(service.RefreshTokenTTL.Seconds()), "/", "", false, true) // set secure=true in prod

	ctx.JSON(http.StatusOK, res)
}

func (c *AuthController) RefreshToken(ctx *gin.Context) {
	refreshToken, err := ctx.Cookie("refreshToken")
	if err != nil || refreshToken == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "missing refresh token"})
		return
	}

	res, newRefreshToken, err := c.authService.RefreshToken(ctx.Request.Context(), refreshToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Set new refresh token in HttpOnly cookie
	ctx.SetCookie("refreshToken", newRefreshToken, int(service.RefreshTokenTTL.Seconds()), "/", "", false, true)

	ctx.JSON(http.StatusOK, res)
}

func (c *AuthController) Logout(ctx *gin.Context) {
	accessToken := ""
	authHeader := ctx.GetHeader("Authorization")
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		accessToken = authHeader[7:]
	}

	refreshToken, _ := ctx.Cookie("refreshToken")

	if err := c.authService.Logout(ctx.Request.Context(), accessToken, refreshToken); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Clear cookie
	ctx.SetCookie("refreshToken", "", -1, "/", "", false, true)
	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}

func (c *AuthController) GetMe(ctx *gin.Context) {
	userIdObj, exists := ctx.Get("userId")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	userId := userIdObj.(uint)
	res, err := c.authService.GetMe(ctx.Request.Context(), userId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, res)
}
