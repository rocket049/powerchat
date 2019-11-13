package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

//db struct
//create table if not exists users (
//	id INTEGER PRIMARY KEY AUTOINCREMENT,
//	name text not null unique,
//	sex integer,
//	birthday DATE not null,
//	desc text,
//	pwdmd5 text not null);

var (
	db      *sql.DB
	dbMutex = &sync.Mutex{}
	dbstore string
)

func dbReconnect() error {
	var err error
	dbMutex.Lock()
	err = db.Ping()
	if err != nil {
		db.Close()
		db, err = sql.Open("sqlite3", filepath.Join(dbstore, "users.db"))
	}
	dbMutex.Unlock()
	return err
}

func init() {
	var err error
	exe1, _ := os.Executable()
	dbstore = filepath.Join(filepath.Dir(exe1), "dbstore")
	os.MkdirAll(dbstore, os.ModePerm)
	db, err = sql.Open("sqlite3", filepath.Join(dbstore, "users.db"))
	if err != nil {
		panic(err)
	}
	_, err = db.Exec("create table if not exists users (id INTEGER PRIMARY KEY AUTOINCREMENT,name text not null unique,sex integer,birthday DATE not null,desc text,pwdmd5 text not null);")
	if err != nil {
		panic(err)
	}
	go func() {
		for {
			time.Sleep(time.Hour * 24)
			dbReconnect()
		}
	}()
}

func dbClose() {
	dbMutex.Lock()
	db.Close()
	dbMutex.Unlock()
}

func updatePasswd(id int64, pwdmd5 string) error {
	dbMutex.Lock()
	defer dbMutex.Unlock()
	_, err := db.Exec("update users set pwdmd5=? where id=?;", pwdmd5, id)
	return err
}

func updateDesc(id int64, desc string) error {
	dbMutex.Lock()
	defer dbMutex.Unlock()
	_, err := db.Exec("update users set desc=? where id=?;", desc, id)
	return err
}

func deleteUser(id int64) error {
	dbMutex.Lock()
	defer dbMutex.Unlock()
	_, err := db.Exec("delete from users where id=?;", id)
	return err
}

//insertUser insert and return id
func insertUser(name string, sex int, birthday, desc string, pwdmd5 string) (int64, error) {
	name1 := strings.TrimSpace(name)
	if len(name) != len(name1) {
		return 0, errors.New("Name format error")
	}
	if len(pwdmd5) == 0 {
		return 0, errors.New("Password length error")
	}
	sql1 := "insert into users(name,sex,birthday,desc,pwdmd5) values(?,?,date(?),?,?);"
	dbMutex.Lock()
	defer dbMutex.Unlock()

	_, err := db.Exec(sql1, name, sex, birthday, desc, pwdmd5)
	if err != nil {
		return 0, err
	}
	sql2 := "select last_insert_rowid() from users"
	row := db.QueryRow(sql2)
	var res int64
	err = row.Scan(&res)

	return res, err
}

type UserInfo struct {
	Id       int64
	Name     string
	Sex      int
	Birthday time.Time
	Desc     string
	Pwdmd5   string
}

type UserBaseInfo struct {
	Id         int64
	Name       string
	Sex        int
	Birthday   time.Time
	Desc       string
	MsgOffline string
}

func getUsersByIds(ids []int64) (map[int64]*UserBaseInfo, error) {
	res := make(map[int64]*UserBaseInfo)
	sql1 := "select id,name,sex,birthday,desc from users where id in ("
	buf := bytes.NewBufferString(sql1)
	lenIds := len(ids)
	for i, v := range ids {
		buf.WriteString(fmt.Sprintf("%d", v))
		if i < (lenIds - 1) {
			buf.WriteString(",")
		}
	}
	buf.WriteString(");")
	dbMutex.Lock()
	rows, err := db.Query(buf.String())
	dbMutex.Unlock()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		r1 := new(UserBaseInfo)
		err := rows.Scan(&r1.Id, &r1.Name, &r1.Sex, &r1.Birthday, &r1.Desc)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		res[r1.Id] = r1
	}
	return res, nil
}

//getUser get user info from id
func getUserById(id int64) (res *UserInfo, err error) {
	res = new(UserInfo)
	sql1 := "select id,name,sex,birthday,desc,pwdmd5 from users where id=?;"
	dbMutex.Lock()
	row := db.QueryRow(sql1, id)
	dbMutex.Unlock()
	err = row.Scan(&res.Id, &res.Name, &res.Sex, &res.Birthday, &res.Desc, &res.Pwdmd5)
	return
}

func getUserByName(name string) (res *UserInfo, err error) {
	res = new(UserInfo)
	sql1 := "select id,name,sex,birthday,desc,pwdmd5 from users where name=?;"
	dbMutex.Lock()
	row := db.QueryRow(sql1, name)
	dbMutex.Unlock()
	err = row.Scan(&res.Id, &res.Name, &res.Sex, &res.Birthday, &res.Desc, &res.Pwdmd5)
	return
}

func searchUsers(key string) (map[int64]*UserBaseInfo, error) {
	res := make(map[int64]*UserBaseInfo)
	sql1 := "select id,name,sex,birthday,desc from users where INSTR(name,?)>0;"
	dbMutex.Lock()
	rows, err := db.Query(sql1, key)
	dbMutex.Unlock()

	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		r1 := new(UserBaseInfo)
		err := rows.Scan(&r1.Id, &r1.Name, &r1.Sex, &r1.Birthday, &r1.Desc)
		if err != nil {
			return nil, err
		}
		res[r1.Id] = r1
	}
	return res, nil
}

const create_frends = "create table if not exists friends (id INTEGER not null unique,name text not null unique,sex integer,birthday DATE not null,desc text,msgOffline text default '');"

func addFriend(uid, fid int64) error {
	if uid == 0 || fid == 0 {
		return errors.New("zero uid.")
	}
	dbMutex.Lock()
	defer dbMutex.Unlock()

	udb, err := sql.Open("sqlite3", filepath.Join(dbstore, fmt.Sprintf("%d", uid)))
	if err != nil {
		return err
	}
	defer udb.Close()
	_, err = udb.Exec(create_frends)
	if err != nil {
		return err
	}
	var (
		id       int64
		name     string
		sex      int
		birthday time.Time
		desc     string
	)
	sql1 := "select id,name,sex,birthday,desc from users where id=?;"
	sql2 := "insert into friends(id,name,sex,birthday,desc) values(?,?,?,?,?);"
	row := db.QueryRow(sql1, fid)
	if row == nil {
		return errors.New("Not exist.")
	}
	row.Scan(&id, &name, &sex, &birthday, &desc)
	_, err = udb.Exec(sql2, id, name, sex, birthday, desc)
	return err
}

func getFriends(uid int64) (map[int64]*UserBaseInfo, error) {
	if uid == 0 {
		return nil, errors.New("Not login.")
	}
	dbMutex.Lock()
	defer dbMutex.Unlock()

	udb, err := sql.Open("sqlite3", filepath.Join(dbstore, fmt.Sprintf("%d", uid)))
	if err != nil {
		return nil, err
	}
	defer udb.Close()
	_, err = udb.Exec(create_frends)
	if err != nil {
		return nil, err
	}
	res := make(map[int64]*UserBaseInfo)
	sql1 := "select id,name,sex,birthday,desc,msgOffline from friends;"
	rows, err := udb.Query(sql1)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		r1 := new(UserBaseInfo)
		err := rows.Scan(&r1.Id, &r1.Name, &r1.Sex, &r1.Birthday, &r1.Desc, &r1.MsgOffline)
		if err != nil {
			return nil, err
		}
		res[r1.Id] = r1
	}
	udb.Exec("update friends set msgOffline='' where msgOffline!='';")
	return res, nil
}

func getStrangers(uid int64) (map[int64]*UserBaseInfo, error) {
	if uid == 0 {
		return nil, errors.New("Not login.")
	}
	dbMutex.Lock()
	defer dbMutex.Unlock()
	udb, err := sql.Open("sqlite3", filepath.Join(dbstore, fmt.Sprintf("%d", uid)))
	if err != nil {
		return nil, err
	}
	defer udb.Close()
	_, err = udb.Exec(create_strangers)
	if err != nil {
		return nil, err
	}
	res := make(map[int64]*UserBaseInfo)
	sql1 := "select id,name,sex,birthday,desc,msgOffline from strangers where msgOffline!='';"
	rows, err := udb.Query(sql1)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		r1 := new(UserBaseInfo)
		err := rows.Scan(&r1.Id, &r1.Name, &r1.Sex, &r1.Birthday, &r1.Desc, &r1.MsgOffline)
		if err != nil {
			return nil, err
		}
		res[r1.Id] = r1
	}
	udb.Exec("update strangers set msgOffline='' where msgOffline!='';")
	return res, nil
}

type offlineMsgData struct {
	Msg       string
	Timestamp string
}

func offlineMsg(from, to int64, msg string) error {
	if from == 0 {
		return errors.New("Not login.")
	}
	if strings.HasPrefix(msg, "TEXT") == false {
		return errors.New("Message Type error.")
	}
	dbMutex.Lock()
	defer dbMutex.Unlock()

	udb, err := sql.Open("sqlite3", filepath.Join(dbstore, fmt.Sprintf("%d", to)))
	if err != nil {
		return err
	}
	defer udb.Close()
	_, err = udb.Exec(create_frends)
	if err != nil {
		return err
	}
	var offMsg = &offlineMsgData{Msg: msg[4:], Timestamp: time.Now().Format("2006-01-02 15:04:05")}
	bv, err := json.Marshal(offMsg)
	if err != nil {
		return err
	}
	var v = string(bv)
	res, _ := udb.Exec("update friends set msgOffline=? where id=?;", v, from)
	n, err := res.RowsAffected()
	log.Println("offline msg update:", n)
	if n == 0 {
		err = strangerMsg(udb, from, v)
	}
	return err
}

const create_strangers = "create table if not exists strangers (id INTEGER not null unique,name text not null unique,sex integer,birthday DATE not null,desc text,msgOffline text default '');"

//call before dbMutex.Lock() , so not need call it inside it
func strangerMsg(udb *sql.DB, from int64, msg string) error {
	udb.Exec(create_strangers)
	//update first
	sql_update := "update strangers set msgOffline=? where id=?"
	res, err := udb.Exec(sql_update, msg, from)
	n, err := res.RowsAffected()
	log.Println("stranger msg update", n)
	if n > 0 {
		return nil
	}
	//insert if fail update
	var (
		id       int64
		name     string
		sex      int
		birthday time.Time
		desc     string
	)
	sql1 := "select id,name,sex,birthday,desc from users where id=?;"
	sql2 := "insert into strangers(id,name,sex,birthday,desc,msgOffline) values(?,?,?,?,?,?);"
	row := db.QueryRow(sql1, from)
	if row == nil {
		return errors.New("Not exist.")
	}
	row.Scan(&id, &name, &sex, &birthday, &desc)
	res, err = udb.Exec(sql2, id, name, sex, birthday, desc, msg)
	n, err = res.RowsAffected()
	log.Println("stranger msg insert", n)
	return err
}

func moveStrangerToFriend(uid, fid int64) error {
	if uid == 0 || fid == 0 {
		return errors.New(fmt.Sprintf("Zero ID ,uid:%d, fid:%d", uid, fid))
	}
	dbMutex.Lock()
	defer dbMutex.Unlock()
	udb, err := sql.Open("sqlite3", filepath.Join(dbstore, fmt.Sprintf("%d", uid)))
	if err != nil {
		return err
	}
	defer udb.Close()
	var (
		id       int64
		name     string
		sex      int
		birthday time.Time
		desc     string
	)
	sql1 := "select id,name,sex,birthday,desc from users where id=?;"
	sql2 := "insert into friends(id,name,sex,birthday,desc) values(?,?,?,?,?);"
	row := db.QueryRow(sql1, fid)
	if row == nil {
		return errors.New("Not exist.")
	}
	row.Scan(&id, &name, &sex, &birthday, &desc)
	_, err = udb.Exec(sql2, id, name, sex, birthday, desc)
	if err == nil {
		_, err = udb.Exec("delete from strangers where id=?", fid)
	}
	return err
}

func removeFriend(uid, fid int64) error {
	if uid == 0 {
		return errors.New(fmt.Sprintf("Zero ID ,uid:%d", uid))
	}
	dbMutex.Lock()
	defer dbMutex.Unlock()
	udb, err := sql.Open("sqlite3", filepath.Join(dbstore, fmt.Sprintf("%d", uid)))
	if err != nil {
		return err
	}
	defer udb.Close()
	udb.Exec("delete from friends where id=?", fid)
	return nil
}

func isFriend(from, to int64) bool {
	if from == 0 || to == 0 {
		return false
	}

	dbMutex.Lock()
	defer dbMutex.Unlock()
	udb, err := sql.Open("sqlite3", filepath.Join(dbstore, fmt.Sprintf("%d", to)))
	if err != nil {
		log.Println(err)
		return false
	}
	defer udb.Close()
	row := udb.QueryRow("select name from friends where id=?;", from)
	var name1 string
	err = row.Scan(&name1)
	if err != nil {
		log.Println(err)
		return false
	} else {
		return true
	}
}
