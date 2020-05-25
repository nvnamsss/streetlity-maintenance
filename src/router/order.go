package router

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"streetlity-maintenance/maintenance"

	"github.com/gorilla/mux"
	"github.com/nvnamsss/goinf/pipeline"
)

func requestOrder(w http.ResponseWriter, req *http.Request) {
	var res Response = Response{Status: true}
	req.ParseForm()
	p := pipeline.NewPipeline()

	vStage := pipeline.NewStage(func() (str struct {
		Users  []string
		Reason string
		Note   string
	}, e error) {
		form := req.PostForm
		users, ok := form["user"]

		if !ok {
			return str, errors.New("user param is missing")
		}

		reasons, ok := form["reason"]
		if !ok {
			return str, errors.New("reason param is missing")
		}

		notes, ok := form["note"]
		if ok {
			str.Note = notes[0]
		}

		str.Reason = reasons[0]
		str.Users = users

		return
	})

	if res.Status {
		user_ids := p.GetString("Users")
		reason := p.GetString("Reason")[0]
		note := p.GetStringFirstOrDefault("Note")

		maintenance.Order(user_ids, reason, note)
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
