package main

import (
	"github.com/subchen/go-log"
	"github.com/tomato3017/mineraltrackerdecoder/cmd"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		log.Error(err)
	}
}
