package commands

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/spf13/cobra"

	shaman "github.com/nanopack/shaman/core/common"
)

var (
	// ResetDomains resets all domains in shaman
	ResetDomains = &cobra.Command{
		Use:   "reset",
		Short: "Reset all domains in shaman",
		Long:  ``,

		Run: resetRecords,
	}
)

func resetRecords(ccmd *cobra.Command, args []string) {
	if jsonString == "" {
		fail("Must pass json string. Try adding `-j`.")
	}

	resources := make([]shaman.Resource, 0)

	err := json.Unmarshal([]byte(jsonString), &resources)
	if err != nil {
		fail("Bad JSON syntax")
	}

	// validate valid values
	jsonBytes, err := json.Marshal(resources)
	if err != nil {
		fail("Bad values for resource")
	}

	res, err := rest("PUT", "/records", bytes.NewBuffer(jsonBytes))
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
