package router

import (
	"log"
	"net/http"
	"streetlity-maintenance/maintenance"
	"streetlity-maintenance/model"
	"streetlity-maintenance/stages"

	"github.com/gorilla/mux"
	"github.com/nvnamsss/goinf/pipeline"
)

func RequestOrder(w http.ResponseWriter, req *http.Request) {
	var res struct {
		Response
		Order model.MaintenanceOrder
	}
	res.Status = true

	p := pipeline.NewPipeline()
	stage := stages.RequestOrderValidate(req)
	p.First = stage
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

func AcceptOrder(w http.ResponseWriter, req *http.Request) {
	var res Response = Response{Status: true}

	p := pipeline.NewPipeline()
	stage := stages.AcceptOrderValidate(req)
	p.First = stage

	res.Error(p.Run())

	if res.Status {
		maintenance_user := p.GetStringFirstOrDefault("MaintenanceUser")
		order_id := p.GetIntFirstOrDefault("OrderId")
		// timestamp := p.GetInt("Timestamp")[0]

		maintenance.Accept(order_id, maintenance_user)
	}

	WriteJson(w, res)
}

func DenyOrder(w http.ResponseWriter, req *http.Request) {
	var res struct {
		Response
		Order model.MaintenanceOrder
	}

	p := pipeline.NewPipeline()
	stage := stages.DenyOrderValidate(req)
	p.First = stage
	res.Error(p.Run())

	if res.Status {
		order_id := p.GetIntFirstOrDefault("OrderId")
		deny_type := p.GetIntFirstOrDefault("DenyType")
		if order, e := maintenance.Deny(order_id, deny_type); e != nil {
			res.Error(e)
		} else {
			res.Order = order
		}
	}

	WriteJson(w, res)
}

func CompleteOrder(w http.ResponseWriter, req *http.Request) {
	var res struct {
		Response
		Order model.MaintenanceOrder
	}
	res.Status = true

	req.ParseForm()
	p := pipeline.NewPipeline()
	stage := stages.IdValidate(req.PostForm, "order_id")
	p.First = stage
	res.Error(p.Run())

	if res.Status {
		order_id := p.GetIntFirstOrDefault("Id")
		if order, e := maintenance.Complete(order_id); e != nil {
			res.Error(e)
		} else {
			res.Order = order
		}
	}

	WriteJson(w, res)
}

func HandleOrder(router *mux.Router) {
	log.Println("[Router]", "Handle order")

	s := router.PathPrefix("/order").Subrouter()
	s.HandleFunc("/request", RequestOrder).Methods("POST")
	s.HandleFunc("/accept", AcceptOrder).Methods("POST")
	s.HandleFunc("/deny", DenyOrder).Methods("POST")
	s.HandleFunc("/complete", CompleteOrder).Methods("POST")
}
