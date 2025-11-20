package service

import (
	"context"
	"errors"
	"TM4/app/model"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MockFileRepository struct {
	MockCreate      func(ctx context.Context, file model.File) (*model.File, error)
	MockGetAll      func(ctx context.Context, search, sortBy, order string, limit, offset int) ([]model.File, error)
	MockCount       func(ctx context.Context, search string) (int64, error) // TAMBAH INI
	MockGetByID     func(ctx context.Context, id primitive.ObjectID) (*model.File, error)
	MockGetByAlumniID func(ctx context.Context, alumniID string) ([]model.File, error)
	MockUpdate      func(ctx context.Context, id primitive.ObjectID, file model.File) error
	MockDelete      func(ctx context.Context, id primitive.ObjectID) error
}

func (m *MockFileRepository) Create(ctx context.Context, file model.File) (*model.File, error) {
	return m.MockCreate(ctx, file)
}

func (m *MockFileRepository) GetAll(ctx context.Context, search, sortBy, order string, limit, offset int) ([]model.File, error) {
	return m.MockGetAll(ctx, search, sortBy, order, limit, offset)
}

func (m *MockFileRepository) Count(ctx context.Context, search string) (int64, error) {
	return m.MockCount(ctx, search)
}

func (m *MockFileRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*model.File, error) {
	return m.MockGetByID(ctx, id)
}

func (m *MockFileRepository) GetByAlumniID(ctx context.Context, alumniID string) ([]model.File, error) {
	return m.MockGetByAlumniID(ctx, alumniID)
}

func (m *MockFileRepository) Update(ctx context.Context, id primitive.ObjectID, file model.File) error {
	return m.MockUpdate(ctx, id, file)
}

func (m *MockFileRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	return m.MockDelete(ctx, id)
}


func TestGetFileByID(t *testing.T) {
	mockRepo := &MockFileRepository{}
	tempDir := t.TempDir()
	svc := FileService{repo: mockRepo, uploadPath: tempDir}
	ctx := context.Background()

	expected := &model.File{
		ID:           primitive.NewObjectID(),
		FileName:     "test.jpg",
		OriginalName: "test.jpg",
		FilePath:     "/path/to/file",
		FileSize:     1024,
		FileType:     "image/jpeg",
		UploadedAt:   time.Now(),
		AlumniID:     primitive.NewObjectID(),
	}

	mockRepo.MockGetByID = func(ctx context.Context, id primitive.ObjectID) (*model.File, error) {
		return expected, nil
	}

	result, err := svc.repo.GetByID(ctx, expected.ID)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if result.FileName != expected.FileName {
		t.Errorf("Expected filename %v, got %v", expected.FileName, result.FileName)
	}
}

func TestGetAllFiles(t *testing.T) {
	mockRepo := &MockFileRepository{}
	tempDir := t.TempDir()
	svc := FileService{repo: mockRepo, uploadPath: tempDir}
	ctx := context.Background()

	mockData := []model.File{
		{
			ID:           primitive.NewObjectID(),
			FileName:     "file1.jpg",
			OriginalName: "file1.jpg",
			FilePath:     "/path/to/file1",
			FileSize:     1024,
			FileType:     "image/jpeg",
			UploadedAt:   time.Now(),
			AlumniID:     primitive.NewObjectID(),
		},
		{
			ID:           primitive.NewObjectID(),
			FileName:     "file2.pdf",
			OriginalName: "file2.pdf",
			FilePath:     "/path/to/file2",
			FileSize:     2048,
			FileType:     "application/pdf",
			UploadedAt:   time.Now(),
			AlumniID:     primitive.NewObjectID(),
		},
	}

	mockRepo.MockGetAll = func(ctx context.Context, search, sortBy, order string, limit, offset int) ([]model.File, error) {
		return mockData, nil
	}

	result, err := svc.repo.GetAll(ctx, "", "uploadedAt", "desc", 10, 0)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(result) != 2 {
		t.Errorf("Expected 2 files, got %d", len(result))
	}
}

func TestCreateFile(t *testing.T) {
	mockRepo := &MockFileRepository{}
	tempDir := t.TempDir()
	svc := FileService{repo: mockRepo, uploadPath: tempDir}
	ctx := context.Background()

	input := model.File{
		FileName:     "test.jpg",
		OriginalName: "test.jpg",
		FilePath:     "/uploads/test.jpg",
		FileSize:     1024,
		FileType:     "image/jpeg",
		UploadedAt:   time.Now(),
		AlumniID:     primitive.NewObjectID(),
	}

	mockRepo.MockCreate = func(ctx context.Context, file model.File) (*model.File, error) {
		file.ID = primitive.NewObjectID()
		return &file, nil
	}

	result, err := svc.repo.Create(ctx, input)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if result.ID.IsZero() {
		t.Errorf("Expected created ID, got zero object ID")
	}
}

func TestUpdateFile(t *testing.T) {
	mockRepo := &MockFileRepository{}
	tempDir := t.TempDir()
	svc := FileService{repo: mockRepo, uploadPath: tempDir}
	ctx := context.Background()

	fileID := primitive.NewObjectID()
	file := model.File{
		ID:       fileID,
		FileName: "updated.jpg",
	}

	mockRepo.MockUpdate = func(ctx context.Context, id primitive.ObjectID, file model.File) error {
		return nil
	}

	err := svc.repo.Update(ctx, fileID, file)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
}

func TestDeleteFile(t *testing.T) {
	mockRepo := &MockFileRepository{}
	tempDir := t.TempDir()
	svc := FileService{repo: mockRepo, uploadPath: tempDir}
	ctx := context.Background()

	fileID := primitive.NewObjectID()

	mockRepo.MockDelete = func(ctx context.Context, id primitive.ObjectID) error {
		return nil
	}

	err := svc.repo.Delete(ctx, fileID)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
}

func TestCountFiles(t *testing.T) {
	mockRepo := &MockFileRepository{}
	tempDir := t.TempDir()
	svc := FileService{repo: mockRepo, uploadPath: tempDir}
	ctx := context.Background()

	mockRepo.MockCount = func(ctx context.Context, search string) (int64, error) {
		return 5, nil
	}

	count, err := svc.repo.Count(ctx, "")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if count != 5 {
		t.Errorf("Expected count 5, got %d", count)
	}
}


func TestGetFileByID_NotFound(t *testing.T) {
	mockRepo := &MockFileRepository{}
	tempDir := t.TempDir()
	svc := FileService{repo: mockRepo, uploadPath: tempDir}
	ctx := context.Background()

	mockRepo.MockGetByID = func(ctx context.Context, id primitive.ObjectID) (*model.File, error) {
		return nil, errors.New("file not found")
	}

	_, err := svc.repo.GetByID(ctx, primitive.NewObjectID())
	if err == nil {
		t.Fatalf("Expected error, got nil")
	}
}

func TestCreateFile_Error(t *testing.T) {
	mockRepo := &MockFileRepository{}
	tempDir := t.TempDir()
	svc := FileService{repo: mockRepo, uploadPath: tempDir}
	ctx := context.Background()

	mockRepo.MockCreate = func(ctx context.Context, file model.File) (*model.File, error) {
		return nil, errors.New("create error")
	}

	_, err := svc.repo.Create(ctx, model.File{})
	if err == nil {
		t.Fatalf("Expected error, got nil")
	}
}

func TestCountFiles_Error(t *testing.T) {
	mockRepo := &MockFileRepository{}
	tempDir := t.TempDir()
	svc := FileService{repo: mockRepo, uploadPath: tempDir}
	ctx := context.Background()

	mockRepo.MockCount = func(ctx context.Context, search string) (int64, error) {
		return 0, errors.New("count error")
	}

	_, err := svc.repo.Count(ctx, "")
	if err == nil {
		t.Fatalf("Expected error, got nil")
	}
}