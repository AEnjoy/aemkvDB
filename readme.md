# aemkvDB

aemkvDB is a simple encapsulation of key values for Redis, JSON, SQLite, 
go map memory key value pairs, and local file storage methods, characterized 
by high performance and fast speed.

Please see the support information to get the support platform.

A very tiny key-value store that can handle any number of inserts at the same time.

At the same time, the program also supports simple handling of key value storage using the cli method.

# Supported Backends:

> Redis(Support link to redis server)
> JSON(Support store to json, support save to local file)
> SQLite/MySQL/postgresql(Support link to mysql, postgresql and sqlite)
> go map(package name is integ)(Using go-building map, do not support save data to local file or database)
> local file(using local file storage)
> buntDB(using buntDB to store data)
> ETCD(support link to etcd server)

# SUPPRT PLATFORM

> Linux x86_x64 
>
> Windows amd64
> 
> MacOS x86_x64

# use guide

Usage

## using by interface

```go
package main

import "github.com/aenjoy/aemkvDB"
import "github.com/aenjoy/aemkvDB/backends/redis" //using redis backends
/*
import(
	"github.com/aenjoy/aemkvDB/backends/buntdb"
	"github.com/aenjoy/aemkvDB/backends/etcd"
	"github.com/aenjoy/aemkvDB/backends/file"
	"github.com/aenjoy/aemkvDB/backends/integ"
	"github.com/aenjoy/aemkvDB/backends/json"
	"github.com/aenjoy/aemkvDB/backends/redis"
	"github.com/aenjoy/aemkvDB/backends/sql"
)
*/
var db aemkvDB.StrApi //- data type support: string
//var db2 aemkvDB.AutoApi //auto detect backend - data type support :any
func main(){
	db = redis.New(aemkvDB.ConfigMemDb{
		Addr:     []string{"127.0.0.1:6379"},
		Password: "123456",
		Database: "",
	})
	db.Set("key", "value")
	print(db.Get("key")) //out value
}
```

## using a db backend

- import package from github.com/aenjoy/aemkvDB/backends/***** 

- import  github.com/aenjoy/aemkvDB , which is aemkvDB config, you could use aemkvDB.ConfigMemDb or aemkvDB.ConfigFileDb or aemkvDB.ConfigDb to create db instance. 
 
other example [view](/example)

### BuntDB backend

```golang
package main

import "github.com/aenjoy/aemkvDB"
import "github.com/aenjoy/aemkvDB/backends/buntdb"

func main() {
	// using buntdb backend
	config := aemkvDB.ConfigMemDb{}
	config.Database=":memory:"
	db := buntdb.New(config)
	db.Set("key", "value")
	db.Set("key2", "value2")
	print(db.Get("key"), db.Get("key2")) //value value2
}

```

### Redis backend

```go
package example

import (
	"context"
	"github.com/aenjoy/aemkvDB"
	"github.com/aenjoy/aemkvDB/backends/redis"
)

func redisDB()  {
	config := aemkvDB.ConfigMemDb{}
	config.Addr = []string{"127.0.0.1:6379"}
	config.Ctx = context.Background()
	config.Password="123456"
	config.Database="0"
	db:=redis.New(config)
	db.Set("key", "value")
	db.Set("key2", "value2")
	print(db.Get("key"), db.Get("key2")) //value value2
}
```

# support method

standard api: get,set,delete

```go
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
```

meanwhile, different backends also support different and additional APIs

such as redis backend:

```
Incr(key string)
Decr(key string)
```

json backend:

```
GetStr(key string)
GetInt(key string)
GetBool(key string)
GetFloat64(key string)
SetStr(key, value string)
SetInt(key string, value int)
SetBool(key string, value bool)
SetFloat64(key string, value float64)
```

# In key-value, the types supported by value:

- buntDB:string
- redis:string,int
- json:string,int,bool,float64
- go map and sync map:any
- local file:string
- sql:string

# config readme:

## ConfigMemDb

```go
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
```

Addr: server address and port

Password: password for redis(if needed)

Database: database name for redis(Usually "0") and buntDB(buntDB mode support ":memory:" and file path)

Ctx: context for redis and etcd

## SqlDB

```go
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
```
Addr: server address

User: user name for mysql,postgresql

Password: password for mysql,postgresql

DataBase: database name for mysql,postgresql,sqlite3

Port: port for mysql,postgresql

Mode: open mode . If mode is 0 or 1, only DataBase is needed, because sqlite3 is a file-based database, DataBase is file path. If mode is 2 or 3, Addr,User,Password,DataBase,Port are needed.

## ConfigStandDb

```go
// ConfigStandDb for json,file,syncMap,map db.
// ->github.com/aenjoy/aemkvDB/backends/json
// ->github.com/aenjoy/aemkvDB/backends/integ
type ConfigStandDb struct {
	Data string //for json,file.
	Mode int    //for json,syncMap,map db. 0:Map 1:syncMap
}
```

Data: file path for json and fileDB

Mode: mapDB mode. 0:go-Map 1:go-sync.Map

# Cli Support

```
COMMANDS:
set		Set an key & value in the store
get		Get a value from the store
delete	Delete a key from the store
:help, h	Shows a list of commands or help for one command
```