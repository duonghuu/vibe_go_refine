package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"go_refine_dashboard_be/internal/di"
	"go_refine_dashboard_be/internal/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
)

func main() {
	// 1. Connect DB (thay đổi dsn theo môi trường)
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "127.0.0.1"
	}
	dbPort := os.Getenv("DB_PORT")
	if dbPort == "" {
		dbPort = "3307"
	}
	dbUser := os.Getenv("DB_USER")
	if dbUser == "" {
		dbUser = "root"
	}
	dbPassword := os.Getenv("DB_PASSWORD")
	if dbPassword == "" {
		dbPassword = "root"
	}
	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "vibe_db"
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbUser, dbPassword, dbHost, dbPort, dbName)
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		log.Fatalf("Failed to connect database: %v", err)
	}

	// 1.1 Connect Redis
	redisHost := os.Getenv("REDIS_HOST")
	if redisHost == "" {
		redisHost = "127.0.0.1"
	}
	redisPort := os.Getenv("REDIS_PORT")
	if redisPort == "" {
		redisPort = "6380"
	}
	redisPassword := os.Getenv("REDIS_PASSWORD")
	if redisPassword == "" {
		redisPassword = "vibe_redis_secret"
	}
	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", redisHost, redisPort),
		Password: redisPassword,
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
	postTypeController := di.InitializePostTypeController(db, redisClient)
	postMediaController := di.InitializePostMediaController(db, redisClient)
	postMetaController := di.InitializePostMetaController(db, redisClient)
	postController := di.InitializePostController(db, redisClient)
	postCategoryController := di.InitializePostCategoryController(db, redisClient)

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
	mediaGroup.Use(middleware.AuthMiddleware(redisClient))
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
		products.GET("/:id", productController.GetProductByID)
		products.POST("", productController.CreateProduct)
		products.PUT("/:id", productController.UpdateProduct)
		products.PUT("/:id/status", productController.UpdateProductStatus)
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

	postTypes := admin.Group("/post-types")
	{
		postTypes.GET("", postTypeController.GetPostTypes)
		postTypes.GET("/:id", postTypeController.GetPostTypeByID)
		postTypes.POST("", postTypeController.CreatePostType)
		postTypes.PUT("/:id", postTypeController.UpdatePostType)
		postTypes.DELETE("/:id", postTypeController.DeletePostType)
	}

	postCategories := admin.Group("/post-categories")
	{
		postCategories.GET("", postCategoryController.GetPostCategories)
		postCategories.GET("/tree", postCategoryController.GetPostCategoryTree)
		postCategories.GET("/:id", postCategoryController.GetPostCategory)
		postCategories.POST("", postCategoryController.CreatePostCategory)
		postCategories.PUT("/:id", postCategoryController.UpdatePostCategory)
		postCategories.DELETE("/:id", postCategoryController.DeletePostCategory)
	}

	posts := admin.Group("/posts")
	{
		posts.GET("", postController.GetPosts)
		posts.POST("", postController.CreatePost)
		posts.DELETE("/:post_id", postController.DeletePost)
	}

	postMedia := admin.Group("/posts/:post_id/media")
	{
		postMedia.GET("", postMediaController.GetPostMedia)
		postMedia.POST("", postMediaController.CreatePostMedia)
		postMedia.PUT("/:collection", postMediaController.SyncPostMedia)
		postMedia.DELETE("/:media_id", postMediaController.DeletePostMedia)
	}

	postMeta := admin.Group("/posts/:post_id/meta")
	{
		postMeta.GET("", postMetaController.GetPostMeta)
		postMeta.PUT("", postMetaController.SyncPostMeta)
		postMeta.DELETE("/:key", postMetaController.DeletePostMeta)
	}

	// 5. Start Server
	serverPort := os.Getenv("PORT")
	if serverPort == "" {
		serverPort = "8080"
	}
	log.Println("Server is running on port " + serverPort + "...")
	if err := r.Run(":" + serverPort); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
