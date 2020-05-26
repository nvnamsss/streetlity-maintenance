package maintenance

import "streetlity-maintenance/model"

//Deny cancel the request
func Deny() (order model.MaintenanceOrder, e error) {
	order = model.MaintenanceOrder{}
	if order, e = model.FindOrder(order); e != nil {
		return
	}

	e = order.Delete()
	return
}
