package main

import (
	"bytes"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/nanopack/shaman/config"
)

func TestMain(m *testing.M) {
	// manually configure
	config.LogLevel = "fatal"
	discard := &bytes.Buffer{}
	shamanTool.SetOutput(discard)

	// set args for shaman
	args := strings.Split("-O 127.0.0.1:8053 -2 none:// -s", " ")
	shamanTool.SetArgs(args)

	// run shaman server
	go main()
	<-time.After(time.Second)

	// run tests
	rtn := m.Run()

	os.Exit(rtn)
}

func TestShowHelp(t *testing.T) {
	config.Server = false
	shamanTool.SetArgs([]string{""})

	shamanTool.Execute()
}

func TestBadConfig(t *testing.T) {
	args := strings.Split("-c /tmp/nowaythisexists list", " ")
	shamanTool.SetArgs(args)

	shamanTool.Execute()
	config.ConfigFile = ""
}

func TestShowVersion(t *testing.T) {
	args := strings.Split("-v", " ")
	shamanTool.SetArgs(args)

	shamanTool.Execute()
	config.Version = false
}

func TestBadCache(t *testing.T) {
	config.L2Connect = "!@#$%^&"
	args := strings.Split("-s", " ")
	shamanTool.SetArgs(args)

	shamanTool.Execute()
	config.L2Connect = "none://"
}

func TestBadDNSListen(t *testing.T) {
	config.L2Connect = "none://"
	config.DnsListen = "127.0.0.1:53"
	args := strings.Split("-s", " ")
	shamanTool.SetArgs(args)

	go shamanTool.Execute()
	<-time.After(time.Second)

	// port already in use, will fail here
	shamanTool.Execute()
	config.DnsListen = "127.0.0.1:8053"
}
