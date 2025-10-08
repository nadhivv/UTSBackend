package model

import "time"

type Alumni struct {
	ID         int       `json:"id"`
	NIM        string    `json:"nim"`
	Nama       string    `json:"nama"`
	Password   string    `json:"-"`
	Jurusan    string    `json:"jurusan"`
	Angkatan   int       `json:"angkatan"`
	TahunLulus int       `json:"tahun_lulus"`
	Email      string    `json:"email"`
	NoTelepon  string    `json:"no_telepon"`
	Alamat     string    `json:"alamat"`
	Role       string    `json:"role"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	IsDeleted  bool      `json:"-"`
}

type CreateAlumni struct {
	NIM        string `json:"nim"`
	Nama       string `json:"nama"`
	Jurusan    string `json:"jurusan"`
	Angkatan   int    `json:"angkatan"`
	TahunLulus int    `json:"tahun_lulus"`
	Email      string `json:"email"`
	NoTelepon  string `json:"no_telepon"`
	Alamat     string `json:"alamat"`
	Role       string `json:"role"`
	Password   string `json:"password"`
}

type UpdateAlumni struct {
	Nama       string `json:"nama"`
	Jurusan    string `json:"jurusan"`
	Angkatan   int    `json:"angkatan"`
	TahunLulus int    `json:"tahun_lulus"`
	Email      string `json:"email"`
	NoTelepon  string `json:"no_telepon"`
	Alamat     string `json:"alamat"`
	Role       string `json:"role"`
	Password   string `json:"password"`
}
