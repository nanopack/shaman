package commands

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/spf13/cobra"
)

var (
	// AddDomain adds a domain to shaman
	AddDomain = &cobra.Command{
		Use:   "add",
		Short: "Add a domain to shaman",
		Long:  ``,

		Run: addRecord,
	}
)

func addRecord(ccmd *cobra.Command, args []string) {
	if jsonString != "" {
		err := json.Unmarshal([]byte(jsonString), &resource)
		if err != nil {
			fail("Bad JSON syntax")
		}
	} else {
		if record.Address == "" {
			// warn if record.Address is empty - doesn't apply to jsonString
			fail("Missing address for record. Try adding `-A`")
		}
		resource.Records = append(resource.Records, record)
	}

	if resource.Domain == "" {
		fail("Domain must be specified. Try adding `-d`.")
	}

	jsonBytes, err := json.Marshal(resource)
	if err != nil {
		fail("Bad values for resource")
	}

	res, err := rest("POST", "/records", bytes.NewBuffer(jsonBytes))
	if err != nil {
		fail("Could not contact shaman - %v", err)
	}

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fail("Could not read shaman's response - %v", err)
	}

	fmt.Print(string(b))
}
