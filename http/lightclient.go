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
	"bytes"
	"context"
	"fmt"

	"github.com/pkg/errors"

	"github.com/attestantio/go-eth2-client/api"
	"github.com/attestantio/go-eth2-client/spec/deneb"
)

// LightClientBootstrap provides the light client bootstrap of a given block ID.
func (s *Service) LightClientBootstrap(ctx context.Context, opts *api.LightClientBootstrapOpts) (
	*api.Response[*deneb.LightClientBootstrap],
	error,
) {
	url := fmt.Sprintf("/eth/v1/beacon/light_client/bootstrap/%s", opts.Block)
	resp, err := s.get(ctx, url, &opts.Common)
	if err != nil {
		return nil, errors.Wrap(err, "failed to request beacon light client bootstrap")
	}
	if resp == nil {
		return nil, nil
	}

	//if err := json.NewDecoder(bytes.NewReader(resp.body)).Decode(&resp); err != nil {
	data, metadata, err := decodeJSONResponse(bytes.NewReader(resp.body), deneb.LightClientBootstrap{})
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse beacon light client bootstrap")
	}

	return &api.Response[*deneb.LightClientBootstrap]{
		Data:     &data,
		Metadata: metadata,
	}, nil
}

// LightClientUpdates provides the light client update instances in the sync committee period range [startPeriod, startPeriod + count]
func (s *Service) LightClientUpdates(ctx context.Context, opts *api.LightClientUpdatesOpts) (
	*api.Response[[]*deneb.LightClientUpdate],
	error,
) {
	url := fmt.Sprintf("/eth/v1/beacon/light_client/updates?start_period=%d&count=%d", opts.StartPeriod, opts.Count)
	resp, err := s.get(ctx, url, &opts.Common)
	if err != nil {
		return nil, errors.Wrap(err, "failed to request beacon light client updates")
	}
	if resp == nil {
		return nil, nil
	}

	data, metadata, err := decodeJSONResponse(bytes.NewReader(resp.body), []*deneb.LightClientUpdate{})
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse beacon light client updates")
	}

	return &api.Response[[]*deneb.LightClientUpdate]{
		Data:     data,
		Metadata: metadata,
	}, nil
}

// LightClientFinalityUpdate provides the latest light client finality_update
func (s *Service) LightClientFinalityUpdate(ctx context.Context, opts *api.CommonOpts) (
	*api.Response[*deneb.LightClientFinalityUpdate],
	error,
) {
	resp, err := s.get(ctx, "/eth/v1/beacon/light_client/finality_update", opts)
	if err != nil {
		return nil, errors.Wrap(err, "failed to request beacon light client finality_update")
	}
	if resp == nil {
		return nil, nil
	}

	data, metadata, err := decodeJSONResponse(bytes.NewReader(resp.body), deneb.LightClientFinalityUpdate{})
	if err != nil {
		//var resp lightClientFinalityUpdateJSON
		//if err := json.NewDecoder(bytes.NewReader(resp.body)).Decode(&resp); err != nil {
		return nil, errors.Wrap(err, "failed to parse beacon light client finality_update")
	}

	return &api.Response[*deneb.LightClientFinalityUpdate]{
		Data:     &data,
		Metadata: metadata,
	}, nil
}

// LightClientOptimisticUpdate provides the latest light client optimistic_update
func (s *Service) LightClientOptimisticUpdate(ctx context.Context, opts *api.CommonOpts) (
	*api.Response[*deneb.LightClientOptimisticUpdate],
	error,
) {
	resp, err := s.get(ctx, "/eth/v1/beacon/light_client/optimistic_update", opts)
	if err != nil {
		return nil, errors.Wrap(err, "failed to request beacon light client optimistic_update")
	}
	if resp == nil {
		return nil, nil
	}

	data, metadata, err := decodeJSONResponse(bytes.NewReader(resp.body), deneb.LightClientOptimisticUpdate{})
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse beacon light client optimistic_update")
	}

	return &api.Response[*deneb.LightClientOptimisticUpdate]{
		Data:     &data,
		Metadata: metadata,
	}, nil
}
