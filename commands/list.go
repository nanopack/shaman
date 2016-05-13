package commands

import (
	"fmt"
	"io/ioutil"

	"github.com/spf13/cobra"
)

var (
	// ListDomains lists all domains in shaman
	ListDomains = &cobra.Command{
		Use:   "list",
		Short: "List all domains in shaman",
		Long:  ``,

		Run: listRecords,
	}
)

func listRecords(ccmd *cobra.Command, args []string) {
	var query string
	if full {
		query = "?full=true"
	}

	res, err := rest("GET", fmt.Sprintf("/records%v", query), nil)
	if err != nil {
		fail("Could not contact shaman - %v", err)
	}

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fail("Could not read shaman's response - %v", err)
	}

	fmt.Print(string(b))
}
