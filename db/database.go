package db

import (
	"fmt"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql" // import _ for gorm

	"github.com/jatgam/wishlist-api/config"
)

var DB *gorm.DB

func Connect(dbConf *config.DBConfig) *gorm.DB {
	connectString := fmt.Sprintf("%s:%s@(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbConf.User, dbConf.Password, dbConf.Hostname, dbConf.Database)
	db, err := gorm.Open("mysql", connectString)
	if err != nil {
		fmt.Println("db err: ", err)
		os.Exit(1)
	}
	db.LogMode(true)
	DB = db
	return DB
}

func GetDB() *gorm.DB {
	return DB
}
