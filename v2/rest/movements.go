package rest

import (
	"github.com/trever-io/bitfinex-api-go/pkg/models/common"
	"github.com/trever-io/bitfinex-api-go/pkg/models/movement"
)

type MovementService struct {
	requestFactory
	Synchronous
}

func (s *MovementService) Movements() (*movement.Snapshot, error) {
	req, err := s.requestFactory.NewAuthenticatedRequest(common.PermissionRead, "movements/hist")
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
