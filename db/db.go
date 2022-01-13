package db

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	SqlDB *gorm.DB
)

func SetUp() {
	databaseConfig := viper.GetStringMap("database")
	username := databaseConfig["username"]
	password := databaseConfig["password"]
	name := databaseConfig["name"]
	maxIdleConns := databaseConfig["max_idle_conns"]
	maxOpenConns := databaseConfig["max_open_conns"]
	connMaxLifetime := databaseConfig["conn_max_lifetime"]
	dsn := fmt.Sprintf("%s:%s@tcp(localhost:3306)/%s?charset=utf8mb4&parseTime=True&loc=Local", username.(string), password.(string), name.(string))

	var openError error
	SqlDB, openError = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if openError != nil {
		fmt.Println("init db error", openError)
	}
	fmt.Println("sqldn: ", SqlDB)
	db, err := SqlDB.DB()
	if err != nil {
		fmt.Println("get db error: ", err)
	}
	db.SetMaxIdleConns(maxIdleConns.(int))
	db.SetMaxOpenConns(maxOpenConns.(int))
	db.SetConnMaxLifetime(time.Duration(connMaxLifetime.(int)) * time.Millisecond)
}
