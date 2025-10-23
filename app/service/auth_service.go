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

func (s *AuthService) Login(c *fiber.Ctx) error {
	var req model.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "invalid request"})
	}

	alumni, err := s.repo.GetByEmail(context.Background(), req.Email)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"message": "invalid email or password"})
	}	

	// cek password
	if !utils.CheckPasswordHash(req.Password, alumni.Password) {
		return c.Status(401).JSON(fiber.Map{"message": "invalid email or password"})
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
