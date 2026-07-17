package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"go_refine_dashboard_be/internal/domain/auth/valueobjects"
	"go_refine_dashboard_be/internal/infrastructure/repository"
	"go_refine_dashboard_be/internal/usecases/auth/dto"

	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrUserInactive       = errors.New("user is inactive")
	ErrInvalidToken       = errors.New("invalid token")
)

const (
	AccessTokenTTL  = 30 * time.Minute
	RefreshTokenTTL = 7 * 24 * time.Hour
	JWTSecret       = "my-super-secret-key" // Should use env var in prod
)

type AuthService interface {
	Login(ctx context.Context, req dto.LoginRequest) (*dto.TokenResponse, string, error)
	RefreshToken(ctx context.Context, oldToken string) (*dto.TokenResponse, string, error)
	Logout(ctx context.Context, accessToken, refreshToken string) error
	GetMe(ctx context.Context, userId uint) (*dto.UserProfile, error)
}

type authService struct {
	userRepo    repository.UserRepository
	redisClient *redis.Client
}

func NewAuthService(userRepo repository.UserRepository, redisClient *redis.Client) AuthService {
	return &authService{
		userRepo:    userRepo,
		redisClient: redisClient,
	}
}

func generateJti() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

func generateTokens(userId uint, role string) (string, string, string, error) {
	jti := generateJti()
	claims := valueobjects.JwtCustomClaims{
		UserID: userId,
		Role:   role,
		Jti:    jti,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(AccessTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, err := token.SignedString([]byte(JWTSecret))
	if err != nil {
		return "", "", "", err
	}

	// For refresh token, we can also use JWT or just a random string.
	// We'll use JWT for consistency to carry the jti and userId.
	rtClaims := valueobjects.JwtCustomClaims{
		UserID: userId,
		Role:   role,
		Jti:    jti, // shared JTI for the session
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(RefreshTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	rtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	refreshToken, err := rtToken.SignedString([]byte(JWTSecret))
	if err != nil {
		return "", "", "", err
	}

	return accessToken, refreshToken, jti, nil
}

func (s *authService) saveSessionToRedis(ctx context.Context, userId uint, jti string) error {
	key := fmt.Sprintf("auth:rf:%d:%s", userId, jti)
	return s.redisClient.SetEx(ctx, key, "valid", RefreshTokenTTL).Err()
}

func (s *authService) removeSessionFromRedis(ctx context.Context, userId uint, jti string) error {
	key := fmt.Sprintf("auth:rf:%d:%s", userId, jti)
	return s.redisClient.Del(ctx, key).Err()
}

func (s *authService) Login(ctx context.Context, req dto.LoginRequest) (*dto.TokenResponse, string, error) {
	user, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, "", ErrInvalidCredentials
	}

	if user.Status != "ACTIVE" {
		return nil, "", ErrUserInactive
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, "", ErrInvalidCredentials
	}

	accessToken, refreshToken, jti, err := generateTokens(user.ID, user.Role)
	if err != nil {
		return nil, "", err
	}

	if err := s.saveSessionToRedis(ctx, user.ID, jti); err != nil {
		return nil, "", err
	}

	// Update last login
	now := time.Now()
	user.LastLogin = &now
	_ = s.userRepo.Update(ctx, user)

	return &dto.TokenResponse{
		AccessToken: accessToken,
		ExpiresIn:   int(AccessTokenTTL.Seconds()),
		User: &dto.UserProfile{
			ID:    user.ID,
			Email: user.Email,
			Name:  user.Name,
			Role:  user.Role,
		},
	}, refreshToken, nil
}

func (s *authService) parseToken(tokenStr string) (*valueobjects.JwtCustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &valueobjects.JwtCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(JWTSecret), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*valueobjects.JwtCustomClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, ErrInvalidToken
}

func (s *authService) RefreshToken(ctx context.Context, oldRefreshToken string) (*dto.TokenResponse, string, error) {
	claims, err := s.parseToken(oldRefreshToken)
	if err != nil {
		return nil, "", ErrInvalidToken
	}

	// Check if jti exists in Redis
	key := fmt.Sprintf("auth:rf:%d:%s", claims.UserID, claims.Jti)
	exists, err := s.redisClient.Exists(ctx, key).Result()
	if err != nil || exists == 0 {
		return nil, "", ErrInvalidToken
	}

	// Revoke old session
	_ = s.removeSessionFromRedis(ctx, claims.UserID, claims.Jti)

	// Generate new tokens
	accessToken, newRefreshToken, newJti, err := generateTokens(claims.UserID, claims.Role)
	if err != nil {
		return nil, "", err
	}

	// Save new session
	if err := s.saveSessionToRedis(ctx, claims.UserID, newJti); err != nil {
		return nil, "", err
	}

	return &dto.TokenResponse{
		AccessToken: accessToken,
		ExpiresIn:   int(AccessTokenTTL.Seconds()),
	}, newRefreshToken, nil
}

func (s *authService) Logout(ctx context.Context, accessToken, refreshToken string) error {
	// 1. Revoke refresh token
	if rtClaims, err := s.parseToken(refreshToken); err == nil {
		_ = s.removeSessionFromRedis(ctx, rtClaims.UserID, rtClaims.Jti)
	}

	// 2. Blacklist access token
	atClaims, err := s.parseToken(accessToken)
	if err == nil {
		ttl := time.Until(atClaims.ExpiresAt.Time)
		if ttl > 0 {
			key := fmt.Sprintf("auth:bl:at:%s", atClaims.Jti) // using jti to blacklist
			_ = s.redisClient.SetEx(ctx, key, "revoked", ttl).Err()
		}
	}
	return nil
}

func (s *authService) GetMe(ctx context.Context, userId uint) (*dto.UserProfile, error) {
	user, err := s.userRepo.FindByID(ctx, userId)
	if err != nil {
		return nil, err
	}
	return &dto.UserProfile{
		ID:    user.ID,
		Email: user.Email,
		Name:  user.Name,
		Role:  user.Role,
	}, nil
}
