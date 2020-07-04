package model

import (
	"errors"
	"log"
)

func GetById(tablename string, id interface{}, ref interface{}) (e error) {
	db := Db.Table(tablename).Where("id=?", id).First(ref)
	e = db.Error

	if db.RowsAffected == 0 {
		e := errors.New("record was not found")
		log.Println("[Database]", "get", tablename, e.Error())
	}

	return
}
