package config

import (
	"fmt"
	"log"
	"os"

	"github.com/gusrylmubarok/mygram-backend/src/domain"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDB() *gorm.DB {
	var (
		env      = os.Getenv("APP_ENV")
		host     = os.Getenv("DB_PG_HOST")
		user     = os.Getenv("DB_PG_USER")
		password = os.Getenv("DB_PG_PASSWD")
		dbname   = os.Getenv("DB_PG_NAME")
		port     = os.Getenv("DB_PG_PORT")
		timeZone = os.Getenv("APP_TIMEZONE")
		dsn      = ""
		db       *gorm.DB
		err      error
	)

	if env == "production" {
		dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=require TimeZone=%s", host, user, password, dbname, port, timeZone)
	} else {
		dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=%s", host, user, password, dbname, port, timeZone)
	}

	if db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{FullSaveAssociations: true}); err != nil {
		log.Fatal("Error connecting to database: ", err)
	}

	if err = db.AutoMigrate(&domain.User{}, &domain.Photo{}, &domain.Comment{}, &domain.SocialMedia{}); err != nil {
		log.Fatal("Error migrating database: ", err.Error())
	}

	return db
}
