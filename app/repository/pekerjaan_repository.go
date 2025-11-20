package repository

import (
	"TM4/app/model"
	// "database/sql"
	"context"
	"time"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)


type IPekerjaanRepository interface {
	GetAll(ctx context.Context, search, sortBy, order string, limit, offset int) ([]model.Pekerjaan, error)
	Count(ctx context.Context, search string) (int64, error)
	GetByID(ctx context.Context, id string) (*model.Pekerjaan, error)
	GetByAlumniID(ctx context.Context, alumniID primitive.ObjectID) ([]model.Pekerjaan, error)
	Create(ctx context.Context, req *model.CreatePekerjaan) (*mongo.InsertOneResult, error)
	Update(ctx context.Context, id string, req *model.UpdatePekerjaan) error
	Delete(ctx context.Context, id string) error
	SoftDelete(ctx context.Context, id string) error
	SoftDeleteBulk(ctx context.Context) error
	Restore(ctx context.Context, id string) error
	Trash(ctx context.Context, search, sortBy, order string, limit, offset int) ([]model.Pekerjaan, error)
	CountTrash(ctx context.Context, search string) (int64, error)
	HardDelete(ctx context.Context) error
}

type PekerjaanRepository struct {
	collection *mongo.Collection
}

func NewPekerjaanRepository(db *mongo.Database) IPekerjaanRepository {
	return &PekerjaanRepository{
		collection: db.Collection("pekerjaan_alumni"),
	}
}

func (r *PekerjaanRepository) GetAll(ctx context.Context, search, sortBy, order string, limit, offset int) ([]model.Pekerjaan, error) {
	filter := bson.M{"isdeleted": false}

	if search != "" {
		filter["$or"] = []bson.M{
			{"nama_perusahaan": bson.M{"$regex": search, "$options": "i"}},
			{"posisi_jabatan": bson.M{"$regex": search, "$options": "i"}},
			{"bidang_industri": bson.M{"$regex": search, "$options": "i"}},
			{"lokasi_kerja": bson.M{"$regex": search, "$options": "i"}},
		}
	}

	sortOrder := 1
	if order == "desc" {
		sortOrder = -1
	}

	opts := options.Find().
		SetSort(bson.D{{Key: sortBy, Value: sortOrder}}).
		SetLimit(int64(limit)).
		SetSkip(int64(offset))

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var list []model.Pekerjaan
	if err := cursor.All(ctx, &list); err != nil {
		return nil, err
	}
	return list, nil
}

func (r *PekerjaanRepository) Count(ctx context.Context, search string) (int64, error) {
	filter := bson.M{"isdeleted": false}
	if search != "" {
		filter["$or"] = []bson.M{
			{"nama_perusahaan": bson.M{"$regex": search, "$options": "i"}},
			{"posisi_jabatan": bson.M{"$regex": search, "$options": "i"}},
			{"bidang_industri": bson.M{"$regex": search, "$options": "i"}},
			{"lokasi_kerja": bson.M{"$regex": search, "$options": "i"}},
		}
	}
	return r.collection.CountDocuments(ctx, filter)
}

func (r *PekerjaanRepository) GetByID(ctx context.Context, id string) (*model.Pekerjaan, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var pekerjaan model.Pekerjaan
	err = r.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&pekerjaan)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &pekerjaan, nil
}

func (r *PekerjaanRepository) GetByAlumniID(ctx context.Context, alumniID primitive.ObjectID) ([]model.Pekerjaan, error) {
	filter := bson.M{"alumni_id": alumniID, "isdeleted": false}
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var list []model.Pekerjaan
	if err := cursor.All(ctx, &list); err != nil {
		return nil, err
	}
	return list, nil
}

func (r *PekerjaanRepository) Create(ctx context.Context, req *model.CreatePekerjaan) (*mongo.InsertOneResult, error) {
	now := time.Now()

	alumniID, err := primitive.ObjectIDFromHex(req.AlumniID)
	if err != nil {
		return nil, fmt.Errorf("invalid alumni_id: %v", err)
	}

	doc := bson.M{
		"alumni_id":             alumniID,
		"nama_perusahaan":       req.NamaPerusahaan,
		"posisi_jabatan":        req.PosisiJabatan,
		"bidang_industri":       req.BidangIndustri,
		"lokasi_kerja":          req.LokasiKerja,
		"gaji_range":            req.GajiRange,
		"tanggal_mulai_kerja":   req.TanggalMulaiKerja,
		"tanggal_selesai_kerja": req.TanggalSelesaiKerja,
		"status_pekerjaan":      req.StatusPekerjaan,
		"deskripsi_pekerjaan":   req.DeskripsiPekerjaan,
		"created_at":            now,
		"updated_at":            now,
		"isdeleted":             false,
	}

	return r.collection.InsertOne(ctx, doc)
}

func (r *PekerjaanRepository) Update(ctx context.Context, id string, req *model.UpdatePekerjaan) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	update := bson.M{"$set": bson.M{
		"nama_perusahaan":       req.NamaPerusahaan,
		"posisi_jabatan":        req.PosisiJabatan,
		"bidang_industri":       req.BidangIndustri,
		"lokasi_kerja":          req.LokasiKerja,
		"gaji_range":            req.GajiRange,
		"tanggal_mulai_kerja":   req.TanggalMulaiKerja,
		"tanggal_selesai_kerja": req.TanggalSelesaiKerja,
		"status_pekerjaan":      req.StatusPekerjaan,
		"deskripsi_pekerjaan":   req.DeskripsiPekerjaan,
		"updated_at":            time.Now(),
	}}
	_, err = r.collection.UpdateOne(ctx, bson.M{"_id": objID}, update)
	return err
}

func (r *PekerjaanRepository) Delete(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = r.collection.DeleteOne(ctx, bson.M{"_id": objID})
	return err
}

func (r *PekerjaanRepository) SoftDelete(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	update := bson.M{"$set": bson.M{"isdeleted": true, "updated_at": time.Now()}}
	_, err = r.collection.UpdateOne(ctx, bson.M{"_id": objID}, update)
	return err
}

func (r *PekerjaanRepository) SoftDeleteBulk(ctx context.Context) error {
	update := bson.M{"$set": bson.M{"isdeleted": true, "updated_at": time.Now()}}
	_, err := r.collection.UpdateMany(ctx, bson.M{"isdeleted": false}, update)
	return err
}

func (r *PekerjaanRepository) Restore(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	update := bson.M{"$set": bson.M{"isdeleted": false, "updated_at": time.Now()}}
	_, err = r.collection.UpdateOne(ctx, bson.M{"_id": objID}, update)
	return err
}

func (r *PekerjaanRepository) Trash(ctx context.Context, search, sortBy, order string, limit, offset int) ([]model.Pekerjaan, error) {
	filter := bson.M{"isdeleted": true}

	if search != "" {
		filter["$or"] = []bson.M{
			{"nama_perusahaan": bson.M{"$regex": search, "$options": "i"}},
			{"posisi_jabatan": bson.M{"$regex": search, "$options": "i"}},
			{"bidang_industri": bson.M{"$regex": search, "$options": "i"}},
			{"lokasi_kerja": bson.M{"$regex": search, "$options": "i"}},
		}
	}

	sortOrder := 1
	if order == "desc" {
		sortOrder = -1
	}

	opts := options.Find().
		SetSort(bson.D{{Key: sortBy, Value: sortOrder}}).
		SetLimit(int64(limit)).
		SetSkip(int64(offset))

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var list []model.Pekerjaan
	if err := cursor.All(ctx, &list); err != nil {
		return nil, err
	}
	return list, nil
}

func (r *PekerjaanRepository) CountTrash(ctx context.Context, search string) (int64, error) {
	filter := bson.M{"isdeleted": true}

	if search != "" {
		filter["$or"] = []bson.M{
			{"nama_perusahaan": bson.M{"$regex": search, "$options": "i"}},
			{"posisi_jabatan": bson.M{"$regex": search, "$options": "i"}},
			{"bidang_industri": bson.M{"$regex": search, "$options": "i"}},
			{"lokasi_kerja": bson.M{"$regex": search, "$options": "i"}},
		}
	}

	return r.collection.CountDocuments(ctx, filter)
}

func (r *PekerjaanRepository) HardDelete(ctx context.Context) error {
	_, err := r.collection.DeleteMany(ctx, bson.M{"isdeleted": true})
	return err
}
