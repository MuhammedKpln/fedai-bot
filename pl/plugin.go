package main

import (
	"fmt"
	"io"
	Cx "muhammedkpln/fedai/context"
	C "muhammedkpln/fedai/core"
	"muhammedkpln/fedai/shared"
	S "muhammedkpln/fedai/shared"
	"net/http"
	"os"
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

	os.Remove(fmt.Sprintf("./pl/%s", file))

	database := C.GetDatabase()
	var plugin C.Plugin
	database.Where(C.Plugin{Name: file, Url: *Payload.Payload}).Attrs(C.Plugin{
		Name: file,
		Url:  *Payload.Payload,
	}).Delete(&plugin)

	go Cx.EditMessage(Cx.InfoMessage(fmt.Sprintf("Deleted %s, restarting in 5 seconds...", file)), message)

	time.Sleep(5 * time.Second)

	os.Exit(0)
}

func AddPlugin(message *shared.PluginRunOptions, Payload S.RegexpMatches) {
	splittedUrl := strings.Split(*Payload.Payload, "/")
	file := splittedUrl[len(splittedUrl)-1]

	go Cx.EditMessage(Cx.InfoMessage(fmt.Sprintf("Downloading %s...", file)), message)

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	resp, err := client.Get(*Payload.Payload)
	defer resp.Body.Close()

	if err != nil {
		panic(err)
	}

	go Cx.EditMessage(Cx.InfoMessage(fmt.Sprintf("Writing %s...", file)), message)

	out, err := os.Create(fmt.Sprintf("./pl/%s", file))
	defer out.Close()
	n, err := io.Copy(out, resp.Body)

	if err != nil {
		panic(err)
	}

	database := C.GetDatabase()
	var plugin C.Plugin
	database.Where(C.Plugin{Name: file, Url: *Payload.Payload}).Attrs(C.Plugin{
		Name: file,
		Url:  *Payload.Payload,
	}).FirstOrCreate(&plugin)

	// database.Create(C.Plugin{
	// 	Name: file,
	// 	Url:  *Payload.Payload,
	// })

	go Cx.EditMessage(Cx.SuccessMessage("Download Complete, restarting in 5 seconds..."), message)

	fmt.Printf("Written %s bytes", n)

	time.Sleep(5 * time.Second)

	os.Exit(0)
}
