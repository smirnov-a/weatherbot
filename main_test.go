package main

import (
	"flag"
	"os"
	"testing"
)

func TestMainFunction(t *testing.T) {

	crontabFile := flag.String("crontab", "nonexistent_file", "path to crontab file")
	err := checkCronTabFile(*crontabFile)
	if err == nil {
		t.Error("expected error, got nil")
	}
	if !os.IsNotExist(err) {
		t.Errorf("expected not-exist error, got %q", err.Error())
	}
}
