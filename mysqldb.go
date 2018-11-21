package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"sync"
)

var dbs *sql.DB
var onceInitMysql sync.Once

func mysqlInit(host, port, user, passwd string) error {
	// db, err := sql.Open("mysql", "root:15801250037@tcp(172.16.165.129:3306)/ntrip?charset=utf8")
	// setd = fmt.Sprintf("%s:%s@tcp(%s:%s)/ntrip?charset=utf8", user, passwd, host, port)
	var err error = nil
	onceInitMysql.Do(func() {
		dbs, err = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/ntrip", user, passwd, host, port))
		// dbs = db
	})

	return err
}

func testmysql() {
	rows, err := dbs.Query("SELECT loginname,password FROM rover")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("2")
	for rows.Next() {
		var loginname string
		var password string
		rows.Scan(&loginname, &password)
		fmt.Println(loginname, password)
	}

}

func mountpointVer(name, password string) bool {
	var rows *sql.Rows
	var err error
	defer func() {
		if rows != nil {
			rows.Close()
		}

	}()
	// rows, _ := dbs.db.Query("SELECT id FROM mountpoint where name=\"%s\" and password=\"%s\"", name, password)
	sqls := fmt.Sprintf("SELECT id FROM mountpoint where name=\"%s\" and password=\"%s\"", name, password)
	// fmt.Println(sqls)
	// rows, _ := dbs.db.Query("SELECT name,password FROM mountpoint")
	rows, err = dbs.Query(sqls)
	if err != nil {
		return false
	}
	return rows.Next()
}

func clientVer(name, password string) bool {
	var rows *sql.Rows
	var err error
	defer func() {
		if rows != nil {
			rows.Close()
		}

	}()
	sqls := fmt.Sprintf("SELECT id FROM rover where loginname=\"%s\" and password=\"%s\"", name, password)
	// fmt.Println(sqls)
	rows, err = dbs.Query(sqls)
	if err != nil {
		return false
	}
	return rows.Next()
}

func setStatus(typee, name, st string) bool {
	var rows *sql.Rows
	var rows1 *sql.Rows
	var err error
	sqls := ""
	sqls2 := ""

	defer func() {
		if rows != nil {
			rows.Close()
		}
		if rows1 != nil {
			rows1.Close()
		}
	}()
	switch typee {
	case "rover":
		sqls = fmt.Sprintf("update rover set status=\"%s\" where loginname=\"%s\"", st, name)
	case "mountpoint":
		sqls = fmt.Sprintf("update mountpoint set status=\"%s\" where name=\"%s\" ", st, name)
	case "all":
		sqls = fmt.Sprintf("update rover set status=\"n\"  ")
		sqls2 = fmt.Sprintf("update mountpoint set status=\"n\"  ")
	default:
		return false

	}
	if typee == "all" {
		rows1, _ = dbs.Query(sqls2)
	}
	// fmt.Println(sqls)
	rows, err = dbs.Query(sqls)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}
