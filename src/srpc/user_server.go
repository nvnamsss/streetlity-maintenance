package srpc

import (
	"net/http"
	"net/url"
	"streetlity-maintenance/config"
)

func RequestNotify(values url.Values) (res *http.Response, e error) {
	host := "http://" + config.Config.UserHost + "/user/notify/"

	res, e = http.PostForm(host, values)

	return
}
