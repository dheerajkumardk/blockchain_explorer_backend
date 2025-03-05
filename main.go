package main

import (
	"log"

	"github.com/dheerajkumardk/blockchain_explorer_backend/database"
	"github.com/dheerajkumardk/blockchain_explorer_backend/listener"
	"github.com/gofiber/fiber/v2"
)

func main() {
	// create app
	app := fiber.New()

	// init database
	database.InitDatabase()

	// start block listener
	go listener.SubscribeBlocks()

	// start the server
	log.Fatal(app.Listen(":3000"))
}
