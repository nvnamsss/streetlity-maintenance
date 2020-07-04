package maintenance

import "streetlity-maintenance/model/order"

//Complete mark the request as completed
func Complete(order_id int64) (o order.MaintenanceOrder, e error) {
	o.Id = order_id
	if o, e = order.FindOrder(o); e != nil {
		return
	}

	o.Status = order.Completed
	e = o.Save()
	return
}
