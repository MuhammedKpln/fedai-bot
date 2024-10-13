package core

import waLog "go.mau.fi/whatsmeow/util/log"

func AppLog() waLog.Logger {
	var appLog = waLog.Stdout("APP", "", true)

	return appLog
}
