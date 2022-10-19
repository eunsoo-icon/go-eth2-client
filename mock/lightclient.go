// Copyright © 2020 Attestant Limited.
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

package mock

import (
	"context"

	api "github.com/attestantio/go-eth2-client/api/v1"
	"github.com/attestantio/go-eth2-client/spec/altair"
	"github.com/attestantio/go-eth2-client/spec/phase0"
)

// LightClientBootstrap provides the light client bootstrap of a given block ID.
func (s *Service) LightClientBootstrap(ctx context.Context, blockID string) (*api.LightClientBootstrap, error) {
	return &api.LightClientBootstrap{
		Header:                     &phase0.BeaconBlockHeader{},
		CurrentSyncCommittee:       &altair.SyncCommittee{},
		CurrentSyncCommitteeBranch: [][32]byte{},
	}, nil
}

// LightClientFinalityUpdate provides the light client finality_update
func (s *Service) LightClientFinalityUpdate(ctx context.Context) (*api.LightClientFinalityUpdate, error) {
	return &api.LightClientFinalityUpdate{
		AttestedHeader:  &phase0.BeaconBlockHeader{},
		FinalizedHeader: &phase0.BeaconBlockHeader{},
		FinalityBranch:  [][32]byte{},
		SyncAggregate:   &altair.SyncAggregate{},
	}, nil
}

// LightClientOptimisticUpdate provides the light client optimistic_update
func (s *Service) LightClientOptimisticUpdate(ctx context.Context) (*api.LightClientOptimisticUpdate, error) {
	return &api.LightClientOptimisticUpdate{
		AttestedHeader: &phase0.BeaconBlockHeader{},
		SyncAggregate:  &altair.SyncAggregate{},
	}, nil
}
