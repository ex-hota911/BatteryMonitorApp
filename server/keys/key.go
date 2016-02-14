package keys

import (
	"encoding/json"
	"io/ioutil"
)

var (
	WebClientId            = ""
	AndroidReleaseClientId = ""
	AndroidDebugClientId   = ""

	GcmApplicationKey = ""

	MyNexus5x = ""
)

type MySecretKeys struct {
	WebClientId            string `json:"web_client_id"`
	AndroidReleaseClientId string `json:"android_release_client_id"`
	AndroidDebugClientId   string `json:"android_debug_client_id"`
	GcmApplicationKey      string `json:"gcm_application_key"`
	MyNexus5x              string `json:"my_nexus_5x"`
}

func Init() error {
	b, err := ioutil.ReadFile("secret.json")
	if err != nil {
		b, err = ioutil.ReadFile("secret.dev.json")
		if err != nil {
			return err
		}
	}
	var keys MySecretKeys
	if err = json.Unmarshal(b, &keys); err != nil {
		return err
	}

	WebClientId = keys.WebClientId
	AndroidReleaseClientId = keys.AndroidReleaseClientId
	AndroidDebugClientId = keys.AndroidDebugClientId
	GcmApplicationKey = keys.GcmApplicationKey
	MyNexus5x = keys.MyNexus5x

	return nil
}
