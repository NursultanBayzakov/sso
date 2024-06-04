package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"grpc-service-ref/internal/lib/logger/handlers/slogpretty"
	"log"
	"log/slog"
	"os"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		connStr = "postgres://auth:auth@localhost:5433/postgres?sslmode=disable"
		log.Println("DATABASE_URL not set, using default connection string")
	}

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	fmt.Println("Successfully connected to the database")

	// Uncomment and configure the following when ready to use the full application
	/*
	   cfg := config.MustLoad()
	   log := setupLogger(cfg.Env)

	   application := app.New(log, cfg.GRPC.Port, cfg.StoragePath, cfg.TokenTTL)

	   go func() {
	       application.GRPCServer.MustRun()
	   }()

	   // Graceful shutdown
	   stop := make(chan os.Signal, 1)
	   signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	   <-stop

	   application.GRPCServer.Stop()
	   log.Info("Gracefully stopped")
	*/
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = setupPrettySlog()
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}

func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)
	return slog.New(handler)
}
