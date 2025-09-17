package krap

import (
	"database/sql"
	"fmt"

	"github.com/roidaradal/rdb"
)

func (rq *Request) StartTransaction(numSteps int) error {
	if rq.DB == nil {
		rq.AddLog("No DB connection")
		return errNoDBConn
	}
	dbtx, err := rq.DB.Begin()
	if err != nil {
		rq.AddLog("Failed to start transaction")
		return err
	}
	rq.DBTx = dbtx
	rq.txSteps = make([]rdb.Query, 0, numSteps)
	rq.Checker = rdb.AssertNothing() // default result checker
	return nil
}

func (rq *Request) CommitTransaction() error {
	if rq.DB == nil {
		rq.AddLog("No DB connection")
		return errNoDBConn
	}
	if rq.DBTx == nil {
		rq.AddLog("No DB transaction")
		return errNoDBTx
	}
	err := rq.DBTx.Commit()
	if err != nil {
		// Log txSteps
		for _, query := range rq.txSteps {
			rq.AddLog("Query: " + rdb.QueryString(query))
		}
		rq.AddLog("Transaction commit failed")
		return fmt.Errorf("dbtx commit error: %w", err)
	}
	return nil
}

func (rq *Request) AddTxStep(query rdb.Query) {
	rq.txSteps = append(rq.txSteps, query)
}

func TxExec(rqtx *Request, q rdb.Query) error {
	rqtx.AddTxStep(q)
	_, err := rdb.ExecTx(q, rqtx.DBTx, rqtx.Checker)
	return err
}

func TxInsertRows[T any](rqtx *Request, items []*T, table string) error {
	numItems := len(items)
	rows := make([]map[string]any, numItems)
	for i, item := range items {
		rows[i] = rdb.ToRow(item)
	}
	q := rdb.NewInsertRowsQuery(table)
	q.Rows(rows)
	rqtx.AddTxStep(q)

	checker := rdb.AssertRowsAffected(numItems)
	_, err := rdb.ExecTx(q, rqtx.DBTx, checker)
	if err != nil {
		rqtx.AddFmtLog("Failed to insert %d rows to %s", numItems, table)
		return err
	}
	return nil
}

func TxInsertRow[T any](rqtx *Request, item *T, table string) error {
	_, err := txInsertRow(rqtx, item, table)
	return err
}

func TxInsertRowID[T any](rqtx *Request, item *T, table string) (uint, error) {
	result, err := txInsertRow(rqtx, item, table)
	if err != nil {
		return 0, err
	}
	return getInsertID(rqtx, result, table)
}

func txInsertRow[T any](rqtx *Request, item *T, table string) (*sql.Result, error) {
	q := rdb.NewInsertRowQuery(table)
	q.Row(rdb.ToRow(item))
	rqtx.AddTxStep(q)

	result, err := rdb.ExecTx(q, rqtx.DBTx, rqtx.Checker)
	if err != nil {
		rqtx.AddFmtLog("Failed to insert row to %s", table)
		return nil, err
	}
	return result, nil
}
