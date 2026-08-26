package main

import (
	"log"

	postTypeEntity "go_refine_dashboard_be/internal/domain/posttype/entity"

	"gorm.io/gorm"
)

func SeedPostTypes(db *gorm.DB) {
	postTypes := []postTypeEntity.PostType{
		{Code: "NEWS", Name: "Tin tức", Status: "ACTIVE", SortOrder: 1},
		{Code: "SERVICE", Name: "Dịch vụ", Status: "ACTIVE", SortOrder: 2},
	}

	for i := range postTypes {
		var existing postTypeEntity.PostType
		result := db.Where("code = ?", postTypes[i].Code).First(&existing)
		if result.Error == nil {
			log.Printf("Post type %s already exists, skipping.", postTypes[i].Code)
			continue
		}

		if result.Error != gorm.ErrRecordNotFound {
			log.Printf("Failed to check post type %s: %v", postTypes[i].Code, result.Error)
			continue
		}

		if err := db.Create(&postTypes[i]).Error; err != nil {
			log.Printf("Failed to seed post type %s: %v", postTypes[i].Code, err)
			continue
		}

		log.Printf("Seeded post type: %s", postTypes[i].Code)
	}
}
