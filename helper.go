package main

import (
	"context"
	"fmt"
)

func ReadErrors(errorChan chan error, expected int, ctx context.Context) error {
	var err error
	var res error
	for i := 0; i < expected; i += 1 {
		err = <-errorChan
		if err != nil && res == nil {
			res = err
		}
	}
	return res
}

func errorf(format string, err error) error {
	if err == nil {
		return err
	}
	return fmt.Errorf(format, err)
}
