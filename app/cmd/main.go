package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
	"hometask/internal/config"
	"hometask/internal/db"
	"hometask/internal/handler"
	"hometask/internal/service"
	"hometask/internal/web"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	ctx := context.Background()
	ctx, _ = signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	conf, err := config.NewConfig()
	if err != nil {
		panic(err)
	}

	loggerCfg := zap.NewProductionConfig()
	level := zap.InfoLevel
	if err := level.Set(conf.App.LogLevel); err != nil {
		panic(err)
	}
	loggerCfg.Level = zap.NewAtomicLevelAt(level)
	zapLogger, err := loggerCfg.Build()
	if err != nil {
		panic(err)
	}
	zap.ReplaceGlobals(zapLogger)

	dbConn, err := sql.Open("pgx", conf.GetDbConnString())
	if err != nil {
		zap.L().Panic("unable to connect to database", zap.String("host", conf.Db.Host))
	}
	defer dbConn.Close()

	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
		JSONEncoder:           json.Marshal,
		JSONDecoder:           json.Unmarshal,
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return ctx.Status(code).JSON(fiber.Map{
				"message": "Not found",
			})
		},
	})

	metricsApp := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})

	app.Use(cors.New())

	prometheus := fiberprometheus.New("hello-revolut")
	prometheus.RegisterAt(metricsApp, "/metrics")
	app.Use(prometheus.Middleware)

	redis := db.NewPostgres(dbConn)
	birthdayService := service.NewBirthdayService(redis)
	birthdayHandler := handler.NewBirthdayHandler(birthdayService)
	helloRouter := web.NewRouter(birthdayHandler)

	app.Get("/health", func(ctx *fiber.Ctx) error {
		return ctx.SendStatus(fiber.StatusOK)
	})

	app.Route("/hello", helloRouter.Route, "hello")

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		<-ctx.Done()
		_ = app.Shutdown()
		_ = metricsApp.Shutdown()
		wg.Done()
	}()

	go func() {
		if err := app.Listen(":" + conf.App.HttpPort); err != nil {
			panic(err)
		}
	}()

	go func() {
		if err := metricsApp.Listen(":" + conf.App.MetricsPort); err != nil {
			panic(err)
		}
	}()
	wg.Wait()
}
