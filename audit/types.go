package audit

import (
	"errors"

	"github.com/roidaradal/rdb/ze"
)

var errMissingSchema = errors.New("schema is not initialized")

var (
	ActionLogs    *ze.Schema[ActionLog]
	BatchLogs     *ze.Schema[BatchLog]
	BatchLogItems *ze.Schema[BatchLogItem]
)

type ActionDetails struct {
	Action  string
	Details string
}

type ActionLog struct {
	ze.CreatedItem
	ActorID    ze.ID  `json:"-"`
	ActorCode_ string `col:"-" json:"ActorCode"`
	ActionDetails
}

type BatchLog struct {
	ze.CodedItem
	ze.CreatedItem
	ActionDetails
}

type BatchLogItem struct {
	ze.CodedItem
	Details string
}
