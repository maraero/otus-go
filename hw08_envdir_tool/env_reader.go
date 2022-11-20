package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"unicode"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

func validateFile(fInfo fs.FileInfo) error {
	fName := fInfo.Name()

	if strings.Contains(fName, "=") {
		return fmt.Errorf("filename %s contains \"=\"", fName)
	}

	if fInfo.Mode().IsRegular() != true {
		return fmt.Errorf("file %s is not regular", fName)
	}

	return nil
}

func readFileFirstLine(file *os.File) (string, error) {
	scanner := bufio.NewScanner(file)
	scanner.Scan()
	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("Can not read file %s: %w", file.Name(), err)
	}

	fContent := scanner.Bytes()
	fContent = bytes.Replace(fContent, []byte("\x00"), []byte("\n"), -1)
	fContent = bytes.TrimRightFunc(fContent, unicode.IsSpace)
	return string(fContent), nil
}

func getFirstLineFromFile(fInfo fs.FileInfo, dir string) (string, error) {
	fName := fInfo.Name()
	err := validateFile(fInfo)
	if err != nil {
		return "", fmt.Errorf("File %s is invalid: %w", fName, err)
	}

	if fInfo.Size() == 0 {
		return "", nil
	}

	fPath := filepath.Join(dir, fName)
	file, err := os.Open(fPath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	line, err := readFileFirstLine(file)
	if err != nil {
		return "", fmt.Errorf("Can not get first line from file %s: %w", fName, err)
	}

	return line, nil
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("Can not read directories: %w", err)
	}

	envs := Environment{}

	for _, file := range files {
		str, err := getFirstLineFromFile(file, dir)
		if err != nil {
			continue
		}

		ev := EnvValue{}

		if str == "" {
			ev.NeedRemove = true
		} else {
			ev.Value = str
		}

		envs[file.Name()] = ev
	}

	return envs, nil
}
