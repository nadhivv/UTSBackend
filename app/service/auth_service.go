package service

import (
	"TM4/app/model"
	"TM4/app/repository"
	"TM4/utils"
	"context"

	"github.com/gofiber/fiber/v2"
)

type AuthService struct {
	repo repository.IAlumniRepository
}

func NewAuthService(repo repository.IAlumniRepository) *AuthService {
	return &AuthService{repo: repo}
}

// Login godoc
// @Summary Login Alumni
// @Description Login menggunakan email dan password untuk mendapatkan JWT token
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body model.LoginRequest true "Login Data"
// @Success 200 {object} model.LoginResponse
// @Failure 400 {object} map[string]interface{} "Invalid Request Format"
// @Failure 401 {object} map[string]interface{} "Invalid Email or Password"
// @Failure 500 {object} map[string]interface{} "Failed to Generate Token"
// @Router /TM4/login [post]
func (s *AuthService) Login(c *fiber.Ctx) error {
	var req model.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "invalid request"})
	}

	alumni, err := s.repo.GetByEmail(context.Background(), req.Email)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"message": "invalid email "})
	}	

	// cek password
	if !utils.CheckPasswordHash(req.Password, alumni.Password) {
		return c.Status(401).JSON(fiber.Map{"message": "invalid password"})
	}

	// generate JWT
	token, err := utils.GenerateToken(model.Alumni{
		ID:       alumni.ID,
		Nama:     alumni.Nama,
		Role:     alumni.Role,
		Email:    alumni.Email,
	})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "failed to generate token"})
	}

	return c.JSON(model.LoginResponse{
		Alumni: *alumni,
		Token:  token,
	})
}
