package integ

import (
	"errors"
)

type MapDB struct {
	opened bool
	Db     map[string]interface{}
}

func (db *MapDB) open() error {
	if db.opened {
		return errors.New("db is opened")
	}
	db.opened = true
	db.Db = make(map[string]interface{})
	return nil
}
func (db *MapDB) Set(key string, value interface{}) error {
	if !db.opened {
		return errors.New("db is not opened")
	}
	db.Db[key] = value
	return nil
}
func (db *MapDB) Get(key string) interface{} {
	if !db.opened {
		return errors.New("db is not opened")
	}
	return db.Db[key]
}
func (db *MapDB) Delete(key string) error {
	if !db.opened {
		return errors.New("db is not opened")
	}
	delete(db.Db, key)
	return nil
}
func (db *MapDB) Close() error {
	if !db.opened {
		return errors.New("db is not opened")
	}
	db.opened = false
	db.Db = make(map[string]interface{})
	return nil
}
func (db *MapDB) IsOpened() bool {
	return db.opened
}
