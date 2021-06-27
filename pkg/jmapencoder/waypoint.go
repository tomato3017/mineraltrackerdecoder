package jmapencoder

import (
	"encoding/json"
	"fmt"
	"github.com/tomato3017/mineraltrackerdecoder/pkg/common"
)

const JmmtPrefix = "MT-"

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

func (w *JMapWaypoint) ToJSON() ([]byte, error) {
	rtnJson, err := json.MarshalIndent(w, "", "    ")
	if err != nil {
		return nil, fmt.Errorf("unable to create json string for waypoint. Err: %w", err)
	}

	return rtnJson, nil
}

func (w *JMapWaypoint) GetFilename() string {
	return fmt.Sprintf("%s.json", w.Id)
}

func resolveWaypointName(name string) string {
	return fmt.Sprintf("%s%s", JmmtPrefix, name)
}

func NewJMapWaypoint(name string, coords common.CoordXYZ) JMapWaypoint {
	resolvedName := resolveWaypointName(name)
	rtnWaypoint := JMapWaypoint{
		Id:         fmt.Sprintf("%s_%s", resolvedName, coords),
		Name:       resolvedName,
		Icon:       "waypoint-normal.png",
		X:          int(coords.X),
		Y:          int(coords.Y),
		Z:          int(coords.Z),
		R:          255,
		G:          0,
		B:          0,
		Enable:     false,
		Type:       "Normal",
		Origin:     "",
		Dimensions: []int{0},
		Persistent: true,
	}

	return rtnWaypoint
}
