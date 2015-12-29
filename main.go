package main

// Main entry point into the shaman program. This starts up the API, caching,
// and DNS servers in their own routines.

// TODO:
//  - handle signals
//  - add logging
//  - test

import (
	"github.com/nanopack/shaman/commands"
)

// main entry point
func main() {
	commands.Shaman.Execute()
}
