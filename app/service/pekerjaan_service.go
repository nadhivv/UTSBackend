package service

import (
	"TM4/app/model"
	"TM4/app/repository"
	"TM4/helper"
	"strconv"
	"strings"
	// "fmt"

	"github.com/gofiber/fiber/v2"
)

type PekerjaanService struct {
	repo *repository.PekerjaanRepository
}

func NewPekerjaanService(repo *repository.PekerjaanRepository) *PekerjaanService {
	return &PekerjaanService{repo: repo}
}

// === GET Pekerjaan by ID ===
func (s *PekerjaanService) GetByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return helper.ResponseJSON(c, 400, "Invalid ID", false, nil)
	}
	pekerjaan, err := s.repo.GetByID(id)
	if err != nil {
		return helper.ResponseJSON(c, 404, "Pekerjaan not found", false, nil)
	}
	return helper.ResponseJSON(c, 200, "Success", true, pekerjaan)
}

// === GET Pekerjaan by Alumni ID ===
func (s *PekerjaanService) GetByAlumniID(c *fiber.Ctx) error {
	alumniID, err := strconv.Atoi(c.Params("alumni_id"))
	if err != nil {
		return helper.ResponseJSON(c, 400, "Invalid Alumni ID", false, nil)
	}
	pekerjaan, err := s.repo.GetByAlumniID(alumniID)
	if err != nil {
		return helper.ResponseJSON(c, 500, err.Error(), false, nil)
	}
	return helper.ResponseJSON(c, 200, "Success", true, pekerjaan)
}

// === CREATE Pekerjaan ===
func (s *PekerjaanService) Create(c *fiber.Ctx) error {
	var input model.CreatePekerjaan
	if err := c.BodyParser(&input); err != nil {
		return helper.ResponseJSON(c, 400, "Invalid body", false, nil)
	}
	id, err := s.repo.Create(input)
	if err != nil {
		return helper.ResponseJSON(c, 500, err.Error(), false, nil)
	}
	return helper.ResponseJSON(c, 201, "Pekerjaan created", true, fiber.Map{"id": id})
}

// === UPDATE Pekerjaan ===
func (s *PekerjaanService) Update(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return helper.ResponseJSON(c, 400, "Invalid ID", false, nil)
	}
	var input model.UpdatePekerjaan
	if err := c.BodyParser(&input); err != nil {
		return helper.ResponseJSON(c, 400, "Invalid body", false, nil)
	}
	if err := s.repo.Update(id, input); err != nil {
		return helper.ResponseJSON(c, 500, err.Error(), false, nil)
	}
	return helper.ResponseJSON(c, 200, "Pekerjaan updated successfully", true, nil)
}

// === DELETE Pekerjaan ===
func (s *PekerjaanService) Delete(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return helper.ResponseJSON(c, 400, "Invalid ID", false, nil)
	}
	if err := s.repo.Delete(id); err != nil {
		return helper.ResponseJSON(c, 500, err.Error(), false, nil)
	}
	return helper.ResponseJSON(c, 200, "Pekerjaan deleted successfully", true, nil)
}

func (s *PekerjaanService) GetAll(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	sortBy := c.Query("sortBy", "id")
	order := c.Query("order", "asc")
	search := c.Query("search", "")

	offset := (page - 1) * limit

	// whitelist kolom sort
	sortWhitelist := map[string]bool{
		"id": true, "alumni_id": true, "nama_perusahaan": true, "posisi_jabatan": true,
		"bidang_industri": true, "lokasi_kerja": true, "gaji_range": true,
		"tanggal_mulai_kerja": true, "tanggal_selesai_kerja": true, "status_pekerjaan": true, "deskripsi_pekerjaan": true, "created_at": true,
	}
	if !sortWhitelist[sortBy] {
		sortBy = "id"
	}
	if strings.ToLower(order) != "desc" {
		order = "asc"
	}

	pekerjaan, err := s.repo.GetAll(search, sortBy, order, limit, offset)
	if err != nil {
		return helper.ResponseJSON(c, 500, "Failed to fetch pekerjaan alumni", false, nil)
	}

	total, err := s.repo.Count(search)
	if err != nil {
		return helper.ResponseJSON(c, 500, "Failed to count pekerjaan alumni", false, nil)
	}

	response := model.PekerjaanResponse{
		Data: pekerjaan,
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

func (s *PekerjaanService) SoftDelete(c *fiber.Ctx) error {
	role := c.Locals("role").(string)
	userID := c.Locals("user_id")
	idStr := c.Query("id")

	if idStr != "" {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{
				"success": false,
				"message": "Invalid ID",
				"data":    nil,
			})
		}

		existingData, err := s.repo.GetByID(id)
		if err != nil {
			return c.Status(404).JSON(fiber.Map{
				"success": false,
				"message": "pekerjaan not found",
				"data":    nil,
			})
		}

		if role != "admin" && existingData.AlumniID != userID {
			return c.Status(403).JSON(fiber.Map{
				"success": false,
				"message": "bukan pekerjaanmu",
				"data":    nil,
			})
		}

		var updateReq model.UpdatePekerjaan
		if err := s.repo.SoftDelete(id, updateReq); err != nil {
			return c.Status(500).JSON(fiber.Map{
				"success": false,
				"message": err.Error(),
				"data":    nil,
			})
		}

		return c.JSON(fiber.Map{
			"success": true,
			"message": "pekerjaan soft deleted",
			"data":    nil,
		})
	}

	if role != "admin" {
		return c.Status(403).JSON(fiber.Map{
			"success": false,
			"message": "unauthorized: admin access required",
			"data":    nil,
		})
	}

	if err := s.repo.SoftDeleteBulk(); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
			"data":    nil,
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "all pekerjaan soft deleted",
		"data":    nil,
	})
}

func (s *PekerjaanService) Trash(c *fiber.Ctx) error {
	pekerjaan, err := s.repo.Trash()
	if err != nil {
		return helper.ResponseJSON(c, 500, err.Error(), false, nil)
	}
	return helper.ResponseJSON(c, 200, "Success", true, pekerjaan)
}

func (s *PekerjaanService) Restore(c *fiber.Ctx) error {
	err := s.repo.Restore()
	if err != nil {
		return helper.ResponseJSON(c, 500, err.Error(), false, nil)
	}
	return helper.ResponseJSON(c, 200, "Semua data berhasil direstore", true, nil)
}

func (s *PekerjaanService) HardDelete(c *fiber.Ctx) error {
	err := s.repo.HardDelete()
	if err != nil {
		return helper.ResponseJSON(c, 500, err.Error(), false, nil)
	}
	return helper.ResponseJSON(c, 200, "Semua data dihapus secara permanen", true, nil)
}