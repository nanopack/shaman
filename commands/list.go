package commands

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/spf13/cobra"

	"github.com/nanopack/shaman/config"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List entries in shaman database",
	Long:  ``,

	Run: list,
}

func list(ccmd *cobra.Command, args []string) {
	var client *http.Client
	if config.Insecure {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client = &http.Client{Transport: tr}
	} else {
		client = http.DefaultClient
	}

	uri := fmt.Sprintf("https://%s:%s/records", config.ApiHost, config.ApiPort)
	fmt.Println(uri)
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
	req.Header.Add("X-NANOBOX-TOKEN", config.ApiToken)
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
