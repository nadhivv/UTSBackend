package service

import (
	"TM4/app/model"
	"TM4/app/repository"
	"TM4/helper"
	"TM4/utils"
	"strconv"
	"strings"
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
)

type AlumniService struct {
	repo repository.IAlumniRepository
}

func NewAlumniService(repo repository.IAlumniRepository) *AlumniService {
	return &AlumniService{repo: repo}
}


// GetByID godoc
// @Summary Get alumni by ID
// @Description Ambil data alumni berdasarkan ID
// @Tags Alumni
// @Accept json
// @Produce json
// @Param id path string true "Alumni ID"
// @Success 200 {object} model.Alumni
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /alumni/{id} [get]
func (s *AlumniService) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	alumni, err := s.repo.GetAlumniByID(ctx, id)
	if err != nil {
		return helper.ResponseJSON(c, 500, "Failed to get alumni", false, nil)
	}
	if alumni == nil {
		return helper.ResponseJSON(c, 404, "Alumni not found", false, nil)
	}

	return helper.ResponseJSON(c, 200, "Success", true, alumni)
}

// Create godoc
// @Summary Create new alumni
// @Description Membuat data alumni baru
// @Tags Alumni
// @Accept json
// @Produce json
// @Param alumni body model.CreateAlumni true "Alumni Data"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /alumni [post]
func (s *AlumniService) Create(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var input model.CreateAlumni
	if err := c.BodyParser(&input); err != nil {
		return helper.ResponseJSON(c, 400, "Invalid body", false, nil)
	}

	hashed, err := utils.HashPassword(input.Password)
	if err != nil {
		return helper.ResponseJSON(c, 500, "Failed to hash password", false, nil)
	}
	input.Password = hashed

	alumni := model.Alumni{
		NIM:        input.NIM,
		Nama:       input.Nama,
		Jurusan:    input.Jurusan,
		Angkatan:   input.Angkatan,
		TahunLulus: input.TahunLulus,
		Email:      input.Email,
		NoTelepon:  input.NoTelepon,
		Alamat:     input.Alamat,
		Role:       input.Role,
		Password:   input.Password,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		IsDeleted:  false,
	}

	result, err := s.repo.CreateAlumni(ctx, &alumni)
	if err != nil {
		return helper.ResponseJSON(c, 500, err.Error(), false, nil)
	}

	return helper.ResponseJSON(c, 201, "Alumni created successfully", true, fiber.Map{
		"id": result.ID.Hex(),
	})
}

// Update godoc
// @Summary Update alumni data
// @Description Mengubah data alumni berdasarkan ID
// @Tags Alumni
// @Accept json
// @Produce json
// @Param id path string true "Alumni ID"
// @Param alumni body model.UpdateAlumni true "Updated Alumni Data"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /alumni/{id} [put]
func (s *AlumniService) Update(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	id := c.Params("id")
	var input model.UpdateAlumni
	if err := c.BodyParser(&input); err != nil {
		return helper.ResponseJSON(c, 400, "Invalid body", false, nil)
	}

	if err := s.repo.UpdateAlumni(ctx, id, &input); err != nil {
		return helper.ResponseJSON(c, 500, err.Error(), false, nil)
	}

	return helper.ResponseJSON(c, 200, "Alumni updated successfully", true, nil)
}

// Delete godoc
// @Summary Delete alumni
// @Description Menghapus data alumni berdasarkan ID
// @Tags Alumni
// @Param id path string true "Alumni ID"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /alumni/{id} [delete]
func (s *AlumniService) Delete(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	id := c.Params("id")
	if err := s.repo.DeleteAlumni(ctx, id); err != nil {
		return helper.ResponseJSON(c, 500, err.Error(), false, nil)
	}
	return helper.ResponseJSON(c, 200, "Alumni deleted successfully", true, nil)
}

// GetAlumni godoc
// @Summary Get list of alumni
// @Description Mengambil daftar alumni dengan pagination, sorting, dan search
// @Tags Alumni
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Limit per page"
// @Param search query string false "Search keyword"
// @Param sortBy query string false "Sort column"
// @Param order query string false "Sort order (asc/desc)"
// @Success 200 {object} model.AlumniResponse
// @Failure 500 {object} map[string]interface{}
// @Router /alumni [get]
func (s *AlumniService) GetAlumni(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	sortBy := c.Query("sortBy", "created_at")
	order := c.Query("order", "asc")
	search := c.Query("search", "")

	offset := (page - 1) * limit

	// whitelist kolom sort (disesuaikan)
	sortWhitelist := map[string]bool{
		"nim": true, "nama": true, "jurusan": true,
		"angkatan": true, "tahun_lulus": true, "email": true,
		"role": true, "created_at": true,
	}
	if !sortWhitelist[sortBy] {
		sortBy = "created_at"
	}
	if strings.ToLower(order) != "desc" {
		order = "asc"
	}

	alumni, err := s.repo.GetAlumni(ctx, search, limit, offset, sortBy, order)
	if err != nil {
		return helper.ResponseJSON(c, 500, "Failed to fetch alumni", false, nil)
	}

	total, err := s.repo.Count(ctx, search)
	if err != nil {
		return helper.ResponseJSON(c, 500, "Failed to count alumni", false, nil)
	}

	response := model.AlumniResponse{
		Data: alumni,
		Meta: model.MetaInfo{
			Page:   page,
			Limit:  limit,
			Total:  int(total),
			Pages:  (int(total) + limit - 1) / limit,
			SortBy: sortBy,
			Order:  order,
			Search: search,
		},
	}
	return c.JSON(response)
}
