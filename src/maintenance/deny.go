package maintenance

import "streetlity-maintenance/model"

//Deny remove the order, notify to the receiver user if it was accepted
func Deny() (order model.MaintenanceOrder, e error) {
	order = model.MaintenanceOrder{}
	if order, e = model.FindOrder(order); e != nil {
		return
	}

	e = order.Delete()
	return
}
