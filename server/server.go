package server

import (
	"context"
	"fmt"
	"kn-assignment/internal/log"
	"kn-assignment/property"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func InitServer() *gin.Engine {
	mode := strings.TrimSpace(property.Get().Gin.Mode)
	switch mode {
	case "release", "debug", "test":
	default:
		mode = "release"
	}
	gin.SetMode(mode)
	return gin.New()
}

func StartServerWithCtx(ctx context.Context, h http.Handler, host string, port string) {
	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", host, port),
		Handler: h,
	}

	go func() {
		log.Infof(ctx, "server running at port: %s", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Errorf(ctx, "error server listen and serve: %v\n", err)
		}
	}()

	gracefulStop := make(chan os.Signal, 1)
	signal.Notify(gracefulStop, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-gracefulStop

	log.Info(ctx, "Shutdown Server ...")

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Errorf(ctx, "Shut down server error: %v\n", err)
	}
	select {
	case <-ctx.Done():
		log.Info(ctx, "timeout of 10 seconds.\n")
	default:
	}
	log.Info(ctx, "Server exiting")
}

func RunMigrations(ctx context.Context) {
	postgres := property.Get().Postgres

	databaseUrl := "postgres://" + postgres.User + ":" + postgres.Password + "@" + postgres.Host + ":" + postgres.Port + "/" + postgres.Database + "?sslmode=disable"

	migrationsPath := "file://migrations"
	log.Info(ctx, databaseUrl)
	if _, err := os.Stat("migrations"); os.IsNotExist(err) {
		log.Fatalf(ctx, "Migrations folder not found: %v", err)
	}
	m, err := migrate.New(
		migrationsPath,
		databaseUrl)
	if err != nil {
		log.Fatalf(ctx, "Failed to create migrate instance: %v", err)
	}

	log.Info(ctx, "Running migrations...")
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf(ctx, "Failed to run migrations: %v", err)
	}
	log.Info(ctx, "Migrations applied successfully")
}
