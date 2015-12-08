package config

// TODO:
//  - read command line arguments for options
//  - read config file for options
//  - test

import (
	"flag"
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
	flag.StringVar(&L1Connect, "l1-connect", "map://127.0.0.1/",
		"Connection string for the l1 cache")
	flag.IntVar(&L1Expires, "l1-expires", 120,
		"TTL for the L1 Cache (0 = never expire)")
	flag.StringVar(&L2Connect, "l2-connect", "map://127.0.0.1/",
		"Connection string for the l2 cache")
	flag.IntVar(&L2Expires, "l2-expires", 0,
		"TTL for the L2 Cache (0 = never expire)")
	flag.IntVar(&TTL, "ttl", 60,
		"Default TTL for DNS records")
	flag.StringVar(&Domain, "domain", "example.com",
		"Parent domain for requests")
	flag.StringVar(&Address, "address", "127.0.0.1:8053",
		"Listen address for DNS requests")
	flag.StringVar(&ApiKey, "api-key", "",
		"Path to SSL key for API access")
	flag.StringVar(&ApiCrt, "api-crt", "",
		"Path to SSL crt for API access")
	flag.StringVar(&ApiToken, "api-token", "",
		"Token for API Access")
	flag.StringVar(&ApiAddress, "api-address", "127.0.0.1:8443",
		"Listen address for the API")
	flag.Parse()
}
