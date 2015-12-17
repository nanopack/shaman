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

// main entry point
func main() {
	// make channel for errors
	errors := make(chan error)
	// Start cache engine, api server, and dns server
	caches.InitCache()
	go func() {
		errors <- caches.StartCache()
	}()
	go func() {
		errors <- api.StartApi()
	}()
	go func() {
		errors <- server.StartServer()
	}()
	// break if any of them return an error
	if err := <-errors; err != nil {
		panic(err)
	}
}
