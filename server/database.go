package server

import (
	"time"

	"appengine"
	"appengine/datastore"
	"appengine/user"
)

type Device struct {
	UserId         string // User.ID
	DeviceId       string // Unique ID for a device
	DeviceName     string // Display name.
	AlertThreshold int32  // 0 - 100.

	// Not stored.
	Batteries []Battery
}

type Battery struct {
	UserId   string    `json:"-"`
	DeviceId string    `json:"-"`
	Battery  int32     `json:"battery"` // 0 - 100.
	Time     time.Time `json:"time"`    // timestamp
}

func userKey(u *user.User, c appengine.Context) *datastore.Key {
	return datastore.NewKey(c, "User", u.Email, 0, nil)
}

func deviceKey(u *user.User, d string, c appengine.Context) *datastore.Key {
	uk := userKey(u, c)
	c.Debugf("%#v", uk)
	return datastore.NewKey(c, "Device", d, 0, uk)
}

func batteryKey(u *user.User, d string, t time.Time, c appengine.Context) *datastore.Key {
	dk := deviceKey(u, d, c)
	c.Debugf("%#v", dk)
	return datastore.NewKey(c, "Battery", "", t.Unix(), dk)
}
