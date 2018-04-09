package server_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/jcelliott/lumber"
	"github.com/miekg/dns"

	"github.com/nanopack/shaman/config"
	"github.com/nanopack/shaman/core"
	sham "github.com/nanopack/shaman/core/common"
	"github.com/nanopack/shaman/server"
)

var nanopack = sham.Resource{Domain: "nanopack.io.", Records: []sham.Record{{Address: "127.0.0.1"}}}

func TestMain(m *testing.M) {
	// manually configure
	config.DnsListen = "127.0.0.1:8053"
	config.Log = lumber.NewConsoleLogger(lumber.LvlInt("FATAL"))

	// start dns server
	go server.Start()
	<-time.After(time.Second)

	// run tests
	rtn := m.Run()

	os.Exit(rtn)
}

func TestDNS(t *testing.T) {
	err := shaman.AddRecord(&nanopack)
	if err != nil {
		t.Errorf("Failed to add record - %v", err)
		t.FailNow()
	}

	r, err := ResolveIt("nanopack.io", dns.TypeA)
	if err != nil {
		t.Errorf("Failed to get record - %v", err)
	}
	if len(r.Answer) == 0 {
		t.Error("No record found")
	}
	if len(r.Answer) > 0 && r.Answer[0].String() != "nanopack.io.\t60\tIN\tA\t127.0.0.1" {
		t.Errorf("Response doesn't match expected - %+q", r.Answer[0].String())
	}

	r, err = ResolveIt("a.b.nanobox.io", dns.TypeA)
	if err != nil {
		t.Errorf("Failed to get record - %v", err)
	}
	if len(r.Answer) != 0 {
		t.Error("Found non-existant record")
	}

	r, err = ResolveIt("nanopack.io", dns.TypeMX, true)
	if err != nil {
		t.Errorf("Failed to get record - %v", err)
	}
	if len(r.Answer) != 0 {
		t.Error("Found non-existant record")
	}
	// test fallback
	config.DnsFallBack = "8.8.8.8:53"
	r, err = ResolveIt("www.google.com", dns.TypeA)
	if err != nil {
		t.Errorf("Failed to get record - %v", err)
	}
	if len(r.Answer) == 0 {
		t.Error("No record found")
	}

	// reset fallback
	config.DnsFallBack = ""
	r, err = ResolveIt("www.google.com", dns.TypeA)
	if len(r.Answer) != 0 {
		t.Error("answer found for unregistered domain when fallback is off.")
	}
}

func ResolveIt(domain string, rType uint16, badop ...bool) (*dns.Msg, error) {
	// root domain if not already
	root(&domain)
	m := new(dns.Msg)
	m.SetQuestion(domain, rType)

	if len(badop) > 0 {
		m.Opcode = dns.OpcodeStatus
	}

	// ask the dns server
	r, err := dns.Exchange(m, config.DnsListen)
	if err != nil {
		return nil, fmt.Errorf("Failed to exchange - %v", err)
	}

	return r, nil
}

func root(domain *string) {
	t := []byte(*domain)
	if len(t) > 0 && t[len(t)-1] != '.' {
		*domain = string(append(t, '.'))
	}
}
