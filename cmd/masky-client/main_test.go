package main

import (
	"testing"

	"github.com/wkj9893/masky/internal/client"
	"github.com/wkj9893/masky/internal/log"
)

func Test_parseArgs(t *testing.T) {
	config = client.Config{}
	tests := [][]string{
		{"--port=3000", "--addr=127.0.0.1:3001", "--allowlan=false"},
		{"--mode=direct", "--log=info", "--allowlan=true"},
		{"--mode=rule", "--log=warn"},
		{"--mode=global", "--log=error"},
	}
	parseArgs(tests[0])
	if config.Port != 3000 || config.Addr != "127.0.0.1:3001" || config.AllowLan != false {
		t.FailNow()
	}
	parseArgs(tests[1])
	if config.Mode != client.DirectMode || config.LogLevel != log.InfoLevel || config.AllowLan != true {
		t.FailNow()
	}
	parseArgs(tests[2])
	if config.Mode != client.RuleMode || config.LogLevel != log.WarnLevel {
		t.FailNow()
	}
	parseArgs(tests[3])
	if config.Mode != client.GlobalMode || config.LogLevel != log.ErrorLevel {
		t.FailNow()
	}
}
