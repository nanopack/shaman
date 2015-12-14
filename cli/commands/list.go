package commands

import (
	// "github.com/nanopack/shaman/cli/config"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List entries in shaman database",
	Long:  ``,

	Run: list,
}

func list(ccmd *cobra.Command, args []string) {

}
