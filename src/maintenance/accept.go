package maintenance

import (
	"errors"
	"log"
	"net/url"
	"strconv"
	"streetlity-maintenance/model/order"
	"streetlity-maintenance/srpc"
	"time"
)

//Accept the order
func Accept(order_id int64, maintenance_user string) (o order.MaintenanceOrder, e error) {
	o = order.MaintenanceOrder{Id: order_id}
	if o, e = order.FindOrder(o); e != nil {
		return
	}

	if o.Status != order.Waiting {
		e = errors.New("This order is not available")
		return
	}

	o.MaintenanceUser = maintenance_user
	o.Timestamp = time.Now().Unix()
	o.Status = order.Accepted
	NotifyAccepted(o)
	e = o.Save()
	return
}

//NotifyAccepted send a notify back to common user to confirm that the order is accepted by some one
func NotifyAccepted(o order.MaintenanceOrder) {
	data_id := "id:" + strconv.FormatInt(o.Id, 10)
	data_action := "action:" + "Accepted"
	data_maintenance_user := "maintenance_user:" + o.MaintenanceUser

	log.Println("[Order]", "Send notify to", o.CommonUser)
	srpc.RequestNotify(url.Values{
		"id":            {o.CommonUser},
		"notify-tittle": {"Your request is accepted"},
		"notify-body":   {"A maintenance accepted your order"},
		"data":          {data_id, data_action, data_maintenance_user},
		"click-action":  {"MaintenanceAcceptNotification"},
	})

	data := make(map[string]string)
	data["id"] = strconv.FormatInt(o.Id, 10)
	data["maintenance_user"] = o.MaintenanceUser
	srpc.SaveNotification(o.CommonUser, "We got a dream", "No longer dreaming", data)
}
