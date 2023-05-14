package bashCmd

import (
	"log"
	"os/exec"
	"strings"
)

func BashCmd(commands []string, logger *log.Logger) {
	for _, command := range commands {
		if logger != nil {
			logger.Print(command)
		}
		t := parseCmd(command)
		var cmd *exec.Cmd
		if len(t) == 1 {
			cmd = exec.Command(t[0])
		} else if len(t) > 1 {
			cmd = exec.Command(t[0], t[1:]...)
		}
		//err := cmd.Run()
		stdOut, err := cmd.CombinedOutput()
		if err != nil {
			if logger != nil {
				logger.Print(err)
			} else {
				log.Print(err)
			}
			break
		}
		if logger != nil {
			logger.Printf("\t%s", stdOut)
		}

	}
}

func parseCmd(cmd string) []string {
	res := strings.Split(cmd, " ")
	return res
}
