// Copyright © 2021 Attestant Limited.
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

package multi

import (
	"context"

	consensusclient "github.com/attestantio/go-eth2-client"
	api "github.com/attestantio/go-eth2-client/api/v1"
)

// LightClientBootstrap provides the light client bootstrap of a given block ID.
func (s *Service) LightClientBootstrap(ctx context.Context, blockID string) (*api.LightClientBootstrap, error) {
	res, err := s.doCall(ctx, func(ctx context.Context, client consensusclient.Service) (interface{}, error) {
		bootstrap, err := client.(consensusclient.LightClientProvider).LightClientBootstrap(ctx, blockID)
		if err != nil {
			return nil, err
		}
		return bootstrap, nil
	}, nil)
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, nil
	}
	return res.(*api.LightClientBootstrap), nil
}

// LightClientUpdates provides the light client updates of a given start_period and count
func (s *Service) LightClientUpdates(ctx context.Context, start, count uint64) ([]*api.LightClientUpdate, error) {
	res, err := s.doCall(ctx, func(ctx context.Context, client consensusclient.Service) (interface{}, error) {
		updates, err := client.(consensusclient.LightClientProvider).LightClientUpdates(ctx, start, count)
		if err != nil {
			return nil, err
		}
		return updates, nil
	}, nil)
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, nil
	}
	return res.([]*api.LightClientUpdate), nil
}

// LightClientFinalityUpdate provides the light client finality_update
func (s *Service) LightClientFinalityUpdate(ctx context.Context) (*api.LightClientFinalityUpdate, error) {
	res, err := s.doCall(ctx, func(ctx context.Context, client consensusclient.Service) (interface{}, error) {
		bootstrap, err := client.(consensusclient.LightClientProvider).LightClientFinalityUpdate(ctx)
		if err != nil {
			return nil, err
		}
		return bootstrap, nil
	}, nil)
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, nil
	}
	return res.(*api.LightClientFinalityUpdate), nil
}

// LightClientOptimisticUpdate provides the light client optimistic_update
func (s *Service) LightClientOptimisticUpdate(ctx context.Context) (*api.LightClientOptimisticUpdate, error) {
	res, err := s.doCall(ctx, func(ctx context.Context, client consensusclient.Service) (interface{}, error) {
		bootstrap, err := client.(consensusclient.LightClientProvider).LightClientOptimisticUpdate(ctx)
		if err != nil {
			return nil, err
		}
		return bootstrap, nil
	}, nil)
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, nil
	}
	return res.(*api.LightClientOptimisticUpdate), nil
}
