// Package commands provides the cli functionality.
// Runnable commands are:
//  add
//  get
//  update
//  delete
//  list
//  reset
package commands

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/spf13/cobra"

	"github.com/nanopack/shaman/config"
	shaman "github.com/nanopack/shaman/core/common"
)

func rest(method string, path string, body io.Reader) (*http.Response, error) {
	uri := fmt.Sprintf("https://%s%s", config.ApiListen, path)

	if config.Insecure {
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}

	req, err := http.NewRequest(method, uri, body)
	if err != nil {
		panic(err)
	}
	req.Header.Add("X-AUTH-TOKEN", config.ApiToken)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		// if requesting `https://` failed, server may have been started with `-i`, try `http://`
		uri = fmt.Sprintf("http://%s%s", config.ApiListen, path)
		req, er := http.NewRequest(method, uri, body)
		if er != nil {
			panic(er)
		}
		req.Header.Add("X-AUTH-TOKEN", config.ApiToken)
		var err2 error
		res, err2 = http.DefaultClient.Do(req)
		if err2 != nil {
			// return original error to client
			return nil, err
		}
	}
	if res.StatusCode == 401 {
		return nil, fmt.Errorf("401 Unauthorized. Please specify api token (-t 'token')")
	}
	return res, nil
}

func fail(format string, args ...interface{}) {
	fmt.Printf(fmt.Sprintf("%v\n", format), args...)
	os.Exit(1)
}

func init() {
	domainFlags(AddDomain)
	DelDomain.Flags().StringVarP(&resource.Domain, "domain", "d", "", "Domain to remove")
	GetDomain.Flags().StringVarP(&resource.Domain, "domain", "d", "", "Domain to get")
	ListDomains.Flags().BoolVarP(&full, "full", "f", false, "Show complete records")
	ResetDomains.Flags().StringVarP(&jsonString, "json", "j", "", "JSON encoded data for domain[s] and record[s]")
	domainFlags(UpdateDomain)
}

var (
	resource   shaman.Resource
	record     shaman.Record
	jsonString string
	full       bool
)

// ResetVars resets the flag vars (used for testing)
func ResetVars() {
	resource = shaman.Resource{}
	record = shaman.Record{}
	jsonString = ""
	full = false
}

func domainFlags(ccmd *cobra.Command) {
	ccmd.Flags().StringVarP(&resource.Domain, "domain", "d", "", "Domain")
	ccmd.Flags().IntVarP(&record.TTL, "ttl", "T", 60, "Record time to live")
	ccmd.Flags().StringVarP(&record.Class, "class", "C", "IN", "Record class")
	ccmd.Flags().StringVarP(&record.RType, "type", "R", "A", "Record type (A, CNAME, MX, etc...)")
	ccmd.Flags().StringVarP(&record.Address, "address", "A", "", "Record address")
	ccmd.Flags().StringVarP(&jsonString, "json", "j", "", "JSON encoded data for domain[s] and record[s]")
}
