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
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"

	"github.com/attestantio/go-eth2-client/api"
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/altair"
	"github.com/attestantio/go-eth2-client/spec/capella"
	"github.com/attestantio/go-eth2-client/spec/deneb"
)

// LightClientBootstrap provides the light client bootstrap of a given block ID.
func (s *Service) LightClientBootstrap(ctx context.Context, opts *api.LightClientBootstrapOpts) (
	*api.Response[*spec.VersionedLCBootstrap],
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

	switch resp.contentType {
	case ContentTypeSSZ:
		return s.lightClientBootstrapFromSSZ(resp)
	case ContentTypeJSON:
		return s.lightClientBootstrapFromJSON(resp)
	default:
		return nil, fmt.Errorf("unhandled content type %v", resp.contentType)
	}
}

func (s *Service) lightClientBootstrapFromSSZ(res *httpResponse) (
	*api.Response[*spec.VersionedLCBootstrap],
	error,
) {
	response := &api.Response[*spec.VersionedLCBootstrap]{
		Data: &spec.VersionedLCBootstrap{
			Version: res.consensusVersion,
		},
		Metadata: metadataFromHeaders(res.headers),
	}

	var err error
	switch res.consensusVersion {
	case spec.DataVersionPhase0:
		err = fmt.Errorf("unsupported version %s", res.consensusVersion)
	case spec.DataVersionAltair:
		response.Data.Altair = &altair.LightClientBootstrap{}
		err = response.Data.Altair.UnmarshalSSZ(res.body)
	case spec.DataVersionBellatrix:
		response.Data.Bellatrix = &altair.LightClientBootstrap{}
		err = response.Data.Bellatrix.UnmarshalSSZ(res.body)
	case spec.DataVersionCapella:
		response.Data.Capella = &capella.LightClientBootstrap{}
		err = response.Data.Capella.UnmarshalSSZ(res.body)
	case spec.DataVersionDeneb:
		response.Data.Deneb = &deneb.LightClientBootstrap{}
		err = response.Data.Deneb.UnmarshalSSZ(res.body)
	default:
		err = fmt.Errorf("unsupported version %s", res.consensusVersion)
	}
	if err != nil {
		return nil, errors.Wrapf(err, "failed to decode %s light client optimistic update", res.consensusVersion)
	}

	return response, nil
}

func (s *Service) lightClientBootstrapFromJSON(res *httpResponse) (
	*api.Response[*spec.VersionedLCBootstrap],
	error,
) {
	response := &api.Response[*spec.VersionedLCBootstrap]{
		Data: &spec.VersionedLCBootstrap{
			Version: res.consensusVersion,
		},
	}

	var err error
	switch res.consensusVersion {
	case spec.DataVersionPhase0:
		err = fmt.Errorf("unsupported version %s", res.consensusVersion)
	case spec.DataVersionAltair:
		response.Data.Altair, response.Metadata, err = decodeJSONResponse(bytes.NewReader(res.body), &altair.LightClientBootstrap{})
	case spec.DataVersionBellatrix:
		response.Data.Bellatrix, response.Metadata, err = decodeJSONResponse(bytes.NewReader(res.body), &altair.LightClientBootstrap{})
	case spec.DataVersionCapella:
		response.Data.Capella, response.Metadata, err = decodeJSONResponse(bytes.NewReader(res.body), &capella.LightClientBootstrap{})
	case spec.DataVersionDeneb:
		response.Data.Deneb, response.Metadata, err = decodeJSONResponse(bytes.NewReader(res.body), &deneb.LightClientBootstrap{})
	default:
		err = fmt.Errorf("unsupported version %s", res.consensusVersion)
	}
	if err != nil {
		return nil, err
	}

	return response, nil
}

// LightClientUpdates provides the light client update instances in the sync committee period range [startPeriod, startPeriod + count]
func (s *Service) LightClientUpdates(ctx context.Context, opts *api.LightClientUpdatesOpts) (
	*api.Response[[]*spec.VersionedLCUpdate],
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

	switch resp.contentType {
	case ContentTypeJSON:
		return s.lightClientUpdatesFromJSON(resp)
	default:
		return nil, fmt.Errorf("unhandled content type %v", resp.contentType)
	}
}

type lcUpdatesJSON struct {
	Data    []interface{} `json:"data"`
	Version string        `json:"version,omitempty"`
}

func (s *Service) lightClientUpdatesFromJSON(res *httpResponse) (
	*api.Response[[]*spec.VersionedLCUpdate],
	error,
) {
	response := &api.Response[[]*spec.VersionedLCUpdate]{
		Metadata: map[string]any{
			"version": string(res.consensusVersion),
		} ,
	}

	//updates, _, err := decodeJSONResponse(bytes.NewReader(res.body), &lcUpdatesJSON{})
	//if err != nil {
	//	return nil, err
	//}

	var err error
	decoded := make(map[string]json.RawMessage)
	if err = json.NewDecoder(bytes.NewReader(res.body)).Decode(&decoded); err != nil {
		return nil, errors.Wrap(err, "failed to parse JSON")
	}

	updates := make([]json.RawMessage, 0)
	for k, v := range decoded {
		if k == "data" {
			if err = json.NewDecoder(bytes.NewReader(v)).Decode(&updates); err != nil {
				return nil, errors.Wrap(err, "failed to parse JSON")
			}
		}
	}

	datas := make([]*spec.VersionedLCUpdate, 0)
	for _, u := range updates {
		d := &spec.VersionedLCUpdate{Version: res.consensusVersion}
		switch d.Version {
		case spec.DataVersionPhase0:
			err = fmt.Errorf("unsupported version %s", res.consensusVersion)
		case spec.DataVersionAltair:
			update := altair.LightClientUpdate{}
			err = json.Unmarshal(u, &update)
			d.Altair = &update
			//d.Altair, response.Metadata, err = decodeJSONResponse(bytes.NewReader(res.body), &altair.LightClientUpdate{})
		case spec.DataVersionBellatrix:
			update := altair.LightClientUpdate{}
			err = json.Unmarshal(u, &update)
			d.Bellatrix = &update
			//d.Bellatrix, response.Metadata, err = decodeJSONResponse(bytes.NewReader(res.body), &altair.LightClientUpdate{})
		case spec.DataVersionCapella:
			update := capella.LightClientUpdate{}
			err = json.Unmarshal(u, &update)
			d.Capella = &update
			//d.Capella, response.Metadata, err = decodeJSONResponse(bytes.NewReader(res.body), &capella.LightClientUpdate{})
		case spec.DataVersionDeneb:
			update := deneb.LightClientUpdate{}
			err = json.Unmarshal(u, &update)
			d.Deneb = &update
			//d.Deneb, response.Metadata, err = decodeJSONResponse(bytes.NewReader(res.body), &deneb.LightClientUpdate{})
		default:
			err = fmt.Errorf("unsupported version %s", res.consensusVersion)
		}
		if err != nil {
			return nil, err
		}
		datas = append(datas, d)
	}
	response.Data = datas

	return response, nil
}

// LightClientFinalityUpdate provides the latest light client finality_update
func (s *Service) LightClientFinalityUpdate(ctx context.Context, opts *api.CommonOpts) (
	*api.Response[*spec.VersionedLCFinalityUpdate],
	error,
) {
	resp, err := s.get(ctx, "/eth/v1/beacon/light_client/finality_update", opts)
	if err != nil {
		return nil, errors.Wrap(err, "failed to request beacon light client finality_update")
	}
	if resp == nil {
		return nil, nil
	}

	switch resp.contentType {
	case ContentTypeSSZ:
		return s.lightClientFinalityUpdateFromSSZ(resp)
	case ContentTypeJSON:
		return s.lightClientFinalityUpdateFromJSON(resp)
	default:
		return nil, fmt.Errorf("unhandled content type %v", resp.contentType)
	}
}

func (s *Service) lightClientFinalityUpdateFromSSZ(res *httpResponse) (
	*api.Response[*spec.VersionedLCFinalityUpdate],
	error,
) {
	response := &api.Response[*spec.VersionedLCFinalityUpdate]{
		Data: &spec.VersionedLCFinalityUpdate{
			Version: res.consensusVersion,
		},
		Metadata: metadataFromHeaders(res.headers),
	}

	var err error
	switch res.consensusVersion {
	case spec.DataVersionPhase0:
		err = fmt.Errorf("unsupported version %s", res.consensusVersion)
	case spec.DataVersionAltair:
		response.Data.Altair = &altair.LightClientFinalityUpdate{}
		err = response.Data.Altair.UnmarshalSSZ(res.body)
	case spec.DataVersionBellatrix:
		response.Data.Bellatrix = &altair.LightClientFinalityUpdate{}
		err = response.Data.Bellatrix.UnmarshalSSZ(res.body)
	case spec.DataVersionCapella:
		response.Data.Capella = &capella.LightClientFinalityUpdate{}
		err = response.Data.Capella.UnmarshalSSZ(res.body)
	case spec.DataVersionDeneb:
		response.Data.Deneb = &deneb.LightClientFinalityUpdate{}
		err = response.Data.Deneb.UnmarshalSSZ(res.body)
	default:
		err = fmt.Errorf("unsupported version %s", res.consensusVersion)
	}
	if err != nil {
		return nil, errors.Wrapf(err, "failed to decode %s light client optimistic update", res.consensusVersion)
	}

	return response, nil
}

func (s *Service) lightClientFinalityUpdateFromJSON(res *httpResponse) (
	*api.Response[*spec.VersionedLCFinalityUpdate],
	error,
) {
	response := &api.Response[*spec.VersionedLCFinalityUpdate]{}

	var err error
	response.Data, response.Metadata, err = versionedLCFinalityUpdateFromJSON(res.consensusVersion, res.body)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func versionedLCFinalityUpdateFromJSON(version spec.DataVersion, data []byte) (
	*spec.VersionedLCFinalityUpdate,
	map[string]any,
	error,
) {
	var err error
	var metadata map[string]any
	update := &spec.VersionedLCFinalityUpdate{
		Version: version,
	}
	switch version {
	case spec.DataVersionPhase0:
		err = fmt.Errorf("unsupported version %s", version)
	case spec.DataVersionAltair:
		update.Altair, metadata, err = decodeJSONResponse(bytes.NewReader(data), &altair.LightClientFinalityUpdate{})
	case spec.DataVersionBellatrix:
		update.Bellatrix, metadata, err = decodeJSONResponse(bytes.NewReader(data), &altair.LightClientFinalityUpdate{})
	case spec.DataVersionCapella:
		update.Capella, metadata, err = decodeJSONResponse(bytes.NewReader(data), &capella.LightClientFinalityUpdate{})
	case spec.DataVersionDeneb:
		update.Deneb, metadata, err = decodeJSONResponse(bytes.NewReader(data), &deneb.LightClientFinalityUpdate{})
	default:
		err = fmt.Errorf("unsupported version %s", version)
	}
	if err != nil {
		return nil, nil, err
	}
	return update, metadata, nil
}

// LightClientOptimisticUpdate provides the latest light client optimistic_update
func (s *Service) LightClientOptimisticUpdate(ctx context.Context, opts *api.CommonOpts) (
	*api.Response[*spec.VersionedLCOptimisticUpdate],
	error,
) {
	resp, err := s.get(ctx, "/eth/v1/beacon/light_client/optimistic_update", opts)
	if err != nil {
		return nil, errors.Wrap(err, "failed to request beacon light client optimistic_update")
	}
	if resp == nil {
		return nil, nil
	}

	switch resp.contentType {
	case ContentTypeSSZ:
		return s.lightClientOptimisticUpdateFromSSZ(resp)
	case ContentTypeJSON:
		return s.lightClientOptimisticUpdateFromJSON(resp)
	default:
		return nil, fmt.Errorf("unhandled content type %v", resp.contentType)
	}
}

func (s *Service) lightClientOptimisticUpdateFromSSZ(res *httpResponse) (
	*api.Response[*spec.VersionedLCOptimisticUpdate],
	error,
) {
	response := &api.Response[*spec.VersionedLCOptimisticUpdate]{
		Data: &spec.VersionedLCOptimisticUpdate{
			Version: res.consensusVersion,
		},
		Metadata: metadataFromHeaders(res.headers),
	}

	var err error
	switch res.consensusVersion {
	case spec.DataVersionPhase0:
		err = fmt.Errorf("unsupported version %s", res.consensusVersion)
	case spec.DataVersionAltair:
		response.Data.Altair = &altair.LightClientOptimisticUpdate{}
		err = response.Data.Altair.UnmarshalSSZ(res.body)
	case spec.DataVersionBellatrix:
		response.Data.Bellatrix = &altair.LightClientOptimisticUpdate{}
		err = response.Data.Bellatrix.UnmarshalSSZ(res.body)
	case spec.DataVersionCapella:
		response.Data.Capella = &capella.LightClientOptimisticUpdate{}
		err = response.Data.Capella.UnmarshalSSZ(res.body)
	case spec.DataVersionDeneb:
		response.Data.Deneb = &deneb.LightClientOptimisticUpdate{}
		err = response.Data.Deneb.UnmarshalSSZ(res.body)
	default:
		err = fmt.Errorf("unsupported version %s", res.consensusVersion)
	}
	if err != nil {
		return nil, errors.Wrapf(err, "failed to decode %s light client optimistic update", res.consensusVersion)
	}

	return response, nil
}

func (s *Service) lightClientOptimisticUpdateFromJSON(res *httpResponse) (
	*api.Response[*spec.VersionedLCOptimisticUpdate],
	error,
) {
	response := &api.Response[*spec.VersionedLCOptimisticUpdate]{}

	var err error
	response.Data, response.Metadata, err = versionedLCOptimisticUpdateFromJSON(res.consensusVersion, res.body)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func versionedLCOptimisticUpdateFromJSON(version spec.DataVersion, data []byte) (
	*spec.VersionedLCOptimisticUpdate,
	map[string]any,
	error,
) {
	var err error
	var metadata map[string]any
	update := &spec.VersionedLCOptimisticUpdate{
		Version: version,
	}
	switch version {
	case spec.DataVersionPhase0:
		err = fmt.Errorf("unsupported version %s", version)
	case spec.DataVersionAltair:
		update.Altair, metadata, err = decodeJSONResponse(bytes.NewReader(data), &altair.LightClientOptimisticUpdate{})
	case spec.DataVersionBellatrix:
		update.Bellatrix, metadata, err = decodeJSONResponse(bytes.NewReader(data), &altair.LightClientOptimisticUpdate{})
	case spec.DataVersionCapella:
		update.Capella, metadata, err = decodeJSONResponse(bytes.NewReader(data), &capella.LightClientOptimisticUpdate{})
	case spec.DataVersionDeneb:
		update.Deneb, metadata, err = decodeJSONResponse(bytes.NewReader(data), &deneb.LightClientOptimisticUpdate{})
	default:
		err = fmt.Errorf("unsupported version %s", version)
	}
	if err != nil {
		return nil, nil, err
	}
	return update, metadata, nil
}
