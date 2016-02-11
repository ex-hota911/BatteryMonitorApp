package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"keys"

	"golang.org/x/net/context"
	"google.golang.org/appengine/urlfetch"
)

// curl -X POST \
// -H "Content-Type:application/json" \
// -H "Authorization:key=AIzaSyBKdv2lEyRlASvsoT7KY2UK58F7sy_FwGk"
// -d '{"data":{"score":"3x1"}, "to": "https://gcm-http.googleapis.com/gcm/send"}' https://gcm-http.googleapis.com/gcm/send

type PostData struct {
	Notification    Notification `json:"notification,omitempty"`
	Data            Data         `json:"data,omitempty"`
	To              string       `json:"to,omitempty"`
	RegistrationIds []string     `json:"registration_ids,omitempty"`
	CollapseKey     string       `json:"collapse_key,omitempty"`
}

type Notification struct {
	Body  string `json:"body"`
	Title string `json:"title"`
	Icon  string `json:"icon"`
	Tag   string `json:"tag,omitempty"`
}

type Data struct {
}

func notifyLowBattery(c context.Context, device string, level int32, to []string) error {
	return notify(c, device+" battery low", fmt.Sprintf("%d%%", level), to)
}

func notify(ctx context.Context, title, body string, to []string) error {
	b, err := json.Marshal(PostData{
		Notification: Notification{
			Title: title,
			Body:  body,
			Icon:  "ic_battery_alert_black",
			Tag:   "battery_low",
		},
		RegistrationIds: to,
		CollapseKey:     "battery_low",
	})

	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", "https://gcm-http.googleapis.com/gcm/send", bytes.NewReader(b))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "key="+keys.GcmApplicationKey)

	c := urlfetch.Client(ctx)
	resp, err := c.Do(req)
	if err != nil {
		return err
	}

	log.Print(resp)

	return nil
}
