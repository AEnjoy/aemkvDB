package example

import (
	"github.com/aenjoy/aemkvDB"
	"github.com/aenjoy/aemkvDB/backends/sql"
)

func mysqlDB() {
	config := aemkvDB.SqlDb{Mode: 2, Addr: "127.0.0.1", User: "root", Password: "123456", DataBase: "test", Port: "3306"} // if database exists
	db := sql.New(config)
	defer db.Close()
	db.Set("key", "value")
	db.Set("key2", "value2")
	print(db.Get("key"), db.Get("key2")) //value value2
}
func postgresDB() {
	config := aemkvDB.SqlDb{Mode: 3, Addr: "127.0.0.1", User: "postgres", Password: "123456", DataBase: "test", Port: "5432"} // if database exists
	db := sql.New(config)
	defer db.Close()
	db.Set("key", "value")
	db.Set("key2", "value2")
	print(db.Get("key"), db.Get("key2")) //value value2
}
func sqLiteDB() {
	db := sql.NewSqlDB(1, "./test.db")
	// or
	// db = sql.New(aemkvDB.SqlDb{Mode:1, DataBase: "./test.db"})
	defer db.Close()
	db.Set("key", "value")
	db.Set("key2", "value2")
	print(db.Get("key"), db.Get("key2")) //value value2
}
