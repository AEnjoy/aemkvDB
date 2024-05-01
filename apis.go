package aemkvDB

import (
	"context"
	"github.com/aenjoy/aemkvDB/backends"
)

type StrApi interface {
	Set(key, value string) error
	Get(key string) string
	Delete(key string) error
	IsOpened() bool
	Close() error
}
type AutoApi interface {
	Set(key string, value interface{}) error
	Get(key string) interface{}
	Delete(key string) error
	IsOpened() bool
	Close() error
}

func SetGlobalMode(flag bool) {
	backends.UsingGlobalDB = flag
}
func SetGlobalContext(c context.Context) {
	backends.GlobalContext = c
}
func SetGlobalSplit(split string) {
	backends.Split = split
}

// ConfigMemDb for redis,etcd,buntDB.
// ->github.com/aenjoy/aemkvDB/backends/buntdb
// ->github.com/aenjoy/aemkvDB/backends/etcd
// ->github.com/aenjoy/aemkvDB/backends/redis
type ConfigMemDb struct {
	Addr     []string        //for redis,etcd. redis only Use Addr[0]
	Password string          //for redis
	Database string          //for redis and buntDB.
	Ctx      context.Context //for redis,etcd.
}

// SqlDb for sqlite3,mysql,postgresql.
// ->github.com/aenjoy/aemkvDB/backends/sql
type SqlDb struct {
	Addr     string //for mysql,postgresql.
	User     string //for mysql,postgresql.
	Password string //for mysql,postgresql.
	DataBase string //for sqlite3,mysql,postgresql.
	Port     string //for mysql,postgresql.
	Mode     int    // Open mode 0:sqlite-New 1:sqlite-Open 2:mysql-Open 3:postgresSql-Open
}

// ConfigStandDb for json,file,syncMap,map db.
// ->github.com/aenjoy/aemkvDB/backends/json
// ->github.com/aenjoy/aemkvDB/backends/integ
type ConfigStandDb struct {
	Data string //for json,file.
	Mode int    //for json,syncMap,map db. 0:Map 1:syncMap
}
