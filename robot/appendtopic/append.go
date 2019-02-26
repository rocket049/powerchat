package main

import (
	"bufio"
	"database/sql"
	"errors"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

var (
	db      *sql.DB
	dbstore string
)

func init() {
	const initSql = "create table if not exists topics(id int unique, title text unique);"
	var err error
	exe1, _ := os.Executable()
	dbstore = filepath.Join(filepath.Dir(exe1), "..", "data")
	os.MkdirAll(dbstore, os.ModePerm)
	db, err = sql.Open("sqlite3", filepath.Join(dbstore, "topics.db"))
	if err != nil {
		panic(err)
	}
	_, err = db.Exec(initSql)
	if err != nil {
		panic(err)
	}
}

func appendTopic(filename string) error {
	const sqlAppend = "insert into topics(id,title) values(?,?);"
	var id1 = filepath.Base(filename)
	fp1, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer fp1.Close()
	var rd1 = bufio.NewReader(fp1)
	line1, _, err := rd1.ReadLine()
	if err != nil {
		return err
	}
	res, err := db.Exec(sqlAppend, id1, line1)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return errors.New("Unknown error")
	}
	return nil
}

func main() {
	if len(os.Args) < 2 {
		log.Panicln("usage: appendtopic <file1 file2 ...>")
		return
	}
	for i := 1; i < len(os.Args); i++ {
		err := appendTopic(os.Args[i])
		if err == nil {
			log.Println("OK", filepath.Base(os.Args[i]))
		} else {
			log.Println("Fail", filepath.Base(os.Args[i]))
		}
	}
}
