package main

import (
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"sync"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

type ProgressWriter struct {
	writer  io.Writer
	current int64
	total   int64
	percent uint8
	ch      chan uint8
}

func (pw *ProgressWriter) Write(p []byte) (n int, err error) {
	n, err = pw.writer.Write(p)

	pw.current += int64(n)
	percent := uint8(math.Round(float64(pw.current) / float64(pw.total) * 100))
	if pw.percent != percent {
		pw.percent = percent
		pw.ch <- pw.percent
	}

	return
}

func Copy(fromPath, toPath string, offset, limit int64) error {
	var dst, src *os.File
	var info os.FileInfo
	var err error

	src, err = os.Open(fromPath)
	if err != nil {
		return err
	}
	defer src.Close()

	info, err = src.Stat()
	if err != nil {
		return err
	}

	size := info.Size()
	if offset > size {
		return ErrOffsetExceedsFileSize
	}
	if !info.Mode().IsRegular() {
		return ErrUnsupportedFile
	}

	if _, err = src.Seek(offset, io.SeekStart); err != nil {
		return err
	}

	if limit == 0 {
		limit = size - offset
	} else {
		limit = min(limit, size-offset)
	}

	limitReader := io.LimitReader(src, limit)

	dst, err = os.Create(toPath)
	if err != nil {
		return err
	}
	defer dst.Close()

	pbChan := make(chan uint8)

	progressDst := &ProgressWriter{
		writer: dst,
		total:  limit,
		ch:     pbChan,
	}

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func(ch chan uint8) {
		for percent := range ch {
			fmt.Printf("\rProgress: %3d%%", percent)
		}

		fmt.Println("\rProgress: 100% - done!")
		wg.Done()
	}(pbChan)

	if _, err := io.Copy(progressDst, limitReader); err != nil {
		return err
	}

	close(pbChan)

	wg.Wait()

	return nil
}
