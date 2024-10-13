package main

import (
	"encoding/json"
	"fmt"
	Context "muhammedkpln/fedai/context"
	C "muhammedkpln/fedai/core"
	S "muhammedkpln/fedai/shared"
	"muhammedkpln/fedai/types"
	"net/http"
	"os"
	"regexp"

	ffmpeg "github.com/u2takey/ffmpeg-go"
	"go.mau.fi/whatsmeow/proto/waE2E"
)

var Plugin S.Plugin = S.Plugin{
	Name:         "Voicy",
	CommandRegex: regexp.MustCompile(".voicy"),
	CommandInfo:  "Translates voice to text",
	CommandFn:    Run,
}

func Run(message *S.PluginRunOptions) {

	if message.IsQuoted {
		if message.QuotedMessage.AudioMessage != nil {
			go Context.EditMessage(Context.InfoMessage("Listening..."), message)

			err := DownloadAudio(message.QuotedMessage)
			if err != nil {
				C.GetClient().Log.Errorf(err.Error())
				panic(err)
			}

			StereoToMonoConverter()
			finalMessage, err := RecognizeAudio()

			if err != nil {
				C.GetClient().Log.Errorf(err.Error())
				panic(err)
			}

			go Context.EditMessage(Context.SuccessMessage(finalMessage), message)
			go Cleanup()

		} else {
			go Context.EditMessage(Context.ErrorMessage("Audio message Required!"), message)
		}

	} else {
		go Context.EditMessage(Context.ErrorMessage("Quote a message"), message)
	}

}

func DownloadAudio(QuotedMessage *waE2E.Message) error {
	client := C.GetClient()

	bytes, err := client.DownloadAny(QuotedMessage)

	if err != nil {
		C.GetClient().Log.Errorf(err.Error())
		return err
	}

	f, err := os.Create("selam-stereo.ogg")
	if err != nil {
		C.GetClient().Log.Errorf(err.Error())
		return err
	}
	defer f.Close()

	_, err = f.Write(bytes)

	if err != nil {
		C.GetClient().Log.Errorf(err.Error())
		return err
	}

	return nil
}

func StereoToMonoConverter() error {
	errs := ffmpeg.Input("./selam-stereo.ogg").
		Output("./selam-mono.ogg", ffmpeg.KwArgs{"format": "ogg", "ac": 1, "ar": "44100"}).
		OverWriteOutput().ErrorToStdOut().Run()

	if errs != nil {
		C.GetClient().Log.Errorf(errs.Error())
		return errs
	}

	return nil
}

func RecognizeAudio() (string, error) {
	client := &http.Client{}

	selamMono, _ := os.Open("./selam-mono.ogg")
	req, err := http.NewRequest("POST", "https://api.wit.ai/dictation", selamMono)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("WITAI_TOKEN")))
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "audio/ogg")
	req.Header.Add("Transfer-Encoding", "chunked")

	resp, err := client.Do(req)

	if err != nil {
		C.GetClient().Log.Errorf(err.Error())

		return "", err
	}
	dec := json.NewDecoder(resp.Body)

	var finalMessage string

	for dec.More() {
		var m types.WitAi
		err := dec.Decode(&m)
		if err != nil {
			C.GetClient().Log.Errorf(err.Error())
			return "", err
		}

		if m.IsFinal {
			finalMessage = m.Text
		}
	}

	return finalMessage, nil
}

func Cleanup() {
	go os.Remove("selam-mono.ogg")
	go os.Remove("selam-stereo.ogg")
}
