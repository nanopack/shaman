package api_test

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/jcelliott/lumber"

	"github.com/nanopack/shaman/api"
	"github.com/nanopack/shaman/config"
	shaman "github.com/nanopack/shaman/core/common"
)

var (
	testResource1 = `{"domain":"google.com","records":[{"type":"A","address":"127.0.0.1"}]}`
	testResource2 = `{"domain":"google.com","records":[{"type":"A","address":"127.0.0.2"}]}`
	badResource   = `{"domain":"google.com","records":[{"type":1,"address":"127.0.0.3"}]}`
	testResource3 = `{"domain":"foogle.com","records":[{"type":"A","address":"127.0.0.4"}]}`
)

func TestMain(m *testing.M) {
	// manually configure
	initialize()

	// start api
	go api.Start()
	<-time.After(time.Second)
	rtn := m.Run()

	os.Exit(rtn)
}

// test put records
func TestPutRecords(t *testing.T) {
	// good request test
	resp, _, err := rest("PUT", "/records", fmt.Sprintf("[%v]", testResource1))
	if err != nil {
		t.Error(err)
	}

	var resources []shaman.Resource
	json.Unmarshal(resp, &resources)

	if len(resources) != 1 {
		t.Errorf("%q doesn't match expected out", resources)
	}

	if len(resources) == 1 &&
		len(resources[0].Records) == 1 &&
		resources[0].Records[0].Address != "127.0.0.1" {
		t.Errorf("%q doesn't match expected out", resources)
	}

	// bad request test
	resp, _, err = rest("PUT", "/records", testResource1)
	if err != nil {
		t.Error(err)
	}

	if !strings.Contains(string(resp), "Bad JSON syntax received in body") {
		t.Errorf("%q doesn't match expected out", resp)
	}

	// clear records
	rest("PUT", "/records", "[]")
}

// todo: "tests should be able to run independent" `go test -v ./api -run TestGet`
// test get records
func TestGetRecords(t *testing.T) {
	body, _, err := rest("GET", "/records", "")
	if err != nil {
		t.Error(err)
	}
	if string(body) != "[]\n" {
		t.Errorf("%q doesn't match expected out", body)
	}
	body, _, err = rest("GET", "/records?full=true", "")
	if err != nil {
		t.Error(err)
	}
	if string(body) != "[]\n" {
		t.Errorf("%q doesn't match expected out", body)
	}
}

// test post records
func TestPostRecord(t *testing.T) {
	// good request test
	resp, _, err := rest("POST", "/records", testResource1)
	if err != nil {
		t.Error(err)
	}

	var resource shaman.Resource
	json.Unmarshal(resp, &resource)

	if resource.Domain != "google.com." {
		t.Errorf("%q doesn't match expected out", resource)
	}

	// bad request test
	resp, _, err = rest("POST", "/records", badResource)
	if err != nil {
		t.Error(err)
	}

	if !strings.Contains(string(resp), "Bad JSON syntax received in body") {
		t.Errorf("%q doesn't match expected out", resp)
	}
}

// test get resource
func TestGetRecord(t *testing.T) {
	// good request test
	resp, _, err := rest("GET", "/records/google.com", "")
	if err != nil {
		t.Error(err)
	}

	var resource shaman.Resource
	json.Unmarshal(resp, &resource)

	if resource.Domain != "google.com." {
		t.Errorf("%q doesn't match expected out", resource)
	}

	// bad request test
	resp, _, err = rest("GET", "/records/not-real.com", "")
	if err != nil {
		t.Error(err)
	}

	if !strings.Contains(string(resp), "failed to find record for domain - 'not-real.com'") {
		t.Errorf("%q doesn't match expected out", resp)
	}
}

// test put records
func TestPutRecord(t *testing.T) {
	// good request test - create(201)
	resp, code, err := rest("PUT", "/records/foogle.com", testResource3)
	if err != nil {
		t.Error(err)
	}
	if code != 201 {
		t.Error("Failed to meet rfc2616 spec, expecting 201")
	}

	var resource shaman.Resource
	json.Unmarshal(resp, &resource)

	if len(resource.Records) == 1 &&
		resource.Records[0].Address != "127.0.0.4" {
		t.Errorf("%q doesn't match expected out", resource)
	}

	// good request test - update
	resp, _, err = rest("PUT", "/records/foogle.com", testResource2)
	if err != nil {
		t.Error(err)
	}

	// verify old resource is gone
	resp, _, err = rest("GET", "/records/foogle.com", "")
	if err != nil {
		t.Error(err)
	}

	if !strings.Contains(string(resp), "failed to find record for domain - 'foogle.com'") {
		t.Errorf("%q doesn't match expected out", resp)
	}

	// bad request test
	resp, _, err = rest("PUT", "/records/not-real.com", badResource)
	if err != nil {
		t.Error(err)
	}

	if !strings.Contains(string(resp), "Bad JSON syntax received in body") {
		t.Errorf("%q doesn't match expected out", resp)
	}
}

// test delete resource
func TestDeleteRecord(t *testing.T) {
	// good request test
	resp, _, err := rest("DELETE", "/records/google.com", "")
	if err != nil {
		t.Error(err)
	}

	if !strings.Contains(string(resp), "{\"msg\":\"success\"}") {
		t.Errorf("%q doesn't match expected out", resp)
	}

	// verify gone
	resp, code, err := rest("GET", "/records/google.com", "")
	if err != nil {
		t.Error(err)
	}

	if code != 404 {
		t.Errorf("%q doesn't match expected out", code)
	}

	// bad request test
	resp, _, err = rest("DELETE", "/records/not-real.com", "")
	if err != nil {
		t.Error(err)
	}

	if !strings.Contains(string(resp), "{\"msg\":\"success\"}") {
		t.Errorf("%q doesn't match expected out", resp)
	}
}

////////////////////////////////////////////////////////////////////////////////
// PRIVS
////////////////////////////////////////////////////////////////////////////////
// hit api and return response body
func rest(method, route, data string) ([]byte, int, error) {
	body := bytes.NewBuffer([]byte(data))
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	uri := fmt.Sprintf("https://%s%s", config.ApiListen, route)

	req, _ := http.NewRequest(method, uri, body)
	req.Header.Add("X-AUTH-TOKEN", config.ApiToken)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, 500, fmt.Errorf("Unable to %v %v - %v", method, route, err)
	}
	defer res.Body.Close()

	if res.StatusCode == 401 {
		return nil, res.StatusCode, fmt.Errorf("401 Unauthorized. Please specify api token (-t 'token')")
	}

	b, err := ioutil.ReadAll(res.Body)

	return b, res.StatusCode, err
}

// manually configure and start internals
func initialize() {
	config.L2Connect = "none://"
	config.ApiListen = "127.0.0.1:1633"
	config.Log = lumber.NewConsoleLogger(lumber.LvlInt("FATAL"))
	config.LogLevel = "FATAL"
}
