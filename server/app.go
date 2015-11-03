package server

import (
	"html/template"
	"net/http"
	"strconv"
	"time"

	"appengine"
	"appengine/datastore"
	"appengine/user"
)

func init() {
	http.HandleFunc("/", root)
	http.HandleFunc("/register", register)
	http.HandleFunc("/battery", battery)
}

func root(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	u := toUser(user.Current(c))
	if u == nil {
		url, err := user.LoginURL(c, r.URL.String())
		if err != nil {
			c.Errorf(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Location", url)
		w.WriteHeader(http.StatusFound)
		return
	}

	q := datastore.NewQuery("Device").Ancestor(userKey(u, c)).Order("DeviceName")
	devices := []Device{}
	keys := []*datastore.Key{}
	var err error
	if keys, err = q.GetAll(c, &devices); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	history := [][]Battery{}
	for _, key := range keys {
		qb := datastore.NewQuery("Battery").Ancestor(key).Order("-__key__")
		h := []Battery{}
		if _, err := qb.GetAll(c, &h); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		history = append(history, h)
	}

	data := map[string]interface{}{
		"User":           u,
		"Devices":        devices,
		"BatteryHistory": history,
	}

	if err := registerTemplate.Execute(w, &data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

var registerTemplate = template.Must(template.ParseFiles("index.html"))

func register(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	u := toUser(user.Current(c))
	if u == nil {
		http.Redirect(w, r, "/", http.StatusFound)
	}

	device := r.FormValue("device_id")
	deviceName := r.FormValue("device_name")

	// Default threshold is 15
	threshold := int32(15)
	if t, err := strconv.Atoi(r.FormValue("alert_threshold")); err == nil {
		threshold = int32(t)
	}

	g := Device{
		UserId:         u.UserId,
		DeviceId:       device,
		DeviceName:     deviceName,
		AlertThreshold: threshold,
	}

	_, err := datastore.Put(c, deviceKey(u, device, c), &g)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

func battery(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	u := toUser(user.Current(c))
	if u == nil {
		http.Redirect(w, r, "/", http.StatusFound)
	}

	device := r.FormValue("device_id")

	battery, err := strconv.Atoi(r.FormValue("battery"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	t := time.Now()

	b := Battery{
		UserId:   u.UserId,
		DeviceId: device,
		Battery:  int32(battery),
		Time:     t,
	}

	_, err = datastore.Put(c, batteryKey(u, device, t, c), &b)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusFound)
}
