package cache

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"net/url"

	"github.com/nanopack/shaman/config"

	consul "github.com/hashicorp/consul/api"
	shaman "github.com/nanopack/shaman/core/common"
)

const prefix = "domains:"

type consulDb struct {
	db *consul.Client
}

func addPrefix(in string) string {
	return prefix + in
}

func (client *consulDb) initialize() error {
	u, err := url.Parse(config.L2Connect)
	if err != nil {
		return err
	}

	consulConfig := consul.DefaultNonPooledConfig()
	consulConfig.Address = u.Host
	consulConfig.Scheme = u.Scheme
	consulC, err := consul.NewClient(consulConfig)
	if err != nil {
		return err
	}
	client.db = consulC
	return nil
}

func (client consulDb) addRecord(resource shaman.Resource) error {
	return client.updateRecord(resource.Domain, resource)
}

func (client consulDb) getRecord(domain string) (*shaman.Resource, error) {
	kvHandler := client.db.KV()
	kvPair, _, err := kvHandler.Get(addPrefix(domain), nil)
	if err != nil {
		return nil, err
	}
	if kvPair == nil {
		return nil, errNoRecordError
	}
	var result shaman.Resource
	err = gob.NewDecoder(bytes.NewReader(kvPair.Value)).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (client consulDb) updateRecord(domain string, resource shaman.Resource) error {
	kvHandler := client.db.KV()
	var buf bytes.Buffer
	err := gob.NewEncoder(&buf).Encode(&resource)
	if err != nil {
		return err
	}

	_, err = kvHandler.Put(&consul.KVPair{
		Key:   addPrefix(domain),
		Value: buf.Bytes(),
	}, nil)
	if err != nil {
		return err
	}

	return nil
}

func (client consulDb) deleteRecord(domain string) error {
	kvHandler := client.db.KV()
	_, err := kvHandler.Delete(domain, nil)
	if err != nil {
		return err
	}

	return nil
}

func (client consulDb) resetRecords(resources []shaman.Resource) error {
	kvHandler := client.db.KV()
	_, err := kvHandler.DeleteTree(prefix, nil)
	if err != nil {
		return err
	}

	for i := range resources {
		err = client.addRecord(resources[i]) // prevents duplicates
		if err != nil {
			return fmt.Errorf("Failed to save records - %v", err)
		}
	}
	return nil
}

func (client consulDb) listRecords() ([]shaman.Resource, error) {
	kvHandler := client.db.KV()
	kvPairs, _, err := kvHandler.List(prefix, nil)
	if err != nil {
		return nil, err
	}

	result := []shaman.Resource{}
	for _, kvPair := range kvPairs {
		var resource shaman.Resource
		err := gob.NewDecoder(bytes.NewReader(kvPair.Value)).Decode(&resource)
		if err != nil {
			return nil, err
		}
		result = append(result, resource)
	}

	return result, nil
}
