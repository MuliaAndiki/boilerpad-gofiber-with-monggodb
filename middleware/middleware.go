package middleware

import (
	"context"
	"os"
	"strings"
	"time"

	"boilerpad/config"
	"boilerpad/models"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

// VerifyToken mirip dengan verifyToken di Express
func VerifyToken(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Access denied. No token provided.",
		})
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	// cek blacklist
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var blacklisted bson.M
	err := config.DB.Collection("blacklist_tokens").FindOne(ctx, bson.M{"token": tokenString}).Decode(&blacklisted)
	if err == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Token sudah tidak valid (logout).",
		})
	}

	// parse token
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil || !token.Valid {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "Invalid or expired token.",
		})
	}

	// simpan user ke context
	c.Locals("user", claims)

	return c.Next()
}
// RequireRole mirip dengan requireRole di Express
func RequireRole(roles []string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user, ok := c.Locals("user").(models.User)
		if !ok {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"message": "Akses ditolak. User tidak ditemukan.",
			})
		}

		for _, role := range roles {
			if user.Role == role {
				return c.Next()
			}
		}

		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "Akses ditolak. Role tidak sesuai.",
		})
	}
}
