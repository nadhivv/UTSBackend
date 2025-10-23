package model

import 	(
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// type Alumni struct {
// 	ID        primitive.ObjectID	`json:"id"`
// 	NIM        string    `json:"nim"`
// 	Nama       string    `json:"nama"`
// 	Password   string    `json:"-"`
// 	Jurusan    string    `json:"jurusan"`
// 	Angkatan   int       `json:"angkatan"`
// 	TahunLulus int       `json:"tahun_lulus"`
// 	Email      string    `json:"email"`
// 	NoTelepon  string    `json:"no_telepon"`
// 	Alamat     string    `json:"alamat"`
// 	Role       string    `json:"role"`
// 	CreatedAt  time.Time `json:"created_at"`
// 	UpdatedAt  time.Time `json:"updated_at"`
// 	IsDeleted  bool      `json:"-"`
// }

// type CreateAlumni struct {
// 	NIM        string `json:"nim"`
// 	Nama       string `json:"nama"`
// 	Jurusan    string `json:"jurusan"`
// 	Angkatan   int    `json:"angkatan"`
// 	TahunLulus int    `json:"tahun_lulus"`
// 	Email      string `json:"email"`
// 	NoTelepon  string `json:"no_telepon"`
// 	Alamat     string `json:"alamat"`
// 	Role       string `json:"role"`
// 	Password   string `json:"password"`
// }

// type UpdateAlumni struct {
// 	Nama       string `json:"nama"`
// 	Jurusan    string `json:"jurusan"`
// 	Angkatan   int    `json:"angkatan"`
// 	TahunLulus int    `json:"tahun_lulus"`
// 	Email      string `json:"email"`
// 	NoTelepon  string `json:"no_telepon"`
// 	Alamat     string `json:"alamat"`
// 	Role       string `json:"role"`
// 	Password   string `json:"password"`
// }


type Alumni struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	NIM        string             `bson:"nim" json:"nim"`
	Nama       string             `bson:"nama" json:"nama"`
	Password   string             `bson:"password,omitempty" json:"-"`
	Jurusan    string             `bson:"jurusan" json:"jurusan"`
	Angkatan   int                `bson:"angkatan" json:"angkatan"`
	TahunLulus int                `bson:"tahun_lulus" json:"tahun_lulus"`
	Email      string             `bson:"email" json:"email"`
	NoTelepon  string             `bson:"no_telepon" json:"no_telepon"`
	Alamat     string             `bson:"alamat" json:"alamat"`
	Role       string             `bson:"role" json:"role"`
	CreatedAt  time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt  time.Time          `bson:"updated_at" json:"updated_at"`
	IsDeleted  bool               `bson:"isdeleted" json:"-"`
}


type CreateAlumni struct {
	NIM        string `bson:"nim" json:"nim"`
	Nama       string `bson:"nama" json:"nama"`
	Jurusan    string `bson:"jurusan" json:"jurusan"`
	Angkatan   int    `bson:"angkatan" json:"angkatan"`
	TahunLulus int    `bson:"tahun_lulus" json:"tahun_lulus"`
	Email      string `bson:"email" json:"email"`
	NoTelepon  string `bson:"no_telepon" json:"no_telepon"`
	Alamat     string `bson:"alamat" json:"alamat"`
	Role       string `bson:"role" json:"role"`
	Password   string `bson:"password,omitempty" json:"password"`
}


type UpdateAlumni struct {
	Nama       string `bson:"nama,omitempty" json:"nama"`
	Jurusan    string `bson:"jurusan,omitempty" json:"jurusan"`
	Angkatan   int    `bson:"angkatan,omitempty" json:"angkatan"`
	TahunLulus int    `bson:"tahun_lulus,omitempty" json:"tahun_lulus"`
	Email      string `bson:"email,omitempty" json:"email"`
	NoTelepon  string `bson:"no_telepon,omitempty" json:"no_telepon"`
	Alamat     string `bson:"alamat,omitempty" json:"alamat"`
	Role       string `bson:"role,omitempty" json:"role"`
	Password   string `bson:"password,omitempty" json:"password"`
}
