package controller

import (
	"socialnet/config"
	"socialnet/repository"

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
	
	// Post interaction operations
	LikePost(c *gin.Context)
	UnlikePost(c *gin.Context)
	
	// Comment operations
	GetComments(c *gin.Context)
	CreateComment(c *gin.Context)
	UpdateComment(c *gin.Context)
	DeleteComment(c *gin.Context)
}
