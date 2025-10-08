package repository

import (
	"TM4/app/model"
	"database/sql"
    "fmt"
    "strings")

type AlumniRepository struct {
	DB *sql.DB
}

func NewAlumniRepository(db *sql.DB) *AlumniRepository {
	return &AlumniRepository{DB: db}
}

func (r *AlumniRepository) GetAlumni(search, sortBy, order string, limit, offset int) ([]model.Alumni, error) {
	query := `SELECT id, nim, nama, jurusan, angkatan, 
				tahun_lulus, email, no_telepon, alamat, role,
				created_at, updated_at
			  FROM alumni
			  WHERE nama ILIKE $1 OR nim ILIKE $1
			  ORDER BY ` + sortBy + ` ` + order + `
			  LIMIT $2 OFFSET $3`

	rows, err := r.DB.Query(query, "%"+search+"%", limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var alumniList []model.Alumni
	for rows.Next() {
		var a model.Alumni
		if err := rows.Scan(&a.ID, &a.NIM, &a.Nama, &a.Jurusan, &a.Angkatan,
			&a.TahunLulus, &a.Email, &a.NoTelepon,
			&a.Alamat, &a.Role, &a.CreatedAt, &a.UpdatedAt); err != nil {
			return nil, err
		}
		alumniList = append(alumniList, a)
	}

	return alumniList, nil
}

func (r *AlumniRepository) Count(search string) (int, error) {
	var total int
	countQuery := `
		SELECT COUNT(*) FROM alumni 
		WHERE nama ILIKE $1 OR jurusan ILIKE $1 OR email ILIKE $1 OR nim ILIKE $1 OR angkatan::text ILIKE $1
	`
	err := r.DB.QueryRow(countQuery, "%"+search+"%").Scan(&total)
	if err != nil && err != sql.ErrNoRows {
		return 0, err
	}
		return total, nil
}



func (r *AlumniRepository) GetByID(id int) (*model.Alumni, error) {
	var a model.Alumni
	err := r.DB.QueryRow(`SELECT id, nim, nama, jurusan, angkatan, tahun_lulus, email, no_telepon, alamat, created_at, updated_at FROM alumni WHERE id=$1`, id).
		Scan(&a.ID, &a.NIM, &a.Nama, &a.Jurusan, &a.Angkatan,
			&a.TahunLulus, &a.Email, &a.NoTelepon, &a.Alamat,
			&a.CreatedAt, &a.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *AlumniRepository) Create(req model.CreateAlumni) (int, error) {
	var id int
	err := r.DB.QueryRow(`
		INSERT INTO alumni (nim, nama, jurusan, angkatan, tahun_lulus, email, no_telepon, alamat, role, password)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10) RETURNING id`,
		req.NIM, req.Nama, req.Jurusan, req.Angkatan, req.TahunLulus,
		req.Email, req.NoTelepon, req.Alamat, req.Role, req.Password).Scan(&id)
	return id, err
}

func (r *AlumniRepository) Update(id int, req model.UpdateAlumni) error {
	_, err := r.DB.Exec(`
		UPDATE alumni 
		SET nama=$1, jurusan=$2, angkatan=$3, tahun_lulus=$4, email=$5, no_telepon=$6, alamat=$7, role=$8, password=$9, updated_at=NOW() 
		WHERE id=$10`,
		req.Nama, req.Jurusan, req.Angkatan, req.TahunLulus, req.Email,
		req.NoTelepon, req.Alamat, req.Role, req.Password, id)
	return err
}

func (r *AlumniRepository) Delete(id int) error {
	_, err := r.DB.Exec(`DELETE FROM alumni WHERE id=$1`, id)
	return err
}

func (r *AlumniRepository) GetByEmail(email string) (*model.Alumni, error) {
	var a model.Alumni
	// pastikan email lowercase
	email = strings.ToLower(email)
	
	err := r.DB.QueryRow(`
		SELECT id, nim, nama, jurusan, angkatan, tahun_lulus, email, no_telepon, alamat, role, password, created_at, updated_at 
		FROM alumni 
		WHERE LOWER(email) = $1`, email).
		Scan(&a.ID, &a.NIM, &a.Nama, &a.Jurusan, &a.Angkatan,
			&a.TahunLulus, &a.Email, &a.NoTelepon, &a.Alamat, &a.Role,
			&a.Password, &a.CreatedAt, &a.UpdatedAt)
	if err != nil {
		// log error sementara untuk debug
		fmt.Println("DEBUG GetByEmail:", email, "error:", err)
		return nil, err
	}
	return &a, nil
}




