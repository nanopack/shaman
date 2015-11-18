package main

// TODO:
//  - handle signals
//  - add logging
//  - test

import (
	// "os"

	// "github.com/nanopack/shaman/api"
	"github.com/nanopack/shaman/caches"
	// "github.com/nanopack/shaman/config"
	"github.com/nanopack/shaman/server"
)

func main() {
	caches.Init()
	err := server.StartServer()
	if err != nil {
		panic(err)
	}
}
