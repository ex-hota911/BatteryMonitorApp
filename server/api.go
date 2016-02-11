package server

import (
	logger "log"

	"keys"

	"github.com/GoogleCloudPlatform/go-endpoints/endpoints"
	"golang.org/x/net/context"

	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
)

type BatteryService struct {
}

type UpdateReq struct {
	Device Device
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
	bs := r.Device.Batteries
	// TODO: It reads and writes N times. Should be optimized.
	for _, b := range bs {
		key := historyKey(u, d.DeviceId, b.Time, c)
		h, err := getHistory(key, c)
		if err != nil {
			log.Debugf(c, "%#v", err)
			return err
		}
		h.Batteries = append(h.Batteries, b)
		_, err = datastore.Put(c, key, h)
		if err != nil {
			log.Debugf(c, "%#v", err)
			return err
		}
	}
	latest := bs[len(bs)-1]

	if latest.Battery <= 15 && !latest.Charging {
		// Notify
		err = notifyLowBattery(c, d.DeviceName, latest.Battery, []string{keys.MyNexus5x})
		if err != nil {
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
