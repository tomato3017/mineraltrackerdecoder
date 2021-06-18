package util

import (
	"github.com/subchen/go-log"
	"io"
)

func SafeClose(closer io.Closer) {
	if err := closer.Close(); err != nil {
		log.Errorf("Unable to close stream. Err: %s", err.Error())
	}
}
