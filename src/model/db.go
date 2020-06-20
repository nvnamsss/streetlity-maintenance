package model

import (
	"fmt"
	"log"
	"streetlity-maintenance/config"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/nvnamsss/goinf/event"
)

var Db *gorm.DB

// var Config Configuration
var OnDisconnect *event.Event = event.NewEvent()
var OnConnected *event.Event = event.NewEvent()

func connect() {
	log.Println("[Database]", "connecting to database")
	connectionString := fmt.Sprintf("%s:%s@tcp(%s)/%s",
		config.Config.Username, config.Config.Password, config.Config.Server, config.Config.Database)
	db, err := gorm.Open("mysql", connectionString)
	Db = db

	if err != nil {
		OnDisconnect.Invoke()
		log.Println(err.Error())
	} else {
		OnConnected.Invoke()
		log.Println("[Database]", "connect success")
	}

}

func reconnect() {
	timer := time.NewTimer(10 * time.Second)
	<-timer.C

	connect()
}

//Connect to the database by default config
func Connect() {
	OnDisconnect.Subscribe(reconnect)
	go connect()
}

func ConnectSync() {
	OnDisconnect.Subscribe(reconnect)
	connect()
}
