package main

import (
	"regexp"
)

type Command struct {
	Regex *regexp.Regexp
	Run   func(...string) (result string, err error)
}

var player = &Player{}

var Commands = []Command{
	Command{
		Regex: regexp.MustCompile("welcome.*reverb"),
		Run: func(...string) (string, error) {
			return "Thanks! It's great to be here today. How may I help you?", nil
		},
	},
	Command{
		Regex: regexp.MustCompile("weather"),
		Run: func(...string) (string, error) {
			return GetForecast()
		},
	},
	Command{
		Regex: regexp.MustCompile("how.*look"),
		Run: func(...string) (string, error) {
			return "you all look quite nice today", nil
		},
	},
	Command{
		Regex: regexp.MustCompile("hello"),
		Run: func(...string) (string, error) {
			return "oh well hello, how are you?", nil
		},
	},
	Command{
		Regex: regexp.MustCompile("play.*WBEZ"),
		Run: func(...string) (string, error) {
			go player.Play("http://wbez.streamguys1.com:80/Pledge_Free.mp3")
			return "Playing WBEZ", nil
		},
	},
	Command{
		Regex: regexp.MustCompile("stop"),
		Run: func(...string) (string, error) {
			player.Stop()
			return "Okay", nil
		},
	},
}
