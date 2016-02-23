package server

import (
	"encoding/json"
	"fmt"
	"html/template"
	"math/rand"
	"net/http"
	"sort"
	"strconv"
	"time"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/user"
)

func init() {
	rand.Seed(time.Now().UnixNano())
	if err := InitKeys(); err != nil {
		panic(err.Error())
	}

	http.HandleFunc("/", root)
	http.HandleFunc("/register", register)
	http.HandleFunc("/battery", battery)
	http.HandleFunc("/api/v1/register", registerApi)
	http.HandleFunc("/api/v1/battery", batteryApi)
}

func root(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	u := toUser(user.Current(c))
	if u == nil {
		url, err := user.LoginURL(c, r.URL.String())
		if err != nil {
			log.Errorf(c, err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Location", url)
		w.WriteHeader(http.StatusFound)
		return
	}
	var err error
	logoutUrl, err := user.LogoutURL(c, r.URL.String())
	if err != nil {
		log.Errorf(c, err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	q := datastore.NewQuery("Device").Ancestor(userKey(u, c)).Order("DeviceName")
	devices := []Device{}
	keys := []*datastore.Key{}
	if keys, err = q.GetAll(c, &devices); err != nil && !isErrFieldMismatch(err) {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	history := [][]Battery{}
	activeDevices := []Device{}
	for i, key := range keys {
		device := devices[i]
		if device.Disabled {
			continue
		}
		populateDeviceId(key, &device)

		qb := datastore.NewQuery("History").Ancestor(key).Order("-__key__").Limit(7)
		h := []History{}
		_, err := qb.GetAll(c, &h)
		if err != nil && !isErrFieldMismatch(err) {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		log.Debugf(c, "%#v", h)
		x := []Battery{}
		for _, his := range h {
			x = append(x, his.Batteries...)
		}
		sort.Sort(ByTime(x))

		activeDevices = append(activeDevices, device)
		history = append(history, x)
	}

	data := map[string]interface{}{
		"User":           u,
		"Devices":        activeDevices,
		"BatteryHistory": history,
		"LogoutUrl":      logoutUrl,
	}

	if err := registerTemplate.Execute(w, &data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// isErrFieldMismatch returns whether err is a datastore.ErrFieldMismatch.
func isErrFieldMismatch(err error) bool {
	_, ok := err.(*datastore.ErrFieldMismatch)
	return ok
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

	disabled := r.FormValue("disabled") == "on"

	g := Device{
		UserId:         u.UserId,
		DeviceId:       device,
		DeviceName:     deviceName,
		AlertThreshold: threshold,
		Disabled:       disabled,
	}

	_, err := datastore.Put(c, deviceKey(u, device, c), &g)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

type RegisterResponse struct {
	ID string `json:"id"`
}

func registerApi(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	u := toUser(user.Current(c))
	if u == nil {
		http.Error(w, "User is not authorized", http.StatusUnauthorized)
		return
	}

	deviceId := "api" + strconv.Itoa(rand.Int())
	d := Device{
		UserId:     u.UserId,
		DeviceId:   deviceId,
		DeviceName: "Unknown",
	}

	_, err := datastore.Put(c, deviceKey(u, deviceId, c), &d)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := RegisterResponse{
		ID: deviceId,
	}

	bytes, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(bytes)
}

func battery(w http.ResponseWriter, r *http.Request) {
	if err := batteryBase(w, r); err == nil {
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

func batteryApi(w http.ResponseWriter, r *http.Request) {
	if err := batteryBase(w, r); err == nil {
		w.Write([]byte{})
	}
}

func batteryBase(w http.ResponseWriter, r *http.Request) error {
	c := appengine.NewContext(r)
	u := toUser(user.Current(c))
	if u == nil {
		http.Error(w, "User is not authorized", http.StatusUnauthorized)
		return nil
	}

	deviceId := r.FormValue("device_id")

	battery, err := strconv.Atoi(r.FormValue("battery"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return err
	}

	charging := r.FormValue("charging") != ""

	t := time.Now()

	b := Battery{
		Battery:  int32(battery),
		Time:     t,
		Charging: charging,
	}

	log.Debugf(c, fmt.Sprintf("%+v", b))
	key := historyKey(u, deviceId, t, c)
	h, err := getHistory(key, c)
	if err != nil && err != datastore.ErrNoSuchEntity {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	h.Batteries = append(h.Batteries, b)
	_, err = datastore.Put(c, key, h)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	// Notify
	if battery <= 15 && !charging {
		key := deviceKey(u, deviceId, c)
		var device Device
		if err = datastore.Get(c, key, &device); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return err
		}
		err = notifyLowBattery(c, device.DeviceName, b.Battery, []string{keys.MyNexus5x})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return err
		}
	}

	return nil
}
