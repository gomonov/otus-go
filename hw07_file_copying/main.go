package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	from, to      string
	limit, offset int64
)

func init() {
	flag.StringVar(&from, "from", "", "file to read from")
	flag.StringVar(&to, "to", "", "file to write to")
	flag.Int64Var(&limit, "limit", 0, "limit of bytes to copy")
	flag.Int64Var(&offset, "offset", 0, "offset in input file")
}

func main() {
	flag.Parse()

	flagError := false

	if from == "" {
		fmt.Fprintln(os.Stderr, "-from flag is required")
		flagError = true
	}

	if to == "" {
		fmt.Fprintln(os.Stderr, "-to flag is required")
		flagError = true
	}

	if limit < 0 {
		fmt.Fprintln(os.Stderr, "-limit flag value should be >= 0")
		flagError = true
	}

	if offset < 0 {
		fmt.Fprintln(os.Stderr, "-offset flag value should be >= 0")
		flagError = true
	}

	if flagError {
		os.Exit(64)
	}

	err := Copy(from, to, offset, limit)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
