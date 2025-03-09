package main

import (
	"log"

	"github.com/dheerajkumardk/blockchain_explorer_backend/database"
	"github.com/dheerajkumardk/blockchain_explorer_backend/listener"
	"github.com/dheerajkumardk/blockchain_explorer_backend/routes"
	"github.com/gofiber/fiber/v2"
)

func main() {
	// create app
	app := fiber.New()

	// init database
	database.InitDatabase()

	// set up routes
	routes.SetupRoutes(app)
	
	// start block listener
	go listener.SubscribeBlocks()

	// start the server
	log.Fatal(app.Listen(":3000"))

	sqlDB, err := database.BlockDB.DB()
	if err != nil {
		log.Fatalf("Failed to get underlying database connection: %v", err)
	}
	defer sqlDB.Close()
}
