package server

import (
	"encoding/json"
	"fmt"
	shotgun_consumer "github.com/ricksilliker/shotgun-event-server/shotgun-consumer"
	"github.com/sirupsen/logrus"
	"io/ioutil"
)

func processEvents(config *ServerConfig, lastEventID int) {
	//var plugins []string
	var continueToProcess bool

	url := "https://brazenanimation.shotgunstudio.com/api/v1"
	auth := shotgun_consumer.AuthenticateAsScript(url, config.ShotgunScriptName, config.ShotgunScriptKey)

	for ok := true; ok; ok = !continueToProcess {
		events, _ := shotgun_consumer.GetNewEvents(url, auth, lastEventID)
		for event := range events {
			//for _, plugin := range plugins {
			//	plugin.process(event)
			//}
			logrus.WithField("eventId", event.ID).Info("Processing event")
			storeEventID(event)
			lastEventID = event.ID
		}
	}
}

func storeEventID(event *shotgun_consumer.EventLogEntry) {
	data, err := json.Marshal(event)
	if err != nil {
		logrus.Error("failed to marshal ShotgunEventEntry data")
	}
	err = ioutil.WriteFile("lastProcessedEvent.txt", data, 0775)
	if err != nil {
		logrus.Error("failed to write out ShotgunEventEntry information")
	}
}

func main() {
	fmt.Println("Hello World")
	config, _ := NewServerConfigFromFile("server.json")
	var lastEventID = 3712449
	processEvents(config, lastEventID)
}