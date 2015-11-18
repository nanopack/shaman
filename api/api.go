package api

// TODO:
//  - parse data to build record to add/update
//  - add logging
//  - test

import (
	"encoding/json"
	"github.com/gorilla/pat"
	"github.com/miekg/dns"
	"github.com/nanopack/shaman/caches"
	// "github.com/nanopack/shaman/config"
	"net/http"
)

func routes() *pat.Router {
	router := pat.New()
	router.Get("/records/{rtype}/{domain}", handleRequest(getRecord))
	router.Post("/records/{rtype}/{domain}", handleRequest(addRecord))
	router.Put("/records/{rtype}/{domain}", handleRequest(updateRecord))
	router.Delete("/records/{rtype}/{domain}", handleRequest(deleteRecord))
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
	record, err := caches.FindRecord(key)
	if err != nil {
		writeBody(err, rw, http.StatusInternalServerError)
	}
	if record == "" {
		rw.WriteHeader(http.StatusNotFound)

	}
	rr, err := dns.NewRR(record)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte(err.Error()))
	}
	err = writeBody(rr, rw, http.StatusOK)
	if err != nil {

	}
}

func addRecord(rw http.ResponseWriter, req *http.Request) {
	rtype := dns.StringToType[req.URL.Query().Get(":rtype")]
	domain := req.URL.Query().Get(":domain")
	key := caches.Key(domain, rtype)
	rr, err := dns.NewRR("")
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte(err.Error()))
	}
	err = caches.AddRecord(key, rr.String())
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte(err.Error()))
	}
}

func updateRecord(rw http.ResponseWriter, req *http.Request) {
	rtype := dns.StringToType[req.URL.Query().Get(":rtype")]
	domain := req.URL.Query().Get(":domain")
	key := caches.Key(domain, rtype)
	rr, err := dns.NewRR("")
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte(err.Error()))
	}
	err = caches.UpdateRecord(key, rr.String())
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte(err.Error()))
	}
}

func deleteRecord(rw http.ResponseWriter, req *http.Request) {
	rtype := dns.StringToType[req.URL.Query().Get(":rtype")]
	domain := req.URL.Query().Get(":domain")
	key := caches.Key(domain, rtype)
	err := caches.RemoveRecord(key)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte(err.Error()))
	}
}
