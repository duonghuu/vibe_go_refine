package main

import (
	"log"

	catEntity "go_refine_dashboard_be/internal/domain/category/entity"

	"gorm.io/gorm"
)

func SeedCategories(db *gorm.DB) {
	categories := []catEntity.Category{
		{Name: "Burger", Slug: "burger", SortOrder: 1, Status: "ACTIVE"},
		{Name: "Pizza", Slug: "pizza", SortOrder: 2, Status: "ACTIVE"},
		{Name: "Đồ uống", Slug: "do-uong", SortOrder: 3, Status: "ACTIVE"},
		{Name: "Món chay", Slug: "mon-chay", SortOrder: 4, Status: "HIDDEN"},
	}

	for i := range categories {
		var existing catEntity.Category
		if err := db.Where("slug = ?", categories[i].Slug).First(&existing).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				db.Create(&categories[i])
				log.Printf("Seeded category: %s", categories[i].Name)
			}
		} else {
			categories[i] = existing
			log.Printf("Category %s already exists, skipping.", categories[i].Name)
		}
	}

	var doUong catEntity.Category
	if err := db.Where("slug = ?", "do-uong").First(&doUong).Error; err == nil {
		children := []catEntity.Category{
			{Name: "Cà phê", Slug: "ca-phe", SortOrder: 1, Status: "ACTIVE", ParentID: &doUong.ID},
			{Name: "Trà trái cây", Slug: "tra-trai-cay", SortOrder: 2, Status: "ACTIVE", ParentID: &doUong.ID},
		}

		for i := range children {
			var existing catEntity.Category
			if err := db.Where("slug = ?", children[i].Slug).First(&existing).Error; err != nil {
				if err == gorm.ErrRecordNotFound {
					db.Create(&children[i])
					log.Printf("Seeded category child: %s", children[i].Name)
				}
			} else {
				log.Printf("Category child %s already exists, skipping.", children[i].Name)
			}
		}
	}
}
