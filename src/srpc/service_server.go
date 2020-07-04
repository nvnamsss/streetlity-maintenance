package srpc

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"streetlity-maintenance/config"
	"strings"
)

func prepareHeader(req *http.Request) {
	req.Header.Set("Auth", config.Config.SixtyDaysKey)
	req.Header.Set("Version", "1.0.0")
}

func GetMaintenanceService(ids ...int64) (res []struct {
	Id      int64   `gorm:"column:id"`
	Lat     float32 `gorm:"column:lat"`
	Lon     float32 `gorm:"column:lon"`
	Note    string  `gorm:"column:note"`
	Address string  `gorm:"column:address"`
	Images  string  `gorm:"column:images"`
	Owner   string  `gorm:"column:owner"`
	Name    string  `gorm:"column:name"`
}, e error) {
	host := "http://" + config.Config.Servicehost + "/user/notify"
	req, _ := http.NewRequest("GET", host, nil)
	prepareHeader(req)

	query := req.URL.Query()
	for _, id := range ids {
		query.Add("id", strconv.FormatInt(id, 10))
	}
	req.URL.RawQuery = query.Encode()

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		log.Println("[SRPC]", "get maintenance", e.Error())
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	json.Unmarshal(body, &res)

	return
}

func SaveNotification(receiver string, title string, content string, data map[string]string) (res struct {
	Status  bool   `json:"Status"`
	Message string `json:"Message"`
}, e error) {
	host := "http://" + config.Config.Servicehost + "/service/maintenance/create"
	values := url.Values{}
	values.Add("name", receiver)
	values.Add("location", "1")
	values.Add("location", "1")
	values.Add("note", title)
	values.Add("address", content)
	values.Add("alt", "nani")

	values.Add("id", receiver)
	values.Add("notify-title", title)
	values.Add("notify-body", content)
	for key, value := range data {
		values.Add("data", key+":"+value)
	}
	req, _ := http.NewRequest("POST", host, strings.NewReader(values.Encode()))
	prepareHeader(req)

	client := &http.Client{}
	resp, e := client.Do(req)

	if e != nil {
		log.Println("[SRPC]", "save notification", e.Error())
		return
	}

	defer resp.Body.Close()
	body, e := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &res)

	return
}
