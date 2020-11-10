package server

import (
	"fmt"
	"github.com/sirupsen/logrus"
)

type ActionWorkRequest struct {
	DCC string `json:"dcc"`
}

func LogArgs(event *EventLogEntry) {
	e := fmt.Sprintf("%#v", event)
	logrus.Info(e)
}

func IngestFile(event *EventLogEntry) {
	if event.EventType != "Slingshot_WorkFile_Created" {
		return
	}

}