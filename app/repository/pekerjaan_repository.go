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

// type PekerjaanRepository struct {
// 	DB *sql.DB
// }

// func NewPekerjaanRepository(db *sql.DB) *PekerjaanRepository {
// 	return &PekerjaanRepository{DB: db}
// }

// func (r *PekerjaanRepository) GetAll(search, sortBy, order string, limit, offset int) ([]model.Pekerjaan, error) {
// 	query := `
// 		SELECT id, alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri, 
// 		       lokasi_kerja, gaji_range, tanggal_mulai_kerja, tanggal_selesai_kerja, 
// 		       status_pekerjaan, deskripsi_pekerjaan, created_at, updated_at 
// 		FROM pekerjaan_alumni
// 		WHERE isdeleted = false 
// 		  AND (nama_perusahaan ILIKE $1 OR posisi_jabatan ILIKE $1)
// 		ORDER BY ` + sortBy + ` ` + order + `
// 		LIMIT $2 OFFSET $3`

// 	rows, err := r.DB.Query(query, "%"+search+"%", limit, offset)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	var pekerjaanList []model.Pekerjaan
// 	for rows.Next() {
// 		var p model.Pekerjaan
// 		if err := rows.Scan(
// 			&p.ID, &p.AlumniID, &p.NamaPerusahaan, &p.PosisiJabatan, &p.BidangIndustri,
// 			&p.LokasiKerja, &p.GajiRange, &p.TanggalMulaiKerja, &p.TanggalSelesaiKerja,
// 			&p.StatusPekerjaan, &p.DeskripsiPekerjaan, &p.CreatedAt, &p.UpdatedAt,
// 		); err != nil {
// 			return nil, err
// 		}
// 		pekerjaanList = append(pekerjaanList, p)
// 	}
// 	return pekerjaanList, nil
// }

// func (r *PekerjaanRepository) Count(search string) (int, error) {
// 	var total int
// 	err := r.DB.QueryRow(`
// 		SELECT COUNT(*) 
// 		FROM pekerjaan_alumni 
// 		WHERE isdeleted = false 
// 		  AND (nama_perusahaan ILIKE $1 OR posisi_jabatan ILIKE $1)`,
// 		"%"+search+"%",
// 	).Scan(&total)

// 	if err != nil {
// 		return 0, err
// 	}
// 	return total, nil
// }

// func (r *PekerjaanRepository) GetByID(id int) (*model.Pekerjaan, error) {
// 	var p model.Pekerjaan
// 	err := r.DB.QueryRow(`
// 		SELECT id, alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri, 
// 		       lokasi_kerja, gaji_range, tanggal_mulai_kerja, tanggal_selesai_kerja, 
// 		       status_pekerjaan, deskripsi_pekerjaan, created_at, updated_at 
// 		FROM pekerjaan_alumni 
// 		WHERE isdeleted = false AND id = $1`,
// 		id,
// 	).Scan(
// 		&p.ID, &p.AlumniID, &p.NamaPerusahaan, &p.PosisiJabatan, &p.BidangIndustri,
// 		&p.LokasiKerja, &p.GajiRange, &p.TanggalMulaiKerja, &p.TanggalSelesaiKerja,
// 		&p.StatusPekerjaan, &p.DeskripsiPekerjaan, &p.CreatedAt, &p.UpdatedAt,
// 	)

// 	if err != nil {
// 		return nil, err
// 	}
// 	return &p, nil
// }

// func (r *PekerjaanRepository) GetByAlumniID(alumniID int) ([]model.Pekerjaan, error) {
// 	query := `
// 		SELECT id, alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri, 
// 		       lokasi_kerja, gaji_range, tanggal_mulai_kerja, tanggal_selesai_kerja, 
// 		       status_pekerjaan, deskripsi_pekerjaan, created_at, updated_at 
// 		FROM pekerjaan_alumni
// 		WHERE alumni_id = $1 
// 		  AND isdeleted = false
// 		ORDER BY created_at DESC`

// 	rows, err := r.DB.Query(query, alumniID)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	var list []model.Pekerjaan
// 	for rows.Next() {
// 		var p model.Pekerjaan
// 		if err := rows.Scan(
// 			&p.ID, &p.AlumniID, &p.NamaPerusahaan, &p.PosisiJabatan, &p.BidangIndustri,
// 			&p.LokasiKerja, &p.GajiRange, &p.TanggalMulaiKerja, &p.TanggalSelesaiKerja,
// 			&p.StatusPekerjaan, &p.DeskripsiPekerjaan, &p.CreatedAt, &p.UpdatedAt,
// 		); err != nil {
// 			return nil, err
// 		}
// 		list = append(list, p)
// 	}
// 	return list, nil
// }

// func (r *PekerjaanRepository) Create(req model.CreatePekerjaan) (int, error) {
// 	var id int
// 	err := r.DB.QueryRow(`
// 		INSERT INTO pekerjaan_alumni (
// 			alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri, 
// 			lokasi_kerja, gaji_range, tanggal_mulai_kerja, tanggal_selesai_kerja, 
// 			status_pekerjaan, deskripsi_pekerjaan
// 		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10) RETURNING id`,
// 		req.AlumniID, req.NamaPerusahaan, req.PosisiJabatan, req.BidangIndustri,
// 		req.LokasiKerja, req.GajiRange, req.TanggalMulaiKerja, req.TanggalSelesaiKerja,
// 		req.StatusPekerjaan, req.DeskripsiPekerjaan).Scan(&id)
// 	return id, err
// }

// func (r *PekerjaanRepository) Update(id int, req model.UpdatePekerjaan) error {
// 	_, err := r.DB.Exec(`
// 		UPDATE pekerjaan_alumni SET 
// 			nama_perusahaan=$1, 
// 			posisi_jabatan=$2, 
// 			bidang_industri=$3, 
// 			lokasi_kerja=$4, 
// 			gaji_range=$5, 
// 			tanggal_mulai_kerja=$6, 
// 			tanggal_selesai_kerja=$7, 
// 			status_pekerjaan=$8, 
// 			deskripsi_pekerjaan=$9, 
// 			updated_at=NOW() 
// 		WHERE id=$10`,
// 		req.NamaPerusahaan, req.PosisiJabatan, req.BidangIndustri,
// 		req.LokasiKerja, req.GajiRange, req.TanggalMulaiKerja, req.TanggalSelesaiKerja,
// 		req.StatusPekerjaan, req.DeskripsiPekerjaan, id)
// 	return err
// }

// func (r *PekerjaanRepository) Delete(id int) error {
// 	_, err := r.DB.Exec(`DELETE FROM pekerjaan_alumni WHERE id=$1`, id)
// 	return err
// }

// func (r *PekerjaanRepository) SoftDelete(id int, req model.UpdatePekerjaan) error {
// 	_, err := r.DB.Exec(`
// 		UPDATE pekerjaan_alumni SET 
// 			isdeleted=true,
// 			updated_at=NOW() 
// 		WHERE id=$1`,
// 		id)
// 	return err
// }

// func (r *PekerjaanRepository) SoftDeleteBulk() error {
// 	query := `UPDATE pekerjaan_alumni SET isdeleted=true, updated_at=NOW() 
// 	WHERE isdeleted=false`
// 	_, err := r.DB.Exec(query)
// 	return err
// }

// func (r *PekerjaanRepository) Trash() ([]model.Pekerjaan, error) {
// 	rows, err := r.DB.Query(`
// 	SELECT alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri,
// 	       lokasi_kerja, gaji_range, tanggal_mulai_kerja, tanggal_selesai_kerja,
// 	       status_pekerjaan, deskripsi_pekerjaan, created_at, updated_at
// 	FROM pekerjaan_alumni
// 	WHERE isdeleted = true
// 	ORDER BY updated_at DESC
// `)

// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	var list []model.Pekerjaan
// 	for rows.Next() {
// 		var p model.Pekerjaan
// 		if err := rows.Scan(&p.AlumniID, &p.NamaPerusahaan, &p.PosisiJabatan, &p.BidangIndustri, &p.LokasiKerja, &p.GajiRange, &p.TanggalMulaiKerja, &p.TanggalSelesaiKerja, &p.StatusPekerjaan, &p.DeskripsiPekerjaan, &p.CreatedAt, &p.UpdatedAt); err != nil {
// 			return nil, err
// 		}
// 		list = append(list, p)
// 	}

// 	return list, nil
// }

// func (r *PekerjaanRepository) Restore() error {
// 	_, err := r.DB.Exec(`
// 		UPDATE pekerjaan_alumni SET isdeleted=false, updated_at=NOW() 
// 		WHERE isdeleted=true`)
// 	return err
// }

// func (r *PekerjaanRepository) HardDelete() error {
// 	_, err := r.DB.Exec(`
// 		DELETE FROM pekerjaan_alumni 
// 		WHERE isdeleted = true`)
// 	return err
// }

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
		collection: db.Collection("pekerjaan_alumni"), //ini namanya collection sama kyk tabel disamain aja
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

	// Convert alumni_id string ke ObjectID
	alumniID, err := primitive.ObjectIDFromHex(req.AlumniID)
	if err != nil {
		return nil, fmt.Errorf("invalid alumni_id: %v", err)
	}

	doc := bson.M{
		"alumni_id":            alumniID,
		"nama_perusahaan":      req.NamaPerusahaan,
		"posisi_jabatan":       req.PosisiJabatan,
		"bidang_industri":      req.BidangIndustri,
		"lokasi_kerja":         req.LokasiKerja,
		"gaji_range":           req.GajiRange,
		"tanggal_mulai_kerja":  req.TanggalMulaiKerja,
		"tanggal_selesai_kerja": req.TanggalSelesaiKerja,
		"status_pekerjaan":     req.StatusPekerjaan,
		"deskripsi_pekerjaan":  req.DeskripsiPekerjaan,
		"created_at":           now,
		"updated_at":           now,
		"isdeleted":            false,
	}

	return r.collection.InsertOne(ctx, doc)
}


func (r *PekerjaanRepository) Update(ctx context.Context, id string, req *model.UpdatePekerjaan) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	update := bson.M{"$set": bson.M{
		"nama_perusahaan":      req.NamaPerusahaan,
		"posisi_jabatan":       req.PosisiJabatan,
		"bidang_industri":      req.BidangIndustri,
		"lokasi_kerja":         req.LokasiKerja,
		"gaji_range":           req.GajiRange,
		"tanggal_mulai_kerja":  req.TanggalMulaiKerja,
		"tanggal_selesai_kerja": req.TanggalSelesaiKerja,
		"status_pekerjaan":     req.StatusPekerjaan,
		"deskripsi_pekerjaan":  req.DeskripsiPekerjaan,
		"updated_at":           time.Now(),
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