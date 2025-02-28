package auth

import (
	"github.com/RCSE2025/backend-go/internal/service"
	"github.com/RCSE2025/backend-go/pkg/api/response"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func ValidateJWT(jwtService service.JWTService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, response.Error("Bearer token missing in request header"))
			return
		}
		if !strings.Contains(authHeader, "Bearer ") {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, response.Error("Bearer token missing in request header"))
			return
		}
		authHeader = strings.Replace(authHeader, "Bearer ", "", -1)
		token, err := jwtService.ValidateToken(authHeader)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, response.Error("invalid token"))
			return
		}
		if !token.Valid {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, response.Error("unauthorized"))
			return
		}
		userId, err := jwtService.GetUserIDByToken(authHeader)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, response.Error(err.Error()))
			return
		}

		userRole := jwtService.GetUserRole(authHeader)
		ctx.Set("user_role", userRole)
		ctx.Set("user_id", userId)

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			ctx.Set("claims", claims)
		} else {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			ctx.Abort()
			return
		}
		ctx.Set("token", authHeader)

		ctx.Next()
	}
}
