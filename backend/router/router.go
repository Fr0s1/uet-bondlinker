package router

import (
	"net/http"
	"socialnet/websocket"
	"time"

	"socialnet/config"
	"socialnet/controller"
	"socialnet/middleware"
	"socialnet/repository"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SetupRouter configures the Gin router
func SetupRouter(db *gorm.DB, cfg *config.Config, hub *websocket.Hub) *gin.Engine {
	// Set Gin mode
	if cfg.Server.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	// Configure CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.Server.CorsOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Content-Length", "Accept-Encoding", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           86400 * time.Second,
	}))

	// Initialize repositories
	repo := repository.NewRepository(db)

	// Initialize controllers
	userController := controller.NewUserController(repo, cfg)
	authController := controller.NewAuthController(repo, cfg)
	fileController := controller.NewFileController(cfg)

	// Initialize post controllers
	postController := controller.NewPostController(repo, cfg)
	postInteractionController := controller.NewPostInteractionController(repo, cfg)
	commentController := controller.NewCommentController(repo, cfg)

	// Initialize search controller
	searchController := controller.NewSearchController(repo, cfg)

	// Initialize message controller
	messageController := controller.NewMessageController(repo, hub, cfg)

	// Initialize notification controller
	notificationController := controller.NewNotificationController(repo, cfg)

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	// API v1 routes
	v1 := r.Group("/api/v1")
	{
		// Auth routes
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authController.Register)
			auth.POST("/login", authController.Login)
			auth.GET("/verify-email", authController.VerifyEmail)
			auth.POST("/forgot-password", authController.ForgotPassword)
			auth.POST("/reset-password", authController.ResetPassword)

			// Protected auth routes
			auth.Use(middleware.AuthMiddleware(cfg))
			auth.POST("/logout", authController.Logout)
			auth.PUT("/change-password", authController.ChangePassword)
		}

		// User routes
		users := v1.Group("/users", middleware.AuthMiddleware(cfg))
		{
			users.GET("", userController.GetUsers)
			users.GET("/:id", userController.GetUser)
			users.GET("/username/:username", userController.GetUserByUsername)
			users.GET("/suggested", userController.GetSuggestedUsers)
			users.PUT("/:id", userController.UpdateUser)
			users.GET("/me", userController.GetCurrentUser)
			users.POST("/fcm-token", userController.SaveFCMToken)
			users.POST("/follow/:id", userController.FollowUser)
			users.DELETE("/follow/:id", userController.UnfollowUser)
			users.GET("/followers", userController.GetFollowers)
			users.GET("/following", userController.GetFollowing)
		}

		// File upload routes
		uploads := v1.Group("/uploads", middleware.AuthMiddleware(cfg))
		{
			uploads.POST("", fileController.UploadFile)
		}

		// Post routes
		posts := v1.Group("/posts", middleware.AuthMiddleware(cfg))
		{
			posts.GET("", postController.GetPosts)
			posts.GET("/:id", postController.GetPost)
			posts.GET("/trending", postController.GetTrending)
			posts.POST("", postController.CreatePost)
			posts.PUT("/:id", postController.UpdatePost)
			posts.DELETE("/:id", postController.DeletePost)
			posts.POST("/:id/like", postInteractionController.LikePost)
			posts.DELETE("/:id/like", postInteractionController.UnlikePost)
			posts.POST("/:id/share", postInteractionController.SharePost)
			posts.GET("/feed", postController.GetFeed)
			posts.GET("/suggested", postController.GetSuggestedPosts)

			// Comment routes (nested under posts)
			posts.GET("/:id/comments", commentController.GetComments)
			posts.POST("/:id/comments", commentController.CreateComment)
			posts.PUT("/comments/:commentId", commentController.UpdateComment)
			posts.DELETE("/comments/:commentId", commentController.DeleteComment)
		}

		// Search routes
		search := v1.Group("/search", middleware.AuthMiddleware(cfg))
		{
			search.GET("", searchController.Search)
			search.GET("/users", searchController.SearchUsers)
			search.GET("/posts", searchController.SearchPosts)
		}

		// Message routes
		conversations := v1.Group("/conversations", middleware.AuthMiddleware(cfg))
		{
			conversations.GET("", messageController.GetConversations)
			conversations.POST("", messageController.CreateConversation)
			conversations.GET("/:id", messageController.GetConversation)
			conversations.GET("/:id/messages", messageController.GetMessages)
			conversations.POST("/:id/messages", messageController.CreateMessage)
			conversations.POST("/:id/read", messageController.MarkConversationAsRead)
		}

		// Notification routes
		notifications := v1.Group("/notifications", middleware.AuthMiddleware(cfg))
		{
			notifications.GET("", notificationController.GetNotifications)
			notifications.GET("/unread-count", notificationController.GetUnreadCount)
			notifications.PUT("/:id/read", notificationController.MarkAsRead)
			notifications.PUT("/read-all", notificationController.MarkAllAsRead)
		}

		// WebSocket endpoint
		v1.GET("/ws", middleware.AuthMiddleware(cfg), websocket.HandleWebSocket(hub))
	}

	return r
}
