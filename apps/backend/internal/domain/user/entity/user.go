package entity

import (
	"time"

	"gorm.io/gorm"
)

// User đại diện cho tài khoản người dùng trong hệ thống
type User struct {
	ID        uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	Email     string         `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	Password  string         `gorm:"type:varchar(255);not null" json:"-"`
	Name      string         `gorm:"type:varchar(255);not null" json:"name"`
	Role      string         `gorm:"type:varchar(20);default:'CUSTOMER';index" json:"role"` // ADMIN, STAFF, CUSTOMER
	Status    string         `gorm:"type:varchar(20);default:'ACTIVE';index" json:"status"` // ACTIVE, INACTIVE
	LastLogin *time.Time     `gorm:"index" json:"lastLogin"`

	// Audit logs & Soft Delete
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (User) TableName() string {
	return "users"
}
