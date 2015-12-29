package config

// TODO:
//  - read command line arguments for options
//  - read config file for options
//  - test

import (
	"github.com/jcelliott/lumber"
)

var (
	Insecure       bool
	L1Connect      string
	L2Connect      string
	L1Expires      int
	L2Expires      int
	Domain         string
	TTL            int
	Host           string
	Port           string
	ApiKey         string
	ApiKeyPassword string
	ApiCrt         string
	ApiToken       string
	ApiHost        string
	ApiPort        string
	LogLevel       string
	LogFile        string
	Log            lumber.Logger
)
