package middlewares

import (
	"fmt"
	"net/http"
	"strings"

	"gotickets/internal/utils"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v5"
)

func JWTMiddleware(jwtSecret string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return utils.SendError(c, http.StatusUnauthorized, "Missing Authorization header", nil)
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				return utils.SendError(c, http.StatusUnauthorized, "Invalid Authorization header format", nil)
			}

			tokenString := parts[1]

			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return []byte(jwtSecret), nil
			})

			if err != nil || !token.Valid {
				return utils.SendError(c, http.StatusUnauthorized, "Invalid or expired token", err)
			}

			// Inject token into context
			c.Set("user", token)
			return next(c)
		}
	}
}

func RoleMiddleware(requiredRoles ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			userToken, ok := c.Get("user").(*jwt.Token)
			if !ok {
				return utils.SendError(c, http.StatusUnauthorized, "Unauthorized", nil)
			}

			claims, ok := userToken.Claims.(jwt.MapClaims)
			if !ok {
				return utils.SendError(c, http.StatusUnauthorized, "Invalid token claims", nil)
			}

			role, ok := claims["role"].(string)
			if !ok {
				return utils.SendError(c, http.StatusForbidden, "Role not found in token", nil)
			}

			for _, requiredRole := range requiredRoles {
				if role == requiredRole {
					return next(c)
				}
			}

			return utils.SendError(c, http.StatusForbidden, "Insufficient permissions", nil)
		}
	}
}
