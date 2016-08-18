package cache

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/nanobox-io/golang-scribble"

	"github.com/nanopack/shaman/config"
	shaman "github.com/nanopack/shaman/core/common"
)

type scribbleDb struct {
	db *scribble.Driver
}

func (self *scribbleDb) initialize() error {
	u, err := url.Parse(config.L2Connect)
	if err != nil {
		return fmt.Errorf("Failed to parse 'l2-connect' - %v", err)
	}
	dir := u.Path
	if dir == "" || dir == "/" {
		config.Log.Debug("Invalid directory, using default '/var/db/shaman'")
		dir = "/var/db/shaman"
	}
	db, err := scribble.New(dir, nil)
	if err != nil {
		config.Log.Fatal("Failed to create db")
		return fmt.Errorf("Failed to create new db at '%v' - %v", dir, err)
	}

	self.db = db
	return nil
}

func (self scribbleDb) addRecord(resource shaman.Resource) error {
	err := self.db.Write("hosts", resource.Domain, resource)
	if err != nil {
		err = fmt.Errorf("Failed to save record - %v", err)
	}
	return err
}

func (self scribbleDb) getRecord(domain string) (*shaman.Resource, error) {
	resource := shaman.Resource{}
	err := self.db.Read("hosts", domain, &resource)
	if err != nil {
		if strings.Contains(err.Error(), "no such file or directory") {
			err = errNoRecordError
		}
		return nil, err
	}
	return &resource, nil
}

func (self scribbleDb) updateRecord(domain string, resource shaman.Resource) error {
	if domain != resource.Domain {
		err := self.deleteRecord(domain)
		if err != nil {
			return fmt.Errorf("Failed to clear current record - %v", err)
		}
	}

	return self.addRecord(resource)
}

func (self scribbleDb) deleteRecord(domain string) error {
	err := self.db.Delete("hosts", domain)
	if err != nil {
		if strings.Contains(err.Error(), "Unable to find") {
			err = nil
		} else {
			err = fmt.Errorf("Failed to delete record - %v", err)
		}
	}
	return err
}

func (self scribbleDb) resetRecords(resources []shaman.Resource) (err error) {
	self.db.Delete("hosts", "")
	for i := range resources {
		err = self.db.Write("hosts", resources[i].Domain, resources[i])
		if err != nil {
			err = fmt.Errorf("Failed to save records - %v", err)
		}
	}
	return err
}

func (self scribbleDb) listRecords() ([]shaman.Resource, error) {
	resources := make([]shaman.Resource, 0)
	values, err := self.db.ReadAll("hosts")
	if err != nil {
		if strings.Contains(err.Error(), "no such file or directory") {
			// if error is about a missing db, return empty array
			return resources, nil
		}
		return nil, err
	}
	for i := range values {
		var resource shaman.Resource
		if err = json.Unmarshal([]byte(values[i]), &resource); err != nil {
			return nil, fmt.Errorf("Bad JSON syntax found in stored body")
		}
		resources = append(resources, resource)
	}
	return resources, nil
}
