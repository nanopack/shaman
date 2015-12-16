package api

// This is a restful interface to manage entries in the DNS database

// TODO:
//  - parse data to build record to add/update
//  - add logging
//  - test

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/gorilla/pat"
	"github.com/miekg/dns"
	nanoauth "github.com/nanobox-io/golang-nanoauth"
	"github.com/nanopack/shaman/caches"
	"github.com/nanopack/shaman/config"
	"net/http"
)

var auth nanoauth.Auth

func StartApi() error {
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
	auth.Header = "X-NANOBOX-TOKEN"
	return auth.ListenAndServeTLS(config.ApiAddress, config.ApiToken, routes())
}

func routes() *pat.Router {
	router := pat.New()
	router.Get("/records/{rtype}/{domain}", handleRequest(getRecord))
	router.Post("/records/{rtype}/{domain}", handleRequest(addRecord))
	router.Put("/records/{rtype}/{domain}", handleRequest(updateRecord))
	router.Delete("/records/{rtype}/{domain}", handleRequest(deleteRecord))
	router.Get("/records", handleRequest(listRecords))
	return router
}

func handleRequest(fn func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		fn(rw, req)
	}
}

func writeBody(v interface{}, rw http.ResponseWriter, status int) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(status)
	rw.Write(b)

	return nil
}

func getRecord(rw http.ResponseWriter, req *http.Request) {
	rtype := dns.StringToType[req.URL.Query().Get(":rtype")]
	domain := req.URL.Query().Get(":domain")
	dns.IsDomainName(domain)
	key := caches.Key(domain, rtype)
	findReturn := make(chan caches.FindReturn)
	findOp := caches.FindOp{Key: key, Resp: findReturn}
	caches.FindOps <- findOp
	findRet := <-findReturn
	err := findRet.Err
	record := findRet.Value
	if err != nil {
		writeBody(err, rw, http.StatusInternalServerError)
		return
	}
	if record == "" {
		writeBody(nil, rw, http.StatusNotFound)
		return
	}
	rr, err := dns.NewRR(record)
	if err != nil {
		writeBody(err, rw, http.StatusInternalServerError)
		return
	}
	err = writeBody(rr, rw, http.StatusOK)
	if err != nil {
		// log error
	}
}

func addRecord(rw http.ResponseWriter, req *http.Request) {
	rtype := req.URL.Query().Get(":rtype")
	domain := req.URL.Query().Get(":domain")
	value := req.FormValue("value")
	ttl := config.TTL
	key := caches.Key(domain, dns.StringToType[rtype])
	rrString := fmt.Sprintf("%s %d IN %s %s", domain, ttl, rtype, value)
	rr, err := dns.NewRR(rrString)
	if err != nil {
		writeBody(err, rw, http.StatusInternalServerError)
		return
	}
	resp := make(chan error)
	addOp := caches.AddOp{Key: key, Value: rr.String(), Resp: resp}
	caches.AddOps <- addOp
	err = <-resp
	if err != nil {
		writeBody(err, rw, http.StatusInternalServerError)
		return
	}
	writeBody(rr, rw, http.StatusOK)
}

func updateRecord(rw http.ResponseWriter, req *http.Request) {
	rtype := req.URL.Query().Get(":rtype")
	domain := req.URL.Query().Get(":domain")
	key := caches.Key(domain, dns.StringToType[rtype])
	ttl := config.TTL
	value := req.FormValue("value")
	rrString := fmt.Sprintf("%s %d IN %s %s", domain, ttl, rtype, value)
	rr, err := dns.NewRR(rrString)
	if err != nil {
		writeBody(err, rw, http.StatusInternalServerError)
		return
	}
	resp := make(chan error)
	updateOp := caches.UpdateOp{Key: key, Value: rr.String(), Resp: resp}
	caches.UpdateOps <- updateOp
	err = <-resp
	if err != nil {
		writeBody(err, rw, http.StatusInternalServerError)
		return
	}
	writeBody(rr, rw, http.StatusOK)
}

func deleteRecord(rw http.ResponseWriter, req *http.Request) {
	rtype := req.URL.Query().Get(":rtype")
	domain := req.URL.Query().Get(":domain")
	key := caches.Key(domain, dns.StringToType[rtype])
	resp := make(chan error)
	removeOp := caches.RemoveOp{Key: key, Resp: resp}
	caches.RemoveOps <- removeOp
	err := <-resp
	if err != nil {
		writeBody(err, rw, http.StatusInternalServerError)
		return
	}
	writeBody(nil, rw, http.StatusOK)
}

func listRecords(rw http.ResponseWriter, req *http.Request) {
	resp := make(chan caches.ListReturn)
	listOp := caches.ListOp{Resp: resp}
	caches.ListOps <- listOp
	listReturn := <-resp
	if listReturn.Err != nil {
		writeBody(listReturn.Err, rw, http.StatusInternalServerError)
		return
	}
	writeBody(listReturn.Values, rw, http.StatusOK)
}
