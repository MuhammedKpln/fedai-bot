package main

import (
	"fmt"
	Cx "muhammedkpln/fedai/context"
	C "muhammedkpln/fedai/core"
	S "muhammedkpln/fedai/shared"
	"regexp"
)

var Plugin S.Plugin = S.Plugin{
	Name:         "Plugins",
	CommandRegex: regexp.MustCompile(".plugins$"),
	CommandInfo:  "List Plugins",
	CommandFn:    Run,
}

func Run(message *S.PluginRunOptions, payload S.RegexpMatches) {
	database := C.GetDatabase()
	var plugins []C.Plugin
	database.Take(&plugins)

	if len(plugins) < 1 {
		go Cx.EditMessage(Cx.InfoMessage("No plugins available."), message)

		return
	}

	var messages string

	for _, plugin := range plugins {
		messages += fmt.Sprintf("_*%s*_ \n *URL*: %s \n _*To Delete*_: `.plugin del %s` \n\n", plugin.Name, plugin.Url, plugin.Url)
	}

	go Cx.EditMessage(messages, message)

}
