//go:build cgo
// +build cgo

package core

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/go-llsqlite/crawshaw/sqlitex"
	db "github.com/zakaria-didah/exatorrent/internal/db"
)

// addUserWithHash inserts a user with a pre-hashed password (from signup request).
func addUserWithHash(username, hashedPassword string, userType int) error {
	if len(username) <= 5 {
		return fmt.Errorf("username too short")
	}
	udb, ok := Engine.UDb.(*db.Sqlite3UserDb)
	if !ok {
		return fmt.Errorf("unsupported user database type")
	}
	udb.Mu.Lock()
	defer udb.Mu.Unlock()
	token, err := uuid.NewRandom()
	if err != nil {
		return fmt.Errorf("uuid error")
	}
	return sqlitex.Exec(udb.Db, `insert into userdb (username,password,token,usertype,createdat) values (?,?,?,?,?);`, nil, username, hashedPassword, token.String(), userType, time.Now().Format(time.RFC3339))
}
