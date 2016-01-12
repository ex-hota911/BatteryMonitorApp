package server

import (
	logger "log"
	"time"

	"github.com/GoogleCloudPlatform/go-endpoints/endpoints"
	"golang.org/x/net/context"

	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
)

var (
	scopes    = []string{endpoints.EmailScope}
	clientIds = []string{webClientId, androidReleaseClientId, androidDebugClientId, endpoints.APIExplorerClientID}
	audiences = []string{webClientId, androidReleaseClientId, androidDebugClientId}
)

type BatteryService struct {
}

type UpdateReq struct {
	Device Device
}

type History struct {
	Level int32     `json:"level"`
	Time  time.Time `json:"timestamp"`
}

func (s *BatteryService) Update(c context.Context, r *UpdateReq) error {
	u, err := getCurrentUser(c)
	if err != nil {
		return err
	}

	// Store Device
	d := r.Device
	log.Debugf(c, "%#v", deviceKey(u, d.DeviceId, c))
	_, err = datastore.Put(c, deviceKey(u, d.DeviceId, c), &d)
	if err != nil {
		log.Debugf(c, "%#v", err)
		return err
	}

	// Store Histories
	for _, b := range r.Device.Batteries {
		log.Debugf(c, "%#v", batteryKey(u, d.DeviceId, b.Time, c))
		_, err = datastore.Put(c, batteryKey(u, d.DeviceId, b.Time, c), &b)
		if err != nil {
			log.Debugf(c, "%#v", err)
			return err
		}
	}
	return nil
}

type ReadReq struct {
}

type ReadResp struct {
}

func (s *BatteryService) Read(c context.Context, req *ReadReq) (*ReadResp, error) {
	u, err := getCurrentUser(c)
	if err != nil {
		return nil, err
	}
	_ = u
	resp := new(ReadResp)
	return resp, nil
}

// Hello world
type HelloReq struct {
	Message string
}

type HelloResp struct {
	Response string
}

func (s *BatteryService) Hello(c context.Context, req *HelloReq) (*HelloResp, error) {
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
	api, err := endpoints.RegisterService(service, "", "v1", "Battery API", true)
	if err != nil {
		logger.Fatalf("Register service: %v", err)
	}

	register := func(orig, name, method, path, desc string) {
		m := api.MethodByName(orig)
		if m == nil {
			logger.Fatalf("Missing method %s", orig)
		}
		i := m.Info()
		i.Name, i.HTTPMethod, i.Path, i.Desc = name, method, path, desc
	}

	register("Update", "update", "POST", "battery", "Update battery history.")
	register("Hello", "hello", "GET", "battery", "Hello battery history.")
	endpoints.HandleHTTP()
}
