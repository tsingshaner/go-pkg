package main

import (
	"fmt"
	"os"
	"time"

	"github.com/tsingshaner/go-pkg/log"
	"github.com/tsingshaner/go-pkg/log/console"
	"github.com/tsingshaner/go-pkg/log/helper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Report struct {
	gorm.Model
	ClassID      uint    `gorm:"index"`
	ExperimentID uint    `gorm:"index"`
	Student      string  `gorm:"index,type:varchar(8)"`
	Score        float32 `gorm:""`
	Status       uint    `gorm:""`
	Content      Content `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}

type Content struct {
	ReportID uint   `gorm:"primarykey"`
	Content  string `gorm:"not null"`
}

func main() {
	logger, _ := log.NewSlog(os.Stdout,
		&log.SlogHandlerOptions{
			Level: log.SlogLevelAll,
		},
		func(o *log.Options) {
			o.AddSource = true
		},
	)

	db := connectDB("postgres://qingshaner:123456@localhost:5432/mimo", &gorm.Config{
		Logger: helper.NewGormLogger(logger.WithGroup("orm"), helper.GORMLoggerOptions{
			ParameterizedQueries: false,
		}),
	})

	if err := db.Migrator().AutoMigrate(&Report{}, &Content{}); err != nil {
		logger.Fatal(fmt.Sprintf("migrate error: %v", err))
	} else {
		logger.Info("migrate success")
	}
}

func connectDB(dns string, config *gorm.Config) *gorm.DB {
	var db *gorm.DB
	var err error

	for i := 0; i < 5; i++ {
		db, err = gorm.Open(postgres.Open(dns), config)
		if err == nil {
			console.Info("connected database %s", dns)
			return db
		}

		time.Sleep(time.Second * time.Duration(i*2))
	}

	console.Fatal("%v", err)
	return nil
}
