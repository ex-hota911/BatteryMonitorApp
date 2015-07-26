package server

import (
	"log"
	"net/http"
	"os/user"
	"strconv"
	"time"

	"github.com/GoogleCloudPlatform/go-endpoints/endpoints"

	"appengine/datastore"
)

type BatteryListReq struct {
}

type BatteryListResp struct {
}

type BatteryService struct {
}

type TimestampNanos time.Time

func (bs *BatteryService) List(c endpoints.Context, r *BatteryListReq) (*BatteryListResp, error) {
	return nil, nil
}

type UpdateReq struct {
	Id        string    `json:"id"`
	Histories []History `json:"histories"`
}

type History struct {
	Level           int32 `json:"level"`
	TimestampMillis Time  `json:"timestamp"`
}

func (bs *BatteryService) Update(c endpoints.Context, r *UpdateReq) error {
}

type BatteryRegisterReq struct {
	DeviceId       string `json:"id"`              // Unique ID for a device
	DeviceName     string `json:"name"`            // Display name.
	AlertThreshold int32  `json:"alert_threshold"` // 0 - 100.
}

func (bs *BatteryService) Register(c endpoints.Context, r *BatteryRegisterReq) error {
	u := user.Current(c)
	if u == nil {
		return endpoints.UnauthorizedError
	}

	device := r.FormValue("device_id")
	deviceName := r.FormValue("device_name")

	// Default threshold is 15
	threshold := int32(15)
	if t, err := strconv.Atoi(r.FormValue("alert_threshold")); err == nil {
		threshold = int32(t)
	}

	g := Device{
		UserId:         u.ID,
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

// Greeting is a datastore entity that represents a single greeting.
// It also serves as (a part of) a response of GreetingService.
type Greeting struct {
	Key     *datastore.Key `json:"id" datastore:"-"`
	Author  string         `json:"author"`
	Content string         `json:"content" datastore:",noindex" endpoints:"req"`
	Date    time.Time      `json:"date"`
}

// GreetingsList is a response type of GreetingService.List method
type GreetingsList struct {
	Items []*Greeting `json:"items"`
}

// Request type for GreetingService.List
type GreetingsListReq struct {
	Limit int `json:"limit" endpoints:"d=10"`
}

// GreetingService can sign the guesbook, list all greetings and delete
// a greeting from the guestbook.
type GreetingService struct {
}

// List responds with a list of all greetings ordered by Date field.
// Most recent greets come first.
func (gs *GreetingService) List(c endpoints.Context, r *GreetingsListReq) (*GreetingsList, error) {
	if r.Limit <= 0 {
		r.Limit = 10
	}

	q := datastore.NewQuery("Greeting").Order("-Date").Limit(r.Limit)
	greets := make([]*Greeting, 0, r.Limit)
	keys, err := q.GetAll(c, &greets)
	if err != nil {
		return nil, err
	}

	for i, k := range keys {
		greets[i].Key = k
	}
	return &GreetingsList{greets}, nil
}

// Add adds a greeting.
func (gs *GreetingService) Add(c endpoints.Context, g *Greeting) error {
	k := datastore.NewIncompleteKey(c, "Greeting", nil)
	_, err := datastore.Put(c, k, g)
	return err
}

type Count struct {
	N int `json:"count"`
}

// Count returns the number of greetings.
func (gs *GreetingService) Count(c endpoints.Context) (*Count, error) {
	n, err := datastore.NewQuery("Greeting").Count(c)
	if err != nil {
		return nil, err
	}
	return &Count{n}, nil
}

func init() {
	greetService := &GreetingService{}
	api, err := endpoints.RegisterService(greetService,
		"greeting", "v1", "Greetings API", true)
	if err != nil {
		log.Fatalf("Register service: %v", err)
	}

	register := func(orig, name, method, path, desc string) {
		m := api.MethodByName(orig)
		if m == nil {
			log.Fatalf("Missing method %s", orig)
		}
		i := m.Info()
		i.Name, i.HTTPMethod, i.Path, i.Desc = name, method, path, desc
	}

	register("List", "greets.list", "GET", "greetings", "List most recent greetings.")
	register("Add", "greets.add", "PUT", "greetings", "Add a greeting.")
	register("Count", "greets.count", "GET", "greetings/count", "Count all greetings.")
	endpoints.HandleHTTP()
}
