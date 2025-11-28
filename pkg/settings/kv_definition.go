package settings

import "encoding/json"

type KV interface {
}

type OrderSyncKV struct {
	SuiteToSingle bool `json:"suite_to_single"`
}

type OrderSplitKV struct {
	VirtualGoodsNotAlone bool `json:"virtual_goods_not_alone"`
	FreebiesNotAlone     bool `json:"freebies_not_alone"`
}

type OrderMergeKV struct {
	MatchTagAfter bool `json:"match_tag_after"`
}

type Definition struct {
	Key         string
	Generator   func() KV
	DefaultJson string
}

var (
	OrderSync = &Definition{
		Key:       "OrderSync",
		Generator: func() KV { return &OrderSyncKV{} },
	}
	OrderSplit = &Definition{
		Key:       "OrderSplit",
		Generator: func() KV { return &OrderSplitKV{} },
	}
	OrderMerge = &Definition{
		Key:       "OrderMerge",
		Generator: func() KV { return &OrderMergeKV{} },
	}
)

var Keys = map[string]*Definition{
	OrderSync.Key:  OrderSync,
	OrderSplit.Key: OrderSplit,
	OrderMerge.Key: OrderMerge,
}

func init() {
	var bs []byte
	bs, _ = json.Marshal(OrderSync.Generator())
	OrderSync.DefaultJson = string(bs)

	bs, _ = json.Marshal(OrderSplit.Generator())
	OrderSplit.DefaultJson = string(bs)

	bs, _ = json.Marshal(OrderMerge.Generator())
	OrderMerge.DefaultJson = string(bs)
}
