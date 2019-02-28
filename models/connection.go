package models

import (
	log "github.com/sirupsen/logrus"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/spf13/viper"
)

var DB *gorm.DB
var err error

func Connect() {
	DB, err = gorm.Open("postgres", viper.GetString("postgres_url"))
	if err != nil {
		log.Fatal(err)
	}

	// AutoMigrate will ONLY create tables, missing columns and missing indexes,
	// and WON’T change existing column’s type or delete unused columns to protect your data.
	DB.AutoMigrate(&Relayer{}, &Token{}, &Trade{})
	// DB.LogMode(true)
	log.Info("DB Connected")
}

func Close() error {
	return DB.Close()
}
