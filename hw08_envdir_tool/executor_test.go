package main

import (
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	t.Run("empty cmd", func(t *testing.T) {
		expected := 1
		var cmd []string
		returnCode := RunCmd(cmd, Environment{})
		require.Equal(t, expected, returnCode)
	})

	t.Run("wrong cmd", func(t *testing.T) {
		expected := 1
		cmd := []string{""}
		returnCode := RunCmd(cmd, Environment{})
		require.Equal(t, expected, returnCode)
	})

	t.Run("cmd without flags", func(t *testing.T) {
		expected := 0
		cmd := []string{"ls"}
		returnCode := RunCmd(cmd, Environment{})
		require.Equal(t, expected, returnCode)
	})

	t.Run("cmd with flags", func(t *testing.T) {
		expected := 0
		cmd := []string{"ls", "-l", "-a"}
		returnCode := RunCmd(cmd, Environment{})
		require.Equal(t, expected, returnCode)
	})

	t.Run("cmd with exit code", func(t *testing.T) {
		content := []byte("#!/bin/bash\nexit 5")
		tmpFile, err := os.CreateTemp("/tmp", "example.")
		if err != nil {
			log.Fatal(err)
		}
		os.Chmod(tmpFile.Name(), os.FileMode(0777))
		defer os.Remove(tmpFile.Name())
		if _, err := tmpFile.Write(content); err != nil {
			log.Fatal(err)
		}
		if err := tmpFile.Close(); err != nil {
			log.Fatal(err)
		}
		expected := 5
		cmd := []string{tmpFile.Name()}
		returnCode := RunCmd(cmd, Environment{})
		require.Equal(t, expected, returnCode)
	})
}
