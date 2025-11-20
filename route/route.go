package route

import (
	"TM4/app/repository"
	"TM4/app/service"
	"TM4/middleware"

	"github.com/gofiber/fiber/v2"
	_ "TM4/docs"
	"go.mongodb.org/mongo-driver/mongo"
	"github.com/gofiber/swagger"

)


func SetupRoutes(app *fiber.App, db *mongo.Database) {
	api := app.Group("/TM4")
	app.Get("/swagger/*", swagger.HandlerDefault) // default swagger UI


	// === Auth ===
	alumniRepo := repository.NewAlumniRepository(db)
	authService := service.NewAuthService(alumniRepo)
	api.Post("/login", authService.Login)

	// === Alumni ===
	alumniService := service.NewAlumniService(alumniRepo)
	alumni := api.Group("/alumni")
	alumni.Get("/", middleware.AuthRequired(), alumniService.GetAlumni)
	alumni.Get("/:id", middleware.AuthRequired(), alumniService.GetByID)

	// === Admin ===
	alumni.Post("/", middleware.AuthRequired(), middleware.AdminOnly(), alumniService.Create)
	alumni.Put("/:id", middleware.AuthRequired(), middleware.AdminOnly(), alumniService.Update)
	alumni.Delete("/:id", middleware.AuthRequired(), middleware.AdminOnly(), alumniService.Delete)


	// === Pekerjaan ===
	pekerjaanRepo := repository.NewPekerjaanRepository(db)
	pekerjaanService := service.NewPekerjaanService(pekerjaanRepo)
	pekerjaan := api.Group("/pekerjaan")
	pekerjaan.Get("/", middleware.AuthRequired(), middleware.AdminOnly(), pekerjaanService.GetAll)
	pekerjaan.Get("/trash", middleware.AuthRequired(),pekerjaanService.Trash)
	pekerjaan.Get("/:id", middleware.AuthRequired(), pekerjaanService.GetByID)
	pekerjaan.Put("/softdelete", middleware.AuthRequired(),pekerjaanService.SoftDelete)
	pekerjaan.Put("/restore", middleware.AuthRequired(),pekerjaanService.Restore)
	pekerjaan.Delete("/harddelete", middleware.AuthRequired(),pekerjaanService.HardDelete)
	
	// === Admin ===
	pekerjaan.Get("/alumni/:alumni_id", middleware.AuthRequired(), middleware.AdminOnly(), pekerjaanService.GetByAlumniID)
	pekerjaan.Post("/", middleware.AuthRequired(), middleware.AdminOnly(), pekerjaanService.Create)
	pekerjaan.Put("/:id", middleware.AuthRequired(), middleware.AdminOnly(), pekerjaanService.Update)
	pekerjaan.Delete("/:id", middleware.AuthRequired(), middleware.AdminOnly(), pekerjaanService.Delete)

	
	// route file
	fileRepo := repository.NewFileRepository(db)
	fileService := service.NewFileService(fileRepo, "./uploads")
	file := api.Group("/file")

	// user bisa lihat file
	file.Get("/", middleware.AuthRequired(), fileService.GetAllFiles)
	file.Get("/:id", middleware.AuthRequired(), fileService.GetFileByID)
	file.Get("/alumni/:alumni_id", middleware.AuthRequired(), fileService.GetFilesByAlumniID)

	// admin only upload/update/delete
	file.Post("/", middleware.AuthRequired(), fileService.UploadFile)              
	file.Post("/:alumni_id", middleware.AuthRequired(), middleware.AdminOnly(), fileService.UploadFileAdmin) // admin upload ke alumni lain
	file.Put("/:id", middleware.AuthRequired(), middleware.AdminOnly(), fileService.UpdateFile)
	file.Delete("/:id", middleware.AuthRequired(), middleware.AdminOnly(), fileService.DeleteFile)

}
