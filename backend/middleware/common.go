
package middleware

import (
	"errors"
	"net/http"
	"socialnet/util"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ParseUUIDParam extracts and validates a UUID from URL parameters
func ParseUUIDParam(c *gin.Context, paramName string) (uuid.UUID, error) {
	idStr := c.Param(paramName)
	id, err := uuid.Parse(idStr)
	if err != nil {
		return uuid.Nil, errors.New("invalid UUID format")
	}
	return id, nil
}

// RequireAuthentication ensures a user is authenticated
func RequireAuthentication(c *gin.Context) (uuid.UUID, bool) {
	userIDStr, err := GetUserID(c)
	if err != nil {
		util.RespondWithError(c, http.StatusUnauthorized, util.ErrorMessages.NotAuthenticated)
		return uuid.Nil, false
	}
	
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		util.RespondWithError(c, http.StatusBadRequest, util.ErrorMessages.InvalidUserID)
		return uuid.Nil, false
	}
	
	return userID, true
}

// BindJSON binds and validates JSON input
func BindJSON(c *gin.Context, obj interface{}) bool {
	if err := c.ShouldBindJSON(obj); err != nil {
		util.RespondWithError(c, http.StatusBadRequest, err.Error())
		return false
	}
	return true
}

// BindQuery binds and validates query parameters
func BindQuery(c *gin.Context, obj interface{}) bool {
	if err := c.ShouldBindQuery(obj); err != nil {
		util.RespondWithError(c, http.StatusBadRequest, err.Error())
		return false
	}
	return true
}

// CheckResourceOwnership verifies the authenticated user owns a resource
func CheckResourceOwnership(c *gin.Context, resourceOwnerID uuid.UUID, currentUserID uuid.UUID) bool {
	if resourceOwnerID != currentUserID {
		util.RespondWithError(c, http.StatusForbidden, util.ErrorMessages.NotAuthorized)
		return false
	}
	return true
}

// GetOptionalUserID retrieves user ID if authenticated, but doesn't require it
func GetOptionalUserID(c *gin.Context) *uuid.UUID {
	userIDStr, err := GetUserID(c)
	if err != nil {
		return nil
	}
	
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil
	}
	
	return &userID
}
