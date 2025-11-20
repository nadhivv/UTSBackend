package service

import (
	"TM4/app/model"
	"TM4/app/repository"
	"TM4/helper"
	"context"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)


type PekerjaanService struct {
	repo repository.IPekerjaanRepository
}

func NewPekerjaanService(repo repository.IPekerjaanRepository) *PekerjaanService {
	return &PekerjaanService{repo: repo}
}

// @Summary Get pekerjaan by ID
// @Tags Pekerjaan
// @Param id path string true "Pekerjaan ID"
// @Success 200 {object} model.Pekerjaan
// @Failure 404 {object} map[string]interface{}
// @Router /pekerjaan/{id} [get]	
// === GET Pekerjaan by ID ===
func (s *PekerjaanService) GetByID(c *fiber.Ctx) error {
	idParam := c.Params("id")

	pekerjaan, err := s.repo.GetByID(context.Background(), idParam)
	if err != nil {
		return helper.ResponseJSON(c, 404, "Pekerjaan not found", false, nil)
	}
	return helper.ResponseJSON(c, 200, "Success", true, pekerjaan)
}

// @Summary Get pekerjaan by Alumni ID
// @Tags Pekerjaan
// @Param alumni_id path string true "Alumni ID"
// @Success 200 {array} model.Pekerjaan
// @Router /pekerjaan/alumni/{alumni_id} [get]
// === GET Pekerjaan by Alumni ID ===
func (s *PekerjaanService) GetByAlumniID(c *fiber.Ctx) error {
	alumniIDParam := c.Params("alumni_id")

	// Convert string ke ObjectID
	alumniID, err := primitive.ObjectIDFromHex(alumniIDParam)
	if err != nil {
		return helper.ResponseJSON(c, 400, "Invalid alumni_id", false, nil)
	}

	pekerjaan, err := s.repo.GetByAlumniID(context.Background(), alumniID)
	if err != nil {
		return helper.ResponseJSON(c, 500, err.Error(), false, nil)
	}
	return helper.ResponseJSON(c, 200, "Success", true, pekerjaan)
}

// @Summary Create pekerjaan
// @Tags Pekerjaan
// @Accept json
// @Produce json
// @Param pekerjaan body model.CreatePekerjaan true "Data pekerjaan"
// @Success 201 {object} map[string]interface{}
// @Router /pekerjaan [post]
// === CREATE Pekerjaan ===
func (s *PekerjaanService) Create(c *fiber.Ctx) error {
	var input model.CreatePekerjaan
	if err := c.BodyParser(&input); err != nil {
		return helper.ResponseJSON(c, 400, "Invalid body", false, nil)
	}

	result, err := s.repo.Create(context.Background(), &input)
	if err != nil {
		return helper.ResponseJSON(c, 500, err.Error(), false, nil)
	}

	return helper.ResponseJSON(c, 201, "Pekerjaan created", true, fiber.Map{
		"id": result.InsertedID,
	})
}

// @Summary Update pekerjaan
// @Tags Pekerjaan
// @Param id path string true "Pekerjaan ID"
// @Param pekerjaan body model.UpdatePekerjaan true "Updated Data"
// @Success 200 {object} map[string]interface{}
// @Router /pekerjaan/{id} [put]
// === UPDATE Pekerjaan ===
func (s *PekerjaanService) Update(c *fiber.Ctx) error {
	idParam := c.Params("id")

	var input model.UpdatePekerjaan
	if err := c.BodyParser(&input); err != nil {
		return helper.ResponseJSON(c, 400, "Invalid body", false, nil)
	}

	if err := s.repo.Update(context.Background(), idParam, &input); err != nil {
		return helper.ResponseJSON(c, 500, err.Error(), false, nil)
	}
	return helper.ResponseJSON(c, 200, "Pekerjaan updated successfully", true, nil)
}

// @Summary Delete pekerjaan
// @Tags Pekerjaan
// @Param id path string true "Pekerjaan ID"
// @Success 200 {object} map[string]interface{}
// @Router /pekerjaan/{id} [delete]
// === DELETE Pekerjaan ===
func (s *PekerjaanService) Delete(c *fiber.Ctx) error {
	idParam := c.Params("id")

	if err := s.repo.Delete(context.Background(), idParam); err != nil {
		return helper.ResponseJSON(c, 500, err.Error(), false, nil)
	}
	return helper.ResponseJSON(c, 200, "Pekerjaan deleted successfully", true, nil)
}

// GetAll godoc
// @Summary Get all pekerjaan
// @Description Ambil semua data pekerjaan dengan pagination, sorting & search
// @Tags Pekerjaan
// @Accept json
// @Produce json
// @Param search query string false "Search keyword"
// @Param limit query int false "Limit data per halaman"
// @Param offset query int false "Offset data"
// @Param sortBy query string false "Kolom untuk sorting"
// @Param order query string false "Urutan sorting (asc/desc)"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /TM4/pekerjaan [get]
// === GET ALL Pekerjaan ===
func (s *PekerjaanService) GetAll(c *fiber.Ctx) error {
    search := c.Query("search", "")
    limit, _ := strconv.Atoi(c.Query("limit", "10"))
    offset, _ := strconv.Atoi(c.Query("offset", "0"))
    sortBy := c.Query("sortBy", "created_at")
    order := c.Query("order", "desc")

    data, err := s.repo.GetAll(c.Context(), search, sortBy, order, limit, offset)
    if err != nil {
        return helper.ResponseJSON(c, 500, err.Error(), false, nil)
    }

    total, err := s.repo.Count(c.Context(), search)
    if err != nil {
        return helper.ResponseJSON(c, 500, err.Error(), false, nil)
    }

    return helper.ResponseJSON(c, 200, "Success", true, fiber.Map{
        "data":  data,
        "total": total,
    })
}

// SoftDelete godoc
// @Summary Soft delete pekerjaan
// @Description Mengubah status pekerjaan menjadi terhapus (soft delete)
// @Tags Pekerjaan
// @Accept json
// @Produce json
// @Param id query string true "ID pekerjaan"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /TM4/pekerjaan/softdelete [put]
// === SOFT DELETE ===
func (s *PekerjaanService) SoftDelete(c *fiber.Ctx) error {
	idParam := c.Params("id")

	if err := s.repo.SoftDelete(context.Background(), idParam); err != nil {
		return helper.ResponseJSON(c, 500, err.Error(), false, nil)
	}

	return helper.ResponseJSON(c, 200, "Pekerjaan soft deleted", true, nil)
}

// Trash godoc
// @Summary Get list of soft-deleted pekerjaan
// @Description Menampilkan pekerjaan yang berada di trash (soft deleted)
// @Tags Pekerjaan
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Limit per page"
// @Param search query string false "Search keyword"
// @Param sortBy query string false "Sort column"
// @Param order query string false "Sort order (asc/desc)"
// @Success 200 {object} model.PekerjaanResponse
// @Failure 500 {object} map[string]interface{}
// @Router /TM4/pekerjaan/trash [get]
// === TRASH ===
func (s *PekerjaanService) Trash(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	sortBy := c.Query("sortBy", "updated_at")
	order := c.Query("order", "desc")
	search := c.Query("search", "")
	offset := (page - 1) * limit

	// ambil data pekerjaan yang sudah dihapus
	pekerjaan, err := s.repo.Trash(context.Background(), search, sortBy, order, limit, offset)
	if err != nil {
		return helper.ResponseJSON(c, 500, err.Error(), false, nil)
	}

	// hitung total dokumen yang dihapus
	total, err := s.repo.CountTrash(context.Background(), search)
	if err != nil {
		return helper.ResponseJSON(c, 500, err.Error(), false, nil)
	}

	response := model.PekerjaanResponse{
		Data: pekerjaan,
		Meta: model.MetaInfo{
			Page:   page,
			Limit:  limit,
			Total:  int(total),
			Pages:  int((total + int64(limit) - 1) / int64(limit)),
			SortBy: sortBy,
			Order:  order,
			Search: search,
		},
	}
	return c.JSON(response)
}


// Restore godoc
// @Summary Restore soft-deleted pekerjaan
// @Description Mengembalikan pekerjaan yang telah di-soft delete
// @Tags Pekerjaan
// @Accept json
// @Produce json
// @Param id query string true "ID pekerjaan"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /TM4/pekerjaan/restore [put]
// === RESTORE ===
func (s *PekerjaanService) Restore(c *fiber.Ctx) error {
	idParam := c.Params("id")

	if err := s.repo.Restore(context.Background(), idParam); err != nil {
		return helper.ResponseJSON(c, 500, err.Error(), false, nil)
	}

	return helper.ResponseJSON(c, 200, "Pekerjaan restored successfully", true, nil)
}

// HardDelete godoc
// @Summary Permanently delete all pekerjaan
// @Description Menghapus seluruh data pekerjaan secara permanen (tidak dapat dikembalikan)
// @Tags Pekerjaan
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /TM4/pekerjaan/harddelete [delete]

// === HARD DELETE ===
func (s *PekerjaanService) HardDelete(c *fiber.Ctx) error {
	if err := s.repo.HardDelete(context.Background()); err != nil {
		return helper.ResponseJSON(c, 500, err.Error(), false, nil)
	}

	return helper.ResponseJSON(c, 200, "Pekerjaan permanently deleted", true, nil)
}

