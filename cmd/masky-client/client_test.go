package main

import (
	"fmt"
	"testing"

	"github.com/wkj9893/masky/internal/log"
	"github.com/wkj9893/masky/internal/masky"
)

func Test_parseArgs(t *testing.T) {
	config = masky.ClientConfig{}
	tests := [][]string{
		{"--port=3000", "--addr=127.0.0.1:3001", "--dns=8.8.8.8", "--password=abc", "--allowlan=false"},
		{"--mode=direct", "--log=info", "--allowlan=true"},
		{"--mode=rule", "--log=warn"},
		{"--mode=global", "--log=error"},
	}
	parseArgs(tests[0])
	if config.Port != 3000 || config.Addr != "127.0.0.1:3001" || config.Dns != "8.8.8.8" || config.Password != "abc" || config.AllowLan != false {
		t.FailNow()
	}
	parseArgs(tests[1])
	if config.Mode != masky.DirectMode || config.LogLevel != log.InfoLevel || config.AllowLan != true {
		t.FailNow()
	}
	parseArgs(tests[2])
	if config.Mode != masky.RuleMode || config.LogLevel != log.WarnLevel {
		t.FailNow()
	}
	parseArgs(tests[3])
	if config.Mode != masky.GlobalMode || config.LogLevel != log.ErrorLevel {
		t.FailNow()
	}
}

func Test_parseAddr(t *testing.T) {
	config.AllowLan = true
	if parseAddr() != fmt.Sprintf(":%v", config.Port) {
		t.FailNow()
	}
	config.AllowLan = false
	if parseAddr() != fmt.Sprintf("127.0.0.1:%v", config.Port) {
		t.FailNow()
	}
}
