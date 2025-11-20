package repository

import (
	"TM4/app/model"
	// "database/sql"
	"context"
	"time"
    "strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

)


type IAlumniRepository interface {
	GetAlumni(ctx context.Context, search string, limit, offset int, sortBy, order string) ([]model.Alumni, error)
	GetAlumniByID(ctx context.Context, id string) (*model.Alumni, error)
	GetByEmail(ctx context.Context, email string) (*model.Alumni, error)
	CreateAlumni(ctx context.Context, req *model.Alumni) (*model.Alumni, error)
	UpdateAlumni(ctx context.Context, id string, req *model.UpdateAlumni) error
	DeleteAlumni(ctx context.Context, id string) error
	Count(ctx context.Context, search string) (int64, error)
}


type AlumniRepository struct {
	collection *mongo.Collection
}

func NewAlumniRepository(db *mongo.Database) IAlumniRepository {
	return &AlumniRepository{
		collection: db.Collection("alumni"),
	}
}

func (r *AlumniRepository) Count(ctx context.Context, search string) (int64, error) {
	filter := bson.M{}
	if search != "" {
		filter = bson.M{
			"$or": []bson.M{
				{"nama": bson.M{"$regex": search, "$options": "i"}},
				{"nim": bson.M{"$regex": search, "$options": "i"}},
				{"email": bson.M{"$regex": search, "$options": "i"}},
				{"jurusan": bson.M{"$regex": search, "$options": "i"}},
			},
		}
	}

	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, err
	}

	return int64(count), nil
}


func (r *AlumniRepository) GetAlumni(ctx context.Context, search string, limit, offset int, sortBy, order string) ([]model.Alumni, error) {
	filter := bson.M{}
	if search != "" {
		filter = bson.M{
			"$or": []bson.M{
				{"nama": bson.M{"$regex": search, "$options": "i"}},
				{"nim": bson.M{"$regex": search, "$options": "i"}},
				{"email": bson.M{"$regex": search, "$options": "i"}},
				{"jurusan": bson.M{"$regex": search, "$options": "i"}},
			},
		}
	}

	findOptions := mongoOptions(limit, offset, sortBy, order)

	cursor, err := r.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var alumniList []model.Alumni
	if err := cursor.All(ctx, &alumniList); err != nil {
		return nil, err
	}

	return alumniList, nil
}


func mongoOptions(limit, offset int, sortBy, order string) *options.FindOptions {
	opts := &options.FindOptions{}
	if limit > 0 {
		opts.SetLimit(int64(limit))
		opts.SetSkip(int64(offset))
	}
	if sortBy != "" {
		sortOrder := 1
		if strings.ToLower(order) == "desc" {
			sortOrder = -1
		}
		opts.SetSort(bson.D{{Key: sortBy, Value: sortOrder}})
	}
	return opts
}


func (r *AlumniRepository) GetByEmail(ctx context.Context, email string) (*model.Alumni, error) {
	var alumni model.Alumni
	filter := bson.M{"email": (email)}

	err := r.collection.FindOne(ctx, filter).Decode(&alumni)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &alumni, nil
}


func (r *AlumniRepository) GetAlumniByID(ctx context.Context, id string) (*model.Alumni, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var alumni model.Alumni
	err = r.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&alumni)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &alumni, nil
}


func (r *AlumniRepository) CreateAlumni(ctx context.Context, req *model.Alumni) (*model.Alumni, error) {
	req.ID = primitive.NewObjectID()
	req.CreatedAt = time.Now()
	req.UpdatedAt = time.Now()

	result, err := r.collection.InsertOne(ctx, req)
	if err != nil {
		return nil, err
	}

	req.ID = result.InsertedID.(primitive.ObjectID)
	return req, nil
}


func (r *AlumniRepository) UpdateAlumni(ctx context.Context, id string, req *model.UpdateAlumni) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	update := bson.M{
		"$set": bson.M{
			"nama":        req.Nama,
			"jurusan":     req.Jurusan,
			"angkatan":    req.Angkatan,
			"tahun_lulus": req.TahunLulus,
			"email":       req.Email,
			"no_telepon":  req.NoTelepon,
			"alamat":      req.Alamat,
			"role":        req.Role,
			"password":    req.Password,
			"updated_at":  time.Now(), // langsung diset di query, bukan di struct
		},
	}

	_, err = r.collection.UpdateOne(ctx, bson.M{"_id": objID}, update)
	return err
}


func (r *AlumniRepository) DeleteAlumni(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = r.collection.DeleteOne(ctx, bson.M{"_id": objID})
	return err
}
