package caches

// TODO:
//  - implement caching backends
//  - add logging
//  - test

import (
	"fmt"
	// "github.com/nanopack/shaman/config"
)

type Cacher interface {
	GetRecord(string) (string, error)
	SetRecord(string, string) error
	ReviseRecord(string, string) error
	DeleteRecord(string) error
}

var (
	l1 Cacher
	l2 Cacher
)

func Init() {

}

func Key(domain string, rtype uint16) string {
	return fmt.Sprintf("%d-%s", rtype, domain)
}

func AddRecord(key string, value string) error {
	if l1 != nil {
		err := l1.SetRecord(key, value)
		if err != nil {
			return nil
		}
	}
	if l2 != nil {
		err := l2.SetRecord(key, value)
		if err != nil {
			return nil
		}
	}
	return nil
}

func RemoveRecord(key string) error {
	if l1 != nil {
		err := l1.DeleteRecord(key)
		if err != nil {
			return err
		}
	}
	if l2 != nil {
		err := l2.DeleteRecord(key)
		if err != nil {
			return err
		}
	}
	return nil
}

func UpdateRecord(key string, value string) error {
	if l1 != nil {
		err := l1.ReviseRecord(key, value)
		if err != nil {
			return err
		}
	}
	if l2 != nil {
		err := l2.ReviseRecord(key, value)
		if err != nil {
			return err
		}
	}
	return nil
}

func FindRecord(key string) (string, error) {
	var record string
	if l2 != nil {
		record, err := l2.GetRecord(key)
		if err != nil {
			return record, err
		}
	}
	if record != "" {
		return record, nil
	}
	if l1 != nil {
		record, err := l1.GetRecord(key)
		if record != "" {
			l2.SetRecord(key, record)
			return record, err
		}
	}
	return "", nil
}
