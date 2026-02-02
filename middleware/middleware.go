package middleware

import (
	"net/http"
	"strings"

	generate "ec-platform/tokens"

	"github.com/gin-gonic/gin"
)

// проверяет JWT токен в заголовке Authorization
func Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientToken := c.GetHeader("Authorization")

		if clientToken == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		// Убираем префикс "Bearer " если он есть
		clientToken = strings.TrimPrefix(clientToken, "Bearer ")

		// Валидируем токен
		claims, err := generate.ValidateToken(clientToken)

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Сохраняем данные из токена в контекст
		c.Set("email", claims.Email)
		c.Set("first_name", claims.First_Name)
		c.Set("last_name", claims.Last_Name)

		c.Next()
	}
}
