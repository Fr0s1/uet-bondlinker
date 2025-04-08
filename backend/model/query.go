
package model

// Pagination represents common pagination parameters
type Pagination struct {
	Limit  int `form:"limit,default=10"`
	Offset int `form:"offset,default=0"`
}

// UserFilter represents user filtering parameters
type UserFilter struct {
	Query     string `form:"q"`
	Pagination
}

// PostFilter represents post filtering parameters
type PostFilter struct {
	UserID string `form:"user_id"`
	Pagination
}

// CommentFilter represents comment filtering parameters
type CommentFilter struct {
	PostID string `form:"post_id" binding:"required,uuid"`
	Pagination
}

// FollowFilter represents follow relationship filtering parameters
type FollowFilter struct {
	Pagination
}
