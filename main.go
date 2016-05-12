// Shaman is a small, clusterable, lightweight, api-driven dns server.
//
// Usage
//
// To start shaman as a server, simply run (with administrator privileges):
//
//  shaman -s
//
// For more specific usage information, refer to the help doc `shaman -h`:
//  Usage:
//    shaman [flags]
//    shaman [command]
//
//  Available Commands:
//    add         Add a domain to shaman
//    delete      Remove a domain from shaman
//    list        List all domains in shaman
//    get         Get records for a domain
//    update      Update records for a domain
//    reset       Reset all domains in shaman
//
//  Flags:
//    -C, --api-crt string            Path to SSL crt for API access
//    -k, --api-key string            Path to SSL key for API access
//    -p, --api-key-password string   Password for SSL key
//    -H, --api-listen string         Listen address for the API (ip:port) (default "127.0.0.1:1632")
//    -c, --config-file string        Configuration file to load
//    -O, --dns-listen string         Listen address for DNS requests (ip:port) (default "127.0.0.1:53")
//    -d, --domain string             Parent domain for requests (default ".")
//    -i, --insecure                  Disable tls key checking (client) and listen on http (api)
//    -2, --l2-connect string         Connection string for the l2 cache (default "scribble:///var/db/shaman")
//    -l, --log-level string          Log level to output [fatal|error|info|debug|trace] (default "INFO")
//    -s, --server                    Run in server mode
//    -t, --token string              Token for API Access (default "secret")
//    -T, --ttl int                   Default TTL for DNS records (default 60)
//    -v, --version                   Print version info and exit
//
package main

import (
	"fmt"
	"os"

	"github.com/jcelliott/lumber"
	"github.com/spf13/cobra"

	"github.com/nanopack/shaman/api"
	"github.com/nanopack/shaman/cache"
	"github.com/nanopack/shaman/commands"
	"github.com/nanopack/shaman/config"
	"github.com/nanopack/shaman/server"
)

var (
	// shaman provides the shaman cli/server functionality
	shamanTool = &cobra.Command{
		Use:              "shaman",
		Short:            "shaman - api driven dns server",
		Long:             ``,
		PersistentPreRun: readConfig,
		PreRun:           preFlight,
		Run:              startShaman,
	}
)

// add supported cli commands/flags
func init() {
	shamanTool.AddCommand(commands.AddDomain)
	shamanTool.AddCommand(commands.DelDomain)
	shamanTool.AddCommand(commands.ListDomains)
	shamanTool.AddCommand(commands.GetDomain)
	shamanTool.AddCommand(commands.UpdateDomain)
	shamanTool.AddCommand(commands.ResetDomains)

	config.AddFlags(shamanTool)
}

func main() {
	shamanTool.Execute()
}

func readConfig(ccmd *cobra.Command, args []string) {
	if err := config.LoadConfigFile(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func preFlight(ccmd *cobra.Command, args []string) {
	if config.Version {
		fmt.Printf("shaman %s\n", VERSION)
		os.Exit(0)
	}

	if !config.Server {
		ccmd.HelpFunc()(ccmd, args)
		os.Exit(0)
	}
}

func startShaman(ccmd *cobra.Command, args []string) {
	config.Log = lumber.NewConsoleLogger(lumber.LvlInt(config.LogLevel))

	// initialize cache
	err := cache.Initialize()
	if err != nil {
		config.Log.Fatal(err.Error())
		os.Exit(1)
	}

	// make channel for errors
	errors := make(chan error)

	go func() {
		errors <- api.Start()
	}()
	go func() {
		errors <- server.Start()
	}()

	// break if any of them return an error (blocks exit)
	if err := <-errors; err != nil {
		config.Log.Fatal(err.Error())
		os.Exit(1)
	}
}
