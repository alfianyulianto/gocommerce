package infrastucture

import (
	"fmt"
	"time"

	"github.com/alfianyulianto/gocommerce/config"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewDatabase(cfg *config.Config, log *logrus.Logger) *gorm.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Name)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.WithError(err).Error("Failed to connect to database")
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.WithError(err).Error("Failed to connect to database")
	}

	sqlDB.SetMaxIdleConns(cfg.Database.Pool.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.Database.Pool.MaxOpenConns)
	sqlDB.SetConnMaxIdleTime(time.Duration(cfg.Database.Pool.ConnIdleTime) * time.Minute)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.Database.Pool.ConnLifeTime) * time.Minute)

	return db
}
