package example

/*
use interface example
*/

import (
	"github.com/aenjoy/aemkvDB"
	"github.com/aenjoy/aemkvDB/backends/integ"
	"github.com/aenjoy/aemkvDB/backends/json"
	"github.com/aenjoy/aemkvDB/backends/redis"
)

func redis2() {
	var db aemkvDB.StrApi //- data type support: string
	//var db2 aemkvDB.AutoApi //auto detect backend - data type support :any
	db = redis.New(aemkvDB.ConfigMemDb{
		Addr:     []string{"127.0.0.1:6379"},
		Password: "123456",
		Database: "",
	})
	db.Set("key", "value")
	print(db.Get("key")) //out is value
}
func storeAnyTypeExample() {
	// support json(support store type:string int float64 bool),
	// map, syncMap (but all not support to store to disk).
	var db aemkvDB.AutoApi
	db = integ.New(aemkvDB.ConfigStandDb{
		Mode: 0, //0:map, 1:syncMap
	}).(*integ.MapDB)
	db.Set("key", "value")
	print(db.Get("key")) //out is value
}
func json2() {
	var db aemkvDB.AutoApi //- data type support: string
	//var db2 aemkvDB.AutoApi //auto detect backend - data type support :any
	db = json.New(aemkvDB.ConfigStandDb{Data: "./data.json", Mode: 0})
	db.Set("key", "value")
	print(db.Get("key")) //out is value
}
