package parser

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	models "testForPosgrges"
	"testForPosgrges/iternal/logger"
)

func Parser(path string) ([]models.Directory, error) { // смесь массива и мапы позволяют решить проблему колизий
	dirs := make(map[string]map[string][]string) // map dir - type_cmd - cmd
	file, err := os.Open(path)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Print(err)
		}
	}(file)

	scanner := bufio.NewScanner(file)
	newDir := ""
	typeCmd := "" // храним тип команды и записываем туда все, пока не встретим новую команду.
	for scanner.Scan() {
		str := scanner.Text()
		if strings.TrimSpace(str) != "" {
			if index := strings.Index(str, "- path:"); index != -1 {
				newDir = strings.Trim(str, "- path:")
				newDir = strings.TrimSpace(newDir)
				if _, isExist := dirs[newDir]; !isExist {
					dirs[newDir] = make(map[string][]string)
				}
			} else if index = strings.Index(str, "log_file:"); index != -1 {
				if newDir != "" { // Для каждого файла свой лог
					str += " "
					str = str[strings.Index(str, "log_file:")+9:]
					//str = strings.Trim(scanner.Text(), "log_file:")
					str = strings.TrimSpace(str)
					if _, isExist := dirs[newDir]["log_file"]; isExist == false {
						dirs[newDir]["log_file"] = make([]string, 1, 1)
						dirs[newDir]["log_file"][0] = str
					} else {
						err = errors.New("log file already exist")
						log.Print(err)
						//return nil, err // ????
						// Если задано несколько файлов для лога, используется только первый
					}
				}
			} else if strings.Contains(str, "commands:") || strings.Contains(str, "include_regexp:") || strings.Contains(str, "exclude_regexp:") {
				if strings.Contains(str, "commands:") {
					typeCmd = "commands"
				} else if strings.Contains(str, "include_regexp:") {
					typeCmd = "include_regexp"
				} else if strings.Contains(str, "exclude_regexp:") {
					typeCmd = "exclude_regexp"
				}
			} else if typeCmd != "" {
				str = scanner.Text() + " "
				str = scanner.Text()[strings.Index(str, "-")+1:]
				str = strings.TrimSpace(str)

				if _, isExist := dirs[newDir][typeCmd]; isExist {
					dirs[newDir][typeCmd] = append(dirs[newDir][typeCmd], str)
				} else {
					dirs[newDir][typeCmd] = make([]string, 0, 2)
					dirs[newDir][typeCmd] = append(dirs[newDir][typeCmd], str)
				}
			}
		}
	}

	dirs_on_struct := make([]models.Directory, 0, len(dirs))
	for pathDir, _ := range dirs {
		var tempDir models.Directory
		tempDir.Path = pathDir
		if _, isExist := dirs[pathDir]["log_file"]; isExist {
			tempDir.LogFile = dirs[pathDir]["log_file"][0]
			thread, err := logger.OpenLogThread(tempDir.LogFile)
			if err != nil {
				log.Print(err, " | log file can't create or open")
			} else {
				tempDir.LogThread = log.New(thread, fmt.Sprintf("%s: ", path), log.LstdFlags)
			}
		}

		if _, isExist := dirs[pathDir]["commands"]; isExist {
			tempDir.Command = make([]string, 0, len(dirs[pathDir]["commands"]))
			for _, cmd := range dirs[pathDir]["commands"] {
				tempDir.Command = append(tempDir.Command, cmd)
			}
		}

		if _, isExist := dirs[pathDir]["include_regexp"]; isExist {
			tempDir.IncludeRegexp = make([]string, 0, len(dirs[pathDir]["include_regexp"]))
			for _, cmd := range dirs[pathDir]["include_regexp"] {
				tempDir.IncludeRegexp = append(tempDir.IncludeRegexp, cmd)
			}
		}

		if _, isExist := dirs[pathDir]["exclude_regexp"]; isExist {
			tempDir.ExcludeRegexp = make([]string, 0, len(dirs[pathDir]["exclude_regexp"]))
			for _, cmd := range dirs[pathDir]["exclude_regexp"] {
				tempDir.ExcludeRegexp = append(tempDir.ExcludeRegexp, cmd)
			}
		}
		tempDir.FileHash = make(map[string]string)

		dirs_on_struct = append(dirs_on_struct, tempDir)
	}
	return dirs_on_struct, err
}
