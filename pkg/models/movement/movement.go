package movement

import (
	"fmt"

	"github.com/trever-io/bitfinex-api-go/pkg/convert"
)

type Movement struct {
	Id                      int64
	Currency                string
	CurrencyName            string
	MTSStarted              int64
	MTSUpdated              int64
	Status                  string
	Amount                  float64
	Fees                    float64
	DestinationAddress      string
	TransactionId           string
	WithdrawTransactionNote string
}

type Snapshot struct {
	Snapshot []*Movement
}

func FromRaw(raw []interface{}) (m *Movement, err error) {
	if len(raw) < 22 {
		err = fmt.Errorf("data slice too short for movement: %#v", raw)
		return
	}

	m = &Movement{
		Id:                      convert.I64ValOrZero(raw[0]),
		Currency:                convert.SValOrEmpty(raw[1]),
		CurrencyName:            convert.SValOrEmpty(raw[2]),
		MTSStarted:              convert.I64ValOrZero(raw[5]),
		MTSUpdated:              convert.I64ValOrZero(raw[6]),
		Status:                  convert.SValOrEmpty(raw[9]),
		Amount:                  convert.F64ValOrZero(raw[12]),
		Fees:                    convert.F64ValOrZero(raw[13]),
		DestinationAddress:      convert.SValOrEmpty(raw[16]),
		TransactionId:           convert.SValOrEmpty(raw[20]),
		WithdrawTransactionNote: convert.SValOrEmpty(raw[21]),
	}
	return
}

func SnapshotFromRaw(raw []interface{}) (s *Snapshot, err error) {
	if len(raw) == 0 {
		return &Snapshot{Snapshot: make([]*Movement, 0)}, nil
	}

	ms := make([]*Movement, 0)
	switch raw[0].(type) {
	case []interface{}:
		for _, v := range raw {
			if l, ok := v.([]interface{}); ok {
				m, err := FromRaw(l)
				if err != nil {
					return s, err
				}
				ms = append(ms, m)
			}
		}
	default:
		return s, fmt.Errorf("not an wallet snapshot")
	}

	s = &Snapshot{Snapshot: ms}
	return
}
