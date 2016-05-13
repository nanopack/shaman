package api

import (
	"fmt"
	"net/http"

	"github.com/nanopack/shaman/core"
	sham "github.com/nanopack/shaman/core/common"
)

func createRecord(rw http.ResponseWriter, req *http.Request) {
	var resource sham.Resource
	err := parseBody(req, &resource)
	if err != nil {
		writeBody(rw, req, apiError{err.Error()}, http.StatusBadRequest)
		return
	}

	err = shaman.AddRecord(&resource)
	if err != nil {
		writeBody(rw, req, apiError{err.Error()}, http.StatusInternalServerError)
		return
	}

	writeBody(rw, req, resource, http.StatusOK)
}

func listRecords(rw http.ResponseWriter, req *http.Request) {
	if req.URL.Query().Get("full") == "true" {
		writeBody(rw, req, shaman.ListRecords(), http.StatusOK)
		return
	}

	writeBody(rw, req, shaman.ListDomains(), http.StatusOK)
}

func updateAnswers(rw http.ResponseWriter, req *http.Request) {
	resources := make([]sham.Resource, 0)
	err := parseBody(req, &resources)
	if err != nil {
		writeBody(rw, req, apiError{err.Error()}, http.StatusBadRequest)
		return
	}

	err = shaman.ResetRecords(&resources)
	if err != nil {
		writeBody(rw, req, apiError{err.Error()}, http.StatusInternalServerError)
		return
	}

	writeBody(rw, req, resources, http.StatusOK)
}

func updateRecord(rw http.ResponseWriter, req *http.Request) {
	var resource sham.Resource
	err := parseBody(req, &resource)
	if err != nil {
		writeBody(rw, req, apiError{err.Error()}, http.StatusBadRequest)
		return
	}

	domain := req.URL.Query().Get(":domain")

	if !shaman.Exists(domain) {
		// create resource if not exist
		err = shaman.AddRecord(&resource)
		if err != nil {
			writeBody(rw, req, apiError{err.Error()}, http.StatusInternalServerError)
			return
		}

		// "MUST reply 201"(https://www.w3.org/Protocols/rfc2616/rfc2616-sec9.html)
		writeBody(rw, req, resource, http.StatusCreated)
		return
	}

	err = shaman.UpdateRecord(domain, &resource)
	if err != nil {
		writeBody(rw, req, apiError{err.Error()}, http.StatusInternalServerError)
		return
	}

	writeBody(rw, req, resource, http.StatusOK)
}

func getRecord(rw http.ResponseWriter, req *http.Request) {
	domain := req.URL.Query().Get(":domain")

	resource, err := shaman.GetRecord(domain)
	if err != nil {
		writeBody(rw, req, apiError{fmt.Sprintf("failed to find record for domain - '%v'", domain)}, http.StatusNotFound)
		return
	}

	writeBody(rw, req, resource, http.StatusOK)
}

func deleteRecord(rw http.ResponseWriter, req *http.Request) {
	domain := req.URL.Query().Get(":domain")

	err := shaman.DeleteRecord(domain)
	if err != nil {
		writeBody(rw, req, apiError{err.Error()}, http.StatusInternalServerError)
		return
	}

	writeBody(rw, req, apiMsg{"success"}, http.StatusOK)
}
