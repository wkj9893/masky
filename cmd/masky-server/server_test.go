package main

import (
	"testing"

	"github.com/wkj9893/masky/internal/log"
	"github.com/wkj9893/masky/internal/masky"
)

func Test_parseArgs(t *testing.T) {
	config = masky.ServerConfig{}
	tests := [][]string{
		{"--port=3000", "--password=123"},
		{"--log=info"},
		{"--log=warn"},
		{"--log=error"},
	}
	parseArgs(tests[0])
	if config.Port != 3000 || config.Password != "123" {
		t.FailNow()
	}
	parseArgs(tests[1])
	if config.LogLevel != log.InfoLevel {
		t.FailNow()
	}
	parseArgs(tests[2])
	if config.LogLevel != log.WarnLevel {
		t.FailNow()
	}
	parseArgs(tests[3])
	if config.LogLevel != log.ErrorLevel {
		t.FailNow()
	}
}
