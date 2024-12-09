package main

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	. "github.com/opentreehole/backend/common"
	. "github.com/opentreehole/backend/image_hosting/api"
	. "github.com/opentreehole/backend/image_hosting/config"
	. "github.com/opentreehole/backend/image_hosting/model"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	app := fiber.New(fiber.Config{
		ErrorHandler:          ErrorHandler,
		JSONEncoder:           json.Marshal,
		JSONDecoder:           json.Unmarshal,
		DisableStartupMessage: true,
		BodyLimit:             128 * 1024 * 1024,
	})
	// will catch every panic
	app.Use(recover.New(recover.Config{
		EnableStackTrace: true,
		// StackTraceHandler: StackTraceHandler,
	}))

	// app.Use(common.MiddlewareGetUserID)

	// will catch every HTTP request
	// app.Use(MiddlewareCustomLogger)
	Init()
	router := app.Group("/api")
	router.Post("/uploadImage", UploadImage)

	// get images based on the identifier(exclude the extension)
	// format: http://localhost:8000/api/i/2024/12/06/6288772352016bf28f1a571d0.jpg
	router.Get("/i/:year/:month/:day/:identifier", GetImage)

	go func() {
		log.Println("Server is running on " + Config.HostName)
		err := app.Listen(Config.HostName)
		if err != nil {
			log.Println(err)
		}

	}()

	interrupt := make(chan os.Signal, 1)

	// listen for interrupt signal (Ctrl+C)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-interrupt
	log.Println("Shutting down the server...")
	err := app.Shutdown()
	if err != nil {
		log.Println(err)
	}
}
