package settings

import (
	"encoding"
	"encoding/json"
	"fmt"
)

type KV interface {
	Key() string
	encoding.BinaryMarshaler
	encoding.BinaryUnmarshaler
}

const (
	OrderSync = "OrderSync"
)

type OrderSyncKV struct {
	SyncWaitPaymentOrder bool `json:"sync_wait_payment_order"`
}

func (o *OrderSyncKV) Key() string {
	return OrderSync
}

func (o *OrderSyncKV) MarshalBinary() (data []byte, err error) {
	return json.Marshal(o)
}

func (o *OrderSyncKV) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, o)
}

// StrconvKV str转换为struct
func StrconvKV(k, v string) (KV, error) {
	switch k {
	case OrderSync:
		o := &OrderSyncKV{}
		err := o.UnmarshalBinary([]byte(v))
		if err != nil {
			return nil, fmt.Errorf("settings: fail to UnmarshalBinary OrderSyncKV in StrconvKV: %v", err)
		}
		return o, nil
	default:
		return nil, fmt.Errorf("settings: fail StrconvKV for unknown k '%s'", k)
	}
}
