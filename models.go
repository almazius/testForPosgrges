package models

import (
	"log"
)

type Directory struct {
	Path          string
	LogFile       string
	LogThread     *log.Logger
	Command       []string
	ExcludeRegexp []string
	IncludeRegexp []string
	FileHash      map[string]string
}

var Frequency float64
