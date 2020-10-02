package utils

import (
	"log"
	"os/exec"
)

func Ls(path string) string {
	cmd := exec.Command("ls", path)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
	return string(out)
}
