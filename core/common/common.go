// Package "common" contains common structs used in shaman
package common

import (
	"fmt"

	"github.com/nanopack/shaman/config"
)

// Resource contains the domain name and a slice of its records
type Resource struct {
	Domain  string   `json:"domain"`  // google.com
	Records []Record `json:"records"` // dns records
}

// Record contains dns information
type Record struct {
	TTL     int    `json:"ttl"`     // seconds record may be cached (300)
	Class   string `json:"class"`   // protocol family (IN)
	RType   string `json:"type"`    // dns record type (A)
	Address string `json:"address"` // address domain resolves to (216.58.217.46)
}

// StringSlice returns a slice of strings with dns info, each ready for dns.NewRR
func (self Resource) StringSlice() []string {
	var records []string
	for i := range self.Records {
		records = append(records, fmt.Sprintf("%s %d %s %s %s\n", self.Domain,
			self.Records[i].TTL, self.Records[i].Class,
			self.Records[i].RType, self.Records[i].Address))
	}
	return records
}

// SanitizeDomain ensures the domain ends with a `.`
func SanitizeDomain(domain *string) {
	t := []byte(*domain)
	if len(t) > 0 && t[len(t)-1] != '.' {
		*domain = string(append(t, '.'))
	}
}

// UnsanitizeDomain ensures the domain ends with a `.`
func UnsanitizeDomain(domain *string) {
	t := []byte(*domain)
	if len(t) > 0 && t[len(t)-1] == '.' {
		*domain = string(t[:len(t)-1])
	}
}

// Validate ensures record values are set
func (self *Resource) Validate() {
	SanitizeDomain(&self.Domain)

	for i := range self.Records {
		if self.Records[i].Class == "" {
			self.Records[i].Class = "IN"
		}
		if self.Records[i].TTL == 0 {
			self.Records[i].TTL = config.TTL
		}
		if self.Records[i].RType == "" {
			self.Records[i].RType = "A"
		}
	}
}
