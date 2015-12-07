package api

// This is a restful interface to manage entries in the DNS database

// TODO:
//  - parse data to build record to add/update
//  - add logging
//  - test

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"github.com/gorilla/pat"
	"github.com/miekg/dns"
	"github.com/nanopack/shaman/caches"
	"github.com/nanopack/shaman/config"
	"math/big"
	"net/http"
	"time"
)

type handler struct {
	child http.Handler
}

func (self handler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if req.Header.Get("X-SHAMAN-TOKEN") != config.ApiToken {
		rw.WriteHeader(http.StatusUnauthorized)
		return
	}

	self.child.ServeHTTP(rw, req)
}

func loadKeys() (*tls.Certificate, error) {
	crt, err := tls.LoadX509KeyPair(config.ApiCrt, config.ApiKey)
	return &crt, err
}

func generateKeys() (*tls.Certificate, error) {
	host := fmt.Sprintf("shaman.%s", config.Domain)

	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}

	notBefore := time.Now()

	notAfter := notBefore.Add(365 * 24 * 100 * time.Hour) // 100 years..

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return nil, err
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{config.Domain},
		},
		NotBefore: notBefore,
		NotAfter:  notAfter,

		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	template.DNSNames = append(template.DNSNames, host)

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		return nil, err
	}

	cert := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	key := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)})
	c, err := tls.X509KeyPair(cert, key)
	return &c, err
}

func StartApi() error {
	var cert *tls.Certificate
	var err error
	if config.ApiCrt == "" {
		cert, err = generateKeys()
	} else {
		cert, err = loadKeys()
	}
	if err != nil {
		return err
	}
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{*cert},
	}
	tlsConfig.BuildNameToCertificate()
	tlsListener, err := tls.Listen("tcp", config.ApiAddress, tlsConfig)
	if err != nil {
		return err
	}
	h := routes()
	return http.Serve(tlsListener, handler{child: h})
}

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
	rtype := dns.StringToType[req.URL.Query().Get(":rtype")]
	domain := req.URL.Query().Get(":domain")
	key := caches.Key(domain, rtype)
	// miek.nl. 3600 IN MX 10 mx.miek.nl.
	//
	rr, err := dns.NewRR("")
	// rr := new()
	if err != nil {
		writeBody(err, rw, http.StatusInternalServerError)
		return
	}
	err = caches.AddRecord(key, rr.String())
	if err != nil {
		writeBody(err, rw, http.StatusInternalServerError)
	}
}

func updateRecord(rw http.ResponseWriter, req *http.Request) {
	rtype := dns.StringToType[req.URL.Query().Get(":rtype")]
	domain := req.URL.Query().Get(":domain")
	key := caches.Key(domain, rtype)
	rr, err := dns.NewRR("")
	if err != nil {
		writeBody(err, rw, http.StatusInternalServerError)
		return
	}
	err = caches.UpdateRecord(key, rr.String())
	if err != nil {
		writeBody(err, rw, http.StatusInternalServerError)
		return
	}
}

func deleteRecord(rw http.ResponseWriter, req *http.Request) {
	rtype := dns.StringToType[req.URL.Query().Get(":rtype")]
	domain := req.URL.Query().Get(":domain")
	key := caches.Key(domain, rtype)
	err := caches.RemoveRecord(key)
	if err != nil {
		writeBody(err, rw, http.StatusInternalServerError)
	}
}
