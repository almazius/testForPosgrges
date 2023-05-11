package bashCmd

import (
	"log"
	"os/exec"
	"strings"
)

func BashCmd(commands []string) {
	for _, command := range commands {
		t := parseCmd(command)
		var cmd *exec.Cmd
		if len(t) == 1 {
			cmd = exec.Command(t[0])
		} else if len(t) > 1 {
			cmd = exec.Command(t[0], t[1:]...)
		}
		err := cmd.Run()
		if err != nil {
			log.Print(err)
			break
		}
	}
}

func parseCmd(cmd string) []string {
	res := strings.Split(cmd, " ")
	return res
}
