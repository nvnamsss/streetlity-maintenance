package main

import (
	"net/http"
	"streetlity-maintenance/model"

	"github.com/gorilla/mux"
)

var Router *mux.Router = mux.NewRouter()
var Server http.Server

func main() {
	model.Connect()
}
