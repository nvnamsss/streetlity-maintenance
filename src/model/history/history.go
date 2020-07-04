package history

import (
	"log"
	"streetlity-maintenance/model"
	"time"
)

type MaintenanceHistory struct {
	Id              int64
	MaintenanceUser string `gorm:"column:maintenance_user"`
	CommonUser      string `gorm:"column:common_user"`
	Timestamp       int64  `gorm:"type:datetime"`
	Reason          string `gorm:"column:reason"`
	Note            string `gorm:"column:note"`
}

const HistoryTableName = "maintenance_history"

func (MaintenanceHistory) TableName() string {
	return HistoryTableName
}

func CreateHistory(h MaintenanceHistory) (e error) {
	h.Timestamp = time.Now().Unix()
	if e = model.Db.Create(&h).Error; e != nil {
		log.Println("[Database]", "Adding new history:", e.Error())
	}

	return
}

func RemoveHistory(h MaintenanceHistory) (e error) {
	if e = model.Db.Delete(h).Error; e != nil {
		log.Println("[Database]", "Removing history:", e.Error())
	}

	return
}

func RemoveHistoriesById(ids ...int64) (e error) {
	for _, id := range ids {
		if e = model.Db.Where("id=?", id).Delete(&MaintenanceHistory{}).Error; e != nil {
			log.Println("[Database]", "remove M history", e.Error())
		}
	}

	return
}

func HistoriesByMUser(mUser string) (histories []MaintenanceHistory, e error) {
	if e = model.Db.Where("maintenance_user=?", mUser).Find(&histories).Error; e != nil {
		log.Println("[Database]", "query M history by M user", e.Error())
	}

	return
}

func FindMaintenanceHistory(embryo MaintenanceHistory) (history MaintenanceHistory, e error) {
	history = embryo
	if e = model.Db.Find(&history).Error; e != nil {
		log.Println("[Database]", "find maintenance history", e.Error())
	}

	return
}

func HistoryById(id int64) (h MaintenanceHistory, e error) {
	if e := model.Db.Find(&h, id).Error; e != nil {
		log.Println("[Database]", "Maintenance history with id:", id, ":", e.Error())
	}

	return
}

func UpdateHistory(id int64, maintenanceUser string, timestamp int64) error {
	h, e := HistoryById(id)

	if e != nil {
		return e
	}

	h.MaintenanceUser = maintenanceUser
	h.Timestamp = timestamp

	if e = model.Db.Save(&h).Error; e != nil {
		log.Println("[Database]", "Update maintenance history with id:", id, ":", e.Error())
	}

	return e
}
