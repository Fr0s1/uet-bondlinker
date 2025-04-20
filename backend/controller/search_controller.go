package controller

import (
	"net/http"

	"socialnet/config"
	"socialnet/middleware"
	"socialnet/model"
	"socialnet/repository"
	"socialnet/util"

	"github.com/gin-gonic/gin"
)

// SearchController handles search-related requests
type SearchController struct {
	repo *repository.Repository
	cfg  *config.Config
}

// NewSearchController creates a new SearchController
func NewSearchController(repo *repository.Repository, cfg *config.Config) *SearchController {
	return &SearchController{
		repo: repo,
		cfg:  cfg,
	}
}

// SearchUsers searches for users by name or username
func (sc *SearchController) SearchUsers(c *gin.Context) {
	var filter model.SearchFilter
	if !middleware.BindQuery(c, &filter) {
		return
	}

	currentUserID := middleware.GetOptionalUserID(c)

	// Search users in database
	users, err := sc.repo.User.SearchUsers(filter.Query, filter.Pagination)
	if err != nil {
		util.RespondWithError(c, http.StatusInternalServerError, "Failed to search users")
		return
	}

	if users, err = sc.repo.User.FillFollowingInfo(currentUserID, users); err != nil {
		util.RespondWithError(c, http.StatusInternalServerError, "Failed to fetch following info")
		return
	}

	util.RespondWithSuccess(c, http.StatusOK, "Users found", users)
}

// SearchPosts searches for posts by content
func (sc *SearchController) SearchPosts(c *gin.Context) {
	var filter model.SearchFilter
	if !middleware.BindQuery(c, &filter) {
		return
	}

	currentUserID := middleware.GetOptionalUserID(c)

	// Search posts in database
	posts, err := sc.repo.Post.SearchPosts(filter.Query, filter.Pagination)
	if err != nil {
		util.RespondWithError(c, http.StatusInternalServerError, "Failed to search posts")
		return
	}

	// If user is authenticated, check if posts are liked
	if posts, err = sc.repo.Post.FillLikeInfo(currentUserID, posts); err != nil {
		util.RespondWithError(c, http.StatusInternalServerError, "Failed to fetch like info")
		return
	}

	util.RespondWithSuccess(c, http.StatusOK, "Posts found", posts)
}

// Search performs a combined search across users and posts
func (sc *SearchController) Search(c *gin.Context) {
	var filter model.SearchFilter
	if !middleware.BindQuery(c, &filter) {
		return
	}

	currentUserID := middleware.GetOptionalUserID(c)

	// Search users
	users, err := sc.repo.User.SearchUsers(filter.Query, filter.Pagination)
	if err != nil {
		util.RespondWithError(c, http.StatusInternalServerError, "Failed to search users")
		return
	}

	if users, err = sc.repo.User.FillFollowingInfo(currentUserID, users); err != nil {
		util.RespondWithError(c, http.StatusInternalServerError, "Failed to fetch following info")
		return
	}

	// Search posts
	posts, err := sc.repo.Post.SearchPosts(filter.Query, filter.Pagination)
	if err != nil {
		util.RespondWithError(c, http.StatusInternalServerError, "Failed to search posts")
		return
	}

	if posts, err = sc.repo.Post.FillLikeInfo(currentUserID, posts); err != nil {
		util.RespondWithError(c, http.StatusInternalServerError, "Failed to fetch like info")
		return
	}

	// Return combined results
	util.RespondWithSuccess(c, http.StatusOK, "Search results", gin.H{
		"users": users,
		"posts": posts,
	})
}
