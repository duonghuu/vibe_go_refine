package dto

// LoginRequest binding JSON từ client gửi lên
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// UserProfile chứa thông tin công khai trả về cho User
type UserProfile struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
	Role  string `json:"role"`
}

// TokenResponse payload trả về cho client
type TokenResponse struct {
	AccessToken string       `json:"accessToken"`
	ExpiresIn   int          `json:"expiresIn"`
	User        *UserProfile `json:"user,omitempty"` // Có thể nil nếu dùng api refresh token
}
