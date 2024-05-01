package redis

import (
	"context"
	"errors"
	"github.com/aenjoy/aemkvDB"
	"github.com/aenjoy/aemkvDB/backends"
	"github.com/redis/go-redis/v9"
	"strconv"
	"time"
)

type Db struct {
	opened bool
	Db     *redis.Client
	ctx    context.Context
}

func (db *Db) open(addr, passwd string, database int) error {
	return db.openWithCTX(addr, passwd, database, context.Background())
}
func (db *Db) IsOpened() bool {
	return db.opened
}
func (db *Db) openWithCTX(addr, passwd string, database int, ctx context.Context) (err error) {
	if db.IsOpened() {
		return errors.New("db is already opened")
	}
	db.ctx = ctx
	db.Db = redis.NewClient(&redis.Options{
		Addr:        addr,
		Password:    passwd,
		DB:          database,
		ReadTimeout: 1 * time.Second,
	})
	_, err = db.Db.Ping(db.ctx).Result()
	if err == nil {
		db.opened = true
	}
	return
}

func (db *Db) checkLink() bool {
	if !db.IsOpened() {
		return false
	}
	err := db.Db.Ping(db.ctx).Err()
	if err != nil {
		return false
	}
	return true
}
func (db *Db) Incr(key string) error {
	if !db.IsOpened() {
		return errors.New("Db is not opened")
	}
	if !db.checkLink() {
		return errors.New("failed to connect to Redis")
	}
	return db.Db.Incr(db.ctx, key).Err()
}
func (db *Db) Decr(key string) error {
	if !db.IsOpened() {
		return errors.New("Db is not opened")
	}
	if !db.checkLink() {
		return errors.New("failed to connect to Redis")
	}
	return db.Db.Decr(db.ctx, key).Err()
}

func (db *Db) Set(key, value string) error {
	if !db.IsOpened() {
		return errors.New("Db is not opened")
	}
	if !db.checkLink() {
		return errors.New("failed to connect to Redis")
	}
	db.Db.Set(db.ctx, key, value, 0)
	return nil
}
func (db *Db) Get(key string) string {
	if !db.IsOpened() {
		return ""
	}
	if !db.checkLink() {
		return ""
	}
	result, err := db.Db.Get(db.ctx, key).Result()
	if err != nil {
		return ""
	}
	return result
}
func (db *Db) Delete(key string) error {
	if !db.IsOpened() {
		return errors.New("Db is not opened")
	}
	if !db.checkLink() {
		return errors.New("failed to connect to Redis")
	}
	return db.Db.Del(db.ctx, key).Err()
}
func (db *Db) Close() error {
	if !db.IsOpened() {
		return errors.New("db is not opened")
	}
	db.opened = false
	return db.Db.Close()
}

var globalDB *Db

// NewRedisDB 创建一个Redis数据库实例。
// 创建一个新的数据库实例。如果无法成功创建或打开数据库，则返回一个空的Db实例。
//
// 参数:
//
//	addr - Redis服务器的地址。
//	passwd - 连接Redis服务器所需的密码。
//	database - 要连接的数据库索引。
//	ctx - 操作的上下文环境。
//
// 返回值:
//
//	Db - 成功时返回的数据库实例，失败返回空的Db实例。
func NewRedisDB(addr, passwd string, database int, ctx context.Context) *Db {
	if backends.UsingGlobalDB {
		if globalDB == nil {
			globalDB = &Db{}
			err := globalDB.openWithCTX(addr, passwd, database, backends.GlobalContext)
			if err != nil {
				return &Db{}
			}
		}
		return globalDB
	} else {
		db := &Db{}
		err := db.openWithCTX(addr, passwd, database, ctx)
		if err != nil {
			return &Db{}
		}
		if globalDB == nil {
			globalDB = db
		}
		return db
	}
}

// New 根据提供的配置信息创建一个新的Db实例。
//
// 参数:
//
//	config - 包含数据库配置信息的对象，如地址、密码、数据库索引和上下文。
//
// 返回值:
//
//	Db - 成功时返回的数据库实例，失败返回空的Db实例。
func New(config aemkvDB.ConfigMemDb) *Db {
	i, err := strconv.Atoi(config.Database)
	if err != nil {
		return &Db{}
	}
	return NewRedisDB(config.Addr[0], config.Password, i, config.Ctx)
}
