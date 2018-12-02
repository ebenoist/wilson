package main

import (
	"regexp"
)

type Command struct {
	Regex *regexp.Regexp
	Run   func(...string) (result string, err error)
}

var Commands = []Command{
	Command{
		Regex: regexp.MustCompile(".*radio"),
		Run: func(...string) (string, error) {
			return "playing wbez", nil
		},
	},
	Command{
		Regex: regexp.MustCompile(`(?i)weather`),
		Run: func(...string) (string, error) {
			return "the weather is great", nil
		},
	},
	Command{
		Regex: regexp.MustCompile(`(?i)jill`),
		Run: func(...string) (string, error) {
			return "hi jill-- want to play with me?", nil
		},
	},
	Command{
		Regex: regexp.MustCompile("stop"),
		Run: func(...string) (string, error) {
			return "stopping", nil
		},
	},
}
