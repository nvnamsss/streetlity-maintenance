package model_test

import (
	"streetlity-maintenance/model"
	"testing"

	"github.com/brianvoe/gofakeit/v5"
)

func TestCreateHistory(t *testing.T) {
	model.ConnectSync()

	for loop := 0; loop < 100; loop++ {
		var h model.MaintenanceHistory
		h.Reason = gofakeit.Sentence(20)
		h.Note = gofakeit.Sentence(8)
		h.CommonUser = gofakeit.FirstName()
		h.MaintenanceUser = gofakeit.FirstName()
		if e := model.CreateMaintenanceHistory(h); e != nil {
			t.Errorf(e.Error())
		}
	}
	t.Log("Completed")
}
