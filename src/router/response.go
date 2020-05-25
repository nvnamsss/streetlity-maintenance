package router

import (
	"encoding/json"
	"log"
	"net/http"
)

type Response struct {
	Status  bool
	Message string
}

//Error validate the data of response by err
func (res *Response) Error(err error) bool {
	if err != nil {
		res.Status = false
		res.Message = err.Error()
		return true
	}

	return false
}

func WriteJson(w http.ResponseWriter, data interface{}) {
	jsonData, jsonErr := json.Marshal(data)

	if jsonErr != nil {
		log.Println(jsonErr)
	}

	w.Write(jsonData)
}
