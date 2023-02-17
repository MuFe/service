package db

import (
	"database/sql"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"mufe_service/camp/xlog"
)

type Tx struct {
	tx *sql.Tx
}

func (s *Tx) Query(sql string, args ...interface{}) (rows *sql.Rows, err error) {
	t := time.Now()
	rows, err = s.tx.Query(sql, args...)
	xlog.DB(false, time.Now().Sub(t), 0, sql, args...)
	if err != nil {
		xlog.ErrorP(err)
	}
	return rows, err
}

func (s *Tx) QueryRow(sql string, args ...interface{}) (result *sql.Row) {
	t := time.Now()
	result = s.tx.QueryRow(sql, args...)
	xlog.DB(false, time.Now().Sub(t), 0, sql, args...)
	return result
}

func (s *Tx) Exec(sql string, args ...interface{}) (result sql.Result, err error) {
	t := time.Now()
	result, err = s.tx.Exec(sql, args...)
	var affected int64
	if err == nil {
		affected, _ = result.RowsAffected()
	} else {
		xlog.ErrorP(err)
	}
	xlog.DB(true, time.Now().Sub(t), affected, sql, args...)
	return result, err
}
