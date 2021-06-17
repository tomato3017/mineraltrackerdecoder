package main

import (
	"github.com/subchen/go-log"
	"mtdecoder/cmd"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		log.Error(err)
	}
}
