package db

import (
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB
var DBStats *gorm.DB


func ConnectDatabase() *gorm.DB{

    password, _ := os.LookupEnv("DB_PASSWORD")
    user, _ := os.LookupEnv("DB_USERNAME")
    port, _ := os.LookupEnv("DB_PORT")
    host, _ := os.LookupEnv("DB_HOST")
    db, _ := os.LookupEnv("DB_DATABASE")
    
    dsn := "host=" + host + " user=" + user + " password="+ password + " dbname=" + db + " port="+ port + " sslmode=disable TimeZone=Asia/Shanghai"
    database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

    if err != nil {
        panic("Failed to connect to database!")
    }

    DB = database
    return DB
}

func ConnectDatabaseStats() *gorm.DB {

    password, _ := os.LookupEnv("DB_PASSWORD")
    user, _ := os.LookupEnv("DB_USERNAME")
    port, _ := os.LookupEnv("DB_PORT")
    host, _ := os.LookupEnv("DB_HOST")
    db, _ := os.LookupEnv("DB_DATABASE")
    
    dsn := "host=" + host + " user=" + user + " password="+ password + " dbname=" + db + " port="+ port + " sslmode=disable TimeZone=Asia/Shanghai"
    database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
        DisableForeignKeyConstraintWhenMigrating: true,
    })
    if err != nil {
        panic("Failed to connect to database!") 
    }

    DBStats = database
    return DBStats
}