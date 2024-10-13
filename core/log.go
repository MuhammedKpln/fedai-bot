package core

import (
	waLog "go.mau.fi/whatsmeow/util/log"
)

var appLog = waLog.Stdout("APP", "INFO", true)

func info(msg string) {
	appLog.Infof(msg)
}

func error(msg string) {
	appLog.Errorf(msg)
}
