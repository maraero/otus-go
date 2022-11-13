package main

import (
	"errors"
	"io"
	"math"
	"os"

	"github.com/cheggaaa/pb/v3"
)

var ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
var ErrSrcFileIsNotSpecified = errors.New("source file is not specified")
var ErrDstFileIsNotSpecified = errors.New("destination file is not specified")
var ErrSrcDirectory = errors.New("can not copy directory")
var ErrNegativeLimit = errors.New("limit can not be negative")

func getSrcFile(fpath string) (*os.File, error) {
	if fpath == "" {
		return nil, ErrSrcFileIsNotSpecified
	}
	file, err := os.Open(fpath)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func getDstFile(fpath string) (*os.File, error) {
	if fpath == "" {
		return nil, ErrDstFileIsNotSpecified
	}
	file, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE, os.FileMode(0o755))
	if err != nil {
		return nil, err
	}
	return file, nil
}

func getSrcFileSize(file *os.File) (int64, error) {
	fi, err := file.Stat()
	if err != nil {
		return 0, err
	}
	if fi.IsDir() {
		return 0, ErrSrcDirectory
	}
	filesize := fi.Size()
	return filesize, nil
}

func validateOffset(offset int64, filesize int64) error {
	if math.Abs((float64(offset))) > float64(filesize) {
		return ErrOffsetExceedsFileSize
	}
	return nil
}

func validateLimit(limit int64) error {
	if limit < 0 {
		return ErrNegativeLimit
	}
	return nil
}

func setReadPointerInFile(file *os.File, offset int64) {
	switch {
	case offset > 0:
		file.Seek(offset, io.SeekCurrent)
	case offset < 0:
		file.Seek(offset, io.SeekEnd)
	default:
		file.Seek(offset, io.SeekStart)
	}
}

func getReaderLimit(filesize int64, limit int64, offset int64) int64 {
	var available int64

	switch {
	case limit == 0 && offset >= 0:
		return filesize - offset
	case limit == 0 && offset < 0:
		return -offset
	case limit > 0 && offset >= 0:
		available = filesize - offset
	default: // limit > 0 && offset < 0
		available = -offset
	}

	if limit > available {
		return available
	}
	return limit
}

func Copy(fromPath, toPath string, offset, limit int64) error {
	srcFile, err := getSrcFile(fromPath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := getDstFile(toPath)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	srcFilesize, err := getSrcFileSize(srcFile)
	if err != nil {
		return err
	}

	err = validateOffset(offset, srcFilesize)
	if err != nil {
		return err
	}

	err = validateLimit(limit)
	if err != nil {
		return err
	}

	setReadPointerInFile(srcFile, offset)
	readerLimit := getReaderLimit(srcFilesize, limit, offset)

	reader := io.LimitReader(srcFile, readerLimit)
	writer := io.Writer(dstFile)
	bar := pb.Full.Start64(readerLimit)
	barReader := bar.NewProxyReader(reader)
	io.Copy(writer, barReader)
	bar.Finish()

	return nil
}
