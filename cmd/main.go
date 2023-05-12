package main

import (
	"fmt"
	"log"
	models "testForPosgrges"
	"testForPosgrges/config"
	"testForPosgrges/iternal/checkDir"
	"testForPosgrges/iternal/parser"
	"testForPosgrges/repository"
	"time"
)

func main() {
	fmt.Printf("Progmar is worked\n")
	viperConf, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}
	conf, err := config.ParseConfig(viperConf)
	if err != nil {
		log.Fatal(err)
	}
	repository.Connect, err = repository.InitPsqlDB(conf)
	if err != nil {
		log.Fatal(err)
	}

	models.Frequency = conf.System.Frequency

	err = repository.Connect.Ping()
	if err != nil {
		log.Fatal(err)
	}

	dirs, err := parser.Parser("test.txt")
	if err != nil {
		log.Fatal(err)
	}

	for i := range dirs {
		go endlessCycle(&dirs[i])
	}

	t := 1
	fmt.Scan(&t)
}

func endlessCycle(dir *models.Directory) {
	for {
		err := checkDir.CheckDir(dir)
		if err != nil {

			//log.Print(err)
			return
		}
		time.Sleep(time.Duration(models.Frequency) * time.Second)
	}
}
