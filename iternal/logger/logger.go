package logger

import (
	"os"
)

// OpenLogThread открывает поток для хранения лога команд
func OpenLogThread(path string) (*os.File, error) {
	thread, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return nil, err
	} else {
		return thread, nil
	}
}
