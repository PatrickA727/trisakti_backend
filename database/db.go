package database

import (
	"fmt"
	"log"
	"os"

	// "github.com/PatrickA727/trisakti-proto/config/db_config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var DB *gorm.DB
var connectionError error

func ConnectDB() (*gorm.DB, error) {
    USER := os.Getenv("DB_USER");
    PASS := os.Getenv("DB_PASS");
    NAME := os.Getenv("DB_NAME");

    dsn := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", USER, PASS, NAME);

    DB, connectionError = gorm.Open(postgres.Open(dsn), &gorm.Config{
        NamingStrategy: schema.NamingStrategy{
            SingularTable: true,
        },
    })
    if connectionError != nil {
        log.Fatal("error connecting to database: ", connectionError)
        return nil, connectionError
    } 

    log.Println("connected to database")

    return DB, nil
}
