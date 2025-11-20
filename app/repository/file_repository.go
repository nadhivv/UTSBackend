package repository

import (
	"context"
	// "errors"
	"TM4/app/model"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)


type IFileRepository interface {
	GetAll(ctx context.Context, search, sortBy, order string, limit, offset int) ([]model.File, error)
	Count(ctx context.Context, search string) (int64, error)
	GetByID(ctx context.Context, id primitive.ObjectID) (*model.File, error)
	GetByAlumniID(ctx context.Context, alumniID string) ([]model.File, error)
	Create(ctx context.Context, file model.File) (*model.File, error)
	Update(ctx context.Context, id primitive.ObjectID, file model.File) error
	Delete(ctx context.Context, id primitive.ObjectID) error
}


type fileRepository struct {
	collection *mongo.Collection
}


func NewFileRepository(db *mongo.Database) IFileRepository {
	return &fileRepository{
		collection: db.Collection("files"),
	}
}

func (r *fileRepository) GetAll(ctx context.Context, search, sortBy, order string, limit, offset int) ([]model.File, error) {
	filter := bson.M{"is_deleted": false}

	if search != "" {
		filter["$or"] = []bson.M{
			{"file_name": bson.M{"$regex": search, "$options": "i"}},
			{"original_name": bson.M{"$regex": search, "$options": "i"}},
		}
	}

	sortOrder := 1
	if order == "DESC" || order == "desc" {
		sortOrder = -1
	}
	if sortBy == "" {
		sortBy = "uploaded_at"
	}

	opts := options.Find().
		SetSort(bson.M{sortBy: sortOrder}).
		SetLimit(int64(limit)).
		SetSkip(int64(offset))

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var list []model.File
	if err := cursor.All(ctx, &list); err != nil {
		return nil, err
	}

	return list, nil
}

func (r *fileRepository) Count(ctx context.Context, search string) (int64, error) {
	filter := bson.M{"is_deleted": false}

	if search != "" {
		filter["$or"] = []bson.M{
			{"file_name": bson.M{"$regex": search, "$options": "i"}},
			{"original_name": bson.M{"$regex": search, "$options": "i"}},
		}
	}

	return r.collection.CountDocuments(ctx, filter)
}

func (r *fileRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*model.File, error) {
	var file model.File
	filter := bson.M{"_id": id, "is_deleted": false}

	err := r.collection.FindOne(ctx, filter).Decode(&file)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &file, nil
}

func (r *fileRepository) GetByAlumniID(ctx context.Context, alumniID string) ([]model.File, error) {
	filter := bson.M{"alumni_id": alumniID, "is_deleted": false}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var list []model.File
	if err := cursor.All(ctx, &list); err != nil {
		return nil, err
	}

	return list, nil
}

func (r *fileRepository) Create(ctx context.Context, file model.File) (*model.File, error) {
	now := time.Now()
	file.UploadedAt = now
	doc := bson.M{
		"file_name":     file.FileName,
		"original_name": file.OriginalName,
		"file_path":     file.FilePath,
		"file_size":     file.FileSize,
		"file_type":     file.FileType,
		"uploaded_at":   file.UploadedAt,
		"is_deleted":    false,
		"alumni_id":     file.AlumniID,
	}

	res, err := r.collection.InsertOne(ctx, doc)
	if err != nil {
		return nil, err
	}

	file.ID = res.InsertedID.(primitive.ObjectID)
	return &file, nil
}

func (r *fileRepository) Update(ctx context.Context, id primitive.ObjectID, file model.File) error {
	filter := bson.M{"_id": id, "is_deleted": false}

	update := bson.M{
		"$set": bson.M{
			"file_name":     file.FileName,
			"original_name": file.OriginalName,
			"file_path":     file.FilePath,
			"file_size":     file.FileSize,
			"file_type":     file.FileType,
			"uploaded_at":   file.UploadedAt,
		},
	}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}

func (r *fileRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	filter := bson.M{"_id": id}

	result, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}