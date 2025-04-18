package middlewares

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/vadgun/gotrelloclone/board-service/infra/config"
	"github.com/vadgun/gotrelloclone/board-service/infra/logger"
	"go.uber.org/zap"

	"slices"

	"github.com/golang-jwt/jwt/v5"
)

func IsRoleAllowed(allowedRoles ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		// Envia un log del endpoint consultado
		logger.Log.Info("Verificando rol en board-service", zap.String("endpoint", ctx.Request.URL.Path), zap.String("ip", ctx.ClientIP()))

		// 1️⃣ Obtener el header Authorization
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Token requerido"})
			ctx.Abort()
			return
		}

		// 2️⃣ Extraer token del header (Formato: "Bearer <token>")
		tokenString := strings.Split(authHeader, " ")
		if len(tokenString) != 2 || tokenString[0] != "Bearer" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Formato de token inválido"})
			ctx.Abort()
			return
		}

		// 3️⃣ Validar el token
		claims := &config.JWTClaims{}
		token, err := jwt.ParseWithClaims(tokenString[1], claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.JWTSecret), nil
		})

		if err != nil || !token.Valid {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Token inválido"})
			ctx.Abort()
			return
		}

		// 4️⃣ Validar el rol
		claims, ok := token.Claims.(*config.JWTClaims)
		if !ok {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "Sin permisos"})
			ctx.Abort()
			return
		}

		userRole := claims.Role
		if !slices.Contains(allowedRoles, userRole) {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "No autorizado"})
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
