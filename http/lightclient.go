// Copyright © 2022 Attestant Limited.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package http

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"

	"github.com/attestantio/go-eth2-client/spec/capella"
	"github.com/attestantio/go-eth2-client/spec/phase0"
)

type lightClientBootstrapJSON struct {
	Data *capella.LightClientBootstrap `json:"data"`
}

// LightClientBootstrap provides the light client bootstrap of a given block ID.
func (s *Service) LightClientBootstrap(ctx context.Context, blockRoot phase0.Root) (*capella.LightClientBootstrap, error) {
	respBodyReader, err := s.get(ctx, fmt.Sprintf("/eth/v1/beacon/light_client/bootstrap/%#x", blockRoot))
	if err != nil {
		return nil, errors.Wrap(err, "failed to request beacon light client bootstrap")
	}
	if respBodyReader == nil {
		return nil, nil
	}

	var resp lightClientBootstrapJSON
	if err := json.NewDecoder(respBodyReader).Decode(&resp); err != nil {
		return nil, errors.Wrap(err, "failed to parse beacon light client bootstrap")
	}

	return resp.Data, nil
}

type lightClientUpdateJSON struct {
	Data *capella.LightClientUpdate `json:"data"`
}

// LightClientUpdates provides the light client update instances in the sync committee period range [startPeriod, startPeriod + count]
func (s *Service) LightClientUpdates(ctx context.Context, startPeriod, count uint64) ([]*capella.LightClientUpdate, error) {
	respBodyReader, err := s.get(
		ctx,
		fmt.Sprintf("/eth/v1/beacon/light_client/updates?start_period=%d&count=%d", startPeriod, count),
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to request beacon light client updates")
	}
	if respBodyReader == nil {
		return nil, nil
	}

	var resp []lightClientUpdateJSON
	if err := json.NewDecoder(respBodyReader).Decode(&resp); err != nil {
		return nil, errors.Wrap(err, "failed to parse beacon light client updates")
	}

	size := len(resp)
	ret := make([]*capella.LightClientUpdate, size, size)
	for i := 0; i < size; i++ {
		ret[i] = resp[i].Data
	}

	return ret, nil
}

type lightClientFinalityUpdateJSON struct {
	Data *capella.LightClientFinalityUpdate `json:"data"`
}

// LightClientFinalityUpdate provides the latest light client finality_update
func (s *Service) LightClientFinalityUpdate(ctx context.Context) (*capella.LightClientFinalityUpdate, error) {
	respBodyReader, err := s.get(ctx, "/eth/v1/beacon/light_client/finality_update")
	if err != nil {
		return nil, errors.Wrap(err, "failed to request beacon light client finality_update")
	}
	if respBodyReader == nil {
		return nil, nil
	}

	var resp lightClientFinalityUpdateJSON
	if err := json.NewDecoder(respBodyReader).Decode(&resp); err != nil {
		return nil, errors.Wrap(err, "failed to parse beacon light client finality_update")
	}

	return resp.Data, nil
}

type lightClientOptimisticUpdateJSON struct {
	Data *capella.LightClientOptimisticUpdate `json:"data"`
}

// LightClientOptimisticUpdate provides the latest light client optimistic_update
func (s *Service) LightClientOptimisticUpdate(ctx context.Context) (*capella.LightClientOptimisticUpdate, error) {
	respBodyReader, err := s.get(ctx, "/eth/v1/beacon/light_client/optimistic_update")
	if err != nil {
		return nil, errors.Wrap(err, "failed to request beacon light client optimistic_update")
	}
	if respBodyReader == nil {
		return nil, nil
	}

	var resp lightClientOptimisticUpdateJSON
	if err := json.NewDecoder(respBodyReader).Decode(&resp); err != nil {
		return nil, errors.Wrap(err, "failed to parse beacon light client optimistic_update")
	}

	return resp.Data, nil
}
