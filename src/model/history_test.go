package model_test

import (
	"streetlity-maintenance/model"
	"streetlity-maintenance/model/history"
	"testing"

	"github.com/brianvoe/gofakeit/v5"
)

func TestCreateHistory(t *testing.T) {
	model.ConnectSync()

	for loop := 0; loop < 100; loop++ {
		var h history.MaintenanceHistory
		h.Reason = gofakeit.Sentence(20)
		h.Note = gofakeit.Sentence(8)
		h.CommonUser = gofakeit.FirstName()
		h.MaintenanceUser = gofakeit.FirstName()
		if e := history.CreateHistory(h); e != nil {
			t.Errorf(e.Error())
		}
	}
	t.Log("Completed")
}
