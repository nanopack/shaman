package server

// This has the handler for the DNS server.

// TODO:
//  - add logging
//  - test

import (
	// "fmt"
	"errors"
	"github.com/miekg/dns"
	"github.com/nanopack/shaman/caches"
	"github.com/nanopack/shaman/config"
)

var (
	invalidDomain = errors.New("Invalid domain")
	notFound      = errors.New("Record was not found")
)

// This receives requests, looks up the result and returns what is found.
func handleDNSLookup(res dns.ResponseWriter, req *dns.Msg) {

	switch req.Opcode {
	case dns.OpcodeQuery:

		message := new(dns.Msg)
		message.SetReply(req)
		message.Compress = false
		message.Answer = make([]dns.RR, 0)

		for _, question := range message.Question {
			record, err := caches.FindRecord(caches.Key(question.Name, question.Qtype))
			if err != nil {
				continue
			}
			if record == "" {
				continue
			}
			entry, err := dns.NewRR(record)
			if err != nil {
				continue
			}
			message.Answer = append(message.Answer, entry)
		}
		res.WriteMsg(message)
	default:
	}
}

// This starts the DNS listener
func StartServer() error {
	dns.HandleFunc(config.Domain, handleDNSLookup)
	udpListener := &dns.Server{Addr: config.Address, Net: "udp"}
	return udpListener.ListenAndServe()
}
