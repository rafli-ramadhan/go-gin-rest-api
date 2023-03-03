package connection

import (
	"fmt"
	log "github.com/forkyid/go-utils/v1/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
	"os"
)

var db *gorm.DB

type DB struct {
	Master *gorm.DB
}

func DBMaster() *gorm.DB {
	if db == nil {
		log.Infof("creating database connection")
		db = dbConnection(os.Getenv("DB_POSTGRES_HOST_MASTER"), os.Getenv("DB_POSTGRES_DATABASE"))
	}
	return db
}

func dbConnection(hostType, dbName string) *gorm.DB {
	port := os.Getenv("DB_POSTGRES_PORT")
	postgresCon := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s",
		hostType,
		port,
		os.Getenv("DB_POSTGRES_USERNAME"),
		dbName,
		os.Getenv("DB_POSTGRES_PASSWORD"),
	)

	db, err := gorm.Open(postgres.Open(postgresCon), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		QueryFields: true,
	})
	if err != nil {
		log.Fatalf(nil, fmt.Sprintf("failed to connect %s on %s:%s", dbName, hostType, port), err)
	}
	log.Infof(fmt.Sprintf("successfully connected to %s on %s:%s", dbName, hostType, port))

	connConfiguration, err := db.DB()
	if err != nil {
		log.Fatalf(nil, fmt.Sprintf("failed to connect %s on %s:%s", dbName, hostType, port), err)
	}

	connConfiguration.SetConnMaxLifetime(5 * time.Minute)
	connConfiguration.SetMaxIdleConns(20)
	connConfiguration.SetMaxOpenConns(200)
	connConfiguration.SetMaxOpenConns(20)

	return db
}
