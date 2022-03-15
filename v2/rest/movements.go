package rest

import (
	"github.com/trever-io/bitfinex-api-go/pkg/models/common"
	"github.com/trever-io/bitfinex-api-go/pkg/models/movement"
)

type MovementService struct {
	requestFactory
	Synchronous
}

func (s *MovementService) Movements(mr *movement.MovementRequest) (*movement.Snapshot, error) {
	bytes, err := mr.ToJSON()
	if err != nil {
		return nil, err
	}
	req, err := s.requestFactory.NewAuthenticatedRequestWithBytes(common.PermissionRead, "movements/hist", bytes)
	if err != nil {
		return nil, err
	}

	raw, err := s.Request(req)
	if err != nil {
		return nil, err
	}

	os, err := movement.SnapshotFromRaw(raw)
	if err != nil {
		return nil, err
	}

	return os, nil
}
