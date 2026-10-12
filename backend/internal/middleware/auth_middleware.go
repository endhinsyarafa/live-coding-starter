package middleware

import (
	"errors"
	"net/http"
	"strings"

	"live-coding-starter/backend/internal/model"
	"live-coding-starter/backend/internal/repository"
	"live-coding-starter/backend/internal/service"
	"live-coding-starter/backend/pkg/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const currentUserKey = "currentUser"

func Auth(tokenService service.TokenService, userRepository repository.UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenValue := parseBearerToken(c.GetHeader("Authorization"))
		if tokenValue == "" {
			response.Error(c, http.StatusUnauthorized, "Authentication required", nil)
			c.Abort()
			return
		}

		userID, err := tokenService.Validate(tokenValue)
		if err != nil {
			response.Error(c, http.StatusUnauthorized, "Invalid or expired token", nil)
			c.Abort()
			return
		}

		user, err := userRepository.FindByID(userID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				response.Error(c, http.StatusUnauthorized, "Authentication required", nil)
			} else {
				response.Error(c, http.StatusInternalServerError, "Internal server error", err.Error())
			}
			c.Abort()
			return
		}

		if !user.IsActive {
			response.Error(c, http.StatusForbidden, "Akun Anda telah dinonaktifkan.", nil)
			c.Abort()
			return
		}

		c.Set(currentUserKey, user)
		c.Next()
	}
}

func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, ok := CurrentUser(c)
		if !ok || user.Role != model.RoleAdmin {
			response.Error(c, http.StatusForbidden, "Anda tidak memiliki akses ke resource ini.", nil)
			c.Abort()
			return
		}

		c.Next()
	}
}

func CurrentUser(c *gin.Context) (*model.User, bool) {
	value, exists := c.Get(currentUserKey)
	if !exists {
		return nil, false
	}

	user, ok := value.(*model.User)
	return user, ok
}

func parseBearerToken(header string) string {
	if header == "" {
		return ""
	}

	parts := strings.SplitN(header, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return ""
	}

	return strings.TrimSpace(parts[1])
}
