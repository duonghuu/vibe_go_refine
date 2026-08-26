package main

import (
	"log"

	postCategoryEntity "go_refine_dashboard_be/internal/domain/postcategory/entity"
	postTypeEntity "go_refine_dashboard_be/internal/domain/posttype/entity"

	"gorm.io/gorm"
)

type postCategorySeed struct {
	TypeCode    string
	Name        string
	Slug        string
	Description string
	SortOrder   int
	Children    []postCategorySeed
}

func SeedPostCategories(db *gorm.DB) {
	categoriesByType := []postCategorySeed{
		{
			TypeCode: "NEWS",
			Children: []postCategorySeed{
				{
					Name:        "Tin tức sản phẩm",
					Slug:        "tin-tuc-san-pham",
					Description: "Các tin tức và cập nhật liên quan đến sản phẩm.",
					SortOrder:   1,
					Children: []postCategorySeed{
						{Name: "Món mới", Slug: "mon-moi", Description: "Các món ăn mới ra mắt.", SortOrder: 1},
						{Name: "Khuyến mãi", Slug: "khuyen-mai", Description: "Thông tin chương trình khuyến mãi.", SortOrder: 2},
					},
				},
				{
					Name:        "Tin tức doanh nghiệp",
					Slug:        "tin-tuc-doanh-nghiep",
					Description: "Các tin tức hoạt động của doanh nghiệp.",
					SortOrder:   2,
					Children: []postCategorySeed{
						{Name: "Hoạt động doanh nghiệp", Slug: "hoat-dong-doanh-nghiep", Description: "Hoạt động và sự kiện doanh nghiệp.", SortOrder: 1},
						{Name: "Tuyển dụng", Slug: "tuyen-dung", Description: "Thông tin tuyển dụng và cơ hội nghề nghiệp.", SortOrder: 2},
					},
				},
			},
		},
		{
			TypeCode: "SERVICE",
			Children: []postCategorySeed{
				{
					Name:        "Dịch vụ đặt món",
					Slug:        "dich-vu-dat-mon",
					Description: "Thông tin về các hình thức đặt món.",
					SortOrder:   1,
					Children: []postCategorySeed{
						{Name: "Giao hàng", Slug: "giao-hang", Description: "Hướng dẫn và chính sách giao hàng.", SortOrder: 1},
						{Name: "Tại cửa hàng", Slug: "tai-cua-hang", Description: "Thông tin dịch vụ tại cửa hàng.", SortOrder: 2},
					},
				},
				{
					Name:        "Dịch vụ khách hàng",
					Slug:        "dich-vu-khach-hang",
					Description: "Thông tin hỗ trợ và chăm sóc khách hàng.",
					SortOrder:   2,
					Children: []postCategorySeed{
						{Name: "Hỗ trợ đơn hàng", Slug: "ho-tro-don-hang", Description: "Hướng dẫn hỗ trợ và xử lý đơn hàng.", SortOrder: 1},
						{Name: "Thành viên", Slug: "thanh-vien", Description: "Quyền lợi và chương trình thành viên.", SortOrder: 2},
					},
				},
			},
		},
	}

	for _, typeCategories := range categoriesByType {
		var postType postTypeEntity.PostType
		if err := db.Where("code = ?", typeCategories.TypeCode).First(&postType).Error; err != nil {
			log.Printf("Skipping post categories for type %s: post type not found: %v", typeCategories.TypeCode, err)
			continue
		}

		for _, category := range typeCategories.Children {
			seedPostCategoryTree(db, typeCategories.TypeCode, category, nil)
		}
	}
}

func seedPostCategoryTree(db *gorm.DB, typeCode string, seed postCategorySeed, parentID *uint) *uint {
	category := postCategoryEntity.PostCategory{
		TypeCode:    typeCode,
		Name:        seed.Name,
		Slug:        seed.Slug,
		ParentID:    parentID,
		Description: seed.Description,
		SortOrder:   seed.SortOrder,
		Status:      "ACTIVE",
	}

	var existing postCategoryEntity.PostCategory
	result := db.Where("slug = ?", seed.Slug).First(&existing)
	if result.Error == nil {
		log.Printf("Post category %s already exists, skipping.", seed.Slug)
		categoryID := existing.ID
		for _, child := range seed.Children {
			seedPostCategoryTree(db, typeCode, child, &categoryID)
		}
		return &categoryID
	}

	if result.Error != gorm.ErrRecordNotFound {
		log.Printf("Failed to check post category %s: %v", seed.Slug, result.Error)
		return nil
	}

	if err := db.Create(&category).Error; err != nil {
		log.Printf("Failed to seed post category %s: %v", seed.Slug, err)
		return nil
	}

	log.Printf("Seeded post category: %s", seed.Slug)
	categoryID := category.ID
	for _, child := range seed.Children {
		seedPostCategoryTree(db, typeCode, child, &categoryID)
	}

	return &categoryID
}
