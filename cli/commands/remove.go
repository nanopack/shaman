package commands

import (
	// "github.com/nanopack/shaman/cli/config"
	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove entry from shaman database",
	Long:  ``,

	Run: remove,
}

func remove(ccmd *cobra.Command, args []string) {

}
