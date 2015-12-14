package commands

import (
	// "github.com/nanopack/shaman/cli/config"
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update entry in shaman database",
	Long:  ``,

	Run: update,
}

func update(ccmd *cobra.Command, args []string) {

}
