package example

import (
	"fiber_zerolog"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/rs/zerolog"
	"testing"
)

func TestExample(t *testing.T) {
	app := fiber.New()
	log.SetLogger(fiber_zerolog.NewLogger())

	app.Get("/", func(c *fiber.Ctx) error {
		log.Info("test logger")
		log.Infow("test logger with fields", "foo", "bar")
		return c.SendString("Hello, World!")
	})

	app.Get("/error", func(c *fiber.Ctx) error {
		log.Error("test error logger")
		return fiber.ErrInternalServerError
	})

	var ctxKey = struct{}{}

	app.Use(func(c *fiber.Ctx) error {
		c.Locals(ctxKey, "ctxValue")
		return c.Next()
	})
	app.Use(fiber_zerolog.InjectLogger(fiber_zerolog.InjectLoggerConfig{
		Logger:           fiber_zerolog.NewLogger(),
		FieldsFromHeader: []string{fiber.HeaderXRequestID},
		FieldsFromLocals: map[string]interface{}{
			"ctxKey": ctxKey,
		},
	}))

	// curl -H "X-Request-ID: 123456" http://localhost:3000/inject
	app.Get("/inject", func(c *fiber.Ctx) error {
		log.Info("test inject logger")
		log.Info(string(c.Request().Header.Peek(fiber.HeaderXRequestID)))
		log.WithContext(c.UserContext()).Info("test inject logger with context")
		zerolog.Ctx(c.UserContext()).Info().Msg("test inject logger with context by zerolog")
		return c.SendString("Hello, World!")
	})

	log.Fatal(app.Listen(":3000"))
}
