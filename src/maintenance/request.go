package maintenance

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/url"
	"strconv"
	"streetlity-maintenance/model/order"
	"streetlity-maintenance/server"
	"streetlity-maintenance/sres"
	"streetlity-maintenance/srpc"
)

func Request(common_user string, maintenance_users []string, reason string, phone string, note string, order_type int) (str struct {
	sres.Response
	Order  order.MaintenanceOrder
	RoomId string
}, e error) {
	if len(maintenance_users) == 0 {
		return str, errors.New("there are no maintenance user that receives the order")
	}
	str.Order.CommonUser = common_user
	str.Order.Reason = reason
	str.Order.Note = note
	str.Order.Status = order.Waiting
	str.Order.Type = order_type
	str.Order.SetReceiver(maintenance_users...)

	if str.Order, e = order.CreateOrder(str.Order); e != nil {
		return
	}

	data_id := "id:" + strconv.FormatInt(str.Order.Id, 10)
	data_common_user := "common_user:" + common_user
	data_reason := "reason:" + reason
	data_note := "note:" + note

	str.RoomId = strconv.FormatInt(str.Order.Id, 10)

	log.Println("[Order]", str)
	log.Println("[Order]", "maintenance users", maintenance_users)
	resp, e := srpc.RequestNotify(url.Values{
		"id":            maintenance_users,
		"notify-tittle": {"Customer is on service"},
		"notify-body":   {"A customer is looking for maintaining"},
		"data":          {data_id, data_common_user, data_reason, data_note},
		"click-action":  {"MaintenanceOrderNotification"},
	})

	if e != nil {
		log.Println("[Order]", "add", e.Error())
		str.Error(e)

		return
	}
	defer resp.Body.Close()

	server.OpenOrderSpace("/" + str.RoomId)
	var res struct {
		Status  bool   `json:"Status"`
		Message string `json:"Message"`
	}

	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &res)

	log.Print("[Request]", res)
	if !res.Status {
		e = errors.New(res.Message)
	}

	return
}
