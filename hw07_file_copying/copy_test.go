package main

import (
	"bytes"
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	resultTmpFile, _ := os.CreateTemp("", "result_*.txt")
	resultTmpFile.Close()

	inputFilePath := "testdata/input.txt"

	tests := []struct {
		offset       int64
		limit        int64
		expectedFile string
	}{
		{offset: 0, limit: 0, expectedFile: "testdata/out_offset0_limit0.txt"},
		{offset: 0, limit: 10, expectedFile: "testdata/out_offset0_limit10.txt"},
		{offset: 0, limit: 1000, expectedFile: "testdata/out_offset0_limit1000.txt"},
		{offset: 100, limit: 1000, expectedFile: "testdata/out_offset100_limit1000.txt"},
		{offset: 6000, limit: 1000, expectedFile: "testdata/out_offset6000_limit1000.txt"},
	}

	for _, tc := range tests {
		t.Run(tc.expectedFile, func(t *testing.T) {
			_ = Copy(inputFilePath, resultTmpFile.Name(), tc.offset, tc.limit)

			expectedData, _ := os.ReadFile(tc.expectedFile)
			copiedData, _ := os.ReadFile(resultTmpFile.Name())

			if !bytes.Equal(expectedData, copiedData) {
				t.Error("Copied file content differs from original!")
			}
		})
	}
}

func TestEmptyFileCopy(t *testing.T) {
	resultTmpFile, _ := os.CreateTemp("", "result_*.txt")
	resultTmpFile.Close()

	emptyTmpFile, _ := os.CreateTemp("", "empty_*.txt")
	emptyTmpFile.Close()

	t.Run("copy empty file", func(t *testing.T) {
		_ = Copy(emptyTmpFile.Name(), resultTmpFile.Name(), 0, 0)

		expectedData, _ := os.ReadFile(emptyTmpFile.Name())
		copiedData, _ := os.ReadFile(resultTmpFile.Name())

		if !bytes.Equal(expectedData, copiedData) {
			t.Error("Copied file content differs from original!")
		}
	})
}

func TestErrorCopy(t *testing.T) {
	resultTmpFile, _ := os.CreateTemp("", "result_*.txt")
	resultTmpFile.Close()

	t.Run("ErrUnsupportedFile", func(t *testing.T) {
		err := Copy("/dev/urandom", resultTmpFile.Name(), 0, 0)

		require.Truef(t, errors.Is(err, ErrUnsupportedFile), "actual error %q", err)
	})

	inputFilePath := "testdata/input.txt"
	t.Run("ErrOffsetExceedsFileSize", func(t *testing.T) {
		err := Copy(inputFilePath, resultTmpFile.Name(), 1000000, 0)

		require.Truef(t, errors.Is(err, ErrOffsetExceedsFileSize), "actual error %q", err)
	})
}
