package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/subchen/go-log"
	"github.com/tomato3017/mineraltrackerdecoder/pkg/mtdecoder"
	"io"
	"os"
	"path"
)

var options struct {
	debug          bool
	outputCoords   bool
	minDistance    int
	filterDistance bool
}

var (
	rootCmd = &cobra.Command{
		Use:   "parse",
		Short: "Parse Mineral Tracker files",
		Long:  "Parses the mineral tracker file passed in.",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
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

	p.BoolVarP(&options.filterDistance, "filterdistance", "f", false, "Filters waypoints very close together")
	p.IntVarP(&options.minDistance, "mindistance", "d", 50, "Filter the waypoints, only showing ones greater then specified distance")
}

//MT Decoder Runner
func runCMD(filename string) {
	if options.debug {
		log.Default.Level = log.DEBUG
	}
	log.Info("Running MT Decoder")

	log.Debug("Opening file")
	file, err := os.Open(path.Join(filename))
	if err != nil {
		log.Panic(err.Error())
	}
	defer file.Close()

	decoder, err := mtdecoder.NewMTDecoder(file)
	if err != nil {
		log.Panic(err.Error())
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

	if options.filterDistance {
		entries = filterEntries(entries)
	}

	if options.outputCoords {
		for _, entry := range entries {
			fmt.Println(entry)
		}
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

//func PrintEntry(entry mtdecoder.MTEntry) {
//	str := []byte(fmt.Sprintf("%s: %d,%d\n", entry.Name, entry.CoordX, entry.CoordZ))
//	var _ = str
//	fmt.Printf("%s: %d,%d\n", entry.Name, entry.CoordX, entry.CoordZ)
//}
