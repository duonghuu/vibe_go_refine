package main

import (
	"log"

	"go_refine_dashboard_be/internal/di"
	"go_refine_dashboard_be/internal/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	// 1. Connect DB (thay đổi dsn theo môi trường)
	dsn := "root:root@tcp(127.0.0.1:3306)/vibe_go_refine?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect database: %v", err)
	}

	// 2. Initialize Dependency Injection
	productController := di.InitializeProductController(db)
	mediaController := di.InitializeMediaController(db)
	categoryController := di.InitializeCategoryController(db)

	// 3. Setup Router
	r := gin.Default()
	r.Static("/uploads", "./public/uploads")

	// 4. API Routes
	api := r.Group("/api/v1")
	// Media Routes
	mediaGroup := api.Group("/media")
	{
		mediaGroup.POST("/upload", mediaController.UploadFile)
		mediaGroup.DELETE("/:id", mediaController.DeleteFile)
	}

	admin := api.Group("/admin")
	admin.Use(middleware.AuthMiddleware())
	admin.Use(middleware.RoleMiddleware("ADMIN", "STAFF"))

	products := admin.Group("/products")
	{
		products.GET("", productController.GetProducts)
		products.POST("", productController.CreateProduct)
		products.PUT("/:id", productController.UpdateProduct)
		products.PATCH("/:id/status", productController.UpdateProductStatus)
		products.DELETE("/:id", productController.DeleteProduct)
		products.POST("/bulk-delete", productController.BulkDelete)
		products.PATCH("/bulk-status", productController.BulkStatus)
	}

	categories := admin.Group("/categories")
	{
		categories.GET("", categoryController.GetCategories)
		categories.GET("/:id", categoryController.GetCategoryByID)
		categories.POST("", categoryController.CreateCategory)
		categories.PUT("/:id", categoryController.UpdateCategory)
		categories.DELETE("/:id", categoryController.DeleteCategory)
		categories.PATCH("/reorder", categoryController.ReorderCategories)
	}

	// 5. Start Server
	log.Println("Server is running on port 8080...")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
