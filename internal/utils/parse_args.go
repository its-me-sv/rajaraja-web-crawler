package utils

import (
	"errors"
	"strconv"
)

const Usage = `usage: rajaraja-web-crawler <baseURL> <maxConcurrency> <maxPages>
example: rajaraja-web-crawler https://example.com 5 100`

var (
	ErrHelp  = errors.New("help requested")
	ErrUsage = errors.New("invalid no. of arguments")
)

func ParseArgs(args []string) (baseURL string, maxConcurrency, maxPages int, err error) {
	for _, arg := range args[1:] {
		if arg == "-h" || arg == "--help" {
			err = ErrHelp
			return
		}
	}
	if len(args) != 4 {
		err = ErrUsage
		return
	}

	baseURL = args[1]

	maxConcurrency, err = strconv.Atoi(args[2])
	if err != nil {
		err = errors.New(`invalid value for argument "maxConcurrency"`)
		return
	}

	maxPages, err = strconv.Atoi(args[3])
	if err != nil {
		err = errors.New(`invalid value for argument "maxPages"`)
		return
	}

	return
}
