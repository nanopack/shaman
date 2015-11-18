package config

// TODO:
//  - read command line arguments for options
//  - read config file for options
//  - test

// import (
// 	"flag"
// 	"os"
// )

var (
	L1Connect string
	L2Connect string
	Domain    string
	Address   string
)

func init() {
	Domain = "example.com"
	Address = "127.0.0.1:8053"
}
