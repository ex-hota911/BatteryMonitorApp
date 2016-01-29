package server

import (
	"errors"
	"time"

	"github.com/GoogleCloudPlatform/go-endpoints/endpoints"
	"golang.org/x/net/context"

	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/user"
)

const (
	DATE_FORMAT = "20060102"
)

type User struct {
	UserId string `datastore:"-"` // User.ID
}

type Device struct {
	UserId         string `datastore:"-"` // User.ID
	DeviceId       string `datastore:"-"` // Unique ID for a device
	DeviceName     string // Display name.
	AlertThreshold int32  // 0 - 100.

	// For API
	Batteries []Battery `database:"-"`
}

type Battery struct {
	Time    time.Time `json:"time"`    // timestamp
	Battery int32     `json:"battery"` // 0 - 100.
}

type History struct {
	Batteries []Battery
}

func userKey(u *User, c context.Context) *datastore.Key {
	return datastore.NewKey(c, "User", u.UserId, 0, nil)
}

func deviceKey(u *User, d string, c context.Context) *datastore.Key {
	uk := userKey(u, c)
	log.Debugf(c, "%#v", uk)
	return datastore.NewKey(c, "Device", d, 0, uk)
}

func batteryKey(u *User, d string, t time.Time, c context.Context) *datastore.Key {
	dk := deviceKey(u, d, c)
	log.Debugf(c, "%#v", dk)
	return datastore.NewKey(c, "Battery", "", t.Unix(), dk)
}

func historyKey(u *User, d string, t time.Time, c context.Context) *datastore.Key {
	dk := deviceKey(u, d, c)
	return datastore.NewKey(c, "History", toDate(t), 0, dk)
}

func getHistory(key *datastore.Key, c context.Context) (*History, error) {
	h := new(History)
	err := datastore.Get(c, key, h)
	if err == datastore.ErrNoSuchEntity {
		err = nil
	}
	return h, err
}

// toDate is the StringID of historyKey.
func toDate(t time.Time) string {
	return t.UTC().Format(DATE_FORMAT)
}

func populateKey(k *datastore.Key, b *Battery) {
	b.Time = time.Unix(k.IntID(), 0)
}

// getCurrentUser retrieves a user associated with the request.
// If there's no user (e.g. no auth info present in the request) returns
// an "unauthorized" error.
func getCurrentUser(c context.Context) (*User, error) {
	u, err := endpoints.CurrentUser(c, scopes, audiences, clientIds)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, errors.New("Unauthorized: Please, sign in.")
	}
	log.Debugf(c, "Current user: %#v", u)
	return toUser(u), nil
}

func toUser(u *user.User) *User {
	if u == nil {
		return nil
	}
	return &User{UserId: u.Email}
}
