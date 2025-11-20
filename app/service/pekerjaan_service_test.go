package service

import (
	"context"
	"errors"
	"TM4/app/model"
	"testing"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MockPekerjaanRepository struct {
	MockGetByID       func(ctx context.Context, id string) (*model.Pekerjaan, error)
	MockGetByAlumniID func(ctx context.Context, alumniID primitive.ObjectID) ([]model.Pekerjaan, error)
	MockCreate        func(ctx context.Context, input *model.CreatePekerjaan) (*mongo.InsertOneResult, error)
	MockUpdate        func(ctx context.Context, id string, input *model.UpdatePekerjaan) error
	MockDelete        func(ctx context.Context, id string) error
	MockGetAll        func(ctx context.Context, search, sortBy, order string, limit, offset int) ([]model.Pekerjaan, error)
	MockCount         func(ctx context.Context, search string) (int64, error)
	MockSoftDelete    func(ctx context.Context, id string) error
	MockSoftDeleteBulk func(ctx context.Context) error // TAMBAH INI
	MockTrash         func(ctx context.Context, search, sortBy, order string, limit, offset int) ([]model.Pekerjaan, error)
	MockCountTrash    func(ctx context.Context, search string) (int64, error)
	MockRestore       func(ctx context.Context, id string) error
	MockHardDelete    func(ctx context.Context) error
}

func (m *MockPekerjaanRepository) GetByID(ctx context.Context, id string) (*model.Pekerjaan, error) {
	return m.MockGetByID(ctx, id)
}

func (m *MockPekerjaanRepository) GetByAlumniID(ctx context.Context, alumniID primitive.ObjectID) ([]model.Pekerjaan, error) {
	return m.MockGetByAlumniID(ctx, alumniID)
}

func (m *MockPekerjaanRepository) Create(ctx context.Context, input *model.CreatePekerjaan) (*mongo.InsertOneResult, error) {
	return m.MockCreate(ctx, input)
}

func (m *MockPekerjaanRepository) Update(ctx context.Context, id string, input *model.UpdatePekerjaan) error {
	return m.MockUpdate(ctx, id, input)
}

func (m *MockPekerjaanRepository) Delete(ctx context.Context, id string) error {
	return m.MockDelete(ctx, id)
}

func (m *MockPekerjaanRepository) GetAll(ctx context.Context, search, sortBy, order string, limit, offset int) ([]model.Pekerjaan, error) {
	return m.MockGetAll(ctx, search, sortBy, order, limit, offset)
}

func (m *MockPekerjaanRepository) Count(ctx context.Context, search string) (int64, error) {
	return m.MockCount(ctx, search)
}

func (m *MockPekerjaanRepository) SoftDelete(ctx context.Context, id string) error {
	return m.MockSoftDelete(ctx, id)
}


func (m *MockPekerjaanRepository) SoftDeleteBulk(ctx context.Context) error {
	return m.MockSoftDeleteBulk(ctx)
}

func (m *MockPekerjaanRepository) Trash(ctx context.Context, search, sortBy, order string, limit, offset int) ([]model.Pekerjaan, error) {
	return m.MockTrash(ctx, search, sortBy, order, limit, offset)
}

func (m *MockPekerjaanRepository) CountTrash(ctx context.Context, search string) (int64, error) {
	return m.MockCountTrash(ctx, search)
}

func (m *MockPekerjaanRepository) Restore(ctx context.Context, id string) error {
	return m.MockRestore(ctx, id)
}

func (m *MockPekerjaanRepository) HardDelete(ctx context.Context) error {
	return m.MockHardDelete(ctx)
}

func TestGetByID_Pekerjaan(t *testing.T) {
	mockRepo := &MockPekerjaanRepository{}
	svc := PekerjaanService{repo: mockRepo}
	ctx := context.Background()

	expected := &model.Pekerjaan{
		ID:             primitive.NewObjectID(),
		PosisiJabatan:  "Software Engineer",
		NamaPerusahaan: "Tech Company",
	}

	mockRepo.MockGetByID = func(ctx context.Context, id string) (*model.Pekerjaan, error) {
		return expected, nil
	}

	result, err := svc.repo.GetByID(ctx, "123")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if result.PosisiJabatan != expected.PosisiJabatan {
		t.Errorf("Expected PosisiJabatan %v, got %v", expected.PosisiJabatan, result.PosisiJabatan)
	}
}

func TestGetByAlumniID_Pekerjaan(t *testing.T) {
	mockRepo := &MockPekerjaanRepository{}
	svc := PekerjaanService{repo: mockRepo}
	ctx := context.Background()

	alumniID := primitive.NewObjectID()
	mockData := []model.Pekerjaan{
		{ID: primitive.NewObjectID(), PosisiJabatan: "Backend Developer", AlumniID: alumniID},
		{ID: primitive.NewObjectID(), PosisiJabatan: "Frontend Developer", AlumniID: alumniID},
	}

	mockRepo.MockGetByAlumniID = func(ctx context.Context, id primitive.ObjectID) ([]model.Pekerjaan, error) {
		return mockData, nil
	}

	result, err := svc.repo.GetByAlumniID(ctx, alumniID)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(result) != 2 {
		t.Errorf("Expected 2 pekerjaan, got %d", len(result))
	}
}

func TestCreatePekerjaan(t *testing.T) {
	mockRepo := &MockPekerjaanRepository{}
	svc := PekerjaanService{repo: mockRepo}
	ctx := context.Background()

	input := &model.CreatePekerjaan{
		PosisiJabatan:  "Software Engineer",
		NamaPerusahaan: "Tech Corp",
	}

	mockRepo.MockCreate = func(ctx context.Context, req *model.CreatePekerjaan) (*mongo.InsertOneResult, error) {
		return &mongo.InsertOneResult{InsertedID: primitive.NewObjectID()}, nil
	}

	result, err := svc.repo.Create(ctx, input)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if result.InsertedID == nil {
		t.Errorf("Expected inserted ID, got nil")
	}
}

func TestUpdatePekerjaan(t *testing.T) {
	mockRepo := &MockPekerjaanRepository{}
	svc := PekerjaanService{repo: mockRepo}
	ctx := context.Background()

	mockRepo.MockUpdate = func(ctx context.Context, id string, req *model.UpdatePekerjaan) error {
		return nil
	}

	err := svc.repo.Update(ctx, "123", &model.UpdatePekerjaan{PosisiJabatan: "Updated Position"})
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
}

func TestDeletePekerjaan(t *testing.T) {
	mockRepo := &MockPekerjaanRepository{}
	svc := PekerjaanService{repo: mockRepo}
	ctx := context.Background()

	mockRepo.MockDelete = func(ctx context.Context, id string) error {
		return nil
	}

	err := svc.repo.Delete(ctx, "123")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
}

func TestSoftDeleteBulk_Pekerjaan(t *testing.T) {
	mockRepo := &MockPekerjaanRepository{}
	svc := PekerjaanService{repo: mockRepo}
	ctx := context.Background()

	mockRepo.MockSoftDeleteBulk = func(ctx context.Context) error {
		return nil
	}

	err := svc.repo.SoftDeleteBulk(ctx)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
}


func TestGetByID_Pekerjaan_NotFound(t *testing.T) {
	mockRepo := &MockPekerjaanRepository{}
	svc := PekerjaanService{repo: mockRepo}
	ctx := context.Background()

	mockRepo.MockGetByID = func(ctx context.Context, id string) (*model.Pekerjaan, error) {
		return nil, errors.New("not found")
	}

	_, err := svc.repo.GetByID(ctx, "missing-id")
	if err == nil {
		t.Fatalf("Expected error, got nil")
	}
}

func TestCreatePekerjaan_Error(t *testing.T) {
	mockRepo := &MockPekerjaanRepository{}
	svc := PekerjaanService{repo: mockRepo}
	ctx := context.Background()

	mockRepo.MockCreate = func(ctx context.Context, req *model.CreatePekerjaan) (*mongo.InsertOneResult, error) {
		return nil, errors.New("insert error")
	}

	_, err := svc.repo.Create(ctx, &model.CreatePekerjaan{})
	if err == nil {
		t.Fatalf("Expected error, got nil")
	}
}

func TestSoftDeleteBulk_Pekerjaan_Error(t *testing.T) {
	mockRepo := &MockPekerjaanRepository{}
	svc := PekerjaanService{repo: mockRepo}
	ctx := context.Background()

	mockRepo.MockSoftDeleteBulk = func(ctx context.Context) error {
		return errors.New("bulk delete error")
	}

	err := svc.repo.SoftDeleteBulk(ctx)
	if err == nil {
		t.Fatalf("Expected error, got nil")
	}
}