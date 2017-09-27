// Package server contains logic to handle DNS requests.
package server

import (
	"fmt"
	"strings"

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
			answers := answerQuestion(question.Qtype, strings.ToLower(question.Name))
			if len(answers) > 0 {
				for i := range answers {
					message.Answer = append(message.Answer, answers[i])
				}
			} else {
				// If there are no records, go back through and search for SOA records
				for _, question := range message.Question {
					answers := answerQuestion(dns.TypeSOA, strings.ToLower(question.Name))
					for i := range answers {
						message.Ns = append(message.Ns, answers[i])
					}
				}
			}
		}
		if len(message.Answer) == 0 && len(message.Ns) == 0 {
			message.Rcode = dns.RcodeNameError
		}
	default:
		message = message.SetRcode(req, dns.RcodeNotImplemented)
	}
	res.WriteMsg(message)
}

// answerQuestion returns resource record answers for the domain in question
func answerQuestion(qtype uint16, name ...string) []dns.RR {
	answers := make([]dns.RR, 0)
	qName := name[len(name)-1] // either `len` every time, or use var

	// get the resource (check memory, cache, and (todo:) upstream)
	r, err := shaman.GetRecord(qName)
	if err != nil {
		config.Log.Trace("Failed to get records for '%s' - %v", qName, err)
	}

	// validate the records and append correct type to answers[]
	for _, record := range r.StringSlice() {
		entry, err := dns.NewRR(record)
		if err != nil {
			config.Log.Debug("Failed to create RR from record - %v", err)
			continue
		}
		entry.Header().Name = name[0]
		if entry.Header().Rrtype == qtype || qtype == dns.TypeANY {
			answers = append(answers, entry)
		}
	}

	// recursively resolve if no records found (essentially provides wildcard
	// registration support)
	if len(answers) == 0 {
		qName = stripSubdomain(qName)
		if len(qName) > 0 {
			config.Log.Trace("Checking again with '%v'", qName)
			return answerQuestion(qtype, name[0], qName)
		}
	}

	return answers
}

// stripSubdomain strips off the subbest domain, returning the domain (won't return TLD)
func stripSubdomain(name string) string {
	words := 3 // assume rooted domain (end with '.')
	// handle edge case of unrooted domain
	t := []byte(name)
	if len(t) > 0 && t[len(t)-1] != '.' {
		words = 2
	}

	config.Log.Trace("Stripping subdomain from '%v'", name)
	names := strings.Split(name, ".")

	// prevent searching for just 'com.' (["domain", "com", ""])
	if len(names) > words {
		return strings.Join(names[1:], ".")
	}
	return ""
}
