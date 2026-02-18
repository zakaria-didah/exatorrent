//go:build cgo
// +build cgo

package db

import (
	"sync"

	sqlite "github.com/go-llsqlite/crawshaw"
	"github.com/go-llsqlite/crawshaw/sqlitex"
)

type CategoryDb interface {
	Open(string)
	Close()
	Set(infohash, category string) error
	Get(infohash string) string
	GetAll() map[string]string
}

type SqliteCategoryDb struct {
	Db *sqlite.Conn
	mu sync.Mutex
}

func (db *SqliteCategoryDb) Open(fp string) {
	var err error
	db.Db, err = sqlite.OpenConn(fp, 0)
	if err != nil {
		DbL.Fatalln(err)
	}
	err = sqlitex.ExecScript(db.Db, `create table if not exists categories (infohash text unique, category text);`)
	if err != nil {
		DbL.Fatalln(err)
	}
}

func (db *SqliteCategoryDb) Close() {
	db.mu.Lock()
	defer db.mu.Unlock()
	if err := db.Db.Close(); err != nil {
		DbL.Fatalln(err)
	}
}

func (db *SqliteCategoryDb) Set(infohash, category string) error {
	db.mu.Lock()
	defer db.mu.Unlock()
	if category == "" {
		return sqlitex.Exec(db.Db, `delete from categories where infohash=?;`, nil, infohash)
	}
	return sqlitex.Exec(db.Db, `insert or replace into categories (infohash, category) values (?, ?);`, nil, infohash, category)
}

func (db *SqliteCategoryDb) Get(infohash string) string {
	var cat string
	db.mu.Lock()
	defer db.mu.Unlock()
	_ = sqlitex.Exec(db.Db, `select category from categories where infohash=?;`,
		func(stmt *sqlite.Stmt) error {
			cat = stmt.GetText("category")
			return nil
		}, infohash)
	return cat
}

func (db *SqliteCategoryDb) GetAll() map[string]string {
	ret := make(map[string]string)
	db.mu.Lock()
	defer db.mu.Unlock()
	_ = sqlitex.Exec(db.Db, `select infohash, category from categories;`,
		func(stmt *sqlite.Stmt) error {
			ret[stmt.GetText("infohash")] = stmt.GetText("category")
			return nil
		})
	return ret
}
