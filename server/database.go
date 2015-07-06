package guestbook

import (
	"time"

	"appengine"
	"appengine/datastore"
	"appengine/user"
)

type Greeting struct {
	Author  string
	Content string
	Date    time.Time
}

type Device struct {
	UserId         string // User.ID
	DeviceId       string // Unique ID for a device
	DeviceName     string // Display name.
	AlertThreshold int32  // 0 - 100.
}

type Battery struct {
	UserId   string
	DeviceId string
	Battery  int32     // 0 - 100.
	Date     time.Time // timestamp
}

// guestbookKey returns the key used for all guestbook entries.
func guestbookKey(c appengine.Context) *datastore.Key {
	// The string "default_guestbook" here could be varied to have multiple guestbooks.
	return datastore.NewKey(c, "Guestbook", "default_guestbook", 0, nil)
}

func userKey(u user.User, c appengine.Context) *datastore.Key {
	// The string "default_guestbook" here could be varied to have multiple guestbooks.
	return datastore.NewKey(c, "User", u.ID, 0, nil)
}
