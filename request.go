package krap

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/roidaradal/fn/clock"
	"github.com/roidaradal/fn/net"
	"github.com/roidaradal/rdb"
)

var (
	dbConnParams *net.SQLConnParams = nil
	dbConn       *sql.DB            = nil
	rqDivider    string             = strings.Repeat("-", 30)
)

type Request struct {
	DB      *sql.DB
	DBTx    *sql.Tx
	Checker rdb.QueryResultChecker
	logs    []string
	txSteps []rdb.Query
}

type RequestOrigin struct {
	BrowserInfo *string
	IPAddress   *string
}

func SetSQLConnParams(params *net.SQLConnParams) {
	dbConnParams = params
}

func NewRequest(name string, args ...any) (*Request, error) {
	if len(args) > 0 {
		name = fmt.Sprintf(name, args...)
	}
	if dbConn == nil {
		var err error
		dbConn, err = net.NewSQLConnection(dbConnParams)
		if err != nil {
			return nil, err
		}
	}
	rq := &Request{}
	rq.DB = dbConn
	rq.logs = []string{nowLog(name)}
	return rq, nil
}

func (rq Request) Output() string {
	return fmt.Sprintf("%s\n%s\n%s", rqDivider, strings.Join(rq.logs, "\n"), rqDivider)
}

func (rq *Request) AddLog(message string) {
	rq.logs = append(rq.logs, nowLog(message))
}

func (rq *Request) AddFmtLog(format string, args ...any) {
	rq.AddLog(fmt.Sprintf(format, args...))
}

func (rq *Request) AddDurationLog(start time.Time) {
	rq.AddFmtLog("Time: %v", time.Since(start))
}

func (rq *Request) AddErrorLog(err error) {
	rq.AddFmtLog("Error: %s", err.Error())
}

func nowLog(message string) string {
	return fmt.Sprintf("%s %s", clock.DateTimeNow(), message)
}
