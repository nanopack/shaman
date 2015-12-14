package commands

import (
	// "github.com/nanopack/shaman/cli/config"
	"github.com/spf13/cobra"
)

var showCmd = &cobra.Command{
	Use:   "show",
	Short: "Show entry in shaman database",
	Long:  ``,

	Run: show,
}

func show(ccmd *cobra.Command, args []string) {

}
