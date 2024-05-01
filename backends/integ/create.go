package integ

import (
	"github.com/aenjoy/aemkvDB"
	"github.com/aenjoy/aemkvDB/backends"
)

var globalDB *MapDB
var globaSynclDB *MapSyncDB

// NewMapDB 创建一个新的MapDB实例。
// 返回一个MapDB实例或空实例（在打开数据库时遇到错误的情况下）。
func NewMapDB() *MapDB {
	if backends.UsingGlobalDB {
		// 如果全局数据库尚未初始化，则初始化之
		if globalDB == nil {
			globalDB = &MapDB{}
			err := globalDB.open()
			if err != nil {
				// 打开失败，返回空实例
				return &MapDB{}
			}
		}
		// 返回全局数据库实例
		return globalDB
	} else {
		// 创建一个新的MapDB实例
		db := &MapDB{}
		// 尝试打开新创建的数据库实例
		err := db.open()
		if err != nil {
			// 打开失败，返回空实例
			return &MapDB{}
		}
		if globalDB == nil {
			globalDB = db
		}
		return db
	}
}

// NewSyncMapDB 创建一个新的MapSyncDB实例。
// 返回一个MapSyncDB实例指针或空指针（在打开数据库时遇到错误的情况下）。
func NewSyncMapDB() *MapSyncDB {
	if backends.UsingGlobalDB {
		if globaSynclDB == nil {
			globaSynclDB = &MapSyncDB{}
			// 尝试打开全局同步数据库
			err := globaSynclDB.open()
			if err != nil {
				return nil
			}
		}
		return globaSynclDB
	} else {
		db := &MapSyncDB{}
		err := db.open()
		if err != nil {
			// 打开失败
			return &MapSyncDB{}
		}
		if globaSynclDB == nil {
			globaSynclDB = db
		}
		return db
	}
}

// New 是一个根据配置创建对应数据库实例的函数。
// 参数：
//
//	config aemkvDB.ConfigStandDb - 数据库配置信息，包含Mode配置。
//
// 返回值：
//
//	interface{} - 根据配置模式返回不同类型的数据库实例。模式0返回 NewMapDB()，模式1返回 NewSyncMapDB()。
//	如果没有匹配的模式，返回nil。
func New(config aemkvDB.ConfigStandDb) interface{} {
	switch config.Mode {
	case 0:

		return NewMapDB()
	case 1:
		// 模式为1时，返回一个同步的MapDB实例
		return NewSyncMapDB()
	}
	return nil
}
