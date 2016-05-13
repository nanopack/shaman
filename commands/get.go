package commands

import (
	"fmt"
	"io/ioutil"

	"github.com/spf13/cobra"
)

var (
	// GetDomain gets records for a domain
	GetDomain = &cobra.Command{
		Use:   "get",
		Short: "Get records for a domain",
		Long:  ``,

		Run: getResource,
	}
)

func getResource(ccmd *cobra.Command, args []string) {
	if resource.Domain == "" {
		fail("Domain must be specified. Try adding `-d`.")
	}

	res, err := rest("GET", fmt.Sprintf("/records/%v", resource.Domain), nil)
	if err != nil {
		fail("Could not contact shaman - %v", err)
	}

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fail("Could not read shaman's response - %v", err)
	}

	fmt.Print(string(b))
}
