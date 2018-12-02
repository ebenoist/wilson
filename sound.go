package main

import (
	"bytes"
	"encoding/binary"
	"log"
	"sync"

	"github.com/ebenoist/wilson/memfile"
	"github.com/go-audio/audio"
	"github.com/go-audio/wav"
	"github.com/gordonklaus/portaudio"
)

type Recorder struct {
	stream        *portaudio.Stream
	data          []int16
	recording     bool
	lastRecording []int
	lock          *sync.RWMutex
}

func NewRecorder(
	stream *portaudio.Stream,
	data []int16,
) *Recorder {
	return &Recorder{
		stream,
		data,
		false,
		make([]int, 1024),
		&sync.RWMutex{},
	}
}

func (s *Recorder) StartRecording() {
	s.lock.Lock()
	s.recording = true
	s.lock.Unlock()
}

func (s *Recorder) StopRecording() []byte {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.recording = false

	f := s.save(s.lastRecording)

	s.lastRecording = []int{}

	return f
}

func (s *Recorder) save(d []int) []byte {
	f := &memfile.File{}
	enc := wav.NewEncoder(f, 16000, 16, 1, 1)

	a := &audio.IntBuffer{
		Format: &audio.Format{
			NumChannels: 1,
			SampleRate:  16000,
		},
		Data: d,
	}

	enc.Write(a)

	err := enc.Close()
	if err != nil {
		log.Printf("got an error on write - %s", err)
	}

	return f.Bytes()
}

func (s *Recorder) IsRecording() bool {
	s.lock.RLock()
	defer s.lock.RUnlock()

	return s.recording
}

// Read is the implementation of the io.Reader interface.
func (s *Recorder) Read(p []byte) (int, error) {
	s.stream.Read()

	buf := &bytes.Buffer{}
	for _, v := range s.data {
		binary.Write(buf, binary.LittleEndian, v)
		if s.IsRecording() {
			s.lastRecording = append(s.lastRecording, int(v))
		}
	}

	copy(p, buf.Bytes())
	return len(p), nil
}
