package tracker

import (
	"log"
	"time"
)

func TrackTime(msg string) (string, time.Time) {
	return msg, time.Now()
}

func Duration(msg string, start time.Time) {
	log.Printf("%v: %s\n", msg, time.Since(start))
}
