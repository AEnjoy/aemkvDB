package file

/*
Cow
*/
import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"github.com/aenjoy/aemkvDB"
	"github.com/aenjoy/aemkvDB/backends"
	"log"
	"net/url"
	"os"
	"strings"
	"sync"
)

type FileDB struct {
	path   string
	opened bool
	sync.Mutex
	file *os.File
}

// Create a new database connection
func (db *FileDB) open(path string) error {
	//db = &FileDB{}
	db.path = path
	f, err := os.OpenFile(path, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	db.file = f
	db.opened = true
	return nil
}
func (db *FileDB) IsOpened() bool {
	return db.opened
}

// Insert the given Key-Value pair into to the FileDB. Auto save
func (db *FileDB) Set(key, value string) error {
	if !db.IsOpened() {
		return errors.New("db is not opened")
	}
	db.Lock()
	defer db.Unlock()
	check := db.Get(key)
	if check != "" {
		// fmt.Println("Reached key already found")
		err := db.delete(key)
		if err != nil {
			log.Fatal(err)
		}

	}
	key = url.QueryEscape(key)
	value = url.QueryEscape(value)
	s := fmt.Sprintf("%s:%s\n", key, value)
	_, err := db.file.Write([]byte(s))
	if err != nil {
		return err
	}
	return nil
}

// Get the value for the particular key given
func (db *FileDB) Get(key string) string {
	if !db.IsOpened() {
		return ""
	}
	db.file.Seek(0, 0)
	var value string
	sc := bufio.NewScanner(db.file)
	for sc.Scan() {
		text := sc.Text()
		// fmt.Println(text)
		splitText := strings.Split(text, ":")
		keyText, _ := url.QueryUnescape(splitText[0])
		if keyText == key {
			value, _ = url.QueryUnescape(splitText[1])
			break
		}
	}
	return value
}

func (db *FileDB) delete(key string) error {
	db.file.Seek(0, 0)

	var bs []byte
	buf := bytes.NewBuffer(bs)

	sc := bufio.NewScanner(db.file)
	for sc.Scan() {
		text := sc.Text()
		// fmt.Println(text)
		splitText := strings.Split(text, ":")
		keyText, _ := url.QueryUnescape(splitText[0])
		if keyText != key {
			_, err := buf.Write(sc.Bytes())
			if err != nil {
				return err
			}
			_, err = buf.WriteString("\n")
			if err != nil {
				return err
			}
		}
	}

	err := os.WriteFile(db.path, buf.Bytes(), 0666)
	if err != nil {
		return err
	}
	return nil
}

// Delete the particular Key-Value pair
func (db *FileDB) Delete(key string) error {
	if !db.IsOpened() {
		return errors.New("Db is not opened")
	}
	db.Lock()
	defer db.Unlock()
	db.file.Seek(0, 0)

	var bs []byte
	buf := bytes.NewBuffer(bs)

	sc := bufio.NewScanner(db.file)
	for sc.Scan() {
		text := sc.Text()
		// fmt.Println(text)
		splitText := strings.Split(text, ":")
		keyText, _ := url.QueryUnescape(splitText[0])
		if keyText != key {
			_, err := buf.Write(sc.Bytes())
			if err != nil {
				return err
			}
			_, err = buf.WriteString("\n")
			if err != nil {
				return err
			}
		}
	}

	err := os.WriteFile(db.path, buf.Bytes(), 0666)
	if err != nil {
		return err
	}
	return nil
}

// Append data to the value of the given key
func (db *FileDB) Append(key, value string) error {
	if !db.IsOpened() {
		return errors.New("Db is not opened")
	}
	old := db.Get(key)
	newVal := old + value
	err := db.Set(key, newVal)
	if err != nil {
		return err
	}
	return nil
}

func (db *FileDB) save() bool {
	return true
}

func (db *FileDB) Close() error {
	if db.IsOpened() {
		err := db.file.Close()
		if err != nil {
			return errors.New("close file fail")
		}
		db.opened = false
		return nil
	} else {
		return nil
	}
}

var globalDB *FileDB

// NewFileDB 创建一个新的FileDB实例。
// 参数:
//
//	data - 用于数据库初始化的参数（文件路径）。
//
// 返回值:
//
//	返回一个指向FileDB实例的指针
func NewFileDB(data string) *FileDB {
	if backends.UsingGlobalDB {
		if globalDB == nil {
			globalDB = &FileDB{}
			globalDB.open(data)
		}
		return globalDB
	} else {
		db := &FileDB{}
		err := db.open(data)
		if err != nil {
			return &FileDB{}
		}
		if globalDB == nil {
			globalDB = db
		}
		return db
	}
}

// New 创建一个新的FileDB实例。
// 参数:
//
//	config - 用于数据库初始化的参数。
//
// 返回值:
//
//	返回一个指向FileDB实例的指针
func New(config aemkvDB.ConfigStandDb) *FileDB {
	return NewFileDB(config.Data)
}
