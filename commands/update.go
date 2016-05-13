package commands

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/spf13/cobra"
)

var (
	// UpdateDomain updates records for a domain
	UpdateDomain = &cobra.Command{
		Use:   "update",
		Short: "Update records for a domain",
		Long:  ``,

		Run: updateRecord,
	}
)

func updateRecord(ccmd *cobra.Command, args []string) {
	if jsonString != "" {
		err := json.Unmarshal([]byte(jsonString), &resource)
		if err != nil {
			fail("Bad JSON syntax")
		}
	}

	if resource.Domain == "" {
		fail("Domain must be specified. Try adding `-d`.")
	}

	resource.Records = append(resource.Records, record)

	// validate valid values
	jsonBytes, err := json.Marshal(resource)
	if err != nil {
		fail("Bad values for resource")
	}

	res, err := rest("PUT", fmt.Sprintf("/records/%v", resource.Domain), bytes.NewBuffer(jsonBytes))
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
