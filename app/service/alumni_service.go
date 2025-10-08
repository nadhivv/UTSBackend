package service

import (
	"TM4/app/model"
	"TM4/app/repository"
	"TM4/helper"
	"TM4/utils"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type AlumniService struct {
	repo *repository.AlumniRepository
}

func NewAlumniService(repo *repository.AlumniRepository) *AlumniService {
	return &AlumniService{repo: repo}
}


func (s *AlumniService) GetByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return helper.ResponseJSON(c, 400, "Invalid ID", false, nil)
	}
	alumni, err := s.repo.GetByID(id)
	if err != nil {
		return helper.ResponseJSON(c, 404, "Alumni not found", false, nil)
	}
	return helper.ResponseJSON(c, 200, "Success", true, alumni)
}

func (s *AlumniService) Create(c *fiber.Ctx) error {
	var input model.CreateAlumni
	if err := c.BodyParser(&input); err != nil {
		return helper.ResponseJSON(c, 400, "Invalid body", false, nil)
	}

	hashed, err := utils.HashPassword(input.Password)
	if err != nil {
		return helper.ResponseJSON(c, 500, "Failed to hash password", false, nil)
	}
	input.Password = hashed

	id, err := s.repo.Create(input)
	if err != nil {
		return helper.ResponseJSON(c, 500, err.Error(), false, nil)
	}
	return helper.ResponseJSON(c, 201, "Alumni created", true, fiber.Map{"id": id})
}


func (s *AlumniService) Update(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return helper.ResponseJSON(c, 400, "Invalid ID", false, nil)
	}
	var input model.UpdateAlumni
	if err := c.BodyParser(&input); err != nil {
		return helper.ResponseJSON(c, 400, "Invalid body", false, nil)
	}
	if err := s.repo.Update(id, input); err != nil {
		return helper.ResponseJSON(c, 500, err.Error(), false, nil)
	}
	return helper.ResponseJSON(c, 200, "Alumni updated successfully", true, nil)
}

func (s *AlumniService) Delete(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return helper.ResponseJSON(c, 400, "Invalid ID", false, nil)
	}
	if err := s.repo.Delete(id); err != nil {
		return helper.ResponseJSON(c, 500, err.Error(), false, nil)
	}
	return helper.ResponseJSON(c, 200, "Alumni deleted successfully", true, nil)
}

func (s *AlumniService) GetAlumni(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	sortBy := c.Query("sortBy", "id")
	order := c.Query("order", "asc")
	search := c.Query("search", "")

	offset := (page - 1) * limit

	// whitelist kolom sort
	sortWhitelist := map[string]bool{
		"id": true, "nim": true, "nama": true, "jurusan": true,
		"angkatan": true, "tahun_lulus": true, "email": true,
		"no_telepon": true, "alamat": true, "role": true, "created_at": true,
	}
	if !sortWhitelist[sortBy] {
		sortBy = "id"
	}
	if strings.ToLower(order) != "desc" {
		order = "asc"
	}

	alumni, err := s.repo.GetAlumni(search, sortBy, order, limit, offset)
	if err != nil {
		return helper.ResponseJSON(c, 500, "Failed to fetch alumni", false, nil)
	}

	total, err := s.repo.Count(search)
	if err != nil {
		return helper.ResponseJSON(c, 500, "Failed to count alumni", false, nil)
	}

	response := model.AlumniResponse{
		Data: alumni,
		Meta: model.MetaInfo{
			Page:   page,
			Limit:  limit,
			Total:  total,
			Pages:  (total + limit - 1) / limit,
			SortBy: sortBy,
			Order:  order,
			Search: search,
		},
	}
	return c.JSON(response)
}
