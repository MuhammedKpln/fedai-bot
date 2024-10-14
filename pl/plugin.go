package main

import (
	"fmt"
	Cx "muhammedkpln/fedai/context"
	C "muhammedkpln/fedai/core"
	"muhammedkpln/fedai/shared"
	S "muhammedkpln/fedai/shared"
	"os"
	"path"
	"regexp"
	"strings"
	"time"
)

var Plugin S.Plugin = S.Plugin{
	Name:         "Plugin Manager",
	CommandRegex: regexp.MustCompile(`^\.plugin (add|del) (https:\/\/[a-zA-Z0-9\-\.]+(?:\/[a-zA-Z0-9\-._~:\/?#@!$&'()*+,;=]*)?\.so)$`),
	CommandInfo:  "Handles plugins",
	CommandFn:    Run,
}

func Run(message *shared.PluginRunOptions, Payload S.RegexpMatches) {
	if Payload.Action == nil && Payload.Payload == nil {
		go Cx.EditMessage(Cx.ErrorMessage("Action or Payload is missing!"), message)
		return
	}

	var action string = *Payload.Action

	switch action {
	case "add":
		go AddPlugin(message, Payload)
		break

	case "del":
		go DelPlugin(message, Payload)
		break

	}

}

func DelPlugin(message *shared.PluginRunOptions, Payload S.RegexpMatches) {
	splittedUrl := strings.Split(*Payload.Payload, "/")
	file := splittedUrl[len(splittedUrl)-1]
	filePath := path.Join("pl", file)
	err := C.DeletePlugin(filePath, file)

	if err != nil {
		go Cx.EditMessage(Cx.ErrorMessage(fmt.Sprintf("Could not uninstall the plugin... %s", err)), message)
		return
	}

	go Cx.EditMessage(Cx.InfoMessage(fmt.Sprintf("Deleted %s, restarting in 5 seconds...", file)), message)

	time.Sleep(5 * time.Second)
	os.Exit(0)
}

func AddPlugin(message *shared.PluginRunOptions, Payload S.RegexpMatches) {
	splittedUrl := strings.Split(*Payload.Payload, "/")
	file := splittedUrl[len(splittedUrl)-1]

	go Cx.EditMessage(Cx.InfoMessage(fmt.Sprintf("Downloading %s...", file)), message)

	err := C.AddPlugin(file, *Payload.Payload)

	if err != nil && err.Error() == "pl_exists" {
		go Cx.EditMessage(Cx.InfoMessage(fmt.Sprintf("%s: You already have this plugin installed.", file)), message)
		return
	}

	if err != nil {
		go Cx.EditMessage(Cx.ErrorMessage(fmt.Sprintf("Could not install the plugin... %s", err)), message)
		return
	}

	go Cx.EditMessage(Cx.SuccessMessage("Download Complete, restarting in 5 seconds..."), message)

	time.Sleep(5 * time.Second)
	os.Exit(0)
}
