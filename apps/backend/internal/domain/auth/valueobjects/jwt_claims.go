package valueobjects

import "github.com/golang-jwt/jwt/v5"

// JwtCustomClaims định nghĩa Payload trong JWT Token
type JwtCustomClaims struct {
	UserID uint   `json:"userId"`
	Role   string `json:"role"`
	Jti    string `json:"jti"` // JWT ID - Dùng để whitelist/blacklist trên Redis
	jwt.RegisteredClaims
}
