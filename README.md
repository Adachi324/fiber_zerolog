# fiber_zerolog
This makes it easy to inject a custom zerologger to your fiber framework

How to install it 
```shell
go get -u github.com/Adachi324/fiber_zerolog
```

## How to use this repo
It's very easy to replace your original fiber logger, and you don't need to modify any other log print code.

```go
// if you want to use a default logger
    logger := fiber_zerolog.NewLogger()
    log.SetLogger(logger)
    app := fiber.New()
    app.Get("/", func(c *fiber.Ctx) error {
		log.Info("test logger")
        log.Infow("test logger with fields", "foo", "bar")
		return c.SendString("Hello, World!")
	})
    app.Listen(":3000")
```
```shell
 # curl http://localhost:3000
 # outputs
 {"level":"info","time":"2024-01-24T17:00:16+08:00","message":"[test logger]"}
 {"level":"info","foo":"bar","time":"2024-01-24T17:01:27+08:00","message":"test logger with fields"}
```

Or if you want to use a custom zerologger, you can do this:
```go
    logger := fiber_zerolog.NewLogger()
    logger.SetLogger(zerolog.New(os.Stderr).With().Timestamp().Logger())
    log.SetLogger(logger)
    app := fiber.New()
    app.Get("/", func(c *fiber.Ctx) error {
        log.Info("test logger")
        log.Infow("test logger with fields", "foo", "bar")
        return c.SendString("Hello, World!")
    })
    app.Listen(":3000")
```
```shell
 # curl http://localhost:3000
 # outputs
 {"level":"info","time":"2024-01-24T17:00:16+08:00","message":"[test logger]"}
 {"level":"info","foo":"bar","time":"2024-01-24T17:01:27+08:00","message":"test logger with fields"}
```
You will get the similar outputs as above.


## Middleware
A injection middleware is provided, you can use it to trace some metadata like traceId or something else
And you can get a zerologger from the ctx of fiber

```go
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
```
```shell
{"level":"info","time":"2024-01-24T17:11:28+08:00","message":"[test inject logger]"}
{"level":"info","time":"2024-01-24T17:11:28+08:00","message":"[123456]"}
{"level":"info","X-Request-ID":"123456","ctxKey":"ctxValue","time":"2024-01-24T17:11:28+08:00","message":"[test inject logger with context]"}
{"level":"info","X-Request-ID":"123456","ctxKey":"ctxValue","time":"2024-01-24T17:11:28+08:00","message":"test inject logger with context by zerolog"}
```