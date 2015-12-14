package commands

import (
	"github.com/nanopack/shaman/cli/config"
	"github.com/spf13/cobra"
)

var (
	ShamanCli = &cobra.Command{
		Use:   "",
		Short: "",
		Long:  ``,
	}
)

func init() {
	ShamanCli.PersistentFlags().StringVarP(&config.AuthToken, "auth", "A", "", "Shaman auth token")
	ShamanCli.PersistentFlags().StringVarP(&config.Host, "host", "H", "127.0.0.1", "Shaman hostname/IP")
	ShamanCli.PersistentFlags().IntVarP(&config.Port, "port", "p", 8443, "Shaman admin port")

	ShamanCli.AddCommand(addCmd)
	ShamanCli.AddCommand(removeCmd)
	ShamanCli.AddCommand(showCmd)
	ShamanCli.AddCommand(updateCmd)
	ShamanCli.AddCommand(listCmd)
}
