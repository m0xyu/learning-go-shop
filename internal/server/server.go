package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/m0xyu/learning-go-shop/docs"
	"github.com/m0xyu/learning-go-shop/internal/config"
	"github.com/m0xyu/learning-go-shop/internal/services"
	"github.com/rs/zerolog"
	"gorm.io/gorm"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Server struct {
	config         *config.Config
	db             *gorm.DB
	logger         *zerolog.Logger
	authService    *services.AuthService
	productService *services.ProductService
	uploadService  *services.UploadService
	userService    *services.UserService
	cartService    *services.CartService
	orderService   *services.OrderService
}

func New(
	ctg *config.Config,
	db *gorm.DB,
	logger *zerolog.Logger,
	authService *services.AuthService,
	productService *services.ProductService,
	userService *services.UserService,
	uploadService *services.UploadService,
	cartService *services.CartService,
	orderService *services.OrderService,
) *Server {
	return &Server{
		config:         ctg,
		db:             db,
		logger:         logger,
		authService:    authService,
		productService: productService,
		userService:    userService,
		uploadService:  uploadService,
		cartService:    cartService,
		orderService:   orderService,
	}
}

func (s *Server) SetupRoutes() *gin.Engine {
	router := gin.Default()

	// Add middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(s.crosMiddleware())

	router.GET("/health", s.healthCheck)

	// Add documentation routes
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.StaticFile("/api-docs", "./docs/rapidoc.html")

	router.Static("/uploads", s.config.Upload.Path)

	router.GET("/playground", s.playgroundHandler())
	router.GET("/playground/public", s.playgroundPublicHandler())
	router.GET("/playground/protected", s.playgroundProtectedHandler())

	graphqlPublic := router.Group("/graphql/public")
	graphqlPublic.Use(s.graphqlMiddleware())
	graphqlPublic.POST("/", s.graphQLHandler())

	graphqlProtected := router.Group("/graphql")
	graphqlProtected.Use(s.authMiddleware())
	graphqlProtected.Use(s.graphqlMiddleware())
	graphqlProtected.POST("/", s.graphQLHandler())

	api := router.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{ //nolint:gocritic // Ignore "too many methods" warning for auth handlers
			auth.POST("/register", s.register)
			auth.POST("/login", s.login)
			auth.POST("/refresh", s.refreshToken)
			auth.POST("/logout", s.logout)
		}

		protected := api.Group("/")
		protected.Use(s.authMiddleware())
		{
			// User Routes
			users := protected.Group("/users")
			{
				userRoutes := users
				userRoutes.GET("/profile", s.getProfile)
				userRoutes.PUT("/profile", s.updateProfile)
			}

			// Category Routes
			categories := protected.Group("/categories")
			{
				categoryRoute := categories
				categoryRoute.POST("/", s.adminMiddleware(), s.createCategory)
				categoryRoute.PUT("/:id", s.adminMiddleware(), s.updateCategory)
				categoryRoute.DELETE("/:id", s.adminMiddleware(), s.deleteCategory)
			}

			// Product Routes
			products := protected.Group("/products")
			{
				productRoute := products
				productRoute.POST("/", s.adminMiddleware(), s.createProduct)
				productRoute.PUT("/:id", s.adminMiddleware(), s.updateProduct)
				productRoute.DELETE("/:id", s.adminMiddleware(), s.deleteProduct)
				productRoute.POST("/:id/image", s.adminMiddleware(), s.uploadProductImage)
			}

			// Cart Routes
			cart := protected.Group("/cart")
			{
				cartRoute := cart
				cartRoute.GET("/", s.getCart)
				cartRoute.POST("/items", s.addToCart)
				cartRoute.PUT("/items/:id", s.updateCartItem)
				cartRoute.DELETE("/items/:id", s.removeFromCart)
			}

			// Order Routes
			order := protected.Group("/orders")
			{
				orderRoute := order
				orderRoute.POST("/", s.createOrder)
				orderRoute.GET("/", s.getOrders)
				orderRoute.GET("/:id", s.getOrder)
			}
		}

		// Public Routes
		api.GET("/categories", s.getCategories)
		api.GET("/products", s.getProducts)
		api.GET("/products/:id", s.getProduct)
	}

	return router
}

func (s *Server) healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (s *Server) crosMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Method", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
