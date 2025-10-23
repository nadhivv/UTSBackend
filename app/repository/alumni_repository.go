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

// type AlumniRepository struct {
// 	DB *sql.DB
// }

// func NewAlumniRepository(db *sql.DB) *AlumniRepository {
// 	return &AlumniRepository{DB: db}
// }

// func (r *AlumniRepository) GetAlumni(search, sortBy, order string, limit, offset int) ([]model.Alumni, error) {
// 	query := `SELECT id, nim, nama, jurusan, angkatan, 
// 				tahun_lulus, email, no_telepon, alamat, role,
// 				created_at, updated_at
// 			  FROM alumni
// 			  WHERE nama ILIKE $1 OR nim ILIKE $1
// 			  ORDER BY ` + sortBy + ` ` + order + `
// 			  LIMIT $2 OFFSET $3`

// 	rows, err := r.DB.Query(query, "%"+search+"%", limit, offset)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	var alumniList []model.Alumni
// 	for rows.Next() {
// 		var a model.Alumni
// 		if err := rows.Scan(&a.ID, &a.NIM, &a.Nama, &a.Jurusan, &a.Angkatan,
// 			&a.TahunLulus, &a.Email, &a.NoTelepon,
// 			&a.Alamat, &a.Role, &a.CreatedAt, &a.UpdatedAt); err != nil {
// 			return nil, err
// 		}
// 		alumniList = append(alumniList, a)
// 	}

// 	return alumniList, nil
// }

// func (r *AlumniRepository) Count(search string) (int, error) {
// 	var total int
// 	countQuery := `
// 		SELECT COUNT(*) FROM alumni 
// 		WHERE nama ILIKE $1 OR jurusan ILIKE $1 OR email ILIKE $1 OR nim ILIKE $1 OR angkatan::text ILIKE $1
// 	`
// 	err := r.DB.QueryRow(countQuery, "%"+search+"%").Scan(&total)
// 	if err != nil && err != sql.ErrNoRows {
// 		return 0, err
// 	}
// 		return total, nil
// }



// func (r *AlumniRepository) GetByID(id int) (*model.Alumni, error) {
// 	var a model.Alumni
// 	err := r.DB.QueryRow(`SELECT id, nim, nama, jurusan, angkatan, tahun_lulus, email, no_telepon, alamat, created_at, updated_at FROM alumni WHERE id=$1`, id).
// 		Scan(&a.ID, &a.NIM, &a.Nama, &a.Jurusan, &a.Angkatan,
// 			&a.TahunLulus, &a.Email, &a.NoTelepon, &a.Alamat,
// 			&a.CreatedAt, &a.UpdatedAt)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &a, nil
// }

// func (r *AlumniRepository) Create(req model.CreateAlumni) (int, error) {
// 	var id int
// 	err := r.DB.QueryRow(`
// 		INSERT INTO alumni (nim, nama, jurusan, angkatan, tahun_lulus, email, no_telepon, alamat, role, password)
// 		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10) RETURNING id`,
// 		req.NIM, req.Nama, req.Jurusan, req.Angkatan, req.TahunLulus,
// 		req.Email, req.NoTelepon, req.Alamat, req.Role, req.Password).Scan(&id)
// 	return id, err
// }

// func (r *AlumniRepository) Update(id int, req model.UpdateAlumni) error {
// 	_, err := r.DB.Exec(`
// 		UPDATE alumni 
// 		SET nama=$1, jurusan=$2, angkatan=$3, tahun_lulus=$4, email=$5, no_telepon=$6, alamat=$7, role=$8, password=$9, updated_at=NOW() 
// 		WHERE id=$10`,
// 		req.Nama, req.Jurusan, req.Angkatan, req.TahunLulus, req.Email,
// 		req.NoTelepon, req.Alamat, req.Role, req.Password, id)
// 	return err
// }

// func (r *AlumniRepository) Delete(id int) error {
// 	_, err := r.DB.Exec(`DELETE FROM alumni WHERE id=$1`, id)
// 	return err
// }

// func (r *AlumniRepository) GetByEmail(email string) (*model.Alumni, error) {
// 	var a model.Alumni
// 	// pastikan email lowercase
// 	email = strings.ToLower(email)
	
// 	err := r.DB.QueryRow(`
// 		SELECT id, nim, nama, jurusan, angkatan, tahun_lulus, email, no_telepon, alamat, role, password, created_at, updated_at 
// 		FROM alumni 
// 		WHERE LOWER(email) = $1`, email).
// 		Scan(&a.ID, &a.NIM, &a.Nama, &a.Jurusan, &a.Angkatan,
// 			&a.TahunLulus, &a.Email, &a.NoTelepon, &a.Alamat, &a.Role,
// 			&a.Password, &a.CreatedAt, &a.UpdatedAt)
// 	if err != nil {
// 		// log error sementara untuk debug
// 		fmt.Println("DEBUG GetByEmail:", email, "error:", err)
// 		return nil, err
// 	}
// 	return &a, nil
// }


type IAlumniRepository interface {
	GetAlumni(ctx context.Context, search string, limit, offset int, sortBy, order string) ([]model.Alumni, error)
	GetAlumniByID(ctx context.Context, id string) (*model.Alumni, error)
	GetByEmail(ctx context.Context, email string) (*model.Alumni, error)
	CreateAlumni(ctx context.Context, req *model.Alumni) (*model.Alumni, error)
	UpdateAlumni(ctx context.Context, id string, req *model.UpdateAlumni) error
	DeleteAlumni(ctx context.Context, id string) error
	Count(ctx context.Context, search string) (int, error)
}


type AlumniRepository struct {
	collection *mongo.Collection
}

func NewAlumniRepository(db *mongo.Database) IAlumniRepository {
	return &AlumniRepository{
		collection: db.Collection("alumni"),
	}
}

func (r *AlumniRepository) Count(ctx context.Context, search string) (int, error) {
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

	return int(count), nil
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
	filter := bson.M{"email": strings.ToLower(email)}

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
