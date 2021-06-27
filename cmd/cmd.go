package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/subchen/go-log"
	"github.com/tomato3017/mineraltrackerdecoder/pkg/jmapencoder"
	"github.com/tomato3017/mineraltrackerdecoder/pkg/mtdecoder"
	"github.com/tomato3017/mineraltrackerdecoder/pkg/util"
	"io"
	"os"
	"path"
	"path/filepath"
)

var options struct {
	debug              bool
	outputCoords       bool
	minDistance        int
	noFilterDistance   bool
	exportToJourneymap bool
	journeymapDir      string
}

var jmapvars struct {
	resolvedPath string
}

var (
	rootCmd = &cobra.Command{
		Use:   "parse",
		Short: "Parse Mineral Tracker files",
		Long:  "Parses the mineral tracker file passed in.",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if options.debug {
				log.Default.Level = log.DEBUG
			}

			if err := checkFlags(); err != nil {
				log.Error(err.Error())
				log.Fatal(cmd.Help().Error())
			}
			runCMD(args[0])
		},
	}
)

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	p := rootCmd.PersistentFlags()
	p.BoolVar(&options.debug, "debug", false, "Set Debug mode")
	p.BoolVar(&options.outputCoords, "output", true, "Toggle if we output the coords to stdout")

	p.BoolVarP(&options.noFilterDistance, "nofilterdistance", "n", false, "Disables filtering waypoints very close together")
	p.IntVarP(&options.minDistance, "mindistance", "d", 50, "Filter the waypoints, only showing ones greater then specified distance")

	p.BoolVar(&options.exportToJourneymap, "journeymapexport", false, "Enables journeymap export")
	p.StringVar(&options.journeymapDir, "journeymapdir", "", "Journeymap directory")
}

func checkFlags() error {
	if options.exportToJourneymap {
		if options.journeymapDir == "" {
			return fmt.Errorf("--journeymapdir not defined")
		}

		//resolve the path
		absPath, err := filepath.Abs(options.journeymapDir)
		if err != nil {
			return fmt.Errorf("unable to resolve journeymap dir path. Err: %w", err)
		}
		jmapvars.resolvedPath = absPath
		log.Debugf("journeymap resolved path is %s", jmapvars.resolvedPath)
	}

	return nil
}

//MT Decoder Runner
func runCMD(filename string) {
	log.Info("Running MT Decoder")

	log.Debug("Opening file")
	file, err := os.Open(path.Join(filename))
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}
	defer util.SafeClose(file)

	decoder, err := mtdecoder.NewMTDecoder(file)
	if err != nil {
		log.Panic(err.Error())
		os.Exit(1)
	}

	entries := make([]mtdecoder.MTEntry, 0)
	for {
		entry, err := decoder.Scan()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				log.Fatal(err.Error())
			}
		}
		entries = append(entries, entry)
	}

	if !options.noFilterDistance {
		entries = filterEntries(entries)
	}

	if options.outputCoords {
		for _, entry := range entries {
			fmt.Println(entry)
		}
	}

	if options.exportToJourneymap {
		jme, err := jmapencoder.NewJMapEncodeWriter(options.journeymapDir)
		if err != nil {
			log.Errorf("unable to encode jmap writer. Err: %w", err)
		}

		if err := jme.WriteOutJmapEntries(entries, true); err != nil {
			log.Fatal("unable to write out entries to jmap. Err: %w", err)
		}
		log.Info("All waypoints written to jm waypoint files!")
	}
}

func filterEntries(entries []mtdecoder.MTEntry) []mtdecoder.MTEntry {
	rtnEntries := make([]mtdecoder.MTEntry, 0)

	for _, entry := range entries {
		closeEntry, found := getCloseEntry(entry, rtnEntries, options.minDistance)
		if found {
			log.Debugf("Entry: %s is too close to Entry: %s. Discarding!", entry, closeEntry)
			continue
		}

		rtnEntries = append(rtnEntries, entry)
	}

	return rtnEntries

}

func getCloseEntry(entry mtdecoder.MTEntry, entryList []mtdecoder.MTEntry, distance int) (mtdecoder.MTEntry, bool) {
	for _, mtEntry := range entryList {
		if mtEntry.Name == entry.Name && mtEntry.DistanceTo(entry) < float64(distance) {
			return mtEntry, true
		}
	}

	return mtdecoder.MTEntry{}, false
}
