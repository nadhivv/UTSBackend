package repository

import (
	"TM4/app/model"
	"database/sql"
)

type PekerjaanRepository struct {
	DB *sql.DB
}

func NewPekerjaanRepository(db *sql.DB) *PekerjaanRepository {
	return &PekerjaanRepository{DB: db}
}

func (r *PekerjaanRepository) GetAll(search, sortBy, order string, limit, offset int) ([]model.Pekerjaan, error) {
	query := `
		SELECT id, alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri, 
		       lokasi_kerja, gaji_range, tanggal_mulai_kerja, tanggal_selesai_kerja, 
		       status_pekerjaan, deskripsi_pekerjaan, created_at, updated_at 
		FROM pekerjaan_alumni
		WHERE isdeleted = false 
		  AND (nama_perusahaan ILIKE $1 OR posisi_jabatan ILIKE $1)
		ORDER BY ` + sortBy + ` ` + order + `
		LIMIT $2 OFFSET $3`

	rows, err := r.DB.Query(query, "%"+search+"%", limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pekerjaanList []model.Pekerjaan
	for rows.Next() {
		var p model.Pekerjaan
		if err := rows.Scan(
			&p.ID, &p.AlumniID, &p.NamaPerusahaan, &p.PosisiJabatan, &p.BidangIndustri,
			&p.LokasiKerja, &p.GajiRange, &p.TanggalMulaiKerja, &p.TanggalSelesaiKerja,
			&p.StatusPekerjaan, &p.DeskripsiPekerjaan, &p.CreatedAt, &p.UpdatedAt,
		); err != nil {
			return nil, err
		}
		pekerjaanList = append(pekerjaanList, p)
	}
	return pekerjaanList, nil
}

func (r *PekerjaanRepository) Count(search string) (int, error) {
	var total int
	err := r.DB.QueryRow(`
		SELECT COUNT(*) 
		FROM pekerjaan_alumni 
		WHERE isdeleted = false 
		  AND (nama_perusahaan ILIKE $1 OR posisi_jabatan ILIKE $1)`,
		"%"+search+"%",
	).Scan(&total)

	if err != nil {
		return 0, err
	}
	return total, nil
}

func (r *PekerjaanRepository) GetByID(id int) (*model.Pekerjaan, error) {
	var p model.Pekerjaan
	err := r.DB.QueryRow(`
		SELECT id, alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri, 
		       lokasi_kerja, gaji_range, tanggal_mulai_kerja, tanggal_selesai_kerja, 
		       status_pekerjaan, deskripsi_pekerjaan, created_at, updated_at 
		FROM pekerjaan_alumni 
		WHERE isdeleted = false AND id = $1`,
		id,
	).Scan(
		&p.ID, &p.AlumniID, &p.NamaPerusahaan, &p.PosisiJabatan, &p.BidangIndustri,
		&p.LokasiKerja, &p.GajiRange, &p.TanggalMulaiKerja, &p.TanggalSelesaiKerja,
		&p.StatusPekerjaan, &p.DeskripsiPekerjaan, &p.CreatedAt, &p.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *PekerjaanRepository) GetByAlumniID(alumniID int) ([]model.Pekerjaan, error) {
	query := `
		SELECT id, alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri, 
		       lokasi_kerja, gaji_range, tanggal_mulai_kerja, tanggal_selesai_kerja, 
		       status_pekerjaan, deskripsi_pekerjaan, created_at, updated_at 
		FROM pekerjaan_alumni
		WHERE alumni_id = $1 
		  AND isdeleted = false
		ORDER BY created_at DESC`

	rows, err := r.DB.Query(query, alumniID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.Pekerjaan
	for rows.Next() {
		var p model.Pekerjaan
		if err := rows.Scan(
			&p.ID, &p.AlumniID, &p.NamaPerusahaan, &p.PosisiJabatan, &p.BidangIndustri,
			&p.LokasiKerja, &p.GajiRange, &p.TanggalMulaiKerja, &p.TanggalSelesaiKerja,
			&p.StatusPekerjaan, &p.DeskripsiPekerjaan, &p.CreatedAt, &p.UpdatedAt,
		); err != nil {
			return nil, err
		}
		list = append(list, p)
	}
	return list, nil
}

func (r *PekerjaanRepository) Create(req model.CreatePekerjaan) (int, error) {
	var id int
	err := r.DB.QueryRow(`
		INSERT INTO pekerjaan_alumni (
			alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri, 
			lokasi_kerja, gaji_range, tanggal_mulai_kerja, tanggal_selesai_kerja, 
			status_pekerjaan, deskripsi_pekerjaan
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10) RETURNING id`,
		req.AlumniID, req.NamaPerusahaan, req.PosisiJabatan, req.BidangIndustri,
		req.LokasiKerja, req.GajiRange, req.TanggalMulaiKerja, req.TanggalSelesaiKerja,
		req.StatusPekerjaan, req.DeskripsiPekerjaan).Scan(&id)
	return id, err
}

func (r *PekerjaanRepository) Update(id int, req model.UpdatePekerjaan) error {
	_, err := r.DB.Exec(`
		UPDATE pekerjaan_alumni SET 
			nama_perusahaan=$1, 
			posisi_jabatan=$2, 
			bidang_industri=$3, 
			lokasi_kerja=$4, 
			gaji_range=$5, 
			tanggal_mulai_kerja=$6, 
			tanggal_selesai_kerja=$7, 
			status_pekerjaan=$8, 
			deskripsi_pekerjaan=$9, 
			updated_at=NOW() 
		WHERE id=$10`,
		req.NamaPerusahaan, req.PosisiJabatan, req.BidangIndustri,
		req.LokasiKerja, req.GajiRange, req.TanggalMulaiKerja, req.TanggalSelesaiKerja,
		req.StatusPekerjaan, req.DeskripsiPekerjaan, id)
	return err
}

func (r *PekerjaanRepository) Delete(id int) error {
	_, err := r.DB.Exec(`DELETE FROM pekerjaan_alumni WHERE id=$1`, id)
	return err
}

func (r *PekerjaanRepository) SoftDelete(id int, req model.UpdatePekerjaan) error {
	_, err := r.DB.Exec(`
		UPDATE pekerjaan_alumni SET 
			isdeleted=true,
			updated_at=NOW() 
		WHERE id=$1`,
		id)
	return err
}

func (r *PekerjaanRepository) SoftDeleteBulk() error {
	query := `UPDATE pekerjaan_alumni SET isdeleted=true, updated_at=NOW() 
	WHERE isdeleted=false`
	_, err := r.DB.Exec(query)
	return err
}

func (r *PekerjaanRepository) Trash() ([]model.Pekerjaan, error) {
	rows, err := r.DB.Query(`
	SELECT alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri,
	       lokasi_kerja, gaji_range, tanggal_mulai_kerja, tanggal_selesai_kerja,
	       status_pekerjaan, deskripsi_pekerjaan, created_at, updated_at
	FROM pekerjaan_alumni
	WHERE isdeleted = true
	ORDER BY updated_at DESC
`)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.Pekerjaan
	for rows.Next() {
		var p model.Pekerjaan
		if err := rows.Scan(&p.AlumniID, &p.NamaPerusahaan, &p.PosisiJabatan, &p.BidangIndustri, &p.LokasiKerja, &p.GajiRange, &p.TanggalMulaiKerja, &p.TanggalSelesaiKerja, &p.StatusPekerjaan, &p.DeskripsiPekerjaan, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		list = append(list, p)
	}

	return list, nil
}

func (r *PekerjaanRepository) Restore() error {
	_, err := r.DB.Exec(`
		UPDATE pekerjaan_alumni SET isdeleted=false, updated_at=NOW() 
		WHERE isdeleted=true`)
	return err
}

func (r *PekerjaanRepository) HardDelete() error {
	_, err := r.DB.Exec(`
		DELETE FROM pekerjaan_alumni 
		WHERE isdeleted = true`)
	return err
}
