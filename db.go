package krap

import (
	"database/sql"
	"errors"

	"github.com/roidaradal/rdb"
)

var (
	errNoDBConn       = errors.New("no db connection")
	errNoDBTx         = errors.New("no db tx")
	errNoRowsInserted = errors.New("db no rows inserted")
	errNoLastInsertID = errors.New("no last insert id")
)

func DBInsertRow[T any](rq *Request, item *T, table string) error {
	_, err := dbInsertRow(rq, item, table)
	return err
}

func DBInsertRowID[T any](rq *Request, item *T, table string) (uint, error) {
	result, err := dbInsertRow(rq, item, table)
	if err != nil {
		return 0, err
	}
	return getInsertID(rq, result, table)
}

func dbInsertRow[T any](rq *Request, item *T, table string) (*sql.Result, error) {
	q := rdb.NewInsertRowQuery(table)
	q.Row(rdb.ToRow(item))
	result, err := rdb.Exec(q, rq.DB)
	if err != nil {
		rq.AddFmtLog("Failed to insert row to %s", table)
		return nil, err
	}
	if rdb.RowsAffected(result) != 1 {
		rq.AddFmtLog("No row inserted to %s", table)
		return nil, errNoRowsInserted
	}
	return result, nil
}

func getInsertID(rq *Request, result *sql.Result, table string) (uint, error) {
	id, ok := rdb.LastInsertID(result)
	if !ok {
		rq.AddFmtLog("Failed to get insert ID from %s", table)
		return 0, errNoLastInsertID
	}
	return id, nil
}
