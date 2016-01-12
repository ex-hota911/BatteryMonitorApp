package server

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	"golang.org/x/net/context"
	"google.golang.org/appengine/urlfetch"
)

// curl -X POST \
// -H "Content-Type:application/json" \
// -H "Authorization:key=AIzaSyBKdv2lEyRlASvsoT7KY2UK58F7sy_FwGk"
// -d '{"data":{"score":"3x1"}, "to": "https://gcm-http.googleapis.com/gcm/send"}' https://gcm-http.googleapis.com/gcm/send

type PostData struct {
	Notification Notification `json:"notification"`
	Data         Data         `json:"data"`
	To           string       `json:"to"`
}

type Notification struct {
	Body  string `json:"body"`
	Title string `json:"title"`
	Icon  string `json:"icon"`
}

type Data struct {
	Device string `json:"device"`
	Level  int    `json:"level"`
}

func notify(ctx context.Context) error {
	body, err := json.Marshal(PostData{
		Notification: Notification{},
		Data:         Data{},
		To:           myNexus5x,
	})

	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", "https://gcm-http.googleapis.com/gcm/send", bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "key="+gcmApplicationKey)

	c := urlfetch.Client(ctx)
	resp, err := c.Do(req)
	if err != nil {
		return err
	}

	log.Print(resp)

	return nil
}
