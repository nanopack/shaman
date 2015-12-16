package commands

import (
	"crypto/tls"
	"fmt"
	"github.com/nanopack/shaman/cli/config"
	"github.com/spf13/cobra"
	"io/ioutil"
	"net/http"
	"os"
)

var showCmd = &cobra.Command{
	Use:   "show",
	Short: "Show entry in shaman database",
	Long:  ``,

	Run: show,
}

func show(ccmd *cobra.Command, args []string) {
	if len(args) != 2 {
		fmt.Fprintln(os.Stderr, "Missing arguments: Needs record type and domain")
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
	fmt.Println("rtype:", rtype, "domain:", domain)

	uri := fmt.Sprintf("https://%s:%d/records/%s/%s", config.Host, config.Port, rtype, domain)
	fmt.Println(uri)
	req, err := http.NewRequest("GET", uri, nil)
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
