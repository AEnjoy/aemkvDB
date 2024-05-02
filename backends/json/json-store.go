package json

import (
	"errors"
	"github.com/aenjoy/aemkvDB"
	"github.com/aenjoy/aemkvDB/backends"
	"github.com/bytedance/sonic"
	"github.com/bytedance/sonic/ast"
	"os"
	"strconv"
)

type JsonDB struct {
	opened bool
	Db     ast.Node
	path   string
}

// open mode: 0: OpenFile, data is json-file path. 1: CreateFile, data is new json-file path. 2:Using JSON data, data is a json string
func (db *JsonDB) open(mode int, data string) error {
	if db.opened {
		return errors.New("db is opened")
	}
	switch mode {
	case 0:
		file, err := os.ReadFile(data)
		db.path = data
		if err != nil {
			return err
		}
		db.Db, err = sonic.Get(file)
		if err != nil {
			return err
		}
		db.opened = true
		return nil
	case 1:
		if err := os.WriteFile(data, []byte("{}"), 0666); err != nil {
			return err
		}
		db.Db, _ = sonic.Get([]byte("{}"))
		db.path = data
		db.opened = true
		return nil
	case 2:
		dbd, e := sonic.Get([]byte(data))
		if e != nil {
			return e
		}
		db.opened = true
		db.Db = dbd
		return nil
	}
	return errors.New("mode error")
}
func (db *JsonDB) IsOpened() bool {
	return db.opened
}
func (db *JsonDB) GetStr(key string) string {
	if !db.IsOpened() {
		return ""
	}
	value, err := db.Db.Get(key).String()
	if err != nil {
		return ""
	}
	return value
}
func (db *JsonDB) GetInt(key string) int64 {
	if !db.IsOpened() {
		return 0
	}
	value, err := db.Db.Get(key).Int64()
	if err != nil {
		return 0
	}
	return value
}
func (db *JsonDB) GetBool(key string) bool {
	if !db.IsOpened() {
		return false
	}
	value, err := db.Db.Get(key).Bool()
	if err != nil {
		return false
	}
	return value
}
func (db *JsonDB) GetFloat64(key string) float64 {
	if !db.IsOpened() {
		return 0
	}
	value, err := db.Db.Get(key).Float64()
	if err != nil {
		return 0
	}
	return value
}
func (db *JsonDB) SetStr(key, value string) error {
	if !db.IsOpened() {
		return errors.New("db is not opened")
	}
	exists := db.Db.Get(key).Exists()
	if !exists {
		db.Db.Add(ast.NewString(value))
		return nil
	}
	_, err := db.Db.Set(key, ast.NewString(value))
	if err != nil {
		return err
	}
	return nil
}
func (db *JsonDB) SetInt(key, value string) error {
	if !db.IsOpened() {
		return errors.New("db is not opened")
	}
	exists := db.Db.Get(key).Exists()
	if !exists {
		db.Db.Add(ast.NewNumber(value))
		return nil
	}
	_, err := db.Db.Set(key, ast.NewNumber(value))
	if err != nil {
		return err
	}
	return nil
}
func (db *JsonDB) SetFloat64(key string, value float64) error {
	if !db.IsOpened() {
		return errors.New("db is not opened")
	}
	exists := db.Db.Get(key).Exists()
	if !exists {
		db.Db.Add(ast.NewNumber(strconv.FormatFloat(value, 'f', -1, 64)))
		return nil
	}
	_, err := db.Db.Set(key, ast.NewNumber(strconv.FormatFloat(value, 'f', -1, 64)))
	if err != nil {
		return err
	}
	return nil
}
func (db *JsonDB) SetBool(key string, value bool) error {
	if !db.IsOpened() {
		return errors.New("db is not opened")
	}
	exists := db.Db.Get(key).Exists()
	if !exists {
		db.Db.Add(ast.NewBool(value))
		return nil
	}
	_, err := db.Db.Set(key, ast.NewBool(value))
	if err != nil {
		return err
	}
	return nil
}
func (db *JsonDB) saveDB(path string) bool {
	json, err := db.Db.MarshalJSON()
	if err != nil {
		return false
	}
	return os.WriteFile(path, json, 0666) == nil
}

func (db *JsonDB) Set(key string, value interface{}) (err error) {
	if !db.IsOpened() {
		return errors.New("db is not opened")
	}
	switch value.(type) {
	case float64:
		err = db.SetFloat64(key, value.(float64))
	case int:
		err = db.SetInt(key, strconv.Itoa(value.(int)))
	case bool:
		err = db.SetBool(key, value.(bool))
	case string:
		err = db.SetStr(key, value.(string))
	}
	return
}
func (db *JsonDB) Delete(key string) error {
	if !db.IsOpened() {
		return errors.New("db is not opened")
	}
	exists := db.Db.Get(key).Exists()
	if !exists {
		return errors.New("key not exists")
	}
	var jsonData map[string]interface{}
	data, err := db.Db.MarshalJSON()
	if err != nil {
		return err
	}
	if err = sonic.Unmarshal(data, &jsonData); err != nil {
		return err
	}
	delete(jsonData, key)
	data, err = sonic.Marshal(jsonData)
	if err != nil {
		return err
	}
	db.Db, err = sonic.Get(data)
	if err != nil {
		return err
	}
	return nil
}
func (db *JsonDB) Get(key string) (retval interface{}) {
	if !db.IsOpened() {
		return
	}
	node := db.Db.Get(key)
	switch node.Type() {
	case ast.V_STRING:
		retval, _ = node.String()
		return
	case ast.V_NUMBER:
		r, _ := node.Int64()
		retval = strconv.FormatInt(r, 10)
	default:
		return
	}
	return
}

/*func (db *JsonDB) Set(key,value string)error  {
	return db.SetStr(key,value)
}*/

func (db *JsonDB) Close() (err error) {
	if !db.IsOpened() {
		return errors.New("db is not opened or has been closed")
	}
	if db.saveDB(db.path) {
		db.opened = false
		return
	} else {
		return errors.New("save and close db failed")
	}
}

var globalDB *JsonDB

// NewJsonDB
// mode:
// 0: OpenFile, data is json-file path.
// 1: CreateFile, data is new json-file path.
// 2:Using JSON data, data is a json string
//
// mode: 操作模式，
//
//	0: 打开文件，data为json文件路径。
//	1: 创建文件，data为新json文件的路径。
//	2: 使用JSON数据，data为一个json字符串。
//
// data: 根据mode不同，代表不同的数据输入。
// 返回一个初始化后的JsonDB实例。
func NewJsonDB(mode int, data string) *JsonDB {
	if backends.UsingGlobalDB {
		if globalDB == nil {
			globalDB = &JsonDB{}
			err := globalDB.open(mode, data)
			if err != nil {
				return &JsonDB{}
			}
		}
		return globalDB
	} else {
		db := &JsonDB{}
		err := db.open(mode, data)
		if err != nil {
			return &JsonDB{}
		}
		if globalDB == nil {
			globalDB = db
		}
		return db
	}
}

// New 通过 aemkvDB.ConfigStandDb 配置创建一个新的JsonDB实例。
// 这是NewJsonDB的另一种调用方式，主要用于配置初始化。
func New(config aemkvDB.ConfigStandDb) *JsonDB {
	// 使用配置信息调用NewJsonDB创建实例
	return NewJsonDB(config.Mode, config.Data)
}
