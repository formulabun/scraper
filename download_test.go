package main

import (
	"net/url"
	"testing"
)

func TestUrl(t *testing.T) {
	path, err := url.JoinPath("", "file")
	if err != nil {
		t.Errorf("Did expect to be able to join empty path")
	}
	res, err := url.Parse(path)
	if err != nil {
		t.Errorf("Did expect to be able to parse only path")
	}
	stringUrl := res.String()
	if stringUrl != "file" {
		t.Errorf("res.String() = %s but expected %s", res.String(), "file")
	}
}
