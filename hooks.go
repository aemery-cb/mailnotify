package main

import (
	"os"
	"os/exec"
)

func ExecuteHook(hookCommand string) error {
	commandsplt := append([]string{"sh", "-c"}, hookCommand)
	cmd := exec.Command(commandsplt[0], commandsplt[1:]...)
	cmd.Stdout = os.Stdout
	return cmd.Run()

}
