package buntdb

import (
	"errors"
	"github.com/aenjoy/aemkvDB"
	"github.com/aenjoy/aemkvDB/backends"
	"github.com/tidwall/buntdb"
)

type BuntDB struct {
	Db     *buntdb.DB
	opened bool
}
type MemDB struct {
	BuntDB
}

// New a database.
func (db *MemDB) newMemDB() error {
	if !db.IsOpened() {
		return db.open(":memory")
	}
	return errors.New("db is opened")
}

func (db *BuntDB) open(data string) (err error) {
	db.Db, err = buntdb.Open(data)
	if err == nil {
		db.opened = true
	}
	return
}
func (db *BuntDB) Set(key, value string) error {
	if !db.IsOpened() {
		return errors.New("db is not opened")
	}
	return db.Db.Update(func(tx *buntdb.Tx) error {
		_, _, err := tx.Set(key, value, nil)
		if err != nil {
			return err
		}
		return nil
	})
}
func (db *BuntDB) Get(key string) (ret string) {
	if !db.IsOpened() {
		return
	}
	err := db.Db.View(func(tx *buntdb.Tx) error {
		v, err := tx.Get(key)
		if err != nil {
			return err
		}
		ret = v
		return nil
	})
	if err != nil {
		return ""
	}
	return
}
func (db *BuntDB) Delete(key string) error {
	if !db.IsOpened() {
		return errors.New("db is not opened")
	}
	return db.Db.Update(func(tx *buntdb.Tx) error {
		_, err := tx.Delete(key)
		if err != nil {
			return err
		}
		return nil
	})
}
func (db *BuntDB) Close() error {
	db.opened = false
	return db.Db.Close()
}
func (db *BuntDB) IsOpened() bool {
	return db.opened
}

var globalDB *BuntDB

// NewBuntDB
// Open a database. data: ":memory:" or database file path.
// if you use MemDB, please use ":memory:".
// if you wish open a file, please set data to file path.
func NewBuntDB(data string) *BuntDB {
	if backends.UsingGlobalDB {
		if globalDB == nil {
			globalDB = &BuntDB{}
			err := globalDB.open(data)
			if err != nil {
				return &BuntDB{}
			}
		}
		return globalDB
	} else {
		db := &BuntDB{}
		err := db.open(data)
		if err != nil {
			return &BuntDB{}
		}
		if globalDB == nil {
			globalDB = db
		}
		return db
	}
}

// New 是一个根据配置创建对应数据库实例的函数。
//
// 参数：
//
//	config aemkvDB.ConfigMemDb - 数据库配置信息。Host[]、密码、数据库名称、上下文。
//
// 返回值：
//
//	BuntDB - 若果成功，返回一个数据库实例，否则返回空数据库实例。
func New(config aemkvDB.ConfigMemDb) *BuntDB {
	return NewBuntDB(config.Database)
}
