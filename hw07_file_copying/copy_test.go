package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	INPUT  = "input"
	OUTPUT = "output"
)

func TestCopyWithExistingFiles(t *testing.T) {
	cases := []struct {
		offset int64
		limit  int64
	}{
		{offset: 0, limit: 0},
		{offset: 0, limit: 10},
		{offset: 0, limit: 1000},
		{offset: 0, limit: 10000},
		{offset: 100, limit: 1000},
		{offset: 6000, limit: 1000},
	}
	fromPath := "testdata/input.txt"

	for _, c := range cases {
		c := c
		expected := fmt.Sprintf("testdata/out_offset%d_limit%d.txt", c.offset, c.limit)
		toPath := fmt.Sprintf("testdata/output_offset%d_limit%d.txt", c.offset, c.limit)
		t.Run(expected, func(t *testing.T) {
			err := Copy(fromPath, toPath, c.offset, c.limit)
			defer os.Remove(toPath)
			require.NoError(t, err)

			toContent, err := os.ReadFile(toPath)
			require.NoError(t, err)
			expectedContent, err := os.ReadFile(expected)
			require.NoError(t, err)
			require.Equal(t, string(toContent), string(expectedContent))
		})
	}
}

func TestEmptySrcFilepath(t *testing.T) {
	err := Copy("", "test", 0, 0)
	require.Equal(t, err, ErrSrcFileIsNotSpecified)
}

func TestCanNotOpenSrcFile(t *testing.T) {
	err := Copy("missging_file.txt", "", 0, 0)
	require.Error(t, err)
}

func TestEmptyDstFilepath(t *testing.T) {
	err := Copy("testdata/input.txt", "", 0, 0)
	require.Equal(t, err, ErrDstFileIsNotSpecified)
}

func TestCanNotCopyDirectory(t *testing.T) {
	dstFile := createTmpFile(t, "", OUTPUT)
	defer os.Remove(dstFile.Name())
	err := Copy("testdata", dstFile.Name(), 0, 0)
	require.Equal(t, err, ErrSrcDirectory)
}

func TestOffsetExceedsFileSize(t *testing.T) {
	t.Run("positive", func(t *testing.T) {
		srcFile := createTmpFile(t, "a", INPUT)
		dstFile := createTmpFile(t, "a", OUTPUT)
		defer os.Remove(srcFile.Name())
		defer os.Remove(dstFile.Name())
		err := Copy(srcFile.Name(), dstFile.Name(), 2, 0)
		require.Equal(t, err, ErrOffsetExceedsFileSize)
	})

	t.Run("negative", func(t *testing.T) {
		fileIn, fileOut := createTestingFilePair(t, "a", "a")
		defer func() {
			defer os.Remove(fileIn.Name())
			defer os.Remove(fileOut.Name())
		}()
		err := Copy(fileIn.Name(), fileOut.Name(), -2, 0)
		require.Equal(t, err, ErrOffsetExceedsFileSize)
	})
}

func TestNegativeLimit(t *testing.T) {
	fileIn, fileOut := createTestingFilePair(t, "a", "a")
	defer func() {
		defer os.Remove(fileIn.Name())
		defer os.Remove(fileOut.Name())
	}()
	err := Copy(fileIn.Name(), fileOut.Name(), 0, -1)
	require.Equal(t, err, ErrNegativeLimit)
}

func TestOffsetLimitCombinations(t *testing.T) {
	cases := []struct {
		name   string
		limit  int64
		offset int64
		in     string
		out    string
	}{
		{name: "limit=0, offset=0", limit: 0, offset: 0, in: "1234567890", out: "1234567890"},
		{name: "limit=0, offset>0", limit: 0, offset: 5, in: "1234567890", out: "67890"},
		{name: "limit=0, offset<0", limit: 0, offset: -3, in: "1234567890", out: "890"},
		{name: "limit>0, offset=0", limit: 3, offset: 0, in: "1234567890", out: "123"},
		{name: "limit>0, offset>0", limit: 3, offset: 3, in: "1234567890", out: "456"},
		{name: "limit>0, offset<0", limit: 3, offset: -5, in: "1234567890", out: "678"},
	}

	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			fileIn, fileOut := createTestingFilePair(t, c.in, c.out)
			defer func() {
				defer os.Remove(fileIn.Name())
				defer os.Remove(fileOut.Name())
			}()
			err := Copy(fileIn.Name(), fileOut.Name(), c.offset, c.limit)
			require.NoError(t, err)
			outContent, err := os.ReadFile(fileOut.Name())
			require.NoError(t, err)
			require.Equal(t, c.out, string(outContent))
		})
	}
}

func createTmpFile(t *testing.T, content string, postfix string) *os.File {
	t.Helper()
	name := strings.ReplaceAll(t.Name()+"_"+postfix, "/", "_")
	f, err := os.CreateTemp("", name)
	if err != nil {
		log.Fatal(err)
	}
	if _, err = f.WriteString(content); err != nil {
		log.Fatal()
	}
	if err = f.Close(); err != nil {
		log.Fatal(err)
	}
	return f
}

func createTestingFilePair(t *testing.T, contentIn string, contentOut string) (fileIn *os.File, fileOut *os.File) {
	t.Helper()
	fileIn = createTmpFile(t, contentIn, INPUT)
	fileOut = createTmpFile(t, contentOut, OUTPUT)
	return fileIn, fileOut
}
