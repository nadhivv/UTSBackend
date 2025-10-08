package database

import	(
	"database/sql"
	"log"
	"os"
	"fmt"
	_ "github.com/lib/pq"

)

var	DB	*sql.DB

func ConnectDB()  {
    dsn := fmt.Sprintf(
        "host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
        os.Getenv("DB_HOST"),
        os.Getenv("DB_USER"),
        os.Getenv("DB_PASSWORD"),
        os.Getenv("DB_NAME"),
        os.Getenv("DB_PORT"),
    )

    var err error
    DB, err = sql.Open("postgres", dsn)
    if err != nil {
        log.Fatal("Error opening DB:", err)
    }

    if err = DB.Ping(); err != nil {
        log.Fatal("Error connecting to DB:", err)
    }

    var version string
    if err := DB.QueryRow("SELECT version()").Scan(&version); err != nil {
        log.Fatal("Error checking DB version:", err)
    }

    log.Println("Connected to PostgreSQL:", version)
}