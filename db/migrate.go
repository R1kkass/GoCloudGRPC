package db

import (
	"os"

	Model "github.com/R1kkass/GoCloudGRPC/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Migration() {

	password, _ := os.LookupEnv("DB_PASSWORD")
	user, _ := os.LookupEnv("DB_USERNAME")
	port, _ := os.LookupEnv("DB_PORT")
	host, _ := os.LookupEnv("DB_HOST")
	db, _ := os.LookupEnv("DB_DATABASE")

	dsn := "host=" + host + " user=" + user + " password=" + password + " dbname=" + db + " port=" + port + " sslmode=disable TimeZone=Asia/Shanghai"
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Apply migration
	database.AutoMigrate(&Model.User{}, &Model.Folder{}, &Model.File{}, &Model.Accesses{}, &Model.Status{}, &Model.RequestAccess{}, &Model.Keys{}, &Model.KeysSecondary{}, &Model.Chat{}, &Model.ChatUser{}, &Model.Message{}, &Model.SavedKeys{}, &Model.UnReadedMessage{}, &Model.ChatFile{})
	database.Migrator()
}
