package main

import (
	"OtterAnalytics/config"
	"OtterAnalytics/handlers"
	"OtterAnalytics/internal/app"
	"OtterAnalytics/models"
	"OtterAnalytics/pkg/errors"
	"OtterAnalytics/pkg/pgsql"
	"context"
	"gorm.io/gorm"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

func mustInitDB() *gorm.DB {
	db, err := pgsql.ConnectPostgres()
	errors.Must(err, "Error initializing DB")
	return db
}

func mustMigrateDB(db *gorm.DB) {
	for _, model := range models.AllModels {
		if !db.Migrator().HasTable(model) {
			if err := db.AutoMigrate(model); err != nil {
				errors.Must(err, "Error migrating table")
			}
		}
	}
}

func closeDB(db *gorm.DB) {
	sqlDB, err := db.DB()
	if err != nil {
		errors.Normal(err, "Error getting DB connection for closing")
		return
	}
	if err := sqlDB.Close(); err != nil {
		errors.Normal(err, "Error closing DB connection")
	}
}

func closeListener(listener net.Listener) {
	if err := listener.Close(); err != nil {
		errors.Normal(err, "Error closing listener")
	}
}

func acceptConnections(ctx context.Context, listener net.Listener, handler *handlers.Handler) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			select {
			case <-ctx.Done():
				return
			default:
				log.Println("Error accepting connection:", err)
				continue
			}
		}
		go handler.HandleConnection(conn)
	}
}

func main() {
	cfg := config.LoadConfig()

	db := mustInitDB()
	defer closeDB(db)
	mustMigrateDB(db)

	handler := handlers.NewHandler(db, cfg)
	app.InitAppRoutes()

	addr := cfg.Host + ":" + strconv.Itoa(cfg.Port)
	listener, err := net.Listen("tcp", addr)
	errors.Must(err, "Error listening on port "+strconv.Itoa(cfg.Port))
	defer closeListener(listener)
	log.Println("Server is listening on port " + strconv.Itoa(cfg.Port) + " ...")

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	go acceptConnections(ctx, listener, handler)

	<-ctx.Done()
	log.Println("Shutting down server gracefully...")
}
