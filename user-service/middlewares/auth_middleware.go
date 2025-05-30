package middlewares

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/vadgun/gotrelloclone/user-service/infra/config"
	"github.com/vadgun/gotrelloclone/user-service/infra/logger"
	"go.uber.org/zap"
)

// AuthMiddleware protege rutas verificando el token JWT.
func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		logger.Log.Info("Verificando Token", zap.String("endpoint", ctx.Request.URL.Path), zap.String("ip", ctx.ClientIP()))

		// Obtener el header Authorization
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Token requerido"})
			ctx.Abort()
			return
		}

		// Extraer token del header (Formato: "Bearer <token>")
		tokenString := strings.Split(authHeader, " ")
		if len(tokenString) != 2 || tokenString[0] != "Bearer" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Formato de token inválido"})
			ctx.Abort()
			return
		}

		// Validar el token
		claims := &config.JWTClaims{}
		token, err := jwt.ParseWithClaims(tokenString[1], claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.JWTSecret), nil
		})

		if err != nil || !token.Valid {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Token inválido"})
			ctx.Abort()
			return
		}

		// Guardar los claims en el contexto
		ctx.Set("userID", claims.UserID)
		ctx.Next()
	}
}
