package main

const helpText = `All Help
Set SystemOption:  Using ":set "
Default:
	Key:	    Value:  Explain:
	backends	  2     cli default using mapDB backends(0~10) Notice: change backends with saveFlag=0 will potentially lead to data loss
	saveFlag      1     do not save data when exited(0,1,false,true)
	addr         [ ]	server backends address
	password	  ""    database passsword
	user		  ""	database user name
	database	  "" 	database name
	port         6379	database port
	filepath	  ""	file path of FileDB, BuntDB(File), JsonDB SqlLiteDB and BuntDB(File)

SystemOptionExplain:
-backends:(0~10)
	0 BuntDB(can't store data in the disk')
	1 EtcdDB
	2 FileDB
	3 MapDB(can't store data in the disk')
	4 SyncMapDB(can't store data in the disk')
	5 JsonDB
	6 RedisDB
	7 SqliteDB
	8 MySqlDB
	9 postgresSqlDB
	10 BuntDB(File)
	example: 
		cli > :set backends 10

-saveFlag:(0,1,false,true)
	0 do not save data when exited
	1 save data when exited
	example: 
		cli > :set saveFlag 0  :set saveFlag true
-addr: [ ] server address. if you use etcd, it allow use multiple address. Notice, if you use etcd, port is in here.
	example:
		cli > :set addr 10.1.1.2:2379 10.1.1.3:2379 10.1.1.4:2379 ... for etcd
		cli > :set addr 10.2.1.1 for redis,mysql and postgresSql.

`
