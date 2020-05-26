package maintenance

import "streetlity-maintenance/model"

//Complete mark the request as completed
func Complete() (order model.MaintenanceOrder, e error) {
	order = model.MaintenanceOrder{}
	if order, e = model.FindOrder(order); e != nil {
		return
	}

	order.Status = model.Completed
	e = order.Save()
	return
}
