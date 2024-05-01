package integ

import (
	"errors"
	"sync"
)

type MapSyncDB struct {
	opened bool
	Db     sync.Map
}

func (db *MapSyncDB) open() error {
	if db.opened {
		return errors.New("db is already opened")
	}
	db.Db = sync.Map{}
	db.opened = true
	return nil
}
func (db *MapSyncDB) IsOpened() bool {
	return db.opened
}
func (db *MapSyncDB) Set(key string, value interface{}) error {
	if !db.IsOpened() {
		return errors.New("db is not opened")
	}
	db.Db.Store(key, value)
	return nil
}
func (db *MapSyncDB) Get(key string) interface{} {
	if !db.IsOpened() {
		return nil
	}
	v, ok := db.Db.Load(key)
	if !ok {
		return nil
	} else {
		return v
	}
}
func (db *MapSyncDB) Delete(key string) error {
	if !db.IsOpened() {
		return errors.New("db is not opened")
	}
	db.Db.Delete(key)
	return nil
}
func (db *MapSyncDB) Close() error {
	if !db.IsOpened() {
		return errors.New("db is not opened")
	}
	db.opened = false
	db.Db = sync.Map{}
	return nil
}
