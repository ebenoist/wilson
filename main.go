package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/brentnd/go-snowboy"
	"github.com/gordonklaus/portaudio"
)

func main() {
	inputChannels := 1
	outputChannels := 0
	sampleRate := 16000
	framesPerBuffer := make([]int16, 1024)

	// initialize the audio recording interface
	err := portaudio.Initialize()
	if err != nil {
		fmt.Errorf("Error initialize audio interface: %s", err)
		return
	}
	defer portaudio.Terminate()

	fmt.Println("Got here")

	// open the sound input for the microphone
	stream, err := portaudio.OpenDefaultStream(
		inputChannels,
		outputChannels,
		float64(sampleRate),
		len(framesPerBuffer),
		framesPerBuffer,
	)

	if err != nil {
		log.Printf("Error open default audio stream: %s", err)
		return
	}
	defer stream.Close()
	fmt.Println("Got here 2")

	// open the snowboy detector
	d := snowboy.NewDetector(os.Args[1])
	defer d.Close()

	sound := NewRecorder(stream, framesPerBuffer)
	d.HandleFunc(snowboy.NewHotword(os.Args[2], 0.5), func(string) {
		fmt.Println("Start Recording")
		sound.StartRecording()
	})

	d.HandleSilenceFunc(1500*time.Millisecond, func(string) {
		fmt.Println("Silence detected")
		if sound.IsRecording() {
			fmt.Println("Stop Recording")
			b := sound.StopRecording()
			GetText(b)
		}
	})

	sr, nc, bd := d.AudioFormat()
	fmt.Printf("sample rate=%d, num channels=%d, bit depth=%d\n", sr, nc, bd)

	err = stream.Start()
	if err != nil {
		fmt.Errorf("Error on stream start: %s", err)
		return
	}

	d.ReadAndDetect(sound)
}
