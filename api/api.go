// Package "api" provides a restful interface to manage entries in the DNS database.
package api

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/pat"
	nanoauth "github.com/nanobox-io/golang-nanoauth"

	"github.com/nanopack/shaman/config"
)

type (
	apiError struct {
		ErrorString string `json:"err"`
	}
	apiMsg struct {
		MsgString string `json:"msg"`
	}
)

var (
	auth         nanoauth.Auth
	badJson      = errors.New("Bad JSON syntax received in body")
	bodyReadFail = errors.New("Body Read Failed")
)

// Start starts shaman's http api
func Start() error {
	var cert *tls.Certificate
	var err error
	if config.ApiCrt == "" {
		cert, err = nanoauth.Generate("shaman.nanobox.io")
	} else {
		cert, err = nanoauth.Load(config.ApiCrt, config.ApiKey, config.ApiKeyPassword)
	}
	if err != nil {
		return err
	}
	auth.Certificate = cert
	auth.Header = "X-AUTH-TOKEN"

	config.Log.Info("Shaman listening on https://%v", config.ApiListen)

	// todo: handle config.Insecure

	return fmt.Errorf("API stopped - %v", auth.ListenAndServeTLS(config.ApiListen, config.ApiToken, routes()))
}

func routes() *pat.Router {
	router := pat.New()

	router.Delete("/records/{domain}", deleteRecord) // delete resource
	router.Put("/records/{domain}", updateRecord)    // reset resource's records
	router.Get("/records/{domain}", getRecord)       // return resource's records

	router.Post("/records", createRecord) // add a resource
	router.Get("/records", listRecords)   // return all domains
	router.Put("/records", updateAnswers) // reset all resources

	return router
}

func writeBody(rw http.ResponseWriter, req *http.Request, v interface{}, status int) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	// print the error only if there is one
	var msg map[string]string
	json.Unmarshal(b, &msg)

	var errMsg string
	if msg["error"] != "" {
		errMsg = msg["error"]
	}

	config.Log.Debug("%s %d %s %s %s", req.RemoteAddr, status, req.Method, req.RequestURI, errMsg)

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(status)
	rw.Write(append(b, byte('\n')))

	return nil
}

// parseBody parses the json body into v
func parseBody(req *http.Request, v interface{}) error {

	// read the body
	b, err := ioutil.ReadAll(req.Body)
	if err != nil {
		config.Log.Error(err.Error())
		return bodyReadFail
	}
	defer req.Body.Close()

	// parse body and store in v
	err = json.Unmarshal(b, v)
	if err != nil {
		return badJson
	}

	return nil
}
