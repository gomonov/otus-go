package main

import (
	"os"
)

func main() {
	args := os.Args
	envDir := args[1]
	environment, err := ReadDir(envDir)
	if err != nil {
		panic(err)
	}

	code := RunCmd(args[2:], environment)
	os.Exit(code)
}
