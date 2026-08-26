package main

import (
	"encoding/json"
	"fmt"
	"log"

	postEntity "go_refine_dashboard_be/internal/domain/post/entity"
	seoEntity "go_refine_dashboard_be/internal/domain/seometa/entity"

	"gorm.io/gorm"
)

const seoPostEntityType = "post"

func SeedPostSEOMeta(db *gorm.DB) {
	var posts []postEntity.Post
	if err := db.Order("id ASC").Find(&posts).Error; err != nil {
		log.Printf("Failed to load posts for SEO seeder: %v", err)
		return
	}

	for _, post := range posts {
		var existing seoEntity.SEOMeta
		result := db.Where("entity_type = ? AND entity_id = ?", seoPostEntityType, post.ID).First(&existing)
		if result.Error == nil {
			log.Printf("SEO meta for post %d already exists, skipping.", post.ID)
			continue
		}

		if result.Error != gorm.ErrRecordNotFound {
			log.Printf("Failed to check SEO meta for post %d: %v", post.ID, result.Error)
			continue
		}

		metaTitle := fmt.Sprintf("%s | TechBite", post.Title)
		metaDescription := fmt.Sprintf("Đọc bài viết %s trên TechBite.", post.Title)
		keywords := fmt.Sprintf("%s, %s, TechBite", post.TypeCode, post.Title)
		canonicalURL := fmt.Sprintf("https://techbite.example/posts/%s", post.Slug)
		ogType := "article"
		twitterCard := "summary_large_image"
		robots := "index,follow"

		schemaJSON, err := json.Marshal(struct {
			Context  string `json:"@context"`
			Type     string `json:"@type"`
			Headline string `json:"headline"`
			URL      string `json:"url"`
		}{
			Context:  "https://schema.org",
			Type:     "Article",
			Headline: post.Title,
			URL:      canonicalURL,
		})
		if err != nil {
			log.Printf("Failed to generate schema JSON for post %d: %v", post.ID, err)
			continue
		}

		meta := seoEntity.SEOMeta{
			EntityType:         seoPostEntityType,
			EntityID:           post.ID,
			MetaTitle:          &metaTitle,
			MetaDescription:    &metaDescription,
			MetaKeywords:       &keywords,
			CanonicalURL:       &canonicalURL,
			OGTitle:            &metaTitle,
			OGDescription:      &metaDescription,
			OGType:             &ogType,
			TwitterTitle:       &metaTitle,
			TwitterDescription: &metaDescription,
			TwitterCard:        &twitterCard,
			Robots:             &robots,
			SchemaJSON:         schemaJSON,
		}

		if err := db.Create(&meta).Error; err != nil {
			log.Printf("Failed to seed SEO meta for post %d: %v", post.ID, err)
			continue
		}

		log.Printf("Seeded SEO meta for post: %d", post.ID)
	}
}
