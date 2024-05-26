package db

import (
	"github.com/yunsonggo/helper/types"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"io"
	"log"
	"os"
	"time"
)

func NewGormDB(conf *types.Sql) (*gorm.DB, error) {
	var newLogger logger.Interface
	if !conf.StdPrint {
		newLogger = logger.New(log.New(io.Discard, "", log.LstdFlags), logger.Config{})
	} else {
		newLogger = logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				SlowThreshold: time.Second,
				LogLevel:      logger.Info,
				Colorful:      true,
			},
		)
	}
	dsn := conf.MysqlDSN()
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		return nil, err
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxIdleConns(conf.MaxIdleConn)
	sqlDB.SetMaxOpenConns(conf.MaxOpenConn)
	sqlDB.SetConnMaxLifetime(time.Hour * time.Duration(conf.ConnMaxLifeHour))
	sqlDB.SetConnMaxIdleTime(time.Hour * time.Duration(conf.ConnMaxLifeHour))
	if err = sqlDB.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
