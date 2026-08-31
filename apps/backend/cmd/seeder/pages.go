package main

import (
	"log"

	pageEntity "go_refine_dashboard_be/internal/domain/page/entity"
	userEntity "go_refine_dashboard_be/internal/domain/user/entity"

	"gorm.io/gorm"
)

// SeedPages creates the static pages used by the public website and admin UI.
// Slugs are used as the idempotency key so this seeder is safe to run repeatedly.
func SeedPages(db *gorm.DB) {
	var author userEntity.User
	if err := db.Where("email = ?", "admin@techbite.com").First(&author).Error; err != nil {
		log.Printf("Skipping pages: seed author not found: %v", err)
		return
	}

	pages := []pageEntity.Page{
		{
			Title:   "Giới thiệu TechBite",
			Slug:    "gioi-thieu",
			Content: "TechBite mang đến những món ăn ngon, nguyên liệu chất lượng và trải nghiệm đặt món tiện lợi cho mọi khách hàng.",
			Status:  pageEntity.PageStatusPublished,
		},
		{
			Title:   "Thực đơn",
			Slug:    "thuc-don",
			Content: "Khám phá các món burger, pizza, đồ uống và món chay được phục vụ mỗi ngày tại TechBite.",
			Status:  pageEntity.PageStatusPublished,
		},
		{
			Title:   "Chính sách giao hàng",
			Slug:    "chinh-sach-giao-hang",
			Content: "TechBite tiếp nhận đơn hàng trực tuyến và giao món đến khu vực phục vụ trong thời gian dự kiến.",
			Status:  pageEntity.PageStatusPublished,
		},
		{
			Title:   "Điều khoản sử dụng",
			Slug:    "dieu-khoan-su-dung",
			Content: "Vui lòng đọc các điều khoản sử dụng trước khi đặt món và sử dụng dịch vụ của TechBite.",
			Status:  pageEntity.PageStatusDraft,
		},
	}

	for i := range pages {
		pages[i].AuthorID = author.ID

		var existing pageEntity.Page
		result := db.Where("slug = ?", pages[i].Slug).First(&existing)
		if result.Error == nil {
			log.Printf("Page %s already exists, skipping.", pages[i].Slug)
			continue
		}
		if result.Error != gorm.ErrRecordNotFound {
			log.Printf("Failed to check page %s: %v", pages[i].Slug, result.Error)
			continue
		}
		if err := db.Create(&pages[i]).Error; err != nil {
			log.Printf("Failed to seed page %s: %v", pages[i].Slug, err)
			continue
		}
		log.Printf("Seeded page: %s", pages[i].Slug)
	}
}
