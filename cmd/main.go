package main

import (
	"fmt"
	"log"
	"os"
	models "testForPosgrges"
	"testForPosgrges/config"
	"testForPosgrges/iternal/checkDir"
	"testForPosgrges/iternal/parser"
	"testForPosgrges/repository"
	"time"
)

func main() {
	fmt.Printf("Progmar is worked\n")
	if len(os.Args) != 2 {
		log.Fatal("you dont write config file on program arguments")
	}
	configFile := os.Args[1]
	//загрузка конфига из файла config.yml
	viperConf, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}
	conf, err := config.ParseConfig(viperConf)
	if err != nil {
		os.Exit(1)
	}

	// создание подключения к бд, оно будет активно все время работы программы
	repository.Connect, err = repository.InitPsqlDB(conf)
	if err != nil {
		log.Fatal(err)
	}

	models.Frequency = conf.System.Frequency

	dirs, err := parser.Parser(configFile)
	if err != nil {
		log.Fatal(err)
	}

	for i := range dirs {
		go endlessCycle(&dirs[i])
	}

	c := make(chan int)
	<-c
}

// endlessCycle Функция, реализующая бесконечный цикл
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
