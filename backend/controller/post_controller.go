
package controller

import (
	"github.com/gin-gonic/gin"
)

// PostControllerInterface encapsulates all post controllers
type PostControllerInterface interface {
	// Post core operations
	GetPosts(c *gin.Context)
	GetPost(c *gin.Context)
	CreatePost(c *gin.Context)
	UpdatePost(c *gin.Context)
	DeletePost(c *gin.Context)
	GetFeed(c *gin.Context)
	GetTrending(c *gin.Context)
	GetSuggestedPosts(c *gin.Context)

	// Post interaction operations
	LikePost(c *gin.Context)
	UnlikePost(c *gin.Context)
	SharePost(c *gin.Context)

	// Comment operations
	GetComments(c *gin.Context)
	CreateComment(c *gin.Context)
	UpdateComment(c *gin.Context)
	DeleteComment(c *gin.Context)
}

// SearchControllerInterface encapsulates search operations
type SearchControllerInterface interface {
	SearchUsers(c *gin.Context)
	SearchPosts(c *gin.Context)
	Search(c *gin.Context) // Combined search across different types
}
