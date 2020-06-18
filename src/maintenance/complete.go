package maintenance

import "streetlity-maintenance/model"

//Complete mark the request as completed
func Complete(order_id int64) (order model.MaintenanceOrder, e error) {
	order.Id = order_id
	if order, e = model.FindOrder(order); e != nil {
		return
	}

	order.Status = model.Completed
	e = order.Save()
	return
}
