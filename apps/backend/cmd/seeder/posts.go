package main

import (
	"fmt"
	"log"

	postEntity "go_refine_dashboard_be/internal/domain/post/entity"
	postCategoryEntity "go_refine_dashboard_be/internal/domain/postcategory/entity"
	postTypeEntity "go_refine_dashboard_be/internal/domain/posttype/entity"
	userEntity "go_refine_dashboard_be/internal/domain/user/entity"

	"gorm.io/gorm"
)

const postsPerCategory = 6

func SeedPosts(db *gorm.DB) {
	var author userEntity.User
	if err := db.Where("email = ?", "admin@techbite.com").First(&author).Error; err != nil {
		log.Printf("Skipping posts: seed author not found: %v", err)
		return
	}

	for _, typeCode := range []string{"NEWS", "SERVICE"} {
		var postType postTypeEntity.PostType
		if err := db.Where("code = ?", typeCode).First(&postType).Error; err != nil {
			log.Printf("Skipping posts for type %s: post type not found: %v", typeCode, err)
			continue
		}

		var categories []postCategoryEntity.PostCategory
		if err := db.Where("type_code = ?", typeCode).Order("id ASC").Find(&categories).Error; err != nil {
			log.Printf("Failed to load post categories for type %s: %v", typeCode, err)
			continue
		}

		for _, category := range categories {
			seedPostsForCategory(db, typeCode, category, author.ID)
		}
	}
}

func seedPostsForCategory(db *gorm.DB, typeCode string, category postCategoryEntity.PostCategory, authorID uint) {
	for index := 1; index <= postsPerCategory; index++ {
		post := postEntity.Post{
			TypeCode:   typeCode,
			Title:      fmt.Sprintf("%s - Bài viết số %d", category.Name, index),
			Slug:       fmt.Sprintf("%s-bai-viet-%02d", category.Slug, index),
			Content:    fmt.Sprintf("Nội dung bài viết số %d thuộc danh mục %s.", index, category.Name),
			AuthorID:   authorID,
			CategoryID: &category.ID,
		}

		var existing postEntity.Post
		result := db.Where("slug = ?", post.Slug).First(&existing)
		if result.Error == nil {
			log.Printf("Post %s already exists, skipping.", post.Slug)
			continue
		}

		if result.Error != gorm.ErrRecordNotFound {
			log.Printf("Failed to check post %s: %v", post.Slug, result.Error)
			continue
		}

		if err := db.Create(&post).Error; err != nil {
			log.Printf("Failed to seed post %s: %v", post.Slug, err)
			continue
		}

		log.Printf("Seeded post: %s", post.Slug)
	}
}
