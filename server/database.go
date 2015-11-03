package server

import (
	"errors"
	"time"

	"github.com/GoogleCloudPlatform/go-endpoints/endpoints"

	"appengine"
	"appengine/datastore"
	"appengine/user"
)

type User struct {
	UserId string
}

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

func userKey(u *User, c appengine.Context) *datastore.Key {
	return datastore.NewKey(c, "User", u.UserId, 0, nil)
}

func deviceKey(u *User, d string, c appengine.Context) *datastore.Key {
	uk := userKey(u, c)
	c.Debugf("%#v", uk)
	return datastore.NewKey(c, "Device", d, 0, uk)
}

func batteryKey(u *User, d string, t time.Time, c appengine.Context) *datastore.Key {
	dk := deviceKey(u, d, c)
	c.Debugf("%#v", dk)
	return datastore.NewKey(c, "Battery", "", t.Unix(), dk)
}

// getCurrentUser retrieves a user associated with the request.
// If there's no user (e.g. no auth info present in the request) returns
// an "unauthorized" error.
func getCurrentUser(c endpoints.Context) (*User, error) {
	u, err := endpoints.CurrentUser(c, scopes, audiences, clientIds)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, errors.New("Unauthorized: Please, sign in.")
	}
	c.Debugf("Current user: %#v", u)
	return toUser(u), nil
}

func toUser(u *user.User) *User {
	if u == nil {
		return nil
	}
	return &User{UserId: u.Email}
}
