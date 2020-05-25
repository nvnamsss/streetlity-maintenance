package router

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"streetlity-maintenance/maintenance"
	"streetlity-maintenance/srpc"

	"github.com/gorilla/mux"
	"github.com/nvnamsss/goinf/pipeline"
)

func requestOrder(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	p := pipeline.NewPipeline()

	vStage := pipeline.NewStage(func() (str struct{ ServiceId []int64 }, e error) {
		form := req.PostForm
		ids, ok := form["id"]

		if !ok {
			return str, errors.New("id param is missing")
		}

		for _, id := range ids {
			v, e := strconv.ParseInt(id, 10, 64)
			if e != nil {
				continue
			}

			str.ServiceId = append(str.ServiceId, v)
		}

		return
	})

	if res.Status {
		service_ids := p.GetInt("ServiceId")

		services := model.MaintenanceByIds(service_ids...)
		fmt.Println(services)
		ids := []string{}
		for _, s := range services {
			ids = append(ids, s.Owner)
		}

		maintenance.Order()
		resp, err := srpc.RequestNotify(url.Values{
			"id":            ids,
			"notify-tittle": {"Customer is on service"},
			"notify-body":   {"A customer is looking for maintaning"},
			"data":          {"score:sss", "id:1"},
		})

		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(resp)
	}
	p.First = vStage
	res.Error(p.Run())
}

func acceptOrder(w http.ResponseWriter, req *http.Request) {
	var res Response = Response{Status: true}

	req.ParseForm()
	p := pipeline.NewPipeline()
	vStage := pipeline.NewStage(func() (str struct {
		User      string
		OrderId   int64
		Timestamp int64
	}, e error) {
		form := req.PostForm
		_, ok := form["user"]
		if !ok {
			return str, errors.New("user param is missing")
		}

		timestamps, ok := form["timestamp"]
		if !ok {
			return str, errors.New("timestamp param is missing")
		}

		ids, ok := form["orderId"]
		if !ok {
			return str, errors.New("oderId param is missing")
		}

		id, e := strconv.ParseInt(ids[0], 10, 64)

		if e != nil {
			return str, errors.New("cannot parse orderId to int")
		}

		timestamp, e := strconv.ParseInt(timestamps[0], 10, 64)

		if e != nil {
			return str, errors.New("cannot parse timestamp to int")
		}

		str.User = form["user"][0]
		str.OrderId = id
		str.Timestamp = timestamp

		return
	})

	p.First = vStage

	res.Error(p.Run())

	if res.Status {
		user := p.GetString("User")[0]
		id := p.GetInt("OrderId")[0]
		timestamp := p.GetInt("Timestamp")[0]

		maintenance.Accept()
	}

	WriteJson(res)
}

func denyOrder(w http.ResponseWriter, req *http.Request) {
	var res Response = Response{Status: true}

	WriteJson(res)
}

func completeOrder(w http.ResponseWriter, req *http.Request) {
	var res Response = Response{Status: true}

	WriteJson(res)
}

func HandleOrder(router *mux.Router) {
	log.Println("[Router]", "Handle order")

	s := router.PathPrefix("/order").Subrouter()
}
