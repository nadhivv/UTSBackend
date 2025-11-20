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
	"log"
	"context"
	"time"
	"TM4/route"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(context.TODO())

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	db := client.Database("alumni") 

	app := fiber.New()

	route.SetupRoutes(app, db)

	log.Fatal(app.Listen(":3000"))

}