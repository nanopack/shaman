// Package shaman contains the logic to add/remove DNS entries.
package shaman

// todo: atomic C.U.D.

import (
	"fmt"

	"github.com/nanopack/shaman/cache"
	"github.com/nanopack/shaman/config"
	sham "github.com/nanopack/shaman/core/common"
)

// Answers is the cached collection of dns records
var Answers map[string]sham.Resource

func init() {
	Answers = make(map[string]sham.Resource, 0)
}

// GetRecord returns a resource for the specified domain
func GetRecord(domain string) (sham.Resource, error) {
	sham.SanitizeDomain(&domain)

	resource, ok := Answers[domain]
	// if domain not cached in memory...
	if !ok {
		// fetch from cache
		record, err := cache.GetRecord(domain)
		if record == nil {
			return resource, fmt.Errorf("Failed to find domain - %v", err)
		}
		// update local cache
		config.Log.Debug("Cache differs from local, updating...")
		Answers[domain] = *record
	}

	return Answers[domain], nil
}

// ListDomains returns a list of all known domains
func ListDomains() []string {
	domains := make([]string, 0)

	for _, record := range ListRecords() {
		sham.UnsanitizeDomain(&record.Domain)
		domains = append(domains, record.Domain)
	}

	return domains
}

// ListRecords returns all known domains
func ListRecords() []sham.Resource {
	if cache.Exists() {
		// get from cache
		stored, _ := cache.ListRecords()
		if len(Answers) != len(stored) {
			config.Log.Debug("Cache differs from local, updating...")
			ResetRecords(&stored, true)
		}
	}

	resources := make([]sham.Resource, 0)
	for _, v := range Answers {
		resources = append(resources, v)
	}

	return resources
}

// DeleteRecord deletes the resource(domain)
func DeleteRecord(domain string) error {
	sham.SanitizeDomain(&domain)

	// update cache
	config.Log.Trace("Removing record from persistent cache...")
	err := cache.DeleteRecord(domain)
	if err != nil {
		return err
	}

	// todo: atomic
	delete(Answers, domain)

	// otherwise, be idempotent and report it was deleted...
	return nil
}

// AddRecord adds a record to a resource(domain)
func AddRecord(resource *sham.Resource) error {
	resource.Validate()
	domain := resource.Domain

	// todo: atomic
	_, ok := Answers[domain]
	if ok {
		config.Log.Trace("Domain is in local cache")
		// if we have the domain registered...
		for k := range Answers[domain].Records {
			for j := range resource.Records {
				// check if the record exists...
				if resource.Records[j].RType == Answers[domain].Records[k].RType &&
					resource.Records[j].Address == Answers[domain].Records[k].Address &&
					resource.Records[j].Class == Answers[domain].Records[k].Class {
					// if so, skip...
					config.Log.Trace("Record exists in local cache, skipping...")
					goto next
				}
			}
			// otherwise, add the record
			config.Log.Trace("Record not in local cache, adding...")
			resource.Records = append(resource.Records, Answers[domain].Records[k])
		next:
		}
	}

	// store in cache
	config.Log.Trace("Saving record to persistent cache...")
	err := cache.AddRecord(resource)
	if err != nil {
		return err
	}

	// add the resource to the list of knowns
	Answers[domain] = *resource

	return nil
}

// Exists returns whether or not that domain exists
func Exists(domain string) bool {
	sham.SanitizeDomain(&domain)
	_, ok := Answers[domain]
	return ok
}

// UpdateRecord updates a record to a resource(domain)
func UpdateRecord(domain string, resource *sham.Resource) error {
	resource.Validate()
	sham.SanitizeDomain(&domain)

	// in case of some update to domain name...
	if domain != resource.Domain {
		// delete old domain
		err := DeleteRecord(domain)
		if err != nil {
			return fmt.Errorf("Failed to clean up old domain - %v", err)
		}
	}

	// store in cache
	config.Log.Trace("Updating record in persistent cache...")
	err := cache.UpdateRecord(domain, resource)
	if err != nil {
		return err
	}

	// set new resource to domain
	// todo: atomic
	Answers[resource.Domain] = *resource

	return nil
}

// ResetRecords resets all answers. If any nocache has any values, caching is skipped
func ResetRecords(resources *[]sham.Resource, nocache ...bool) error {
	for i := range *resources {
		(*resources)[i].Validate()
	}

	// new map to clear current answers
	answers := make(map[string]sham.Resource)

	for i := range *resources {
		answers[(*resources)[i].Domain] = (*resources)[i]
	}

	if len(nocache) == 0 {
		// store in cache
		config.Log.Trace("Resetting records in persistent cache...")
		err := cache.ResetRecords(resources)
		if err != nil {
			return err
		}
	}

	// reset the answers
	// todo: atomic
	Answers = answers

	return nil
}
