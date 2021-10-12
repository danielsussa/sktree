package defaultdb

import (
	"github.com/dgraph-io/badger/v3"
	"log"
)

type BadgerDB struct {
	db *badger.DB
}

func (dmp BadgerDB) Find(key string) (string, bool) {
	value := ""
	err := dmp.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			return err
		}
		// Alternatively, you could also use item.ValueCopy().
		valCopy, err := item.ValueCopy(nil)
		if err != nil {
			return err
		}
		value = string(valCopy)
		return nil
	})
	if err != nil {
		return "", false
	}
	return value, true
}

func (dmp BadgerDB) Add(key, value string) error {
	err := dmp.db.Update(func(txn *badger.Txn) error {
		err := txn.Set([]byte(key), []byte(value))
		return err
	})
	return err
}

func NewBadgerDB(filename string) BadgerDB {
	// Open the Badger database located in the /tmp/badger directory.
	// It will be created if it doesn't exist.
	options := badger.DefaultOptions(filename)
	options.Logger = nil
	db, err := badger.Open(options)
	if err != nil {
		log.Fatal(err)
	}
	return BadgerDB{db: db}
}
