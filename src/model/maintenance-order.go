package model

import (
	"log"
	"regexp"

	"github.com/jinzhu/gorm"
)

type MaintenanceOrder struct {
	Id         int64  `gorm:"column:id"`
	CommonUser string `gorm:"column:common_user"`
	Timestamp  int64  `gorm:"type:datetime",gorm:"column:time"`
	Receiver   string `gorm:"column:receiver"`
	Reason     string `gorm:"column:reason"`
	Note       string `gorm:"column:note"`
	db         *gorm.DB
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
	if dbc := Db.Create(&rs); dbc.Error != nil {
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

func (order MaintenanceOrder) Save() {
	order.db.Save(&order)
}

func (order MaintenanceOrder) Delete() {
	order.db.Delete(&order)
}
