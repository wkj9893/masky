package main

import (
	"testing"

	"github.com/wkj9893/masky/internal/relay"
)

func Test_parseArgs(t *testing.T) {
	config = relay.Config{}
	tests := [][]string{
		{"--port=3000"},
		{"--addr=127.0.0.1:3000,:4000"},
		{"--addr=:5000, :6000"},
	}
	parseArgs(tests[0])
	if config.Port != 3000 {
		t.FailNow()
	}
	parseArgs(tests[1])
	if config.Addrs[0] != "127.0.0.1:3000" || config.Addrs[1] != ":4000" {
		t.FailNow()
	}
	parseArgs(tests[2])
	if config.Addrs[2] != ":5000" || config.Addrs[3] != ":6000" {
		t.FailNow()
	}
}
