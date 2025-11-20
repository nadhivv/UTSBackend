package service

import (
	"context"
	"errors"
	"TM4/app/model"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MockAlumniRepository struct {
	MockGetAlumni     func(ctx context.Context, search string, limit, offset int, sortBy, order string) ([]model.Alumni, error)
	MockGetAlumniByID func(ctx context.Context, id string) (*model.Alumni, error)
	MockGetByEmail    func(ctx context.Context, email string) (*model.Alumni, error)
	MockCreateAlumni  func(ctx context.Context, req *model.Alumni) (*model.Alumni, error)
	MockUpdateAlumni  func(ctx context.Context, id string, req *model.UpdateAlumni) error
	MockDeleteAlumni  func(ctx context.Context, id string) error
	MockCount         func(ctx context.Context, search string) (int64, error)
}

func (m *MockAlumniRepository) GetAlumni(ctx context.Context, search string, limit, offset int, sortBy, order string) ([]model.Alumni, error) {
	return m.MockGetAlumni(ctx, search, limit, offset, sortBy, order)
}

func (m *MockAlumniRepository) GetAlumniByID(ctx context.Context, id string) (*model.Alumni, error) {
	return m.MockGetAlumniByID(ctx, id)
}

func (m *MockAlumniRepository) GetByEmail(ctx context.Context, email string) (*model.Alumni, error) {
	return m.MockGetByEmail(ctx, email)
}

func (m *MockAlumniRepository) CreateAlumni(ctx context.Context, req *model.Alumni) (*model.Alumni, error) {
	return m.MockCreateAlumni(ctx, req)
}

func (m *MockAlumniRepository) UpdateAlumni(ctx context.Context, id string, req *model.UpdateAlumni) error {
	return m.MockUpdateAlumni(ctx, id, req)
}

func (m *MockAlumniRepository) DeleteAlumni(ctx context.Context, id string) error {
	return m.MockDeleteAlumni(ctx, id)
}

func (m *MockAlumniRepository) Count(ctx context.Context, search string) (int64, error) {
	return m.MockCount(ctx, search)
}


func TestGetAlumni(t *testing.T) {
	mockRepo := &MockAlumniRepository{}
	svc := AlumniService{repo: mockRepo}
	ctx := context.Background()

	mockData := []model.Alumni{
		{ID: primitive.NewObjectID(), NIM: "123", Nama: "Janu", Email: "janu@example.com"},
		{ID: primitive.NewObjectID(), NIM: "124", Nama: "Dika", Email: "dika@example.com"},
	}

	mockRepo.MockGetAlumni = func(ctx context.Context, search string, limit, offset int, sortBy, order string) ([]model.Alumni, error) {
		return mockData, nil
	}

	alumni, err := svc.repo.GetAlumni(ctx, "", 10, 0, "nama", "asc")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(alumni) != 2 {
		t.Errorf("Expected 2 alumni, got %d", len(alumni))
	}
}

func TestGetByID(t *testing.T) {
	mockRepo := &MockAlumniRepository{}
	svc := AlumniService{repo: mockRepo}
	ctx := context.Background()

	expected := &model.Alumni{
		ID:          primitive.NewObjectID(),
		NIM:         "123456",
		Nama:        "John Doe",
		Email:       "john@example.com",
		Jurusan:     "Teknik Informatika",
		Angkatan:    2020,
		TahunLulus:  2024,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		IsDeleted:   false,
	}

	mockRepo.MockGetAlumniByID = func(ctx context.Context, id string) (*model.Alumni, error) {
		return expected, nil
	}

	result, err := svc.repo.GetAlumniByID(ctx, "123")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if result.Nama != expected.Nama {
		t.Errorf("Expected name %v, got %v", expected.Nama, result.Nama)
	}
}

func TestGetByEmail(t *testing.T) {
	mockRepo := &MockAlumniRepository{}
	svc := AlumniService{repo: mockRepo}
	ctx := context.Background()

	expected := &model.Alumni{
		ID:    primitive.NewObjectID(),
		Nama:  "John Doe",
		Email: "john@example.com",
	}

	mockRepo.MockGetByEmail = func(ctx context.Context, email string) (*model.Alumni, error) {
		return expected, nil
	}

	result, err := svc.repo.GetByEmail(ctx, "john@example.com")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if result.Email != expected.Email {
		t.Errorf("Expected email %v, got %v", expected.Email, result.Email)
	}
}

func TestCreateAlumni(t *testing.T) {
	mockRepo := &MockAlumniRepository{}
	svc := AlumniService{repo: mockRepo}
	ctx := context.Background()

	input := &model.Alumni{
		NIM:        "123456",
		Nama:       "Janu",
		Email:      "janu@example.com",
		Jurusan:    "Teknik Informatika",
		Angkatan:   2020,
		TahunLulus: 2024,
		Password:   "hashed_password",
	}

	mockRepo.MockCreateAlumni = func(ctx context.Context, alumni *model.Alumni) (*model.Alumni, error) {
		alumni.ID = primitive.NewObjectID()
		alumni.CreatedAt = time.Now()
		alumni.UpdatedAt = time.Now()
		return alumni, nil
	}

	result, err := svc.repo.CreateAlumni(ctx, input)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if result.ID.IsZero() {
		t.Errorf("Expected created ID, got zero object ID")
	}
}

func TestUpdateAlumni(t *testing.T) {
	mockRepo := &MockAlumniRepository{}
	svc := AlumniService{repo: mockRepo}
	ctx := context.Background()

	mockRepo.MockUpdateAlumni = func(ctx context.Context, id string, alumni *model.UpdateAlumni) error {
		return nil
	}

	err := svc.repo.UpdateAlumni(ctx, "123", &model.UpdateAlumni{
		Nama:   "Updated Name",
		Email:  "updated@example.com",
	})
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
}

func TestDeleteAlumni(t *testing.T) {
	mockRepo := &MockAlumniRepository{}
	svc := AlumniService{repo: mockRepo}
	ctx := context.Background()

	mockRepo.MockDeleteAlumni = func(ctx context.Context, id string) error {
		return nil
	}

	err := svc.repo.DeleteAlumni(ctx, "123")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
}

func TestCountAlumni(t *testing.T) {
	mockRepo := &MockAlumniRepository{}
	svc := AlumniService{repo: mockRepo}
	ctx := context.Background()

	mockRepo.MockCount = func(ctx context.Context, search string) (int64, error) {
		return 25, nil
	}

	count, err := svc.repo.Count(ctx, "")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if count != 25 {
		t.Errorf("Expected count 25, got %d", count)
	}
}


func TestGetByID_NotFound(t *testing.T) {
	mockRepo := &MockAlumniRepository{}
	svc := AlumniService{repo: mockRepo}
	ctx := context.Background()

	mockRepo.MockGetAlumniByID = func(ctx context.Context, id string) (*model.Alumni, error) {
		return nil, errors.New("not found")
	}

	_, err := svc.repo.GetAlumniByID(ctx, "missing-id")
	if err == nil {
		t.Fatalf("Expected error, got nil")
	}
}

func TestGetByEmail_NotFound(t *testing.T) {
	mockRepo := &MockAlumniRepository{}
	svc := AlumniService{repo: mockRepo}
	ctx := context.Background()

	mockRepo.MockGetByEmail = func(ctx context.Context, email string) (*model.Alumni, error) {
		return nil, errors.New("email not found")
	}

	_, err := svc.repo.GetByEmail(ctx, "nonexistent@example.com")
	if err == nil {
		t.Fatalf("Expected error, got nil")
	}
}

func TestCreateAlumni_Error(t *testing.T) {
	mockRepo := &MockAlumniRepository{}
	svc := AlumniService{repo: mockRepo}
	ctx := context.Background()

	mockRepo.MockCreateAlumni = func(ctx context.Context, alumni *model.Alumni) (*model.Alumni, error) {
		return nil, errors.New("insert error")
	}

	_, err := svc.repo.CreateAlumni(ctx, &model.Alumni{})
	if err == nil {
		t.Fatalf("Expected error, got nil")
	}
}

func TestCountAlumni_Error(t *testing.T) {
	mockRepo := &MockAlumniRepository{}
	svc := AlumniService{repo: mockRepo}
	ctx := context.Background()

	mockRepo.MockCount = func(ctx context.Context, search string) (int64, error) {
		return 0, errors.New("count error")
	}

	_, err := svc.repo.Count(ctx, "")
	if err == nil {
		t.Fatalf("Expected error, got nil")
	}
}