package server

import (
	"log"
	"time"

	"github.com/GoogleCloudPlatform/go-endpoints/endpoints"

	"appengine/datastore"
)

var (
	scopes    = []string{endpoints.EmailScope}
	clientIds = []string{webClientId, androidReleaseClientId, androidDebugClientId, endpoints.APIExplorerClientID}
	audiences = []string{webClientId, androidReleaseClientId, androidDebugClientId}
)

type BatteryService struct {
}

type UpdateReq struct {
	DeviceId   string
	DeviceName string
	Histories  []History
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

	name := r.DeviceId
	if r.DeviceName != "" {
		name = r.DeviceName
	}

	d := Device{
		UserId:         u.UserId,
		DeviceId:       r.DeviceId,
		DeviceName:     name,
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
			UserId:   u.UserId,
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

type HelloReq struct {
	Message string
}

type HelloResp struct {
	Response string
}

// Hello world
func (s *BatteryService) Hello(c endpoints.Context, req *HelloReq) (*HelloResp, error) {
	u, err := getCurrentUser(c)
	if err != nil {
		return nil, err
	}
	resp := HelloResp{
		Response: "Hello, " + u.UserId + ". Your message is " + req.Message,
	}
	return &resp, nil
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

	register("Update", "update", "POST", "battery", "Update battery history.")
	register("Hello", "hello", "GET", "battery", "Hello battery history.")
	endpoints.HandleHTTP()
}
