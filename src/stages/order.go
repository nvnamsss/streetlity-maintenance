package stages

import (
	"errors"
	"net/http"
	"net/url"
	"strconv"

	"github.com/nvnamsss/goinf/pipeline"
)

func RequestOrderValidate(req *http.Request) *pipeline.Stage {
	req.ParseForm()
	stage := pipeline.NewStage(func() (str struct {
		CommonUser       string
		MaintenanceUsers []string
		Reason           string
		Phone            string
		Note             string
		Type             int
	}, e error) {
		form := req.PostForm
		commonUsers, ok := form["common_user"]

		if !ok {
			return str, errors.New("common_user param is missing")
		}

		maintenanceUsers, ok := form["maintenance_users"]

		if !ok {
			return str, errors.New("maintenance_users param is missing")
		}

		reasons, ok := form["reason"]
		if !ok {
			return str, errors.New("reason param is missing")
		}

		phones, ok := form["phone"]
		if !ok {
			return str, errors.New("phone param is missing")
		}

		notes, ok := form["note"]
		if ok {
			str.Note = notes[0]
		}

		if types, ok := form["type"]; !ok {
			return str, errors.New("type param is missing")
		} else {
			if t, e := strconv.Atoi(types[0]); e != nil {
				return str, errors.New("type param cannot parse to int")
			} else {
				str.Type = t
			}
		}

		str.CommonUser = commonUsers[0]
		str.Reason = reasons[0]
		str.MaintenanceUsers = maintenanceUsers
		str.Phone = phones[0]

		return
	})

	return stage
}

func AcceptOrderValidate(req *http.Request) *pipeline.Stage {
	req.ParseForm()
	stage := pipeline.NewStage(func() (str struct {
		MaintenanceUser string
		OrderId         int64
	}, e error) {
		form := req.PostForm
		users, ok := form["maintenance_user"]
		if !ok {
			return str, errors.New("maintenance_user param is missing")
		}

		ids, ok := form["order_id"]
		if !ok {
			return str, errors.New("order_id param is missing")
		}

		id, e := strconv.ParseInt(ids[0], 10, 64)

		if e != nil {
			return str, errors.New("cannot parse orderId to int")
		}

		str.MaintenanceUser = users[0]
		str.OrderId = id

		return
	})

	return stage
}

func DenyOrderValidate(req *http.Request) *pipeline.Stage {
	req.ParseForm()
	stage := pipeline.NewStage(func() (str struct {
		OrderId  int64
		DenyType int
		Reason   string
	}, e error) {
		form := req.PostForm
		ids, ok := form["order_id"]
		if !ok {
			return str, errors.New("order_id param is missing")
		}

		denies, ok := form["deny_type"]
		if !ok {
			return str, errors.New("deny_type param is missing")
		}

		reasons, ok := form["reason"]
		if ok {
			str.Reason = reasons[0]
		}

		if id, e := strconv.ParseInt(ids[0], 10, 64); e != nil {
			return str, errors.New("order_id cannot parse to int")
		} else {
			str.OrderId = id
		}

		if deny, e := strconv.Atoi(denies[0]); e != nil {
			return str, errors.New("deny_type cannot parse to int")
		} else {
			str.DenyType = deny
		}

		return
	})

	return stage
}

func IdValidate(values url.Values, id_name string) *pipeline.Stage {
	stage := pipeline.NewStage(func() (str struct {
		Id int64
	}, e error) {
		ids, ok := values[id_name]
		if !ok {
			return str, errors.New(id_name + " param is missing")
		}

		if id, e := strconv.ParseInt(ids[0], 10, 64); e != nil {
			return str, errors.New(id_name + " cannot parse to int64")
		} else {
			str.Id = id
		}

		return
	})

	return stage
}
