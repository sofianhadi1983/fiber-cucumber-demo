package database

import (
	"fmt"
	"log"
	"os"

	"github.com/tiketdatasofianhadiant0/bdd-demo/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Dbinstance struct {
	Db *gorm.DB
}

var DB Dbinstance

func ConnectDb() {
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASS")
	dbname := os.Getenv("DB_NAME")
	dbport := os.Getenv("DB_PORT")

	dsn := fmt.Sprintf("host=%v user=%v password='%v' dbname=%v port=%v sslmode=disable TimeZone=Asia/Jakarta", host, user, password, dbname, dbport)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Error),
	})

	if err != nil {
		log.Fatal("Failed to connect to database. \n", err)
	}

	log.Println("connected")
	log.Println("running migrations")
	db.AutoMigrate(&models.Book{})

	DB = Dbinstance{
		Db: db,
	}
}
