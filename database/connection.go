package database

import (
	// "database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
    
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
    "context"
    "time"
)

// yang atas ini buat postgres
// var DB *sql.DB

// func ConnectDB() {
// 	dsn := fmt.Sprintf(
// 		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
// 		os.Getenv("DB_HOST"),
// 		os.Getenv("DB_USER"),
// 		os.Getenv("DB_PASSWORD"),
// 		os.Getenv("DB_NAME"),
// 		os.Getenv("DB_PORT"),
// 	)

// 	var err error
// 	DB, err = sql.Open("postgres", dsn)
// 	if err != nil {
// 		log.Fatal("Error opening DB:", err)
// 	}

// 	if err = DB.Ping(); err != nil {
// 		log.Fatal("Error connecting to DB:", err)
// 	}

// 	var version string
// 	if err := DB.QueryRow("SELECT version()").Scan(&version); err != nil {
// 		log.Fatal("Error checking DB version:", err)
// 	}

// 	log.Println("Connected to PostgreSQL:", version)
// }

// ini untuk mongo
var DB *mongo.Database
var client *mongo.Client

func ConnectDB() {
	// Ambil URI dari environment variable
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
		log.Println("⚠️  Peringatan: MONGO_URI tidak disetel. Menggunakan default:", mongoURI)
	}

	// Ambil nama database dari environment variable
	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "alumni" // default
		log.Println("⚠️  Peringatan: DB_NAME tidak disetel. Menggunakan default:", dbName)
	}

	// Siapkan opsi koneksi
	clientOptions := options.Client().ApplyURI(mongoURI)

	// Koneksi ke MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("❌ Gagal terhubung ke MongoDB: %v", err)
	}

	// Cek koneksi
	if err = client.Ping(ctx, nil); err != nil {
		log.Fatalf("❌ Ping ke MongoDB gagal: %v", err)
	}

	// Inisialisasi database global
	DB = client.Database(dbName)
	fmt.Println("✅ Berhasil terhubung ke MongoDB Database:", dbName)

}

func DisconnectDB() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := client.Disconnect(ctx); err != nil {
		log.Printf("⚠️  Gagal memutus koneksi MongoDB: %v", err)
	} else {
		fmt.Println("🔌 Koneksi MongoDB ditutup dengan aman.")
	}
}