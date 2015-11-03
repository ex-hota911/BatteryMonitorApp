package server

import (
	"errors"
	"log"
	"time"

	"github.com/GoogleCloudPlatform/go-endpoints/endpoints"

	"appengine/datastore"
	"appengine/user"
)

const clientId = "546634630324-mkannoor781g7scn86vodbhol9qss1ev.apps.googleusercontent.com"

var (
	scopes    = []string{endpoints.EmailScope}
	clientIds = []string{clientId, endpoints.APIExplorerClientID}
	audiences = []string{clientId}
)

type BatteryService struct {
}

type UpdateReq struct {
	DeviceId  string
	Histories []History
}

type History struct {
	Level int32     `json:"level"`
	Time  time.Time `json:"timestamp"`
}

func (s *BatteryService) Update(c endpoints.Context, r *UpdateReq) error {
	u, err := getCurrentUser(c)
	if err != nil {
		return err
	}

	d := Device{
		UserId:         u.ID,
		DeviceId:       r.DeviceId,
		DeviceName:     r.DeviceId,
		AlertThreshold: 15,
	}
	c.Debugf("%#v", deviceKey(u, d.DeviceId, c))
	_, err = datastore.Put(c, deviceKey(u, d.DeviceId, c), &d)
	if err != nil {
		c.Debugf("%#v", err)
		return err
	}

	for _, h := range r.Histories {
		b := Battery{
			UserId:   u.ID,
			DeviceId: r.DeviceId,
			Battery:  h.Level,
			Time:     h.Time,
		}
		c.Debugf("%#v", batteryKey(u, b.DeviceId, b.Time, c))
		_, err = datastore.Put(c, batteryKey(u, b.DeviceId, b.Time, c), &b)
		if err != nil {
			c.Debugf("%#v", err)
			return err
		}
	}
	return nil
}

type ReadReq struct {
}

type ReadResp struct {
}

func (s *BatteryService) Read(c endpoints.Context, req *ReadReq) (*ReadResp, error) {
	u, err := getCurrentUser(c)
	if err != nil {
		return nil, err
	}
	_ = u
	resp := new(ReadResp)
	return resp, nil
}

// getCurrentUser retrieves a user associated with the request.
// If there's no user (e.g. no auth info present in the request) returns
// an "unauthorized" error.
func getCurrentUser(c endpoints.Context) (*user.User, error) {
	u, err := endpoints.CurrentUser(c, scopes, audiences, clientIds)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, errors.New("Unauthorized: Please, sign in.")
	}
	c.Debugf("Current user: %#v", u)
	return u, nil
}

func init() {
	service := &BatteryService{}
	api, err := endpoints.RegisterService(service,
		"battery", "v1", "Battery API", true)
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

	register("Update", "battery.update", "POST", "battery", "Update battery history.")
	endpoints.HandleHTTP()
}
