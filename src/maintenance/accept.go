package maintenance

import (
	"errors"
	"streetlity-maintenance/model"
)

//Accept the order
func Accept(order_id int64, maintenance_user string) (order model.MaintenanceOrder, e error) {
	order = model.MaintenanceOrder{Id: order_id}
	if order, e = model.FindOrder(order); e != nil {
		return
	}

	if order.Status != model.Waiting {
		e = errors.New("This order is not available")
		return
	}

	order.Status = model.Accepted

	e = order.Save()
	return
}

func NotifyAccepted(order_id int64) {
	order = model.MaintenanceOrder{Id: order_id}
	if order, e = model 
}
