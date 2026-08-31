package main

import (
	"log"

	mediaEntity "go_refine_dashboard_be/internal/domain/media/entity"
	pageEntity "go_refine_dashboard_be/internal/domain/page/entity"
	userEntity "go_refine_dashboard_be/internal/domain/user/entity"

	"gorm.io/gorm"
)

type pageMediaSeed struct {
	PageSlug     string
	FileName     string
	OriginalName string
	URL          string
	Collection   pageEntity.PageMediaCollection
	SortOrder    int
}

// SeedPageMedia also creates the media rows required by the page_media foreign key.
func SeedPageMedia(db *gorm.DB) {
	var owner userEntity.User
	if err := db.Where("email = ?", "admin@techbite.com").First(&owner).Error; err != nil {
		log.Printf("Skipping page media: seed owner not found: %v", err)
		return
	}

	seeds := []pageMediaSeed{
		{PageSlug: "gioi-thieu", FileName: "page-gioi-thieu-cover.jpg", OriginalName: "gioi-thieu-cover.jpg", URL: "https://images.unsplash.com/photo-1517248135467-4c7edcad34c4?auto=format&fit=crop&q=80&w=1200", Collection: pageEntity.PageMediaCollectionThumbnail},
		{PageSlug: "gioi-thieu", FileName: "page-gioi-thieu-team.jpg", OriginalName: "gioi-thieu-team.jpg", URL: "https://images.unsplash.com/photo-1556761175-b413da4baf72?auto=format&fit=crop&q=80&w=1200", Collection: pageEntity.PageMediaCollectionGallery, SortOrder: 0},
		{PageSlug: "thuc-don", FileName: "page-thuc-don-cover.jpg", OriginalName: "thuc-don-cover.jpg", URL: "https://images.unsplash.com/photo-1504670073073-7a66eab7d7d1?auto=format&fit=crop&q=80&w=1200", Collection: pageEntity.PageMediaCollectionThumbnail},
		{PageSlug: "chinh-sach-giao-hang", FileName: "page-giao-hang-cover.jpg", OriginalName: "giao-hang-cover.jpg", URL: "https://images.unsplash.com/photo-1526367790999-0150786686a2?auto=format&fit=crop&q=80&w=1200", Collection: pageEntity.PageMediaCollectionThumbnail},
	}

	for _, seed := range seeds {
		var page pageEntity.Page
		if err := db.Where("slug = ?", seed.PageSlug).First(&page).Error; err != nil {
			log.Printf("Skipping page media for %s: page not found: %v", seed.PageSlug, err)
			continue
		}

		var media mediaEntity.Media
		result := db.Where("file_name = ?", seed.FileName).First(&media)
		if result.Error == gorm.ErrRecordNotFound {
			media = mediaEntity.Media{
				OriginalName: seed.OriginalName,
				FileName:     seed.FileName,
				MimeType:     "image/jpeg",
				Size:         250000,
				OriginalUrl:  seed.URL,
				ThumbnailUrl: seed.URL,
				MediumUrl:    seed.URL,
				Status:       mediaEntity.MediaStatusAttached,
				OwnerID:      owner.ID,
			}
			if err := db.Create(&media).Error; err != nil {
				log.Printf("Failed to seed media %s: %v", seed.FileName, err)
				continue
			}
			log.Printf("Seeded page media asset: %s", seed.FileName)
		} else if result.Error != nil {
			log.Printf("Failed to check media %s: %v", seed.FileName, result.Error)
			continue
		}

		var existing pageEntity.PageMedia
		result = db.Where("page_id = ? AND media_id = ? AND collection = ?", page.ID, media.ID, seed.Collection).First(&existing)
		if result.Error == nil {
			log.Printf("Page media %s/%s already exists, skipping.", seed.PageSlug, seed.Collection)
			continue
		}
		if result.Error != gorm.ErrRecordNotFound {
			log.Printf("Failed to check page media %s/%s: %v", seed.PageSlug, seed.Collection, result.Error)
			continue
		}

		link := pageEntity.PageMedia{PageID: page.ID, MediaID: media.ID, Collection: seed.Collection, SortOrder: seed.SortOrder}
		if err := db.Create(&link).Error; err != nil {
			log.Printf("Failed to seed page media %s/%s: %v", seed.PageSlug, seed.Collection, err)
			continue
		}
		log.Printf("Seeded page media link: %s/%s", seed.PageSlug, seed.Collection)
	}
}
