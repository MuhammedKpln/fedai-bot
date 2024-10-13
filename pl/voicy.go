package main

import (
	"encoding/json"
	"errors"
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

	if os.Getenv("WITAI_TOKEN") == "" {
		go Context.EditMessage(Context.ErrorMessage("Token ekle."), message)

		return
	}

	if message.IsQuoted && message.QuotedMessage.GetAudioMessage() != nil {
		go Context.EditMessage(Context.InfoMessage("Dinliyorum..."), message)

		err := DownloadAudio(message.QuotedMessage)
		if err != nil {
			go Context.EditMessage(Context.ErrorMessage("Yanlis biseyler oldu, olmadi be usta..."), message)
			Cleanup()

			return
		}
		stereoErr := StereoToMonoConverter()
		if stereoErr != nil {
			go Context.EditMessage(Context.ErrorMessage("Yanlis biseyler oldu, olmadi be usta..."), message)
			Cleanup()

			return
		}

		finalMessage, finalMessageErr := RecognizeAudio()

		if finalMessageErr != nil {
			go Context.EditMessage(Context.ErrorMessage("Yanlis biseyler oldu, olmadi be usta..."), message)
			Cleanup()

			return
		}

		m := S.If(finalMessage == nil, Context.InfoMessage("Birsey duyamadim, sanirim konusmayi daha ögrenememis..."), Context.SuccessMessage(string(*finalMessage)))

		go Context.EditMessage(m, message)
		Cleanup()

	} else {
		go Context.EditMessage(Context.ErrorMessage("Lutfen ses alintila."), message)
	}

}

func DownloadAudio(QuotedMessage *waE2E.Message) error {
	client := C.GetClient()

	bytes, err := client.DownloadAny(QuotedMessage)

	if err != nil {
		client.Log.Errorf(err.Error())
		return err
	}

	f, err := os.Create("selam-stereo.ogg")
	if err != nil {
		client.Log.Errorf(err.Error())
		return err
	}
	defer f.Close()

	_, err = f.Write(bytes)

	if err != nil {
		client.Log.Errorf(err.Error())
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

func RecognizeAudio() (*string, error) {
	client := &http.Client{}
	c := C.GetClient()
	selamMono, _ := os.Open("./selam-mono.ogg")
	req, err := http.NewRequest("POST", "https://api.wit.ai/dictation", selamMono)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("WITAI_TOKEN")))
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "audio/ogg")
	req.Header.Add("Transfer-Encoding", "chunked")

	resp, requestErr := client.Do(req)

	if requestErr != nil || resp.StatusCode != 200 {
		c.Log.Errorf("Voicy: Request Error - Token missin maybe? %s - %s", resp.StatusCode, err)
		return nil, errors.New("Voicy: Request Error - Token missin maybe?")
	}

	dec := json.NewDecoder(resp.Body)

	var finalMessage *string

	for dec.More() {
		var m types.WitAi
		err := dec.Decode(&m)
		if err != nil {
			c.Log.Errorf(m.Text)
			return nil, err
		}

		if m.IsFinal {
			c.Log.Debugf(m.Text)
			finalMessage = &m.Text
		}
	}

	return finalMessage, nil
}

func Cleanup() {
	go os.Remove("selam-mono.ogg")
	go os.Remove("selam-stereo.ogg")
}
