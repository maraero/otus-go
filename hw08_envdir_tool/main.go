package main

import (
	"log"
	"os"
)

const minArgs = 3

func main() {
	args := os.Args
	if len(args) < minArgs {
		log.Fatal("A minimum of three arguments are needed. Example: /path/to/env/dir command arg1 arg2")
	}

	envDir := args[1]
	cmd := args[2:]

	envs, err := ReadDir(envDir)
	if err != nil {
		log.Fatal(err)
	}

	returnCode := RunCmd(cmd, envs)
	os.Exit(returnCode)
}
