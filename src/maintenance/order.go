package maintenance

import (
	"net/url"
	"streetlity-maintenance/model"
	"streetlity-maintenance/srpc"
)

func Order(ids []string, reason string, note string) {
	order := model.MaintenanceOrder{}
	order, e := model.AddOrder(order)

	resp, err := srpc.RequestNotify(url.Values{
		"id":            ids,
		"notify-tittle": {"Customer is on service"},
		"notify-body":   {"A customer is looking for maintaning"},
		"data":          {"id:" + order.Id, "reason:" + reason, "note:" + note},
	})
}

func Accept() {
	order := model.MaintenanceOrder{}
	order, e := model.FindOrder(order)

	order.Save()
}

func Deny() {
	order := model.MaintenanceOrder{}
	order, e := model.FindOrder(order)

	order.Delete()
}

func Complete() {
	order := model.MaintenanceOrder{}
	order, e := model.FindOrder(order)

	order.Save()
}
