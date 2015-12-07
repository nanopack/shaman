package main

// Main entry point into the shaman program. This starts up the API, caching,
// and DNS servers in their own routines.

// TODO:
//  - handle signals
//  - add logging
//  - test

import (
	// "os"

	"github.com/nanopack/shaman/api"
	"github.com/nanopack/shaman/caches"
	// "github.com/nanopack/shaman/config"
	"github.com/nanopack/shaman/server"
)

func main() {
	errors := make(chan error)
	err := caches.Init()
	if err != nil {
		panic(err)
	}
	go func() {
		errors <- api.StartApi()
	}()
	go func() {
		errors <- server.StartServer()
	}()

	if err := <-errors; err != nil {
		panic(err)
	}
}
