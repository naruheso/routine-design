package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"github.com/amical/routine-design/backend/internal/model"
)

// JWTClaims はJWTのペイロードを表す。
type JWTClaims struct {
	UserID string   `json:"userId"`
	Email  string   `json:"email"`
	Name   string   `json:"name"`
	Roles  []string `json:"roles"`
	jwt.RegisteredClaims
}

// AuthMiddleware はJWT認証ミドルウェアを返す。
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, model.ErrorResponse("認証が必要です"))
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.AbortWithStatusJSON(http.StatusUnauthorized, model.ErrorResponse("無効な認証形式です"))
			return
		}

		claims := &JWTClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(getJWTSecret()), nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, model.ErrorResponse("無効なトークンです"))
			return
		}

		// コンテキストにユーザー情報をセット
		c.Set("userID", claims.UserID)
		c.Set("email", claims.Email)
		c.Set("name", claims.Name)
		c.Set("roles", claims.Roles)
		c.Next()
	}
}

// GetUserID はコンテキストからユーザーIDを取得する。
func GetUserID(c *gin.Context) string {
	v, _ := c.Get("userID")
	if id, ok := v.(string); ok {
		return id
	}
	return ""
}

// GetUserRoles はコンテキストからユーザーロール一覧を取得する。
func GetUserRoles(c *gin.Context) []string {
	v, _ := c.Get("roles")
	if roles, ok := v.([]string); ok {
		return roles
	}
	return nil
}

// HasRole は指定ロールを保持しているか判定する。
func HasRole(c *gin.Context, role string) bool {
	roles := GetUserRoles(c)
	for _, r := range roles {
		if r == role {
			return true
		}
	}
	return false
}

func getJWTSecret() string {
	if s := os.Getenv("JWT_SECRET"); s != "" {
		return s
	}
	return "dev-secret-key"
}
