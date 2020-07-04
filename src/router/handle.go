package router

import (
	"streetlity-maintenance/router/remergency"

	"github.com/gorilla/mux"
)

func Handle(router *mux.Router) {
	remergency.HandleEmergency(router)
	HandleOrder(router)
}
