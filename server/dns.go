// Package "server" contains logic to handle DNS requests.
package server

import (
	"fmt"

	"github.com/miekg/dns"

	"github.com/nanopack/shaman/config"
	"github.com/nanopack/shaman/core"
)

// Start starts the DNS listener
func Start() error {
	dns.HandleFunc(".", handlerFunc)
	udpListener := &dns.Server{Addr: config.DnsListen, Net: "udp"}
	config.Log.Info("DNS listening at udp://%v", config.DnsListen)
	return fmt.Errorf("DNS listener stopped - %v", udpListener.ListenAndServe())
}

// handlerFunc receives requests, looks up the result and returns what is found.
func handlerFunc(res dns.ResponseWriter, req *dns.Msg) {
	message := new(dns.Msg)
	switch req.Opcode {
	case dns.OpcodeQuery:
		message.SetReply(req)
		message.Compress = false
		message.Answer = make([]dns.RR, 0)

		for _, question := range message.Question {
			answers := answerQuestion(question)
			for i := range answers {
				message.Answer = append(message.Answer, answers[i])
			}
		}
		if len(message.Answer) == 0 {
			message.Rcode = dns.RcodeNameError
		}
	default:
		message = message.SetRcode(req, dns.RcodeNotImplemented)
	}
	res.WriteMsg(message)
}

// answerQuestion returns resource record answers for the domain in question
func answerQuestion(question dns.Question) []dns.RR {
	answers := make([]dns.RR, 0)
	r, _ := shaman.GetRecord(question.Name)
	records := r.StringSlice()
	// fmt.Printf("Records received - %+q\n", records)
	for _, record := range records {
		entry, err := dns.NewRR(record)
		if err != nil {
			config.Log.Trace("Failed to create RR from record - %v", err)
			continue
		}
		entry.Header().Name = question.Name
		if entry.Header().Rrtype == question.Qtype || question.Qtype == dns.TypeANY {
			answers = append(answers, entry)
		}
	}

	return answers
}
