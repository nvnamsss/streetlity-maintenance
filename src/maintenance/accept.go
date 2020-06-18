package maintenance

import (
	"errors"
	"net/url"
	"strconv"
	"streetlity-maintenance/model"
	"streetlity-maintenance/srpc"
	"time"
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

	order.Timestamp = time.Now().Unix()
	order.Status = model.Accepted
	NotifyAccepted(order)
	e = order.Save()
	return
}

//NotifyAccepted send a notify to other maintenance users that the order is accepted by other
func NotifyAccepted(order model.MaintenanceOrder) {
	receivers := order.GetReceiver()
	data_id := "id:" + strconv.FormatInt(order.Id, 10)
	data_action := "action:" + "Accepted"
	data_message := "message:" + "An order is accepted by other"

	srpc.RequestNotify(url.Values{
		"id":            receivers,
		"notify-tittle": {""},
		"notify-body":   {""},
		"data":          {data_id, data_action, data_message},
	})
}
