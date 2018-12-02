package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/brentnd/go-snowboy"
	"github.com/gordonklaus/portaudio"
)

const (
	inputChannels  = 1
	outputChannels = 1
	sampleRate     = 16000
	silenceDelay   = 1500 * time.Millisecond
)

func main() {
	in := make([]int16, 1024)
	out := make([]int16, 1024)

	err := portaudio.Initialize()
	if err != nil {
		fmt.Errorf("Error initialize audio interface: %s", err)
		return
	}
	defer portaudio.Terminate()

	stream, err := portaudio.OpenDefaultStream(
		inputChannels,
		outputChannels,
		float64(sampleRate),
		len(in),
		in,
		out,
	)

	if err != nil {
		log.Printf("Error open default audio stream: %s", err)
		return
	}
	defer stream.Close()

	svc := NewService(out, stream)
	d := snowboy.NewDetector(os.Args[1])
	defer d.Close()

	sound := NewRecorder(stream, in)
	d.HandleFunc(snowboy.NewHotword(os.Args[2], 0.5), func(string) {
		log.Print("start recording")
		sound.StartRecording()
	})

	d.HandleSilenceFunc(silenceDelay, func(string) {
		log.Println("silence detected")
		if sound.IsRecording() {
			log.Println("stop recording")
			b := sound.StopRecording()
			phrase := svc.GetTranscript(b)

			for _, cmd := range Commands {
				if cmd.Regex.Match([]byte(phrase)) {
					args := cmd.Regex.FindStringSubmatch(phrase)
					msg, _ := cmd.Run(args...)

					svc.Say(msg)
				}
			}
		}
	})

	err = stream.Start()
	if err != nil {
		fmt.Errorf("Error on stream start: %s", err)
		return
	}

	d.ReadAndDetect(sound)
}
