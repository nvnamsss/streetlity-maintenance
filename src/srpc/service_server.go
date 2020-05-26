package srpc

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"streetlity-maintenance/config"
)

func prepareHeader(req *http.Request) {
	req.Header.Set("Auth", config.Config.SixtyDaysKey)
	req.Header.Set("Version", "1.0.0")
}

func RequestMaintenance(ids ...int64) (res []struct {
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
		log.Println(err.Error())
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	json.Unmarshal(body, &res)

	return
}
