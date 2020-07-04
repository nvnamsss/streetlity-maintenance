package emergency

import (
	"log"
	"math"
	"streetlity-maintenance/model"

	"github.com/golang/geo/r2"
	"github.com/jinzhu/gorm"
	"github.com/nvnamsss/goinf/spatial"
)

type EmergencyMaintenance struct {
	Id  string  `gorm:"column:id"`
	Lat float64 `gorm:"column:lat"`
	Lon float64 `gorm:"column:lon"`
}

const EmergencyMaintenanceTableName = "emergency_maintenance"

var eme_maintenance spatial.Quadtree
var map_eme map[string]*spatial.Quadtree

func (EmergencyMaintenance) TableName() string {
	return EmergencyMaintenanceTableName
}

func (e EmergencyMaintenance) GetId() string {
	return e.Id
}

func (e EmergencyMaintenance) Location() r2.Point {
	var location r2.Point
	location.X = e.Lat
	location.Y = e.Lon
	return location
}

func (e EmergencyMaintenance) CreateBound() (b spatial.Bounds) {
	b.X = e.Lat
	b.Y = e.Lon
	b.Height = 0.01
	b.Width = 0.01
	b.Item = e

	return
}

func EmergencyById(id string) (m EmergencyMaintenance, e error) {
	e = model.GetById(EmergencyMaintenanceTableName, id, &m)
	return
}
func AllEmergency() (ms []EmergencyMaintenance, e error) {
	if e = model.Db.Find(&ms).Error; e != nil {
		log.Println("[Database]", e.Error())
	}

	return
}

func CreateEmergency(id string, lat float64, lon float64) (m EmergencyMaintenance, e error) {
	m.Id = id
	m.Lat = lat
	m.Lon = lon

	if e = model.Db.Create(&m).Error; e != nil {
		log.Println("[Database]", "create emergency", e.Error())
	}

	return
}

func RemoveEmergency(id string) (e error) {
	if e = model.Db.Where("id=?", id).Delete(&EmergencyMaintenance{Id: id}).Error; e != nil {
		log.Println("[Database]", "remove emergency", e.Error())
	}
	return
}

func (m *EmergencyMaintenance) AfterSave(scope *gorm.Scope) (err error) {
	if qt, ok := map_eme[m.Id]; ok {
		qt.RemoveItem(m)
	}

	map_eme[m.Id] = eme_maintenance.Insert(m.CreateBound())
	log.Println("[Database]", "a maintenance saved")
	return
}

// func (m EmergencyMaintenance) AfterCreate(scope *gorm.Scope) (e error) {
// 	map_eme[m.Id] = eme_maintenance.Insert(m.CreateBound())
// 	log.Println("[Database]", "New maintennace added")
// 	return
// }

func (m *EmergencyMaintenance) AfterDelete(tx *gorm.DB) (e error) {
	if qt, ok := map_eme[m.Id]; ok {
		log.Println("[Database]", "a maintennace deleted")
		qt.RemoveItem(m)
	}

	return
}

func UpdateEmergencyLocation(id string, lat float64, lon float64) (m EmergencyMaintenance, e error) {
	m, e = EmergencyById(id)
	if e != nil {
		return
	}

	m.Lat = lat
	m.Lon = lon

	if e = model.Db.Save(&m).Error; e != nil {
		log.Println("[Database]", "update emergency", e.Error())
	}
	return
}

func distance(p1 r2.Point, p2 r2.Point) float64 {
	x := math.Pow(p1.X-p2.X, 2)
	y := math.Pow(p1.Y-p2.Y, 2)
	return math.Sqrt(x + y)
}

func EmergenciesInRange(p r2.Point, max_range float64) (ms []EmergencyMaintenance, e error) {
	f := spatial.Bounds{X: p.X, Y: p.Y, Height: max_range, Width: max_range}
	bounds := eme_maintenance.Retrieve(f)
	for _, b := range bounds {
		location := b.Item.Location()
		d := distance(location, p)
		if d < max_range {
			ms = append(ms, b.Item.(EmergencyMaintenance))
		}
	}
	return
}

func LoadEmergency() {
	map_eme = make(map[string]*spatial.Quadtree)
	eme_maintenance.MaxLevels = 100
	eme_maintenance.MaxObjects = 8

	ms, _ := AllEmergency()
	for _, m := range ms {
		map_eme[m.Id] = eme_maintenance.Insert(m.CreateBound())
	}
}

func init() {
	model.OnConnected.Subscribe(LoadEmergency)
	model.OnDisconnect.Subscribe(func() {
		model.OnConnected.Unsubscribe(LoadEmergency)
	})
}
