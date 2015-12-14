package commands

import (
	// "github.com/nanopack/shaman/cli/config"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add entry into shaman database",
	Long:  ``,

	Run: add,
}

func add(ccmd *cobra.Command, args []string) {

}
