package config

import	(
	"log"
	"os"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

)

func LoadEnv()	{
	err	:= godotenv.Load()
	if err	!=	nil	{
		log.Fatal("Error loading.env file")
	}
}

func GetEnv(key string) string {
	return os.Getenv(key)
}