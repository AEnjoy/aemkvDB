package sql

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/aenjoy/aemkvDB"
	"github.com/aenjoy/aemkvDB/backends"
	_ "github.com/mattn/go-sqlite3"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"strings"
)

type SqlDB struct {
	opened bool
	Db     *gorm.DB
}
type KeyValueTable struct {
	gorm.Model
	Key   string //`gorm:"primaryKey"`
	Value string
}

func (db *SqlDB) IsOpened() bool {
	return db.opened
}
func (db *SqlDB) open(mode int, path string) (err error) {
	if db.IsOpened() {
		return errors.New("already opened")
	}
	switch mode {
	case 0:
		database, err := sql.Open("sqlite3", "file:"+path+"?mode=rwc")
		if err != nil {
			return err
		}
		database.Close()
		fallthrough
	case 1:
		db.Db, err = gorm.Open(sqlite.Open(path), &gorm.Config{})
		if err == nil {
			db.opened = true
		} else {
			return
		}
		err = db.Db.AutoMigrate(&KeyValueTable{})
		return
	case 2:
		data := strings.Split(path, backends.Split)
		sprint := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", data[0], data[1], data[2], data[3], data[4])
		db.Db, err = gorm.Open(mysql.Open(sprint), &gorm.Config{})
		if err == nil {
			db.opened = true
		} else {
			return
		}
		err = db.Db.AutoMigrate(&KeyValueTable{})
		return
	case 3:
		data := strings.Split(path, backends.Split)
		sprint := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable", data[0], data[1], data[2], data[3], data[4])
		db.Db, err = gorm.Open(postgres.Open(sprint), &gorm.Config{})
		if err == nil {
			db.opened = true
		} else {
			return
		}
		err = db.Db.AutoMigrate(&KeyValueTable{})
		return
	}
	return errors.New("mode error")
}
func (db *SqlDB) Set(key, value string) error {
	if !db.IsOpened() {
		return errors.New("db is not opened or has been cloesd.")
	}
	result := db.Db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&KeyValueTable{
		Key:   key,
		Value: value,
	})
	return result.Error
}
func (db *SqlDB) Get(key string) string {
	if !db.IsOpened() {
		return ""
	}
	var keyValue KeyValueTable
	if err := db.Db.First(&keyValue, KeyValueTable{Key: key}).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ""
		}
		return ""
	}
	return keyValue.Value
}
func (db *SqlDB) Delete(key string) error {
	if !db.IsOpened() {
		return errors.New("db is not opened or has been cloesd.")
	}
	result := db.Db.Delete(&KeyValueTable{}, "key = ?", key)
	return result.Error
}
func (db *SqlDB) Close() (err error) {
	if !db.IsOpened() {
		return errors.New("db is not opened or has been cloesd.")
	}
	db.opened = false
	return
}

var globalDB *SqlDB

// NewSqlDB 创建或打开一个数据库连接。
//
// Mode: 模式标志，用于指定数据库类型。
//
//	0: 创建并打开一个 sqlite 数据库。
//	1: 打开一个 sqlite 数据库。
//	2: 打开一个 mysql 数据库。
//	3: 打开一个 postgres 数据库。
//
// path: 对于 sqlite，这是数据库文件的路径；对于 mysql 和 postgres，这是连接字符串的组成部分之一。
// 返回值 SqlDB: 返回一个初始化好的 SqlDB 实例，如果初始化失败，则返回一个空的 SqlDB 实例。
func NewSqlDB(mode int, path string) *SqlDB {
	if backends.UsingGlobalDB {
		if globalDB == nil {
			globalDB = &SqlDB{}
			err := globalDB.open(mode, path)
			if err != nil {
				return &SqlDB{}
			}
		}
		return globalDB
	} else {
		db := &SqlDB{}
		err := db.open(mode, path)
		if err != nil {
			return &SqlDB{}
		}
		if globalDB == nil {
			globalDB = db
		}
		return db
	}
}

// New 通过配置信息创建一个新的 SqlDB 实例。
//
// config: 包含数据库连接信息的 aemkvDB.SqlDb 结构体。
//
//	包括 User(用户)、Password(密码)、Addr(地址)、Port(端口)、DataBase(数据库名) 和 Mode(数据库类型)。
//
// 返回值 SqlDB: 返回一个根据配置信息初始化好的 SqlDB 实例
func New(config aemkvDB.SqlDb) *SqlDB {
	//user split password split host split port split dbname
	str := fmt.Sprintf("%s%s%s%s%s%s%s%s%s",
		config.User, backends.Split,
		config.Password, backends.Split,
		config.Addr, backends.Split,
		config.Port, backends.Split,
		config.DataBase)
	if config.Mode == 0 || config.Mode == 1 {
		return NewSqlDB(config.Mode, config.DataBase)
	}
	return NewSqlDB(config.Mode, str)
}
