
package router

import (
	"net/http"

	"socialnet/config"
	"socialnet/controller"
	"socialnet/middleware"
	"socialnet/repository"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SetupRouter configures the Gin router
func SetupRouter(db *gorm.DB, cfg *config.Config) *gin.Engine {
	// Set Gin mode
	if cfg.Server.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	// Configure CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:5173", "http://localhost:8080"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Content-Length", "Accept-Encoding", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           86400,
	}))

	// Initialize repositories
	repo := repository.NewRepository(db)

	// Initialize controllers
	userController := controller.NewUserController(repo, cfg)
	authController := controller.NewAuthController(repo, cfg)

	// Initialize post controllers
	postController := controller.NewPostController(repo, cfg)
	postInteractionController := controller.NewPostInteractionController(repo, cfg)
	commentController := controller.NewCommentController(repo, cfg)
	
	// Initialize search controller
	searchController := controller.NewSearchController(repo, cfg)

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
		}

		// User routes
		users := v1.Group("/users")
		{
			users.GET("", userController.GetUsers)
			users.GET("/:id", userController.GetUser)
			users.GET("/username/:username", userController.GetUserByUsername)

			// Protected routes
			users.Use(middleware.AuthMiddleware(cfg))
			users.GET("/suggested", userController.GetSuggestedUsers)  // New endpoint for suggested users
			users.PUT("/:id", userController.UpdateUser)
			users.GET("/me", userController.GetCurrentUser)
			users.POST("/follow/:id", userController.FollowUser)
			users.DELETE("/follow/:id", userController.UnfollowUser)
			users.GET("/followers", userController.GetFollowers)
			users.GET("/following", userController.GetFollowing)
		}

		// Post routes
		posts := v1.Group("/posts")
		{
			posts.GET("", postController.GetPosts)
			posts.GET("/:id", postController.GetPost)
			posts.GET("/trending", postController.GetTrending)

			// Protected routes
			posts.Use(middleware.AuthMiddleware(cfg))
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
		search := v1.Group("/search")
		{
			search.GET("", searchController.Search)
			search.GET("/users", searchController.SearchUsers)
			search.GET("/posts", searchController.SearchPosts)
		}
	}

	return r
}
