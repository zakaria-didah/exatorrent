//go:build cgo
// +build cgo

package db

import (
	"fmt"
	"sync"
	"time"

	"golang.org/x/crypto/bcrypt"

	sqlite "github.com/go-llsqlite/crawshaw"
	"github.com/go-llsqlite/crawshaw/sqlitex"
)

type SignupRequest struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	Password  string    `json:"-"`
	Message   string    `json:"message"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

type SignupRequestDb interface {
	Open(string)
	Close()
	Add(username, password, message string) error
	GetPending() []*SignupRequest
	Approve(id int64) (*SignupRequest, error)
	Decline(id int64) error
	HasPending(username string) bool
}

type SqliteSignupRequestDb struct {
	Db *sqlite.Conn
	mu sync.Mutex
}

func (db *SqliteSignupRequestDb) Open(fp string) {
	var err error
	db.Db, err = sqlite.OpenConn(fp, 0)
	if err != nil {
		DbL.Fatalln(err)
	}
	err = sqlitex.ExecScript(db.Db, `create table if not exists signuprequests (id integer primary key autoincrement, username text unique, password text, message text, status text default 'pending', createdat text);`)
	if err != nil {
		DbL.Fatalln(err)
	}
}

func (db *SqliteSignupRequestDb) Close() {
	db.mu.Lock()
	defer db.mu.Unlock()
	if err := db.Db.Close(); err != nil {
		DbL.Fatalln(err)
	}
}

func (db *SqliteSignupRequestDb) Add(username, password, message string) error {
	if len(username) <= 5 || len(password) <= 5 {
		return fmt.Errorf("username or password size too small")
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return err
	}
	db.mu.Lock()
	defer db.mu.Unlock()
	return sqlitex.Exec(db.Db, `insert into signuprequests (username,password,message,status,createdat) values (?,?,?,?,?);`, nil, username, string(hashed), message, "pending", time.Now().Format(time.RFC3339))
}

func (db *SqliteSignupRequestDb) GetPending() []*SignupRequest {
	ret := make([]*SignupRequest, 0)
	db.mu.Lock()
	defer db.mu.Unlock()
	_ = sqlitex.Exec(db.Db, `select id,username,message,status,createdat from signuprequests where status='pending' order by createdat desc;`,
		func(stmt *sqlite.Stmt) error {
			var sr SignupRequest
			sr.ID = stmt.ColumnInt64(0)
			sr.Username = stmt.GetText("username")
			sr.Message = stmt.GetText("message")
			sr.Status = stmt.GetText("status")
			t, err := time.Parse(time.RFC3339, stmt.GetText("createdat"))
			if err == nil {
				sr.CreatedAt = t
			}
			ret = append(ret, &sr)
			return nil
		})
	return ret
}

func (db *SqliteSignupRequestDb) Approve(id int64) (*SignupRequest, error) {
	db.mu.Lock()
	defer db.mu.Unlock()
	var sr SignupRequest
	var found bool
	err := sqlitex.Exec(db.Db, `select id,username,password,message from signuprequests where id=? and status='pending';`,
		func(stmt *sqlite.Stmt) error {
			found = true
			sr.ID = stmt.ColumnInt64(0)
			sr.Username = stmt.GetText("username")
			sr.Password = stmt.GetText("password")
			sr.Message = stmt.GetText("message")
			return nil
		}, id)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, fmt.Errorf("signup request not found or already processed")
	}
	err = sqlitex.Exec(db.Db, `delete from signuprequests where id=?;`, nil, id)
	if err != nil {
		return nil, err
	}
	return &sr, nil
}

func (db *SqliteSignupRequestDb) Decline(id int64) error {
	db.mu.Lock()
	defer db.mu.Unlock()
	return sqlitex.Exec(db.Db, `delete from signuprequests where id=?;`, nil, id)
}

func (db *SqliteSignupRequestDb) HasPending(username string) bool {
	var exists bool
	db.mu.Lock()
	defer db.mu.Unlock()
	_ = sqlitex.Exec(db.Db, `select exists(select 1 from signuprequests where username=? and status='pending');`,
		func(stmt *sqlite.Stmt) error {
			exists = stmt.ColumnInt(0) != 0
			return nil
		}, username)
	return exists
}
