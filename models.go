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
	Changed       []string
}

var Frequency float64
