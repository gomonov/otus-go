package main

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	filesEntry, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	env := make(Environment)

	for _, fileEntry := range filesEntry {
		if !fileEntry.Type().IsRegular() || strings.Contains(fileEntry.Name(), "=") {
			continue
		}

		filePath := filepath.Join(dir, fileEntry.Name())
		file, err := os.Open(filePath)
		if err != nil {
			return nil, err
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		if scanner.Scan() {
			line := scanner.Text()
			line = strings.TrimRight(line, " \t")
			line = strings.ReplaceAll(line, "\x00", "\n")
			env[fileEntry.Name()] = EnvValue{
				Value:      line,
				NeedRemove: false,
			}
		} else {
			env[fileEntry.Name()] = EnvValue{
				Value:      "",
				NeedRemove: true,
			}
		}
	}
	return env, nil
}
