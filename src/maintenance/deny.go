package maintenance

import (
	"net/url"
	"strconv"
	"streetlity-maintenance/model"
	"streetlity-maintenance/srpc"
)

//Deny remove the order, notify to the receiver user if it was accepted
func Deny(order_id int64, deny_type int64) (order model.MaintenanceOrder, e error) {
	order.Id = order_id
	if order, e = model.FindOrder(order); e != nil {
		return
	}

	e = order.Delete()
	return
}

//NotifyDenied send a notify to other maintenance users that the order is accepted by other
func NotifyDenied(order model.MaintenanceOrder) {
	receivers := []string{order.CommonUser}
	data_id := "id:" + strconv.FormatInt(order.Id, 10)
	data_action := "action:" + "Denied"
	data_message := "message:" + "order is denied by maintenance user"

	srpc.RequestNotify(url.Values{
		"id":            receivers,
		"notify-tittle": {""},
		"notify-body":   {""},
		"data":          {data_id, data_action, data_message},
	})
}
