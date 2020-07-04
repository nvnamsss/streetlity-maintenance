package maintenance

import (
	"net/url"
	"strconv"
	"streetlity-maintenance/model/order"
	"streetlity-maintenance/srpc"
)

//Deny remove the order, notify to the receiver user if it was accepted
func Deny(order_id int64, decline_type int64) (o order.MaintenanceOrder, e error) {
	o.Id = order_id
	if o, e = order.FindOrder(o); e != nil {
		return
	}

	e = o.Delete()
	return
}

//NotifyDeny send a notify to other maintenance users that the order is denied
func NotifyDeny(o order.MaintenanceOrder) {
	receivers := []string{o.CommonUser}
	data_id := "id:" + strconv.FormatInt(o.Id, 10)
	data_action := "action:" + "Denied"
	data_status := "status:" + strconv.Itoa(int(o.Status))
	data_message := "message:" + "order is declined"

	srpc.RequestNotify(url.Values{
		"id":            receivers,
		"notify-tittle": {"An order is denied"},
		"notify-body":   {"An order is denied"},
		"data":          {data_id, data_action, data_status, data_message},
		"click-action":  {"MaintenanceOrderDenyNotification"},
	})
}
