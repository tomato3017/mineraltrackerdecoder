package mtdecoder

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math"
)

type MTEntry struct {
	CoordX int32
	CoordY int32
	CoordZ int32
	Name   string
}

func (e MTEntry) String() string {
	return fmt.Sprintf("%s: %d,%d\n", e.Name, e.CoordX, e.CoordZ)
}

func (e MTEntry) DistanceTo(entry MTEntry) float64 {
	xdiff := math.Abs(float64(e.CoordX - entry.CoordX))
	zdiff := math.Abs(float64(e.CoordZ - entry.CoordZ))

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
		CoordX: convertHexToCoord(xCoord),
		CoordY: 144,
		CoordZ: convertHexToCoord(zCoord),
		Name:   string(name),
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
	endPos = endPos + len(byteEntryStart)

	return data[headerpos:endPos], nil
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
