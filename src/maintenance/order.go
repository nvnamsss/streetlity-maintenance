package maintenance

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/url"
	"strconv"
	"streetlity-maintenance/model"
	"streetlity-maintenance/server"
	"streetlity-maintenance/sres"
	"streetlity-maintenance/srpc"
)

func Order(common_user string, maintenance_users []string, reason string, note string) (str struct {
	sres.Response
	Order  model.MaintenanceOrder
	RoomId string
}, e error) {
	log.Println("[Router]", "len", len(maintenance_users))
	if len(maintenance_users) == 0 {
		return str, errors.New("there are no maintenance user that receives the order")
	}
	str.Order.CommonUser = common_user
	str.Order.Reason = reason
	str.Order.Note = note
	str.Order.SetReceiver(maintenance_users...)
	log.Println("[Router]", "len", len(maintenance_users))

	if str.Order, e = model.CreateOrder(str.Order); e != nil {
		return
	}

	data_id := "id:" + strconv.FormatInt(str.Order.Id, 10)
	data_user := "user:" + common_user
	data_reason := "reason:" + reason
	data_note := "note:" + note

	str.RoomId = strconv.FormatInt(str.Order.Id, 10)

	log.Println("[Order]", str)
	log.Println("[Order]", "maintenance users", maintenance_users)
	resp, e := srpc.RequestNotify(url.Values{
		"id":            maintenance_users,
		"notify-tittle": {"Customer is on service"},
		"notify-body":   {"A customer is looking for maintaining"},
		"data":          {data_id, data_user, data_reason, data_note},
	})

	if e != nil {
		log.Println("[Order]", "add", e.Error())
		str.Error(e)

		return
	}
	defer resp.Body.Close()

	server.OpenOrderSpace(str.RoomId)
	var res struct {
		Status  bool   `json:"Status"`
		Message string `json:"Message"`
	}

	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &res)

	if !res.Status {
		e = errors.New(res.Message)
	}

	return
}
