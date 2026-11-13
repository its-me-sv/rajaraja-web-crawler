package utils

import (
	"errors"
	"strconv"
)

func ParseArgs(args []string) (baseURL string, maxConcurrency, maxPages int, err error) {
	if len(args) < 2 {
		err = errors.New("insufficient no. of arguments")
		return
	}
	if len(args) > 4 {
		err = errors.New("too many arguments")
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
