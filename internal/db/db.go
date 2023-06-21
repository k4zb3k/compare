package db

import (
	"compare/internal/models"
	"compare/pkg/logging"
	"fmt"
	postgresDriver "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	logger   = logging.GetLogger()
	host     = "localhost"
	port     = "5432"
	user     = "test"
	dbname   = "test_db"
	password = "pass"
	sslmode  = "disable"
)

func GetDbConn() (*gorm.DB, error) {
	connString := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		host, port, user, dbname, password, sslmode)

	db, err := gorm.Open(postgresDriver.Open(connString))
	if err != nil {
		logger.Error("wrong data: %v", err)
		return nil, err
	}
	logger.Infoln("successful DB connection")

	if !db.Migrator().HasTable(&models.HumoPayment{}) {
		if err := db.AutoMigrate(&models.HumoPayment{}); err != nil {
			logger.Errorf("error while auto_migrating Humo_Payment(models): %v", err)
		}

	}

	if !db.Migrator().HasTable(&models.PartnerPayment{}) {
		if err := db.AutoMigrate(&models.PartnerPayment{}); err != nil {
			logger.Errorf("error while auto_migrating Partner_Payment(models): %v", err)
		}
	}

	if !db.Migrator().HasTable(&models.Partner{}) {
		if err := db.AutoMigrate(&models.Partner{}); err != nil {
			logger.Errorf("error while auto_migrating Partner(models): %v", err)
		}
	}

	if !db.Migrator().HasTable(&models.Reestr{}) {
		if err := db.AutoMigrate(&models.Reestr{}); err != nil {
			logger.Errorf("error while auto_migrating Reestr(models): %v", err)
		}
	}

	return db, nil
}
