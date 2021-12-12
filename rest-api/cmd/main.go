package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	log "github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

var (
	persistenceURI = "http://persistence:8001/internal/"
)

func main() {

	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})
	app.Use(cors.New())

	app.Get("/api/getrandom", func(c *fiber.Ctx) error {
		statusCode, data, err := fasthttp.Get(nil, persistenceURI+"getrandom")
		if err != nil {
			log.Error(err)
		}
		if statusCode != fasthttp.StatusOK {
			c.SendStatus(statusCode)
		}
		return c.Send(data)
	})

	app.Get("/api/getall", func(c *fiber.Ctx) error {
		statusCode, data, err := fasthttp.Get(nil, persistenceURI+"getall")
		if err != nil {
			log.Error(err)
		}
		if statusCode != fasthttp.StatusOK {
			c.SendStatus(statusCode)
		}
		return c.Send(data)
	})

	app.Get("/api/deleteall", func(c *fiber.Ctx) error {
		statusCode, _, err := fasthttp.Get(nil, persistenceURI+"deleteall")
		if err != nil {
			log.Error(err)
		}
		if statusCode != fasthttp.StatusOK {
			c.SendStatus(statusCode)
		}
		return c.SendString("Successful database purge!")
	})

	log.Info("REST-API started successfully")

	app.Listen(":8000")
}
