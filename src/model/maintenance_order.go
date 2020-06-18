package model

import (
	"log"
	"regexp"

	"github.com/jinzhu/gorm"
)

type OrderStatus int

const (
	Waiting   = 0
	Accepted  = 1
	Completed = 2
	Denied    = 3
)

type MaintenanceOrder struct {
	Id              int64       `gorm:"column:id"`
	CommonUser      string      `gorm:"column:common_user",json:"common_user"`
	MaintenanceUser string      `gorm:"column:maintenance_user",json:"maintenance_user"`
	Timestamp       int64       `gorm:"type:datetime",gorm:"column:time"`
	Receiver        string      `gorm:"column:receiver",json:"receiver"`
	Reason          string      `gorm:"column:reason",json:"reason"`
	Note            string      `gorm:"column:note",json:"note"`
	Status          OrderStatus `gorm:"column:status"`
	db              *gorm.DB
}

func (MaintenanceOrder) TableName() string {
	return "maintenance_order"
}

func (order MaintenanceOrder) GetReceiver() (receivers []string) {
	reg := regexp.MustCompile(";")
	receivers = reg.Split(order.Receiver, -1)

	return
}

func AddOrder(order MaintenanceOrder) (rs MaintenanceOrder, e error) {
	rs = order
	if e := Db.Create(&rs).Error; e != nil {
		log.Println("[Database]", "add order", e.Error())
	}

	return
}

func FindOrder(embryo MaintenanceOrder) (order MaintenanceOrder, e error) {
	order = embryo
	if e = Db.Find(&order).Error; e != nil {
		log.Println("[Database]", "find order", e.Error())
	}

	order.db = Db
	return
}

func (order MaintenanceOrder) Save() (e error) {
	if e = order.db.Save(&order).Error; e != nil {
		log.Println("[Database]", "save order", e.Error())
	}
	return
}

func (order MaintenanceOrder) Delete() (e error) {
	if e = order.db.Delete(&order).Error; e != nil {
		log.Println("[Database]", "delete order", e.Error())
	}

	return
}
