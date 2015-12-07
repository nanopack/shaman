package config

// TODO:
//  - read command line arguments for options
//  - read config file for options
//  - test

import (
// "flag"
// "os"
)

var (
	L1Connect  string
	L2Connect  string
	L1Expires  int
	L2Expires  int
	Domain     string
	TTL        int
	Address    string
	ApiKey     string
	ApiCrt     string
	ApiToken   string
	ApiAddress string
)

// Initialize configuration
func init() {
	// Defaults
	L1Connect = "map://127.0.0.1/"
	L2Connect = "map://127.0.0.1/"
	Domain = "example.com"
	Address = "127.0.0.1:8053"
	ApiKey = ""
	ApiCrt = ""
	ApiToken = ""
	ApiAddress = "127.0.0.1:8443"
	TTL = 60
	L1Expires = 2 * TTL
	// read config file options
	// read command line options
}
