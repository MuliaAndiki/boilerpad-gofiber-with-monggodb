package controllers

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
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type AuthController struct{}

// REGISTER
func (AuthController) Register(c *fiber.Ctx) error {
	var body struct {
		Fullname string `json:"fullname"`
		Email    string `json:"email"`
		Password string `json:"password"`
		Role     string `json:"role"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	// cek email sudah ada
	var existing models.User
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := config.DB.Collection("users").FindOne(ctx, bson.M{"email": body.Email}).Decode(&existing)
	if err == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Email sudah terdaftar",
		})
	}

	// hash password
	hash, _ := bcrypt.GenerateFromPassword([]byte(body.Password), 12)

	user := models.User{
		ID:        primitive.NewObjectID(),
		Fullname:  body.Fullname,
		Email:     body.Email,
		Password:  string(hash),
		Role:      body.Role,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err = config.DB.Collection("users").InsertOne(ctx, user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal membuat user",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Register berhasil",
		"user": fiber.Map{
			"id":       user.ID.Hex(),
			"fullname": user.Fullname,
			"email":    user.Email,
			"role":     user.Role,
		},
	})
}

// LOGIN
func (AuthController) Login(c *fiber.Ctx) error {
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	// cari user
	var user models.User
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := config.DB.Collection("users").FindOne(ctx, bson.M{"email": body.Email}).Decode(&user)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Email tidak ditemukan",
		})
	}

	// cek password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password)); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Password salah",
		})
	}

	// generate JWT
	claims := jwt.MapClaims{
		"id":    user.ID.Hex(),
		"email": user.Email,
		"role":  user.Role,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := os.Getenv("JWT_SECRET")
	tokenString, _ := token.SignedString([]byte(secret))

	return c.JSON(fiber.Map{
		"message": "Login berhasil",
		"token":   tokenString,
		"user": fiber.Map{
			"id":       user.ID.Hex(),
			"fullname": user.Fullname,
			"email":    user.Email,
			"role":     user.Role,
		},
	})
}

func (AuthController) Logout(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Token tidak ditemukan",
		})
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Simpan token ke blacklist
	_, err := config.DB.Collection("blacklist_tokens").InsertOne(ctx, bson.M{
		"token":     tokenString,
		"createdAt": time.Now(),
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal logout",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Logout berhasil, token di-blacklist",
	})
}