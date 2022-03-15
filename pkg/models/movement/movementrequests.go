package movement

import (
	"encoding/json"
	"fmt"
)

// CancelRequest represents an order cancel request.
// An order can be cancelled using the internal ID or a
// combination of Client ID (CID) and the daten for the given
// CID.
type MovementRequest struct {
	Start *int64 `json:"start,omitempty"`
	End   *int64 `json:"end,omitempty"`
}

func (mr *MovementRequest) ToJSON() ([]byte, error) {
	resp := struct {
		Start *int64 `json:"start,omitempty"`
		End   *int64 `json:"end,omitempty"`
	}{
		Start: mr.Start,
		End:   mr.End,
	}

	return json.Marshal(resp)
}

// MarshalJSON converts the order cancel object into the format required by the
// bitfinex websocket service.
func (mr *MovementRequest) MarshalJSON() ([]byte, error) {
	b, err := mr.ToJSON()
	if err != nil {
		return nil, err
	}
	return []byte(fmt.Sprintf("[0, \"oc\", null, %s]", string(b))), nil
}
