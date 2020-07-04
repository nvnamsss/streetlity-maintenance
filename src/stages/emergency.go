package stages

import (
	"errors"
	"net/http"
	"net/url"
	"strconv"

	"github.com/nvnamsss/goinf/pipeline"
)

func CreateEmergencyValidate(req *http.Request) *pipeline.Stage {
	req.ParseForm()
	stage := pipeline.NewStage(func() (str struct {
		Id  string
		Lat float64
		Lon float64
	}, e error) {
		form := req.PostForm
		if ids, ok := form["id"]; !ok {
			return str, errors.New("id param is missing")
		} else {
			str.Id = ids[0]
		}

		if lats, ok := form["lat"]; !ok {
			return str, errors.New("lat param is missing")
		} else {
			if lat, e := strconv.ParseFloat(lats[0], 64); e != nil {
				return str, errors.New("lat param cannot parse to float")
			} else {
				str.Lat = lat
			}
		}

		if lons, ok := form["lon"]; !ok {
			return str, errors.New("lon param is missing")
		} else {
			if lon, e := strconv.ParseFloat(lons[0], 64); e != nil {
				return str, errors.New("lon param cannot parse to float")
			} else {
				str.Lon = lon
			}
		}
		return
	})

	return stage
}

func UpdateEmergencyLocationValidate(req *http.Request) *pipeline.Stage {
	req.ParseForm()
	stage := pipeline.NewStage(func() (str struct {
		Id  string
		Lat float64
		Lon float64
	}, e error) {
		form := req.PostForm
		if ids, ok := form["id"]; !ok {
			return str, errors.New("id param is missing")
		} else {
			str.Id = ids[0]
		}

		if lats, ok := form["lat"]; !ok {
			return str, errors.New("lat param is missing")
		} else {
			if lat, e := strconv.ParseFloat(lats[0], 64); e != nil {
				return str, errors.New("lat param cannot parse to float")
			} else {
				str.Lat = lat
			}
		}

		if lons, ok := form["lon"]; !ok {
			return str, errors.New("lon param is missing")
		} else {
			if lon, e := strconv.ParseFloat(lons[0], 64); e != nil {
				return str, errors.New("lon param cannot parse to float")
			} else {
				str.Lon = lon
			}
		}
		return
	})

	return stage
}

func IdEmergencyValidate(values url.Values) *pipeline.Stage {
	stage := pipeline.NewStage(func() (str struct {
		Id string
	}, e error) {
		if ids, ok := values["id"]; !ok {
			return str, errors.New("id param is missing")
		} else {
			str.Id = ids[0]
		}
		return
	})
	return stage
}

func EmergencyInRangeValidate(req *http.Request) *pipeline.Stage {
	stage := pipeline.NewStage(func() (str struct {
		Lat   float64
		Lon   float64
		Range float64
	}, e error) {
		query := req.URL.Query()
		if lats, ok := query["lat"]; !ok {
			return str, errors.New("lat param is missing")
		} else {
			if lat, e := strconv.ParseFloat(lats[0], 64); e != nil {
				return str, errors.New("lat param cannot parse to float")
			} else {
				str.Lat = lat
			}
		}

		if lons, ok := query["lon"]; !ok {
			return str, errors.New("lon param is missing")
		} else {
			if lon, e := strconv.ParseFloat(lons[0], 64); e != nil {
				return str, errors.New("lon param cannot parse to float")
			} else {
				str.Lon = lon
			}
		}

		if ranges, ok := query["range"]; !ok {
			return str, errors.New("range param is missing")
		} else {
			if r, e := strconv.ParseFloat(ranges[0], 64); e != nil {
				return str, errors.New("range param cannot parse to float")
			} else {
				str.Range = r
			}
		}
		return
	})
	return stage
}

func OrderEmergencyValidate(req *http.Request) *pipeline.Stage {
	req.ParseForm()
	stage := pipeline.NewStage(func() (str struct {
		CommonUser string
		Emergency  []string
	}, e error) {
		form := req.PostForm
		if cusers, ok := form["common_user"]; !ok {
			return str, errors.New("common_user param is missing")
		} else {
			str.CommonUser = cusers[0]
		}

		if ems, ok := form["emergency"]; !ok {
			return str, errors.New("emergency param is missing")
		} else {
			str.Emergency = ems
		}
		return
	})

	return stage
}
