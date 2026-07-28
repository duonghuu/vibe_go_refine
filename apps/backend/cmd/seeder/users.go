package main

import (
	"log"

	"go_refine_dashboard_be/internal/domain/user/entity"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func SeedUsers(db *gorm.DB) {
	passwordStr := "Password123!"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(passwordStr), 12)
	if err != nil {
		log.Fatalf("Failed to hash password: %v", err)
	}

	users := []entity.User{
		{
			Email:    "admin@techbite.com",
			Password: string(hashedPassword),
			Name:     "System Admin",
			Role:     "ADMIN",
			Status:   "ACTIVE",
		},
		{
			Email:    "staff@techbite.com",
			Password: string(hashedPassword),
			Name:     "Operational Staff",
			Role:     "STAFF",
			Status:   "ACTIVE",
		},
		{
			Email:    "customer1@techbite.com",
			Password: string(hashedPassword),
			Name:     "Customer One",
			Role:     "CUSTOMER",
			Status:   "ACTIVE",
		},
		{
			Email:    "customer2@techbite.com",
			Password: string(hashedPassword),
			Name:     "Customer Two",
			Role:     "CUSTOMER",
			Status:   "INACTIVE",
		},
	}

	for _, user := range users {
		var existingUser entity.User
		result := db.Where("email = ?", user.Email).First(&existingUser)
		if result.Error != nil {
			if result.Error == gorm.ErrRecordNotFound {
				if err := db.Create(&user).Error; err != nil {
					log.Printf("Failed to seed user %s: %v", user.Email, err)
				} else {
					log.Printf("Seeded user: %s", user.Email)
				}
			} else {
				log.Printf("Error checking user %s: %v", user.Email, result.Error)
			}
		} else {
			log.Printf("User %s already exists, skipping.", user.Email)
		}
	}
}
