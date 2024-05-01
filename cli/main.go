package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/aenjoy/aemkvDB"
	"github.com/aenjoy/aemkvDB/backends/buntdb"
	"github.com/aenjoy/aemkvDB/backends/etcd"
	"github.com/aenjoy/aemkvDB/backends/file"
	"github.com/aenjoy/aemkvDB/backends/integ"
	"github.com/aenjoy/aemkvDB/backends/json"
	"github.com/aenjoy/aemkvDB/backends/redis"
	"github.com/aenjoy/aemkvDB/backends/sql"
	"github.com/urfave/cli/v2"
	"os"
	"strconv"
)

var db aemkvDB.StrApi
var db2 aemkvDB.AutoApi
var nowMode = 2

func save() {
	// save data
	// ! 0 3 4  is not support save
	// and 3 4 5 use db2 -> 5 use db2
	if nowMode != 0 && nowMode != 3 && nowMode != 4 && nowMode != 5 && db.IsOpened() && settings["saveFlag"] == "1" {
		//save
		err := db.Close()
		if err != nil {
			fmt.Print("Error: Data save Error")
			//return
		}
	}
	if nowMode == 5 && db2.IsOpened() && settings["saveFlag"] == "1" {
		err := db2.Close()
		if err != nil {
			fmt.Print("Error: Data save Error")
			//return
		}
	}
}

func openDb(mode string) {
	atoi, _ := strconv.Atoi(mode)
	save()
	// create new database backend
	switch atoi {
	case 0:
		db = buntdb.NewBuntDB(":memory:")
	case 1:
		db = etcd.New(aemkvDB.ConfigMemDb{
			Addr:     []string{settings["addr"]},
			Password: settings["password"],
			Database: settings["database"],
			Ctx:      context.Background(),
		})
		if db == nil || !db.IsOpened() {
			fmt.Print("Error: Connect to etcd Server Error")
		}
	case 2:
		db = file.NewFileDB(settings["filepath"])
		if db == nil || !db.IsOpened() {
			fmt.Print("Error: Can't open file database")
		}
	case 3:
		db2 = integ.NewMapDB()
		db = nil
	case 4:
		db2 = integ.NewSyncMapDB()
		db = nil
	case 5:
		db2 = json.NewJsonDB(0, settings["filepath"])
		if db2 == nil || !db2.IsOpened() {
			fmt.Print("Error: Can't open json file")
		}
		db = nil
	case 6:
		db = redis.New(aemkvDB.ConfigMemDb{
			Addr:     []string{settings["addr"]},
			Password: settings["password"],
			Database: settings["database"],
		})
		if db == nil || !db.IsOpened() {
			fmt.Print("Error: Connect to redis Server Error")
		}
	case 7:
		db = sql.NewSqlDB(1, settings["filepath"])
		if db == nil || !db.IsOpened() {
			fmt.Print("Error: Connect to Sqlite Server Error")
		}
	case 8:
		db = sql.New(aemkvDB.SqlDb{
			Mode:     2,
			Addr:     settings["addr"],
			User:     settings["user"],
			Password: settings["password"],
			DataBase: settings["database"],
			Port:     settings["port"],
		})
		if db == nil || !db.IsOpened() {
			fmt.Print("Error: Connect to MySql Server Error")
		}
	case 9:
		db = sql.New(aemkvDB.SqlDb{
			Mode:     3,
			Addr:     settings["addr"],
			User:     settings["user"],
			Password: settings["password"],
			DataBase: settings["database"],
			Port:     settings["port"],
		})
		if db == nil || !db.IsOpened() {
			fmt.Print("Error: Connect to PostgresSql Server Error")
		}
	case 10:
		db = buntdb.NewBuntDB(settings["filepath"])
		if db == nil || !db.IsOpened() {
			fmt.Print("Error: Connect to BuntDB Server Error")
		}
	}

	//

	nowMode = atoi
}

var settings = map[string]string{
	"backends": "2",
	"saveFlag": "1",
	"addr":     "",
	"password": "",
	"user":     "",
	"database": "",
	"port":     "6379",
	"filepath": "",
}

func get(c *cli.Context) error {
	if nowMode != 3 && nowMode != 4 && nowMode != 5 {
		if db.IsOpened() {
			fmt.Printf("Key:%s Value:%s", c.Args().First(), db.Get(c.Args().First()))
		} else {
			fmt.Print("Error: Database is not opened")
			return errors.New("Database is not opened")
		}
	} else {
		if db2.IsOpened() {
			fmt.Printf("Key:%s Value:%s", c.Args().First(), db2.Get(c.Args().First()).(string))
		} else {
			fmt.Print("Error: Database is not opened")
			return errors.New("Database is not opened")
		}
	}
	return nil
}
func set(c *cli.Context) error {
	if nowMode != 3 && nowMode != 4 && nowMode != 5 {
		if db.IsOpened() {
			err := db.Set(c.Args().First(), c.Args().Get(1))
			if err != nil {
				fmt.Print("Error: Set Error")
				return err
			}
			fmt.Print("ok")
		} else {
			fmt.Print("Error: Database is not opened")
			return errors.New("Database is not opened")
		}
	} else {
		if db2.IsOpened() {
			err := db2.Set(c.Args().First(), c.Args().Get(1))
			if err != nil {
				fmt.Print("Error: Set Error")
				return err
			}
			fmt.Print("ok")
		} else {
			fmt.Print("Error: Database is not opened")
			return errors.New("Database is not opened")
		}
	}
	return nil
}
func deleteA(c *cli.Context) error {
	if nowMode != 3 && nowMode != 4 && nowMode != 5 {
		if db.IsOpened() {
			err := db.Delete(c.Args().First())
			if err != nil {
				fmt.Print("Error: Delete Error")
				return err
			}
			fmt.Print("ok")
		} else {
			fmt.Print("Error: Database is not opened")
			return errors.New("Database is not opened")
		}
	} else {
		if db2.IsOpened() {
			err := db2.Delete(c.Args().First())
			if err != nil {
				fmt.Print("Error: Delete Error")
				return err
			}
			fmt.Print("ok")
		} else {
			fmt.Print("Error: Database is not opened")
			return errors.New("Database is not opened")
		}
	}
	return nil
}
func setting(c *cli.Context) error {
	first := c.Args().First()
	second := c.Args().Get(1)
	switch first {
	case "backends":
		atoi, err := strconv.Atoi(second)
		if err != nil {
			fmt.Print("Error: Type Error. Value must be int")
			return err
		}
		if atoi >= 0 && atoi <= 10 {
			settings["backends"] = second
			openDb(second)
			fmt.Print("ok")
			return nil
		} else {
			fmt.Print("Error: Choose Backend Error. Value must be 0~10")
			return errors.New("Choose Backend Error. Value must be 0~10")
		}
	case "saveFlag":
		if second == "0" || second == "false" {
			settings["saveFlag"] = "0"
		} else if second == "1" || second == "true" {
			settings["saveFlag"] = "1"
		} else {
			fmt.Print("Error: Type Error. Value must be 0 or 1")
			return errors.New("Type Error. Value must be 0 or 1")
		}
	case "addr":
		settings["addr"] = second
		fmt.Print("ok")
	case "password":
		settings["password"] = second
		fmt.Print("ok")
	case "user":
		settings["user"] = second
		fmt.Print("ok")
	case "database":
		settings["database"] = second
		fmt.Print("ok")
	case "port":
		atoi, err := strconv.Atoi(second)
		if err != nil {
			fmt.Print("Error: Type Error. Value must be int")
			return err
		}
		if atoi >= 0 && atoi <= 65535 {
			settings["port"] = second
			fmt.Print("ok")
			return nil
		} else {
			fmt.Print("Error: Port Error. Value must be 0~65535")
			return errors.New("Port Error. Value must be 0~65535")
		}
	case "filepath":
		settings["filepath"] = second
		fmt.Print("ok")
	default:
		fmt.Print("Warning: Unsupport System attribute. However, it will still be recorded")
		settings[first] = second
	}
	return nil
}
func help(c *cli.Context) error {
	fmt.Print(helpText)
	return nil
}
func show(c *cli.Context) error {
	first := c.Args().First()
	if first == "all" {
		for key, value := range settings {
			fmt.Printf("%s:%s\n", key, value)
		}
		return nil
	}
	if first == "" {
		return errors.New("no key")
	}
	v, ok := settings[first]
	if ok {
		fmt.Printf("%s:%s\n", first, v)
		return nil
	}
	return errors.New("no key")
}
func quit(c *cli.Context) error {
	save()
	fmt.Print("bye")
	os.Exit(0)
	return nil
}
func load(c *cli.Context) error {
	settings["filepath"] = c.Args().First()
	openDb(strconv.Itoa(nowMode))
	return nil
}

func main() {
	app := cli.NewApp()
	app.Name = "aemkv"
	app.Usage = "aemkv api cli client "
	app.Version = "0.0.1"
	app.Commands = []*cli.Command{
		{
			Name:    "get",
			Aliases: []string{"g"},
			Usage:   "get key",
			Action:  get,
		},
		{
			Name:    "set",
			Aliases: []string{"s"},
			Usage:   "set key value",
			Action:  set,
		},
		{
			Name:    "delete",
			Aliases: []string{"d"},
			Usage:   "delete key",
			Action:  deleteA,
		},
		{
			Name:    ":set",
			Aliases: []string{"S"},
			Usage:   ":set xxx yyy ...",
			Action:  setting,
		},
		{
			Name:    ":quit",
			Aliases: []string{"q", "e"},
			Usage:   ":quit",
			Action:  quit,
		},
		{
			Name:    ":help",
			Aliases: []string{"h"},
			Usage:   ":help",
			Action:  help,
		},
		{
			Name:    ":show",
			Aliases: []string{"sh"},
			Usage:   ":show all/keyName show setting info and value",
			Action:  show,
		},
		{
			Name:    ":load",
			Aliases: []string{"L"},
			Usage:   ":load filepath(.db)",
			Action:  load,
		},
	}
	db = file.NewFileDB("data.db")
	app.Run(os.Args)
}
