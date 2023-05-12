package checkDir

import (
	"bufio"
	"crypto/md5"
	"fmt"
	"hash"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	models "testForPosgrges"
	"testForPosgrges/iternal/bashCmd"
	"testForPosgrges/repository"
)

func CheckDir(dir *models.Directory) error {
	err := OpenDir(dir)
	if err != nil {
		return err
	}

	err = CalcHashForFiles(dir)
	if err != nil {
		return err
	}
	oldHashes, err := repository.GetHashFromRepository(dir)
	if err != nil {
		return err
	}
	if !CompareHashes(*dir, oldHashes) {
		// bashCmd
		fmt.Printf("change in %s\n", dir.Path)

		bashCmd.BashCmd(dir.Command, dir.LogThread)
		err := repository.SetHashFromRepository(dir)
		if err != nil {
			return err
		}
	}
	return nil
}

func OpenDir(dir *models.Directory) error { // получает все адреса фалов в папке
	err := filepath.Walk(dir.Path,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				if CheckRegex(path, dir.IncludeRegexp, dir.ExcludeRegexp, dir.LogThread) {
					dir.FileHash[path] = ""
				}
			}
			return nil
		})
	if err != nil {
		if dir.LogThread != nil {
			dir.LogThread.Print(err)
		} else {
			log.Print(err)
		}
		return err
	}

	//getHasgFromReposiyory()

	return err
}

func CheckRegex(path string, include, exclude []string, logger *log.Logger) bool {
	fileName := path[strings.LastIndex(path, "/")+1:]
	if exclude != nil {
		for i, expr := range exclude {
			if expr != "" {
				res, err := regexp.Match(expr, []byte(fileName)) // >>>>>
				if err != nil {
					if logger != nil {
						logger.Print(err)
					} else {
						log.Print(err)
					}
					exclude[i] = ""
				} else if res == true {
					return false
				}
			}
		}
	}
	if include != nil {
		for i, expr := range include {
			if expr != "" {
				res, err := regexp.Match(expr, []byte(fileName))
				if err != nil {
					if logger != nil {
						logger.Print(err)
					} else {
						log.Print(err)
					}
					include[i] = ""
				} else if res == false {
					return false
				}
			}
		}
	}
	return true
}

func CompareHashes(dir models.Directory, old map[string]string) bool { // true - equal, false = no equal
	res := true
	for file, hash := range dir.FileHash {
		if oldHash, _ := old[file]; oldHash != hash {
			if dir.Changed == nil {
				dir.Changed = make([]string, 0)
			}
			dir.Changed = append(dir.Changed, file)
			res = false
		}
	}
	return res
}

func CalcHashForFiles(dir *models.Directory) error {
	for file, _ := range dir.FileHash {
		thread, err := os.Open(file)
		if err != nil {
			if dir.LogThread != nil {
				dir.LogThread.Print(err)
			} else {
				log.Print(err)
			}
			return err
		}
		dir.FileHash[file] = fmt.Sprintf("%x", getHash(thread))
		err = thread.Close()
		if err != nil {
			if dir.LogThread != nil {
				dir.LogThread.Print(err)
			} else {
				log.Print(err)
			}
			return err
		}
	}
	return nil
}

func getHash(file io.Reader) hash.Hash { // get hash file
	scanner := bufio.NewScanner(file)
	fileHash := md5.New()
	for scanner.Scan() {
		s := scanner.Text()
		io.WriteString(fileHash, s)
		//fileHash.Sum([]byte(s))
	}
	return fileHash
}
