package dto

type CreateUserRequest struct {
	Email    string `json:"email" binding:"required,email,max=255"`
	Password string `json:"password" binding:"required,min=8,max=255"`
	Name     string `json:"name" binding:"required,max=255"`
	Role     string `json:"role" binding:"required,oneof=ADMIN STAFF CUSTOMER"`
	Status   string `json:"status" binding:"omitempty,oneof=ACTIVE INACTIVE"`
}

type UpdateUserRequest struct {
	Name   *string `json:"name" binding:"omitempty,max=255"`
	Role   *string `json:"role" binding:"omitempty,oneof=ADMIN STAFF CUSTOMER"`
	Status *string `json:"status" binding:"omitempty,oneof=ACTIVE INACTIVE"`
}

type UpdateUserStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=ACTIVE INACTIVE"`
}

type BulkUpdateUserStatusRequest struct {
	IDs    []uint `json:"ids" binding:"required,min=1"`
	Status string `json:"status" binding:"required,oneof=ACTIVE INACTIVE"`
}

type ResetPasswordRequest struct {
	NewPassword string `json:"newPassword" binding:"required,min=8,max=255"`
}
