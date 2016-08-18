package cache

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"

	"github.com/nanopack/shaman/config"
	shaman "github.com/nanopack/shaman/core/common"
)

type postgresDb struct {
	pg *sql.DB
}

func (p *postgresDb) connect() error {
	// todo: example: config.DatabaseConnection = "postgres://postgres@127.0.0.1?sslmode=disable"
	db, err := sql.Open("postgres", config.L2Connect)
	if err != nil {
		return fmt.Errorf("Failed to connect to postgres - %v", err)
	}
	err = db.Ping()
	if err != nil {
		return fmt.Errorf("Failed to ping postgres on connect - %v", err)
	}

	p.pg = db
	return nil
}

func (p postgresDb) createTables() error {
	// create records table
	_, err := p.pg.Exec(`
CREATE TABLE IF NOT EXISTS records (
	recordId SERIAL PRIMARY KEY NOT NULL,
	domain   TEXT NOT NULL,
	address  TEXT NOT NULL,
	ttl      INTEGER,
	class    TEXT,
	type     TEXT
)`)
	if err != nil {
		return fmt.Errorf("Failed to create records table - %v", err)
	}

	return nil
}

func (p *postgresDb) initialize() error {
	err := p.connect()
	if err != nil {
		return fmt.Errorf("Failed to create new connection - %v", err)
	}

	// create tables
	err = p.createTables()
	if err != nil {
		return fmt.Errorf("Failed to create tables - %v", err)
	}

	return nil
}

func (p postgresDb) addRecord(resource shaman.Resource) error {
	resources, err := p.listRecords()
	if err != nil {
		return err
	}

	for i := range resources {
		if resources[i].Domain == resource.Domain {
			// if domains match, check address
			for k := range resources[i].Records {
			next:
				for j := range resource.Records {
					// check if the record exists...
					if resource.Records[j].RType == resources[i].Records[k].RType &&
						resource.Records[j].Address == resources[i].Records[k].Address &&
						resource.Records[j].Class == resources[i].Records[k].Class {
						// if so, skip
						config.Log.Trace("Record exists in persistent, skipping...")
						resource.Records = append(resource.Records[:i], resource.Records[i+1:]...)
						goto next
					}
				}
			}
		}
	}

	// add records
	for i := range resource.Records {
		config.Log.Trace("Adding record to database...")
		_, err = p.pg.Exec(fmt.Sprintf(`
INSERT INTO records(domain, address, ttl, class, type)
VALUES('%v', '%v', '%v', '%v', '%v')`,
			resource.Domain, resource.Records[i].Address, resource.Records[i].TTL,
			resource.Records[i].Class, resource.Records[i].RType))
		if err != nil {
			return fmt.Errorf("Failed to insert into records table - %v", err)
		}
	}

	return nil
}

func (p postgresDb) getRecord(domain string) (*shaman.Resource, error) {
	// read from records table
	rows, err := p.pg.Query(fmt.Sprintf("SELECT address, ttl, class, type FROM records WHERE domain = '%v'", domain))
	if err != nil {
		return nil, fmt.Errorf("Failed to select from records table - %v", err)
	}
	defer rows.Close()

	records := make([]shaman.Record, 0, 0)

	// get data
	for rows.Next() {
		rcrd := shaman.Record{}
		err = rows.Scan(&rcrd.Address, &rcrd.TTL, &rcrd.Class, &rcrd.RType)
		if err != nil {
			return nil, fmt.Errorf("Failed to save results into record - %v", err)
		}

		records = append(records, rcrd)
	}

	// check for errors
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("Error with results - %v", err)
	}

	if len(records) == 0 {
		return nil, errNoRecordError
	}

	return &shaman.Resource{Domain: domain, Records: records}, nil
}

func (p postgresDb) updateRecord(domain string, resource shaman.Resource) error {
	// delete old from records
	err := p.deleteRecord(domain)
	if err != nil {
		return fmt.Errorf("Failed to clean old records - %v", err)
	}

	return p.addRecord(resource)
}

func (p postgresDb) deleteRecord(domain string) error {
	_, err := p.pg.Exec(fmt.Sprintf(`DELETE FROM records WHERE domain = '%v'`, domain))
	if err != nil {
		return fmt.Errorf("Failed to delete from records table - %v", err)
	}

	return nil
}

func (p postgresDb) resetRecords(resources []shaman.Resource) error {
	// truncate records table
	_, err := p.pg.Exec("TRUNCATE records")
	if err != nil {
		return fmt.Errorf("Failed to truncate records table - %v", err)
	}
	for i := range resources {
		err = p.addRecord(resources[i]) // prevents duplicates
		if err != nil {
			return fmt.Errorf("Failed to save records - %v", err)
		}
	}
	return nil
}

func (p postgresDb) listRecords() ([]shaman.Resource, error) {
	// read from records table
	rows, err := p.pg.Query("SELECT DISTINCT domain FROM records")
	if err != nil {
		return nil, fmt.Errorf("Failed to select from records table - %v", err)
	}
	defer rows.Close()

	resources := make([]shaman.Resource, 0)

	// get data
	for rows.Next() {
		var domain string
		err = rows.Scan(&domain)
		if err != nil {
			return nil, fmt.Errorf("Failed to save domain - %v", err)
		}
		resource, err := p.getRecord(domain)
		if err != nil {
			return nil, fmt.Errorf("Failed to get record for domain - %v", err)
		}

		resources = append(resources, *resource)
	}

	// check for errors
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("Error with results - %v", err)
	}
	return resources, nil
}
