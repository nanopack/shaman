package commands

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"github.com/nanopack/shaman/cli/config"
	"github.com/spf13/cobra"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add entry into shaman database",
	Long:  ``,

	Run: add,
}

type addBody struct {
	value string
}

func add(ccmd *cobra.Command, args []string) {
	if len(args) != 3 {
		fmt.Fprintln(os.Stderr, "Missing arguments: Needs record type, domain, and value")
		os.Exit(1)
	}
	var client *http.Client
	if config.Insecure {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client = &http.Client{Transport: tr}
	} else {
		client = http.DefaultClient
	}
	rtype := args[0]
	domain := args[1]
	value := args[2]
	fmt.Println("rtype:", rtype, "domain:", domain, "value:", value)
	data := url.Values{}
	data.Set("value", value)

	uri := fmt.Sprintf("https://%s:%d/records/%s/%s?%s", config.Host, config.Port, rtype, domain, data.Encode())
	fmt.Println(uri)
	req, err := http.NewRequest("POST", uri, bytes.NewBufferString(data.Encode()))
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
	req.Header.Add("X-NANOBOX-TOKEN", config.AuthToken)
	res, err := client.Do(req)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
	fmt.Println(string(b))
}
