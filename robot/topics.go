package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"path/filepath"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

var (
	db      *sql.DB
	dbstore string
	dbmutex *sync.Mutex = &sync.Mutex{}
)

func init() {
	const initSql = "create table if not exists topics(id int unique, title text unique);"
	var err error
	dbstore = getRelatePath(".")
	db, err = sql.Open("sqlite3", filepath.Join(dbstore, "topics.db"))
	if err != nil {
		panic(err)
	}
	_, err = db.Exec(initSql)
	if err != nil {
		panic(err)
	}
}

func getPage(n int) string {
	const selectSql = "select id,title from topics order by id asc limit 10 offset ?;"
	dbmutex.Lock()
	res, err := db.Query(selectSql, n*10)
	dbmutex.Unlock()
	if err != nil {
		return ""
	}
	buf := bytes.NewBufferString(fmt.Sprintf("当前页码：p%d\n", n))
	buf.WriteString("请发送问题编号(Send Question ID):\n\n")
	for res.Next() {
		var id1 int
		var title1 string
		err = res.Scan(&id1, &title1)
		if err == nil {
			buf.WriteString(fmt.Sprintf("%d. %s\n", id1, title1))
		}
	}
	buf.WriteString(fmt.Sprintf("\n发送 'p%d' 显示下一页\n(Send 'p%d' to show next page)\n", n+1, n+1))
	return buf.String()
}
