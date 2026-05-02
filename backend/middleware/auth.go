package middleware

import (
	"jukebox-lite/model"
	"jukebox-lite/service"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func Auth(authService *service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(http.StatusUnauthorized, model.APIResponse{Code: 401, Message: "未登录"})
			c.Abort()
			return
		}

		token = strings.TrimPrefix(token, "Bearer ")

		openID, err := authService.ValidateToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, model.APIResponse{Code: 401, Message: "token无效"})
			c.Abort()
			return
		}

		user, err := authService.GetUserByOpenID(openID)
		if err != nil {
			c.JSON(http.StatusUnauthorized, model.APIResponse{Code: 401, Message: "用户不存在"})
			c.Abort()
			return
		}

		c.Set("open_id", openID)
		c.Set("token", token)
		c.Set("user", user)
		c.Set("user_nick", user.NickName)
		c.Set("user_role", user.Role)
		c.Next()
	}
}
