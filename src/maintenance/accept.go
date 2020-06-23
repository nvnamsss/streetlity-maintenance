package maintenance

import (
	"errors"
	"log"
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

	order.Receiver = maintenance_user
	order.Timestamp = time.Now().Unix()
	order.Status = model.Accepted
	NotifyAccepted(order)
	e = order.Save()
	return
}

//NotifyAccepted send a notify back to common user to confirm that the order is accepted by some one
func NotifyAccepted(order model.MaintenanceOrder) {
	data_id := "id:" + strconv.FormatInt(order.Id, 10)
	data_action := "action:" + "Accepted"
	data_receiver := "receiver:" + order.Receiver

	log.Println("[Order]", "Send notify to", order.CommonUser)
	srpc.RequestNotify(url.Values{
		"id":            {order.CommonUser},
		"notify-tittle": {"We got a dream"},
		"notify-body":   {"A dream is became true"},
		"data":          {data_id, data_action, data_receiver},
		"click-action":  {"MaintenanceAcceptNotification"},
	})
}
