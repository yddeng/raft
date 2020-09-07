package raft

import "log"

var debug = 0

func dlog(format string, v ...interface{}) {
	if debug > 0 {
		log.Printf(format, v...)
	}
}
