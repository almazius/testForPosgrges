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
	err := openDir(dir)
	if err != nil {
		return err
	}

	err = saveHashForFiles(dir)
	if err != nil {
		return err
	}
	oldHashes, err := repository.GetHashFromRepository(dir)
	if err != nil {
		return err
	}
	if !compareHashes(*dir, oldHashes) {
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

// openDir получает все адреса фалов в папке
func openDir(dir *models.Directory) error {
	err := filepath.Walk(dir.Path,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				if checkRegex(path, dir.IncludeRegexp, dir.ExcludeRegexp) {
					dir.FileHash[path] = ""
				}
			}
			return nil
		})
	if err != nil {
		log.Print(err)
		return err
	}
	return err
}

// checkRegex проверяет название файлов на наличие exclude и include
func checkRegex(path string, include, exclude []string) bool {
	fileName := path[strings.LastIndex(path, "/")+1:]
	if exclude != nil {
		for i, expr := range exclude {
			if expr != "" {
				res, err := regexp.Match(expr, []byte(fileName))
				if err != nil {
					log.Print(err)
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

					log.Print(err)
					include[i] = ""
				} else if res == false {
					return false
				}
			}
		}
	}
	return true
}

// compareHashes сравнивает старую и новые хеш-суммы файлов. Если они равны true, иначе false
func compareHashes(dir models.Directory, old map[string]string) bool {
	res := true
	for file, fileHash := range dir.FileHash {
		if oldHash, _ := old[file]; oldHash != fileHash {
			res = false
		}
	}
	return res
}

// saveHashForFiles сохраняет хеш-сумму каждого файла в директории
func saveHashForFiles(dir *models.Directory) error {
	for file, _ := range dir.FileHash {
		thread, err := os.Open(file)
		if err != nil {
			log.Print(err)
			return err
		}
		dir.FileHash[file] = fmt.Sprintf("%x", getHash(thread))
		err = thread.Close()
		if err != nil {
			log.Print(err)
			return err
		}
	}
	return nil
}

// getHash вычисляет хеш-сумму файла
func getHash(file io.Reader) hash.Hash {
	scanner := bufio.NewScanner(file)
	fileHash := md5.New()
	for scanner.Scan() {
		s := scanner.Text()
		_, err := io.WriteString(fileHash, s)
		if err != nil {
			log.Print(err)
			return nil
		}
	}
	return fileHash
}
