package main

import (
	"bytes"
	"context"
	"encoding/binary"
	"io"
	"log"

	speech "cloud.google.com/go/speech/apiv1"
	texttospeech "cloud.google.com/go/texttospeech/apiv1"
	"github.com/gordonklaus/portaudio"
	speechpb "google.golang.org/genproto/googleapis/cloud/speech/v1"
	texttospeechpb "google.golang.org/genproto/googleapis/cloud/texttospeech/v1"
)

type SpeechService struct {
	speechClient *speech.Client
	textClient   *texttospeech.Client
	out          []int16
	stream       *portaudio.Stream
}

func NewService(out []int16, stream *portaudio.Stream) *SpeechService {
	ctx := context.Background()

	speechClient, err := speech.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	textClient, err := texttospeech.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	return &SpeechService{
		speechClient,
		textClient,
		out,
		stream,
	}
}

const (
	languageCode = "en-US"
)

func (s *SpeechService) Say(phrase string) {
	req := texttospeechpb.SynthesizeSpeechRequest{
		// Set the text input to be synthesized.
		Input: &texttospeechpb.SynthesisInput{
			InputSource: &texttospeechpb.SynthesisInput_Text{Text: phrase},
		},
		// Build the voice request, select the language code ("en-US") and the SSML
		// voice gender ("neutral").
		Voice: &texttospeechpb.VoiceSelectionParams{
			LanguageCode: languageCode,
			SsmlGender:   texttospeechpb.SsmlVoiceGender_NEUTRAL,
		},
		// Select the type of audio file you want returned.
		AudioConfig: &texttospeechpb.AudioConfig{
			AudioEncoding:   texttospeechpb.AudioEncoding_LINEAR16,
			SampleRateHertz: sampleRate,
		},
	}

	resp, err := s.textClient.SynthesizeSpeech(context.Background(), &req)
	if err != nil {
		log.Print(err)
	}

	r := bytes.NewReader(resp.AudioContent)

	for {
		audio := make([]byte, 2*len(s.out))
		_, err := r.Read(audio)
		if err == io.EOF {
			log.Print("done reading")
			break
		}

		if err != nil {
			log.Printf("failed to read %s", err)
			break
		}

		binary.Read(bytes.NewBuffer(audio), binary.LittleEndian, s.out)
		s.stream.Write()
	}

}

func (s *SpeechService) GetTranscript(data []byte) string {
	resp, err := s.speechClient.Recognize(context.Background(), &speechpb.RecognizeRequest{
		Config: &speechpb.RecognitionConfig{
			Encoding:        speechpb.RecognitionConfig_LINEAR16,
			SampleRateHertz: sampleRate,
			LanguageCode:    languageCode,
		},
		Audio: &speechpb.RecognitionAudio{
			AudioSource: &speechpb.RecognitionAudio_Content{Content: data},
		},
	})
	if err != nil {
		log.Printf("failed to recognize: %v", err)
	}

	var transcript string

	for _, result := range resp.Results {
		for _, alt := range result.Alternatives {
			transcript = alt.Transcript
			log.Printf("\"%v\" (confidence=%3f)\n", alt.Transcript, alt.Confidence)
		}
	}

	return transcript
}
