package router

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"streetlity-maintenance/maintenance"
	"streetlity-maintenance/model"

	"github.com/gorilla/mux"
	"github.com/nvnamsss/goinf/pipeline"
)

func requestOrder(w http.ResponseWriter, req *http.Request) {
	var res struct {
		Response
		Order model.MaintenanceOrder
	}
	res.Status = true

	req.ParseForm()
	p := pipeline.NewPipeline()

	vStage := pipeline.NewStage(func() (str struct {
		CommonUser       string
		MaintenanceUsers []string
		Reason           string
		Note             string
	}, e error) {
		form := req.PostForm
		commonUsers, ok := form["commonUser"]

		if !ok {
			return str, errors.New("commonUser param is missing")
		}

		maintenanceUsers, ok := form["maintenanceUser"]

		if !ok {
			return str, errors.New("maintenanceUser param is missing")
		}

		reasons, ok := form["reason"]
		if !ok {
			return str, errors.New("reason param is missing")
		}

		notes, ok := form["note"]
		if ok {
			str.Note = notes[0]
		}

		str.CommonUser = commonUsers[0]
		str.Reason = reasons[0]
		str.MaintenanceUsers = maintenanceUsers

		return
	})
	p.First = vStage
	res.Error(p.Run())

	if res.Status {
		common_user := p.GetString("CommonUser")[0]
		maintenance_user_ids := p.GetString("MaintenanceUsers")
		reason := p.GetString("Reason")[0]

		note := p.GetStringFirstOrDefault("Note")

		if order, e := maintenance.Order(common_user, maintenance_user_ids, reason, note); e != nil {
			res.Error(e)
		} else {
			log.Println(order)
			res.Order = order
		}

	}

	WriteJson(w, res)
}

func acceptOrder(w http.ResponseWriter, req *http.Request) {
	var res Response = Response{Status: true}

	req.ParseForm()
	p := pipeline.NewPipeline()
	vStage := pipeline.NewStage(func() (str struct {
		MaintenanceUser string
		OrderId         int64
		Timestamp       int64
	}, e error) {
		form := req.PostForm
		_, ok := form["maintenanceUser"]
		if !ok {
			return str, errors.New("maintenanceUser param is missing")
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

		str.MaintenanceUser = form["user"][0]
		str.OrderId = id
		str.Timestamp = timestamp

		return
	})

	p.First = vStage

	res.Error(p.Run())

	if res.Status {
		maintenance_user := p.GetString("MaintenanceUser")[0]
		order_id := p.GetInt("OrderId")[0]
		// timestamp := p.GetInt("Timestamp")[0]

		maintenance.Accept(order_id, maintenance_user)
	}

	WriteJson(w, res)
}

func denyOrder(w http.ResponseWriter, req *http.Request) {
	var res Response = Response{Status: true}

	WriteJson(w, res)
}

func completeOrder(w http.ResponseWriter, req *http.Request) {
	var res Response = Response{Status: true}

	WriteJson(w, res)
}

func HandleOrder(router *mux.Router) {
	log.Println("[Router]", "Handle order")

	s := router.PathPrefix("/order").Subrouter()
	s.HandleFunc("/", requestOrder).Methods("POST")
}
