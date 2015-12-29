package commands

import (
	"github.com/jcelliott/lumber"
	"github.com/spf13/cobra"

	"github.com/nanopack/shaman/api"
	"github.com/nanopack/shaman/caches"
	"github.com/nanopack/shaman/config"
	"github.com/nanopack/shaman/server"
)

var (
	runServer bool
	Shaman    = &cobra.Command{
		Use:   "",
		Short: "",
		Long:  ``,

		Run: func(ccmd *cobra.Command, args []string) {
			if runServer {
				startServer()
				return
			}
			// Show the help if not starting the server
			ccmd.HelpFunc()(ccmd, args)
		},
	}
)

func init() {
	// Shaman.PersistentFlags().StringVarP(&config.AuthToken, "auth", "A", "", "Shaman auth token")
	// Shaman.PersistentFlags().StringVarP(&config.Host, "host", "H", "127.0.0.1", "Shaman hostname/IP")
	// Shaman.PersistentFlags().IntVarP(&config.Port, "port", "p", 8443, "Shaman admin port")
	Shaman.PersistentFlags().BoolVarP(&config.Insecure, "insecure", "i", false, "Disable tls key checking")

	Shaman.Flags().BoolVarP(&runServer, "server", "s", false, "Run in server mode")

	Shaman.Flags().StringVarP(&config.L1Connect, "l1-connect", "1", "map://127.0.0.1/",
		"Connection string for the l1 cache")
	Shaman.Flags().IntVarP(&config.L1Expires, "l1-expires", "e", 120,
		"TTL for the L1 Cache (0 = never expire)")
	Shaman.Flags().StringVarP(&config.L2Connect, "l2-connect", "2", "map://127.0.0.1/",
		"Connection string for the l2 cache")
	Shaman.Flags().IntVarP(&config.L2Expires, "l2-expires", "E", 0,
		"TTL for the L2 Cache (0 = never expire)")
	Shaman.Flags().IntVarP(&config.TTL, "ttl", "T", 60,
		"Default TTL for DNS records")
	Shaman.Flags().StringVarP(&config.Domain, "domain", "d", ".",
		"Parent domain for requests")
	Shaman.Flags().StringVarP(&config.Host, "host", "O", "127.0.0.1",
		"Listen address for DNS requests")
	Shaman.Flags().StringVarP(&config.Port, "port", "o", "8053",
		"Listen port for DNS requests")
	Shaman.Flags().StringVarP(&config.ApiKey, "api-key", "k", "",
		"Path to SSL key for API access")
	Shaman.Flags().StringVarP(&config.ApiKeyPassword, "api-key-password", "p", "",
		"Password for SSL key")
	Shaman.Flags().StringVarP(&config.ApiCrt, "api-crt", "c", "",
		"Path to SSL crt for API access")
	Shaman.PersistentFlags().StringVarP(&config.ApiToken, "api-token", "t", "",
		"Token for API Access")
	Shaman.PersistentFlags().StringVarP(&config.ApiHost, "api-host", "H", "127.0.0.1",
		"Listen address for the API")
	Shaman.PersistentFlags().StringVarP(&config.ApiPort, "api-port", "P", "8443",
		"Listen address for the API")
	Shaman.Flags().StringVarP(&config.LogLevel, "log-level", "L", "INFO",
		"Log level to use")
	Shaman.Flags().StringVarP(&config.LogFile, "log-file", "l", "",
		"Log file (blank = log to console)")

	Shaman.AddCommand(addCmd)
	Shaman.AddCommand(removeCmd)
	Shaman.AddCommand(showCmd)
	Shaman.AddCommand(updateCmd)
	Shaman.AddCommand(listCmd)
}

func startServer() {
	if config.LogFile == "" {
		config.Log = lumber.NewConsoleLogger(lumber.LvlInt(config.LogLevel))
	} else {
		var err error
		config.Log, err = lumber.NewFileLogger(config.LogFile, lumber.LvlInt(config.LogLevel), lumber.ROTATE, 5000, 9, 100)
		if err != nil {
			panic(err)
		}
	}
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
