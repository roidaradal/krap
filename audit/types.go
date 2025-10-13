package audit

import "github.com/roidaradal/rdb/ze"

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
	ze.CreatedItem
	ze.CodedItem
	ActionDetails
}

type BatchLogItem struct {
	ze.CodedItem
	Details string
}
