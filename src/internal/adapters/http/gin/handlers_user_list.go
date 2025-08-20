package ginadp

import (
	"math"
	"net/http"
	"time"

	ginadp "github.com/PixyBoy/jwt-auth-go/internal/adapters/http/gin/dto"
	"github.com/PixyBoy/jwt-auth-go/internal/core/ports"
	"github.com/gin-gonic/gin"
)

// @Summary      List users
// @Description  List users with search & pagination
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        search   query string false "Search by phone"
// @Param        page     query int    false "Page number"
// @Param        per_page query int    false "Items per page (<=200)"
// @Security     BearerAuth
// @Success      200 {object} UsersListResponse
// @Failure      400 {object} ErrorResponse
// @Failure      401 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Router       /v1/users [get]
func ListUsersHandler(userRepo ports.UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		var q ginadp.UsersQuery
		if err := c.ShouldBindQuery(&q); err != nil {
			c.JSON(http.StatusBadRequest, NewError("INVALID_INPUT", "invalid query params"))
			return
		}
		users, total, err := userRepo.List(q.Search, q.Page, q.PerPage)
		if err != nil {
			c.JSON(http.StatusInternalServerError, NewError("INTERNAL", "db error"))
			return
		}
		items := make([]ginadp.UsersListItem, 0, len(users))
		for _, u := range users {
			items = append(items, ginadp.UsersListItem{
				ID:        u.ID,
				Phone:     u.Phone,
				CreatedAt: u.CreatedAt.Format(time.RFC3339),
			})
		}
		totalPages := int64(math.Ceil(float64(total) / float64(q.PerPage)))
		if q.PerPage <= 0 {
			totalPages = 1
		}
		c.JSON(http.StatusOK, ginadp.UsersListResponse{
			Data: items,
			Meta: ginadp.PaginationMeta{
				Page:       q.Page,
				PerPage:    q.PerPage,
				Total:      total,
				TotalPages: totalPages,
			},
		})
	}
}
