package ginadp

import (
	"net/http"
	"strconv"

	"github.com/PixyBoy/jwt-auth-go/internal/core/ports"
	"github.com/gin-gonic/gin"
)

// GetMeHandler godoc
// @Summary Get current user
// @Description Returns details of the logged in user
// @Tags users
// @Security BearerAuth
// @Produce json
// @Success 200 {object} domain.User
// @Failure 401 {object} map[string]string
// @Router /v1/users/me [get]
func GetMeHandler(userRepo ports.UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		uid, _, ok := GetAuthUser(c)
		if !ok {
			c.JSON(http.StatusUnauthorized, NewError("UNAUTHORIZED", "not authorized"))
			return
		}
		u, err := userRepo.FindByID(uid)
		if err != nil {
			c.JSON(http.StatusInternalServerError, NewError("INTERNAL", "db error"))
			return
		}
		if u == nil {
			c.JSON(http.StatusNotFound, NewError("NOT_FOUND", "user not found"))
			return
		}
		c.JSON(http.StatusOK, gin.H{"id": u.ID, "phone": u.Phone, "created_at": u.CreatedAt})
	}
}

// GetUserByIDHandler godoc
// @Summary Get user by ID
// @Description Returns details of a specific user by ID
// @Tags users
// @Security BearerAuth
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} domain.User
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /v1/users/{id} [get]
func GetUserByIDHandler(userRepo ports.UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, NewError("INVALID_INPUT", "invalid id"))
			return
		}
		u, err := userRepo.FindByID(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, NewError("INTERNAL", "db error"))
			return
		}
		if u == nil {
			c.JSON(http.StatusNotFound, NewError("NOT_FOUND", "user not found"))
			return
		}
		c.JSON(http.StatusOK, gin.H{"id": u.ID, "phone": u.Phone, "created_at": u.CreatedAt})
	}
}
