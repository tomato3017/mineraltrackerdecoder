package jmapencoder

import (
	"fmt"
	"github.com/tomato3017/mineraltrackerdecoder/pkg/util"
)

type JMapWaypoint struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	Icon       string `json:"icon"`
	X          int    `json:"x"`
	Y          int    `json:"y"`
	Z          int    `json:"z"`
	R          int    `json:"r"`
	G          int    `json:"g"`
	B          int    `json:"b"`
	Enable     bool   `json:"enable"`
	Type       string `json:"type"`
	Origin     string `json:"origin"`
	Dimensions []int  `json:"dimensions"`
	Persistent bool   `json:"persistent"`
}

type JMapEncodeWriter struct {
	dirPath string
}

func NewJMapEncodeWriter(dirPath string) (JMapEncodeWriter, error) {
	if _, err := util.PathExists(dirPath); err != nil {
		return JMapEncodeWriter{}, fmt.Errorf("unable to open directory. Err: %w", err)
	}

	return JMapEncodeWriter{dirPath: dirPath}, nil
}
