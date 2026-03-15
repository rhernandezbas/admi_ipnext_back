package config

import (
	"fmt"

	"github.com/ipnext/admin-backend/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewDatabase(cfg *config.Config) (*gorm.DB, error) {
	gormCfg := &gorm.Config{}

	if cfg.Server.Env == "development" {
		gormCfg.Logger = logger.Default.LogMode(logger.Info)
	} else {
		gormCfg.Logger = logger.Default.LogMode(logger.Error)
	}

	db, err := gorm.Open(mysql.Open(cfg.Database.DSN()), gormCfg)
	if err != nil {
		return nil, fmt.Errorf("error conectando a MySQL: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("error obteniendo sql.DB: %w", err)
	}

	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(5)

	return db, nil
}
