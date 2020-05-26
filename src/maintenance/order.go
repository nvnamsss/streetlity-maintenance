package maintenance

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/url"
	"strconv"
	"streetlity-maintenance/model"
	"streetlity-maintenance/srpc"
)

func Order(common_user_id string, maintenance_user_ids []string, reason string, note string) (order model.MaintenanceOrder, e error) {
	if len(maintenance_user_ids) == 0 {
		return order, errors.New("there are no maintenance user that receives the order")
	}
	receiver_length := len(maintenance_user_ids)
	order.CommonUser = common_user_id
	order.Reason = reason
	order.Note = note
	order.Receiver = maintenance_user_ids[0]
	for loop := 0; loop < receiver_length; loop++ {
		order.Receiver += ";" + maintenance_user_ids[loop]
	}

	if order, e = model.AddOrder(order); e != nil {
		return
	}

	data_id := "id:" + strconv.FormatInt(order.Id, 10)
	data_user := "user:" + common_user_id
	data_reason := "reason:" + reason
	data_note := "note:" + note
	log.Println(order)
	resp, e := srpc.RequestNotify(url.Values{
		"id":            maintenance_user_ids,
		"notify-tittle": {"Customer is on service"},
		"notify-body":   {"A customer is looking for maintaning"},
		"data":          {data_id, data_user, data_reason, data_note},
	})

	if e != nil {
		log.Println("[Order]", "add", e.Error())
	}
	defer resp.Body.Close()

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
