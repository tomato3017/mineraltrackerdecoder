package mtdecoder

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"github.com/subchen/go-log"
	"io"
	"os"
)

type MTDecoder struct {
	fileData []byte
	pos      int
}

var bytePreheader = []byte{0x00, 0x00, 0x00}
var byteEntryStart = []byte{0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00}
var byteCoordSep = []byte{0x00, 0x13}

func NewMTDecoder(file *os.File) (MTDecoder, error) {
	bReader := bufio.NewReader(file)
	fileInfo, err := file.Stat()
	if err != nil {
		return MTDecoder{}, err
	}

	fileData := make([]byte, fileInfo.Size())
	count, err := bReader.Read(fileData)
	if err != nil {
		return MTDecoder{}, err
	}
	log.Debugf("Got %d bytes", count)

	mt := MTDecoder{
		fileData: fileData,
	}

	if !mt.ValidateHeader() {
		return MTDecoder{}, errors.New("Unable to validate header")
	}

	mt.MovePosToNextEntry(0)

	return mt, nil
}

func (m *MTDecoder) MovePosToNextEntry(offset int) {
	m.pos = GetEntryHeaderPos(m.pos+offset, m.fileData)
}

func (m *MTDecoder) ValidateHeader() bool {
	//Validate the first 3 bytes
	preheader := m.fileData[:3]
	check := bytes.Compare(preheader, bytePreheader)
	if check != 0 {
		return false
	}

	log.Debugf("Header check was %d", check)

	//Check that the first entry header exists
	if !HasValidEntryHeader(m.fileData[4:12]) {
		return false
	}

	return true
}

func (m *MTDecoder) Scan() (MTEntry, error) {
	//Get the next entry bytes
	if m.pos == -1 {
		return MTEntry{}, io.EOF
	}
	var err error
	entryBytes, err := GetMTEntryBytes(m.fileData[m.pos:])
	if err != nil {
		if errors.Is(err, io.EOF) {
			err = io.EOF
		} else {
			return MTEntry{}, fmt.Errorf("unable to find next MT Entry. Err: %w", err)
		}
	}

	//Convert bytes to entry
	nextEntry, err := GetMTEntryFromBytes(entryBytes)
	if err != nil {
		return MTEntry{}, fmt.Errorf("unable to find next MT Entry. Err: %w", err)
	}

	m.MovePosToNextEntry(1)

	return nextEntry, err
}

func (m *MTDecoder) getName(currentBytes []byte) ([]byte, int, error) {
	name := bytes.Index(currentBytes, byteEntryStart)
	if name == -1 {
		return currentBytes, len(currentBytes), io.EOF
	}

	return currentBytes[:name], name, nil
}

func (m *MTDecoder) getCoords(currentBytes []byte) ([]byte, int) {
	coords := bytes.Index(currentBytes, byteCoordSep)
	if coords == -1 {
		return nil, -1
	}
	return currentBytes[8:coords], coords
}
