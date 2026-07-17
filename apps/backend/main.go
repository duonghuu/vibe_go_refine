package main

import (
	"context"
	"log"

	"go_refine_dashboard_be/internal/di"
	"go_refine_dashboard_be/internal/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
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

	// 1.1 Connect Redis
	redisClient := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
	})
	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		log.Printf("Failed to connect redis: %v. Auth might fail.", err)
	}

	// 2. Initialize Dependency Injection
	productController := di.InitializeProductController(db)
	mediaController := di.InitializeMediaController(db)
	categoryController := di.InitializeCategoryController(db)
	userController := di.InitializeUserController(db)
	authController := di.InitializeAuthController(db, redisClient)

	// 3. Setup Router
	r := gin.Default()
	
	// Add CORS middleware
	r.Use(cors.New(cors.Config{
		AllowOriginFunc: func(origin string) bool {
			return true
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS", "HEAD"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization", "Accept", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	r.Static("/uploads", "./public/uploads")

	// 4. API Routes
	api := r.Group("/api/v1")
	// Media Routes
	mediaGroup := api.Group("/media")
	{
		mediaGroup.POST("/upload", mediaController.UploadFile)
		mediaGroup.DELETE("/:id", mediaController.DeleteFile)
	}

	// Auth Routes (Public)
	authGroup := api.Group("/auth")
	{
		authGroup.POST("/login", authController.Login)
		authGroup.POST("/refresh-token", authController.RefreshToken)
	}

	// Auth Protected Routes
	admin := api.Group("/admin")
	admin.Use(middleware.AuthMiddleware(redisClient))
	admin.Use(middleware.RoleMiddleware("ADMIN", "STAFF"))

	// Protected Auth Endpoints
	authProtected := api.Group("/auth")
	authProtected.Use(middleware.AuthMiddleware(redisClient))
	{
		authProtected.GET("/me", authController.GetMe)
		authProtected.POST("/logout", authController.Logout)
	}

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

	users := admin.Group("/users")
	{
		users.GET("", userController.GetUsers)
		users.GET("/:id", userController.GetUserByID)
		users.POST("", userController.CreateUser)
		users.PUT("/:id", userController.UpdateUser)
		users.PATCH("/:id/status", userController.UpdateStatus)
		users.PATCH("/bulk-status", userController.BulkUpdateStatus)
		users.PATCH("/:id/reset-password", userController.ResetPassword)
	}

	// 5. Start Server
	log.Println("Server is running on port 8080...")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
