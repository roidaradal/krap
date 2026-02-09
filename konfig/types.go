package konfig

import "github.com/roidaradal/rdb/ze"

var (
	KVSchema *ze.Schema[KV]
)

type KV struct {
	Key           string `col:"AppKey"`
	Value         string `col:"AppValue"`
	LastUpdatedAt ze.DateTime
}
