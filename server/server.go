package server

// This has the handler for the DNS server.

// TODO:
//  - add logging
//  - test

import (
	"errors"
	"github.com/miekg/dns"
	"github.com/nanopack/shaman/caches"
	"github.com/nanopack/shaman/config"
	"strings"
)

var (
	invalidDomain = errors.New("Invalid domain")
	notFound      = errors.New("Record was not found")
)

func stripSubdomain(name string) string {
	names := strings.SplitN(name, ".", 2)
	if len(names) == 2 {
		return names[1]
	} else {
		return ""
	}
}

func answerQuestion(question dns.Question) []dns.RR {
	answers := make([]dns.RR, 0)
	name := question.Name
	for {
		findReturn := make(chan caches.FindReturn)
		var findOp caches.FindOp
		var key string
		if name != question.Name {
			key = caches.Key("*."+name, question.Qtype)
		} else {
			key = caches.Key(name, question.Qtype)
		}
		findOp = caches.FindOp{Key: key, Resp: findReturn}
		caches.FindOps <- findOp
		findRet := <-findReturn
		err := findRet.Err
		record := findRet.Value
		if err != nil {
			config.Log.Error("error: %s", err)
			continue
		}
		if record != "" {
			entry, err := dns.NewRR(record)
			if err != nil {
				config.Log.Error("error: %s", err)
				continue
			}
			entry.Header().Name = question.Name
			answers = append(answers, entry)
		}
		if len(answers) > 0 {
			break
		}
		name = stripSubdomain(name)
		if len(name) == 0 {
			break
		}
	}
	return answers
}

// This receives requests, looks up the result and returns what is found.
func handlerFunc(res dns.ResponseWriter, req *dns.Msg) {
	switch req.Opcode {
	case dns.OpcodeQuery:

		message := new(dns.Msg)
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
