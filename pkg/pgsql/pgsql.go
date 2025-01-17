package pgsql

import (
	"OtterAnalytics/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"strconv"
	"time"
)

func ConnectPostgres() (*gorm.DB, error) {
	cfg := config.LoadConfig()
	db_cfg := "host=" + cfg.PgsqlHost +
		" port=" + strconv.Itoa(cfg.PgsqlPort) +
		" user=" + cfg.PgsqlUser +
		" dbname=" + cfg.PgsqlDb +
		" password=" + cfg.PgsqlPassword +
		" sslmode=disable"
	db, err := gorm.Open(postgres.Open(db_cfg), &gorm.Config{})
	if err != nil {
		log.Printf("Error connecting to postgres: %v", err)
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Printf("Error getting db instance: %v", err)
		return nil, err
	}
	sqlDB.SetMaxOpenConns(50)
	sqlDB.SetMaxIdleConns(20)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)

	log.Println("Connected to PostgreSQL with connection pool")
	return db, nil
}
