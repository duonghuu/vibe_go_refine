package main

import (
	"log"

	catEntity "go_refine_dashboard_be/internal/domain/category/entity"
	prodEntity "go_refine_dashboard_be/internal/domain/product/entity"

	"gorm.io/gorm"
)

func SeedProducts(db *gorm.DB) {
	var burgerCat, pizzaCat, cafeCat, traCat catEntity.Category
	db.Where("slug = ?", "burger").First(&burgerCat)
	db.Where("slug = ?", "pizza").First(&pizzaCat)
	db.Where("slug = ?", "ca-phe").First(&cafeCat)
	db.Where("slug = ?", "tra-trai-cay").First(&traCat)

	funcPtr := func(f float64) *float64 { return &f }

	products := []prodEntity.Product{
		{
			SKU:         "TB-0001",
			Slug:        "burger-bo-pho-mai",
			Name:        "Burger Bò Phô Mai",
			Description: "Burger bò phô mai thơm ngon",
			CategoryID:  burgerCat.ID,
			Price:       65000,
			SalePrice:   funcPtr(59000),
			Stock:       100,
			SoldCount:   50,
			Status:      "ACTIVE",
			ImageURL:    "https://images.unsplash.com/photo-1568901346375-23c9450c58cd?auto=format&fit=crop&q=80&w=800",
		},
		{
			SKU:         "TB-0002",
			Slug:        "burger-ga-gion",
			Name:        "Burger Gà Giòn",
			Description: "Burger gà chiên giòn tan",
			CategoryID:  burgerCat.ID,
			Price:       55000,
			Stock:       150,
			SoldCount:   120,
			Status:      "ACTIVE",
			ImageURL:    "https://images.unsplash.com/photo-1615719413546-198b25453f85?auto=format&fit=crop&q=80&w=800",
		},
		{
			SKU:         "TB-0003",
			Slug:        "pizza-hai-san",
			Name:        "Pizza Hải Sản Nhiệt Đới",
			Description: "Pizza hải sản ngập tràn topping",
			CategoryID:  pizzaCat.ID,
			Price:       189000,
			SalePrice:   funcPtr(159000),
			Stock:       50,
			SoldCount:   200,
			Status:      "ACTIVE",
			ImageURL:    "https://images.unsplash.com/photo-1565299624946-b28f40a0ae38?auto=format&fit=crop&q=80&w=800",
		},
		{
			SKU:         "TB-0004",
			Slug:        "pizza-pho-mai-mat-ong",
			Name:        "Pizza Phô Mai Mật Ong",
			Description: "Pizza phô mai béo ngậy cùng mật ong",
			CategoryID:  pizzaCat.ID,
			Price:       149000,
			Stock:       0,
			SoldCount:   150,
			Status:      "OUT_OF_STOCK",
			ImageURL:    "https://images.unsplash.com/photo-1513104890138-7c749659a591?auto=format&fit=crop&q=80&w=800",
		},
		{
			SKU:         "TB-0005",
			Slug:        "ca-phe-sua-da",
			Name:        "Cà Phê Sữa Đá",
			Description: "Cà phê sữa đá pha phin truyền thống",
			CategoryID:  cafeCat.ID,
			Price:       29000,
			Stock:       200,
			SoldCount:   500,
			Status:      "ACTIVE",
			ImageURL:    "https://images.unsplash.com/photo-1559802619-33d3bb7eec28?auto=format&fit=crop&q=80&w=800",
		},
		{
			SKU:         "TB-0006",
			Slug:        "tra-dao-cam-sa",
			Name:        "Trà Đào Cam Sả",
			Description: "Trà đào thanh mát giải nhiệt",
			CategoryID:  traCat.ID,
			Price:       45000,
			SalePrice:   funcPtr(39000),
			Stock:       80,
			SoldCount:   300,
			Status:      "ACTIVE",
			ImageURL:    "https://images.unsplash.com/photo-1556679343-c7306c1976bc?auto=format&fit=crop&q=80&w=800",
		},
		{
			SKU:         "TB-0007",
			Slug:        "tra-vai-nhai",
			Name:        "Trà Vải Lài",
			Description: "Trà vải lài thơm lừng",
			CategoryID:  traCat.ID,
			Price:       45000,
			Stock:       100,
			SoldCount:   10,
			Status:      "HIDDEN",
			ImageURL:    "https://images.unsplash.com/photo-1556679343-c7306c1976bc?auto=format&fit=crop&q=80&w=800",
		},
	}

	for i := range products {
		if products[i].CategoryID == 0 {
			log.Printf("Skipping product %s because category is missing", products[i].Name)
			continue
		}
		var existing prodEntity.Product
		if err := db.Where("sku = ?", products[i].SKU).First(&existing).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				db.Create(&products[i])
				log.Printf("Seeded product: %s", products[i].Name)
			}
		} else {
			log.Printf("Product %s already exists, skipping.", products[i].Name)
		}
	}
}
