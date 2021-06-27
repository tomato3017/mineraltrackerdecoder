package mtdecoder

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/tomato3017/mineraltrackerdecoder/pkg/common"
	"io"
	"math"
)

type MTEntry struct {
	Coords common.CoordXYZ
	Name   string
}

func (e MTEntry) String() string {
	return fmt.Sprintf("%s: %d,%d", e.Name, e.Coords.X, e.Coords.Z)
}

func (e MTEntry) DistanceTo(entry MTEntry) float64 {
	xdiff := math.Abs(float64(e.Coords.X - entry.Coords.X))
	zdiff := math.Abs(float64(e.Coords.Z - entry.Coords.Z))

	//Calculate the hypotenuse
	distance := math.Hypot(xdiff, zdiff)

	return distance
}

func convertHexToCoord(bytes []byte) int32 {
	return int32(binary.BigEndian.Uint32(bytes))
}

func GetMTEntryFromBytes(entryBytes []byte) (MTEntry, error) {
	//First validate the header
	if !HasValidEntryHeader(entryBytes) {
		return MTEntry{}, errors.New("no valid header found")
	}

	dataBytes := entryBytes[len(byteEntryStart):]

	xCoord := dataBytes[:4]
	zCoord := dataBytes[4:8]
	name := dataBytes[10:]

	return MTEntry{
		Coords: common.CoordXYZ{
			X: convertHexToCoord(xCoord),
			Y: 144,
			Z: convertHexToCoord(zCoord),
		},
		Name: string(name),
	}, nil
}

func GetMTEntryBytes(data []byte) ([]byte, error) {
	headerpos := GetEntryHeaderPos(0, data)
	if headerpos == -1 {
		return nil, io.EOF
	}

	//Ok get the next entry header so we know where to stop
	endPos := GetEntryHeaderPos(len(byteEntryStart), data)
	if endPos == -1 {
		//No more headers
		return data[headerpos:], io.EOF
	}

	rtnData := data[headerpos:endPos]
	return rtnData, nil
}

func HasValidEntryHeader(data []byte) bool {
	headerpos := GetEntryHeaderPos(0, data)
	return headerpos == 0
}

func GetEntryHeaderPos(startPos int, data []byte) int {
	checkData := data[startPos:]
	loc := bytes.Index(checkData, byteEntryStart)
	if loc == -1 {
		return -1
	}

	return loc + startPos
}
