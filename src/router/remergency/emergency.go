package remergency

import (
	"net/http"
	"streetlity-maintenance/model/emergency"
	"streetlity-maintenance/sres"
	"streetlity-maintenance/stages"

	"github.com/golang/geo/r2"
	"github.com/gorilla/mux"
	"github.com/nvnamsss/goinf/pipeline"
)

func GetEmergency(w http.ResponseWriter, req *http.Request) {
	var res struct {
		sres.Response
		Emergency emergency.EmergencyMaintenance
	}

	res.Response = sres.Response{Status: true}

	p := pipeline.NewPipeline()
	stage := stages.IdEmergencyValidate(req.URL.Query())
	p.First = stage
	res.Error(p.Run())

	if res.Status {
		id := p.GetString("Id")[0]
		if m, e := emergency.EmergencyById(id); e != nil {
			res.Error(e)
		} else {
			res.Emergency = m
		}
	}
	sres.WriteJson(w, res)
}

func EmergenciesInRange(w http.ResponseWriter, req *http.Request) {
	var res struct {
		sres.Response
		Emergencies []emergency.EmergencyMaintenance
	}
	res.Response = sres.Response{Status: true}

	p := pipeline.NewPipeline()
	stage := stages.EmergencyInRangeValidate(req)
	p.First = stage
	res.Error(p.Run())

	if res.Status {
		lat := p.GetFloat("Lat")[0]
		lon := p.GetFloat("Lon")[0]
		r := p.GetFloat("Range")[0]
		p := r2.Point{X: lat, Y: lon}
		if ms, e := emergency.EmergenciesInRange(p, r); e != nil {
			res.Error(e)
		} else {
			res.Emergencies = ms
		}
	}

	sres.WriteJson(w, res)
}

func CreateEmergencyMaintenance(w http.ResponseWriter, req *http.Request) {
	var res struct {
		sres.Response
		Emergency emergency.EmergencyMaintenance
	}
	res.Response = sres.Response{Status: true, Message: "Create emergency maintenance successfully"}

	p := pipeline.NewPipeline()
	stage := stages.CreateEmergencyValidate(req)
	p.First = stage
	res.Error(p.Run())

	if res.Status {
		id := p.GetString("Id")[0]
		lat := p.GetFloat("Lat")[0]
		lon := p.GetFloat("Lon")[0]

		if m, e := emergency.CreateEmergency(id, lat, lon); e != nil {
			res.Error(e)
		} else {
			res.Emergency = m
		}
	}

	sres.WriteJson(w, res)
}

func UpdateEmergencyLocation(w http.ResponseWriter, req *http.Request) {
	var res struct {
		sres.Response
		Emergency emergency.EmergencyMaintenance
	}
	res.Response = sres.Response{Status: true, Message: "Update successfully"}

	p := pipeline.NewPipeline()
	stage := stages.UpdateEmergencyLocationValidate(req)
	p.First = stage
	res.Error(p.Run())

	if res.Status {
		id := p.GetString("Id")[0]
		lat := p.GetFloat("Lat")[0]
		lon := p.GetFloat("Lon")[0]

		if m, e := emergency.UpdateEmergencyLocation(id, lat, lon); e != nil {
			res.Error(e)
		} else {
			res.Emergency = m
		}
	}

	sres.WriteJson(w, res)
}

func RemoveEmergency(w http.ResponseWriter, req *http.Request) {
	var res sres.Response = sres.Response{Status: true, Message: "Remove emergency successfully"}

	req.ParseForm()
	p := pipeline.NewPipeline()
	stage := stages.IdEmergencyValidate(req.URL.Query())
	p.First = stage
	res.Error(p.Run())

	if res.Status {
		id := p.GetString("Id")[0]
		if e := emergency.RemoveEmergency(id); e != nil {
			res.Error(e)
		}
	}

	sres.WriteJson(w, res)
}

func OrderEmergency(w http.ResponseWriter, req *http.Request) {
	var res struct {
		sres.Response
	}
	res.Response = sres.Response{Status: true, Message: "Order emergency successfully"}

	sres.WriteJson(w, res)
}

func HandleEmergency(router *mux.Router) *mux.Router {
	s := router.PathPrefix("/emergency").Subrouter()
	s.HandleFunc("/", GetEmergency).Methods("GET")
	s.HandleFunc("/", CreateEmergencyMaintenance).Methods("POST")
	s.HandleFunc("/", RemoveEmergency).Methods("DELETE")
	s.HandleFunc("/location", UpdateEmergencyLocation).Methods("POST")
	s.HandleFunc("/range", EmergenciesInRange).Methods("GET")
	return s
}
