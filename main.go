package	main
import	(
	"fmt"
	"log"
	"TM4/config"
	"TM4/database"
	"TM4/route"
	"github.com/gofiber/fiber/v2"

)

func main()	{
config.LoadEnv()
  database.ConnectDB()
  defer database.DB.Close()
  log.Println("db connected")

	app := fiber.New()
	route.SetupRoutes(app, database.DB)

	for _, r := range app.GetRoutes() {
    fmt.Println(r.Method, r.Path)
}

	port := config.GetEnv("APP_PORT")
	log.Fatal(app.Listen(fmt.Sprintf(":%s", port)))

}

