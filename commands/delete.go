package commands

import (
	"fmt"
	"io/ioutil"

	"github.com/spf13/cobra"
)

var (
	// DelDomain removes a domain from shaman
	DelDomain = &cobra.Command{
		Use:   "delete",
		Short: "Remove a domain from shaman",
		Long:  ``,

		Run: delRecord,
	}
)

func delRecord(ccmd *cobra.Command, args []string) {
	if resource.Domain == "" {
		fail("Domain must be specified. Try adding `-d`.")
	}

	res, err := rest("DELETE", fmt.Sprintf("/records/%v", resource.Domain), nil)
	if err != nil {
		fail("Could not contact shaman - %v", err)
	}

	// parse response
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fail("Could not read shaman's response - %v", err)
	}

	fmt.Print(string(b))
}
