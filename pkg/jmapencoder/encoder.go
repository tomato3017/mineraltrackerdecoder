package jmapencoder

import (
	"fmt"
	"github.com/subchen/go-log"
	"github.com/tomato3017/mineraltrackerdecoder/pkg/mtdecoder"
	"github.com/tomato3017/mineraltrackerdecoder/pkg/util"
	"io/fs"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

type JMapEncodeWriter struct {
	dirPath string
}

func (jm *JMapEncodeWriter) WriteOutJmapEntries(mtentries []mtdecoder.MTEntry, deleteOldEntries bool) error {
	fileEntries, err := ioutil.ReadDir(jm.dirPath)
	if err != nil {
		return fmt.Errorf("unable to list dir entries. Err: %w", err)
	}

	//delete old entries
	if deleteOldEntries {
		jm.deleteOldEntries(fileEntries)
	}

	jmWaypoints := getJMWaypoints(mtentries)
	for _, waypoint := range jmWaypoints {
		if err := jm.writeOutJmapEntry(waypoint); err != nil {
			log.Errorf("unable to write out waypoint %s. Err: %w Skipping! ", waypoint.Name, err)
			continue
		}
		log.Debugf("Successfully wrote out waypoint %s to file %s", waypoint.Name, waypoint.GetFilename())
	}

	return nil
}

func (jm *JMapEncodeWriter) writeOutJmapEntry(waypoint JMapWaypoint) error {
	waypointFile, err := os.Create(path.Join(jm.dirPath, waypoint.GetFilename()))
	if err != nil {
		return fmt.Errorf("unable to create waypoint file. Err: %w", err)
	}
	defer util.SafeClose(waypointFile)

	jsonWaypoint, err := waypoint.ToJSON()
	if err != nil {
		return fmt.Errorf("unable to json waypoint. Err: %w", err)
	}

	fileBytes, err := waypointFile.Write(jsonWaypoint)
	if err != nil {
		return fmt.Errorf("unable to write waypoint file. Err: %w", err)
	}

	log.Debugf("Wrote %d bytes", fileBytes)

	return nil
}

func (jm *JMapEncodeWriter) deleteOldEntries(fileEntries []fs.FileInfo) {
	for _, entry := range fileEntries {
		if strings.HasPrefix(entry.Name(), JmmtPrefix) && strings.HasSuffix(entry.Name(), ".json") {
			if err := os.Remove(path.Join(jm.dirPath, entry.Name())); err != nil {
				log.Errorf("Unable to delete file %s. Err: %w", entry.Name(), err)
			}
		}
	}
}

func getJMWaypoints(mtentries []mtdecoder.MTEntry) []JMapWaypoint {
	outJMEntries := make([]JMapWaypoint, 0)
	for _, mtentry := range mtentries {
		outJMEntries = append(outJMEntries, NewJMapWaypoint(mtentry.Name, mtentry.Coords))
	}

	return outJMEntries
}

func NewJMapEncodeWriter(dirPath string) (JMapEncodeWriter, error) {
	if _, err := util.PathExists(dirPath); err != nil {
		return JMapEncodeWriter{}, fmt.Errorf("unable to open directory. Err: %w", err)
	}

	return JMapEncodeWriter{dirPath: dirPath}, nil
}
