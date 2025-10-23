// package	main
// import	(
// 	"fmt"
// 	"log"
// 	"TM4/config"
// 	"TM4/database"
// 	"TM4/route"
// 	"github.com/gofiber/fiber/v2"

// )

// func main()	{
// config.LoadEnv()
//   database.ConnectDB()
//   defer database.DB.Close()
//   log.Println("db connected")

// 	app := fiber.New()
// 	route.SetupRoutes(app, database.DB)

// 	for _, r := range app.GetRoutes() {
//     fmt.Println(r.Method, r.Path)
// }

// 	port := config.GetEnv("APP_PORT")
// 	log.Fatal(app.Listen(fmt.Sprintf(":%s", port)))

// }

package main

import (
	"fmt"
	"log"

	"TM4/config"
	"TM4/database"
	"TM4/route"

	"github.com/gofiber/fiber/v2"
)

func main() {
	// Load environment variable
	config.LoadEnv()

	// Koneksi ke MongoDB
	database.ConnectDB()
	defer database.DisconnectDB()
	log.Println("✅ MongoDB connected successfully")

	// Setup Fiber app
	app := fiber.New()
	route.SetupRoutes(app, database.DB)

	// Debug: tampilkan semua route
	for _, r := range app.GetRoutes() {
		fmt.Println(r.Method, r.Path)
	}

	// Jalankan server
	port := config.GetEnv("APP_PORT")
	if port == "" {
		port = "8080"
	}
	log.Fatal(app.Listen(fmt.Sprintf(":%s", port)))
}
