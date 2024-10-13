package main

import (
	"fmt"
	Context "muhammedkpln/fedai/context"
	M "muhammedkpln/fedai/modules"
	S "muhammedkpln/fedai/shared"
	"regexp"
)

var Plugin S.Plugin = S.Plugin{
	Name:         "Help",
	CommandRegex: regexp.MustCompile(".help"),
	CommandInfo:  "Help",
	CommandFn:    Run,
}

func Run(message *S.PluginRunOptions) {
	var helpString string

	for _, plugin := range M.LoadedPlugins {
		helpString += fmt.Sprintf("_*%s*_ \n *About*: %s \n Command: *%s* \n\n", plugin.Name, plugin.CommandInfo, plugin.CommandRegex.String())
	}

	go Context.EditMessage(helpString, message)

}
