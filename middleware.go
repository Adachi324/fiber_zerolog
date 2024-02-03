package fiber_zerolog

import (
	"github.com/gofiber/fiber/v2"
)

type InjectLoggerConfig struct {
	Logger           *Logger
	FieldsFromHeader []string
	FieldsFromLocals map[string]interface{}
}

var DefaultInjectLoggerConfig = InjectLoggerConfig{
	Logger:           NewLogger(),
	FieldsFromHeader: []string{fiber.HeaderXRequestID},
}

// InjectLogger injects logger into fiber context, which can trace your request by your request id or something else,
// and you can get it by zerolog.Ctx(c.UserContext) or log.WithContext(c.UserContext)
// this middleware should be used after your keys are all set
func InjectLogger(cfg ...InjectLoggerConfig) func(ctx *fiber.Ctx) error {
	if len(cfg) == 0 {
		cfg = append(cfg, DefaultInjectLoggerConfig)
	}
	config := cfg[0]
	return func(ctx *fiber.Ctx) error {
		l := config.Logger.logger.With()
		for _, field := range config.FieldsFromHeader {
			val := string(ctx.Request().Header.Peek(field))
			if val != "" {
				l = l.Str(field, val)
			}
		}
		for name, field := range config.FieldsFromLocals {
			val := ctx.Locals(field)
			if val != nil {
				l = l.Interface(name, val)
			}
		}
		ctx.SetUserContext(l.Logger().WithContext(ctx.UserContext()))
		return ctx.Next()
	}
}
