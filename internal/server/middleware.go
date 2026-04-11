package server

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/m0xyu/learning-go-shop/internal/models"
	"github.com/m0xyu/learning-go-shop/internal/utils"
)

func (s *Server) authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Authorizationヘッダーの検証
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.UnauthorizedResponse(c, "Authorization header required")
			c.Abort()
			return
		}

		// Bearerトークンの形式を検証
		tokenParts := strings.Split(authHeader, " ")
		// "Bearer <token>"の形式であることを確認
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			utils.UnauthorizedResponse(c, "Invalid authorization header format")
			c.Abort()
			return
		}

		// トークンの検証とクレームの抽出
		claims, err := utils.ValidateToken(tokenParts[1], s.config.JWT.Secret)
		if err != nil {
			utils.UnauthorizedResponse(c, "Invalid token")
			c.Abort()
			return
		}

		// コンテキストに保存
		c.Set("user_id", claims.UserID)
		c.Set("email", claims.Email)
		c.Set("role", claims.Role)

		c.Next()
	}
}

func (s *Server) adminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists || role != string(models.UserRoleAdmin) {
			utils.ForbiddenResponse(c, "Forbidden")
			c.Abort()
			return
		}
		c.Next()
	}
}
