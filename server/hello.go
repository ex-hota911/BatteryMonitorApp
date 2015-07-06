package guestbook

import (
	"html/template"
	"net/http"

	"appengine"
	"appengine/datastore"
	"appengine/user"
)

func init() {
	http.HandleFunc("/", root)
	http.HandleFunc("/register", register)
}

func root(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	u := user.Current(c)
	if u == nil {
		url, err := user.LoginURL(c, r.URL.String())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Location", url)
		w.WriteHeader(http.StatusFound)
		return
	}

	q := datastore.NewQuery("Device").Filter("UserId =", u.ID).Order("DeviceName")

	devices := make([]Device, 10)
	if _, err := q.GetAll(c, &devices); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data := map[string]interface{}{
		"User":    u,
		"Devices": devices,
	}

	if err := registerTemplate.Execute(w, &data); err != nil {
		//if err := registerTemplate.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

var registerTemplate = template.Must(template.ParseFiles("index.html"))

func register(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	u := user.Current(c)
	if u == nil {
		http.Redirect(w, r, "/", http.StatusFound)
	}

	device := r.FormValue("device_id")
	deviceName := r.FormValue("device_name")
	// Default threshold is 15
	threshold := int32(15)

	strKey := u.String() + ":" + device
	g := Device{
		UserId:         u.ID,
		DeviceId:       device,
		DeviceName:     deviceName,
		AlertThreshold: threshold,
	}

	key := datastore.NewKey(c, "Device", strKey, 0, nil)
	_, err := datastore.Put(c, key, &g)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

func registerDevice(u *user.User, c appengine.Context, w http.ResponseWriter, r *http.Request) {
}
