package server

// This has the handler for the DNS server.

// TODO:
//  - add logging
//  - test

import (
	"errors"
	"fmt"
	"github.com/miekg/dns"
	"github.com/nanopack/shaman/caches"
	"github.com/nanopack/shaman/config"
)

var (
	invalidDomain = errors.New("Invalid domain")
	notFound      = errors.New("Record was not found")
)

// This receives requests, looks up the result and returns what is found.
func handlerFunc(res dns.ResponseWriter, req *dns.Msg) {
	switch req.Opcode {
	case dns.OpcodeQuery:

		message := new(dns.Msg)
		message.SetReply(req)
		message.Compress = false
		message.Answer = make([]dns.RR, 0)

		for _, question := range message.Question {
			findReturn := make(chan caches.FindReturn)
			var findOp caches.FindOp
			// findOp.key = caches.Key(question.Name, question.Qtype)
			// findOp.resp = findReturn
			findOp = caches.FindOp{Key: caches.Key(question.Name, question.Qtype), Resp: findReturn}
			caches.FindOps <- findOp
			findRet := <-findReturn
			err := findRet.Err
			record := findRet.Value
			if err != nil {
				config.Log.Error("error: %s", err)
				continue
			}
			if record == "" {
				// TESTING ONLY!!!
				record = fmt.Sprintf("%s %d %s 127.0.0.1", question.Name, config.TTL, dns.TypeToString[question.Qtype])
				config.Log.Debug("nothing found, setting to 127.0.0.1")
				// continue
			}
			entry, err := dns.NewRR(record)
			if err != nil {
				config.Log.Error("error: %s\n", err)
				continue
			}
			config.Log.Info("record: %s\n", entry)
			message.Answer = append(message.Answer, entry)
		}
		res.WriteMsg(message)
	default:
	}
}

// This starts the DNS listener
func StartServer() error {
	dns.HandleFunc(config.Domain, handlerFunc)
	udpListener := &dns.Server{Addr: config.Address, Net: "udp"}
	return udpListener.ListenAndServe()
}
