package main

import (
	"io"
	"log"
	"os/exec"
	"sync"
	"syscall"
)

type Player struct {
	cmd     *exec.Cmd
	mplayer io.Writer
	sync.Mutex
}

func (p *Player) isPlaying() bool {
	return p.cmd != nil
}

func (p *Player) Play(url string) error {
	p.Lock()
	defer p.Unlock()

	if !p.isPlaying() {
		p.cmd = exec.Command("mplayer", "-slave", "-quiet", url)
		var err error
		p.mplayer, err = p.cmd.StdinPipe()
		if err != nil {
			log.Println(err)
		}
		return p.cmd.Start()
	}

	return nil
}

func (p *Player) Mute() {
	p.Lock()
	defer p.Unlock()

	if p.isPlaying() {
		p.mplayer.Write([]byte("mute\n"))
	}
}

func (p *Player) Stop() error {
	p.Lock()
	defer p.Unlock()

	if p.isPlaying() {
		p.cmd.Process.Signal(syscall.SIGTERM)
		p.cmd.Wait()
		p.cmd = nil
		return nil
	}

	return nil
}
