package main

import (
	"fmt"
	"os/exec"
)

func main() {

	var (
		command        *exec.Cmd
		combinedOutput []byte
		err            error
	)

	command = exec.Command("/bin/bash", "-c", "echo 1 ; echo 2;")

	if combinedOutput, err = command.CombinedOutput(); err != nil {
		fmt.Print(err)
		return
	}

	fmt.Println(string(combinedOutput))

}
