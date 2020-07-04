package router_test

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"streetlity-maintenance/model/order"
	"strings"
	"testing"
)

func TestOrder(t *testing.T) {
	host := "http://localhost:9002/order/"
	client := &http.Client{}

	var orders struct {
		orders []order.MaintenanceOrder `json:"orders"`
	}

	file, fileErr := os.Open("./order_test.json")

	if fileErr != nil {
		log.Panic(fileErr)
	}

	defer file.Close()
	decoder := json.NewDecoder(file)

	_ = decoder.Decode(&orders)

	log.Println(orders)
	for _, order := range orders.orders {
		form := url.Values{
			"commonUser":       {order.CommonUser},
			"maintenanceUsers": order.GetReceiver(),
			"reason":           {order.Reason},
			"note":             {order.Note},
		}

		req, _ := http.NewRequest("POST", host, strings.NewReader(form.Encode()))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		resp, e := client.Do(req)

		if e != nil {
			log.Println(e.Error())
		}

		defer resp.Body.Close()

		var res struct {
			Status  bool   `json:"Status"`
			Message string `json:"Message"`
		}

		body, _ := ioutil.ReadAll(resp.Body)
		json.Unmarshal(body, &res)

		log.Println(res)
	}
}

func TestAcceptOrder(t *testing.T) {

}
