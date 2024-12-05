package main

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	. "github.com/opentreehole/backend/common"
	. "github.com/opentreehole/backend/image_hosting/api"
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
	// app.Use(recover.New(recover.Config{
	// 	EnableStackTrace:  true,
	// 	StackTraceHandler: common.StackTraceHandler,
	// }))
	// app.Use(common.MiddlewareGetUserID)
	// app.Use(common.MiddlewareCustomLogger)
	Init()
	router := app.Group("/api")
	router.Post("/uploadImage", UploadImage)
	// router.Get("/i/:year/:month/:day/:identifier", api.GetImage)
	router.Get("/getImage.png", GetImage)

	go func() {
		err := app.Listen(":8000")
		if err != nil {
			log.Println(err)
		}
	}()

	interrupt := make(chan os.Signal, 1)

	// wait for CTRL-C interrupt
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-interrupt

	// close app
	err := app.Shutdown()
	if err != nil {
		log.Println(err)
	}
}
