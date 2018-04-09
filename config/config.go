// Package config is a central location for configuration options. It also contains
// config file parsing logic.
package config

import (
	"fmt"
	"path/filepath"

	"github.com/jcelliott/lumber"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	ApiDomain          = "shaman.nanobox.io"         // Domain for generating cert (if none passed)
	ApiCrt             = ""                          // Path to SSL crt for API access
	ApiKey             = ""                          // Path to SSL key for API access
	ApiKeyPassword     = ""                          // Password for SSL key
	ApiListen          = "127.0.0.1:1632"            // Listen address for the API (ip:port)
	ApiToken           = "secret"                    // Token for API Access
	Insecure           = false                       // Disable tls key checking (client) and listen on http (server)
	L2Connect          = "scribble:///var/db/shaman" // Connection string for the l2 cache
	TTL            int = 60                          // Default TTL for DNS records
	Domain             = "."                         // Parent domain for requests
	DnsListen          = "127.0.0.1:53"              // Listen address for DNS requests (ip:port)
	DnsFallBack        = ""                          // fallback dns server if record not found in cache, not used if empty

	LogLevel   = "INFO" // Log level to output [fatal|error|info|debug|trace]
	Server     = false  // Run in server mode
	ConfigFile = ""     // Configuration file to load
	Version    = false  // Print version info and exit

	Log lumber.Logger // Central logger for shaman
)

// AddFlags adds the available cli flags
func AddFlags(cmd *cobra.Command) {
	// api
	cmd.Flags().StringVarP(&ApiDomain, "api-domain", "a", ApiDomain, "Domain of generated cert (if none passed)")
	cmd.Flags().StringVarP(&ApiCrt, "api-crt", "C", ApiCrt, "Path to SSL crt for API access")
	cmd.Flags().StringVarP(&ApiKey, "api-key", "k", ApiKey, "Path to SSL key for API access")
	cmd.Flags().StringVarP(&ApiKeyPassword, "api-key-password", "p", ApiKeyPassword, "Password for SSL key")
	cmd.PersistentFlags().StringVarP(&ApiListen, "api-listen", "H", ApiListen, "Listen address for the API (ip:port)")
	cmd.PersistentFlags().StringVarP(&ApiToken, "token", "t", ApiToken, "Token for API Access")
	cmd.PersistentFlags().BoolVarP(&Insecure, "insecure", "i", Insecure, "Disable tls key checking (client) and listen on http (api). Also disables auth-token")

	// dns
	cmd.Flags().StringVarP(&L2Connect, "l2-connect", "2", L2Connect, "Connection string for the l2 cache")
	cmd.Flags().IntVarP(&TTL, "ttl", "T", TTL, "Default TTL for DNS records")
	cmd.Flags().StringVarP(&Domain, "domain", "d", Domain, "Parent domain for requests")
	cmd.Flags().StringVarP(&DnsListen, "dns-listen", "O", DnsListen, "Listen address for DNS requests (ip:port)")
	cmd.Flags().StringVarP(&DnsFallBack, "fallback-dns", "f", DnsFallBack, "Fallback dns server address (ip:port), if not specified fallback is not used")

	// core
	cmd.Flags().StringVarP(&LogLevel, "log-level", "l", LogLevel, "Log level to output [fatal|error|info|debug|trace]")
	cmd.Flags().BoolVarP(&Server, "server", "s", Server, "Run in server mode")
	cmd.PersistentFlags().StringVarP(&ConfigFile, "config-file", "c", ConfigFile, "Configuration file to load")

	cmd.Flags().BoolVarP(&Version, "version", "v", Version, "Print version info and exit")
}

// LoadConfigFile reads the specified config file
func LoadConfigFile() error {
	if ConfigFile == "" {
		return nil
	}

	// Set defaults to whatever might be there already
	viper.SetDefault("api-domain", ApiDomain)
	viper.SetDefault("api-crt", ApiCrt)
	viper.SetDefault("api-key", ApiKey)
	viper.SetDefault("api-key-password", ApiKeyPassword)
	viper.SetDefault("api-listen", ApiListen)
	viper.SetDefault("token", ApiToken)
	viper.SetDefault("insecure", Insecure)
	viper.SetDefault("l2-connect", L2Connect)
	viper.SetDefault("ttl", TTL)
	viper.SetDefault("domain", Domain)
	viper.SetDefault("dns-listen", DnsListen)
	viper.SetDefault("log-level", LogLevel)
	viper.SetDefault("server", Server)
	viper.SetDefault("fallback-dns", DnsFallBack)

	filename := filepath.Base(ConfigFile)
	viper.SetConfigName(filename[:len(filename)-len(filepath.Ext(filename))])
	viper.AddConfigPath(filepath.Dir(ConfigFile))

	err := viper.ReadInConfig()
	if err != nil {
		return fmt.Errorf("Failed to read config file - %v", err)
	}

	// Set values. Config file will override commandline
	ApiDomain = viper.GetString("api-domain")
	ApiCrt = viper.GetString("api-crt")
	ApiKey = viper.GetString("api-key")
	ApiKeyPassword = viper.GetString("api-key-password")
	ApiListen = viper.GetString("api-listen")
	ApiToken = viper.GetString("token")
	Insecure = viper.GetBool("insecure")
	L2Connect = viper.GetString("l2-connect")
	TTL = viper.GetInt("ttl")
	Domain = viper.GetString("domain")
	DnsListen = viper.GetString("dns-listen")
	LogLevel = viper.GetString("log-level")
	Server = viper.GetBool("server")

	return nil
}
