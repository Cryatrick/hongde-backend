package main

import (
	"log"
	"os"

	"hongde_backend/internal/config"
	"hongde_backend/internal/database"
	// "hongde_backend/internal/service"
	"hongde_backend/internal/middleware"
	"hongde_backend/internal/route"
	// "hongde_backend/internal/controller"
)

func init() {
	// First Initialisasi variable variable configuration
	config.InitEnvronment()
	config.InitDatabaseVars()
	config.InitEncryptionVars()

	// Second Initialisasi middleware yang perlu dinyalakan
	middleware.InitValidator()

	// Third Nyalakan koneksi ke database yang diperlukan
	database.OpenMain()
	database.OpenMongo()

	// Fourth Initialisasi variable variable yang bersifat universal dari database agar hanya perlu dijalankan sekali saja

	// Last Initialisasi variable variable khusus dari controller
}

func main() {
	router := route.SetupRouter()
	// if os.Getenv("ENVIRONMENT") == "development" {
	// }
		err := router.Run(":"+os.Getenv("PORT"))
		if err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	// err := router.RunTLS(":"+os.Getenv("PORT"),"server.pem","server.key")
	// if err != nil {
	// 	log.Fatalf("Failed to start server: %v", err)
	// }
}