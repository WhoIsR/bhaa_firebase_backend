package routes

import (
	"github.com/WhoIsR/bhaa_firebase_backend/handlers"
	"github.com/WhoIsR/bhaa_firebase_backend/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	// gin.Default() udah include otomatis fitur Logger & Recovery biar ga gampang crash
	r := gin.Default()

	// CORS Middleware (Penting BANGET CUY biar request dari Flutter/Postman tidak diblokir)
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// Panggil (handlers)
	authHandler := handlers.NewAuthHandler()
	productHandler := handlers.NewProductHandler()

	// Kita bungkus semua rute di dalam grup "/v1"
	v1 := r.Group("/v1")

	// 1. Health check (Buat ngetes server idup apa nggak, tanpa perlu login)
	v1.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "gin-firebase-backend"})
	})

	// 2. Auth routes (Jalur Login)
	auth := v1.Group("/auth")
	auth.POST("/verify-token", authHandler.VerifyToken)

	// 3. Protected routes (Jalur VIP, WAJIB bawa token JWT yang valid)
	protected := v1.Group("")
	protected.Use(middleware.AuthMiddleware()) // Taruh satpam JWT di sini

	// - Jalur Products (Bisa diakses semua user login)
	products := protected.Group("/products")
	products.GET("", productHandler.GetAll)
	products.GET("/:id", productHandler.GetByID)

	// - Jalur Admin Khusus (Buat Nambah, edit, hapus)
	adminProducts := products.Group("")
	adminProducts.Use(middleware.AdminOnly()) // Taruh satpam Admin di sini
	adminProducts.POST("", productHandler.Create)
	adminProducts.PUT("/:id", productHandler.Update)
	adminProducts.DELETE("/:id", productHandler.Delete)

	return r
}
