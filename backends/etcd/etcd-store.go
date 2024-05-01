package etcd

/*
Cow
*/
import (
	"context"
	"errors"
	"github.com/aenjoy/aemkvDB"
	"github.com/aenjoy/aemkvDB/backends"
	cli "go.etcd.io/etcd/client/v3"
	"time"
)

type EtcdDB struct {
	opened bool
	Db     *cli.Client
	kv     cli.KV
	ctx    context.Context
}

func (db *EtcdDB) IsOpened() bool {
	return db.opened
}
func (db *EtcdDB) open(addr []string) error {
	return db.openWithCTX(addr, context.Background())
}
func (db *EtcdDB) openWithCTX(addr []string, ctx context.Context) (err error) {
	if db.IsOpened() {
		return errors.New("db already opened")
	}
	db.Db, err = cli.New(cli.Config{
		Endpoints:   addr,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return
	}
	db.kv = cli.NewKV(db.Db)
	db.opened = true
	db.ctx = ctx
	return
}
func (db *EtcdDB) Set(key, value string) (err error) {
	if !db.IsOpened() {
		return errors.New("db is not opened")
	}
	_, err = db.kv.Put(db.ctx, key, value)
	return
}
func (db *EtcdDB) Get(key string) string {
	if !db.IsOpened() {
		return ""
	}
	get, err := db.kv.Get(db.ctx, key)
	if err != nil {
		return ""
	}
	for _, kv := range get.Kvs {
		return string(kv.Value) //返回第一个value
	}
	return ""
}
func (db *EtcdDB) Delete(key string) (err error) {
	if !db.IsOpened() {
		return errors.New("db is not opened")
	}
	_, err = db.Db.Delete(db.ctx, key)
	return
}
func (db *EtcdDB) Close() (err error) {
	if !db.IsOpened() {
		return errors.New("db is not opened or has been cloesd.")
	}
	err = db.Db.Close()
	if err == nil {
		db.opened = false
	}
	return
}

var globalDB *EtcdDB

// NewEtcdDB 创建一个新的EtcdDB实例。如果使用全局数据库，则尝试获取或创建一个全局实例；否则，创建一个新实例。
//
// addr []string: Etcd服务器的地址列表。
// ctx context.Context: 上下文，用于控制数据库操作的生命周期。
// 返回值 *EtcdDB: 返回初始化好的EtcdDB实例地址。
func NewEtcdDB(addr []string, ctx context.Context) *EtcdDB {
	if backends.UsingGlobalDB {
		if globalDB == nil {
			globalDB = &EtcdDB{}
			err := globalDB.openWithCTX(addr, backends.GlobalContext)
			if err != nil {
				return &EtcdDB{}
			}
		}
		return globalDB
	} else {
		db := &EtcdDB{}
		err := db.openWithCTX(addr, ctx)
		if err != nil {
			return &EtcdDB{}
		}
		if globalDB == nil {
			globalDB = db
		}
		return db
	}
}

// New 通过配置信息创建一个新的EtcdDB实例。
//
// config aemkvDB.ConfigMemDb: 包含数据库配置信息，如地址和上下文。
// 返回值 EtcdDB: 返回初始化好的EtcdDB实例。
func New(config aemkvDB.ConfigMemDb) *EtcdDB {
	return NewEtcdDB(config.Addr, config.Ctx)
}
