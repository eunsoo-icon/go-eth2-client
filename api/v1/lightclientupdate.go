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

package v1

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/pkg/errors"

	"github.com/attestantio/go-eth2-client/spec/altair"
	"github.com/attestantio/go-eth2-client/spec/phase0"
)

// LightClientUpdate is the data providing light client update
type LightClientUpdate struct {
	AttestedHeader          *phase0.BeaconBlockHeader // The beacon block header that is attested to by the sync committee
	NextSyncCommittee       *altair.SyncCommittee     // Next sync committee corresponding to `attested_header`
	NextSyncCommitteeBranch [][]byte                  `ssz-size:"5,32"`
	FinalizedHeader         *phase0.BeaconBlockHeader // The finalized beacon block header attested to by Merkle branch
	FinalityBranch          [][]byte                  `ssz-size:"6,32"`
	SyncAggregate           *altair.SyncAggregate     // Sync committee aggregate signature
	SignatureSlot           phase0.Slot               // Slot at which the aggregate signature was created (untrusted)
}

// lightClientUpdateJSON is the spec representation of the struct.
type lightClientUpdateJSON struct {
	AttestedHeader          *phase0.BeaconBlockHeader `json:"attested_header"`
	NextSyncCommittee       *altair.SyncCommittee     `json:"next_sync_committee"`
	NextSyncCommitteeBranch []string                  `json:"next_sync_committee_branch"`
	FinalizedHeader         *phase0.BeaconBlockHeader `json:"finalized_header"`
	FinalityBranch          []string                  `json:"finality_branch"`
	SyncAggregate           *altair.SyncAggregate     `json:"sync_aggregate"`
	SignatureSlot           string                    `json:"signature_slot"`
}

// MarshalJSON implements json.Marshaler.
func (l *LightClientUpdate) MarshalJSON() ([]byte, error) {
	nb := make([]string, len(l.NextSyncCommitteeBranch))
	for i := range l.NextSyncCommitteeBranch {
		nb[i] = fmt.Sprintf("%#x", l.NextSyncCommitteeBranch[i])
	}
	fb := make([]string, len(l.FinalityBranch))
	for i := range l.FinalityBranch {
		fb[i] = fmt.Sprintf("%#x", l.FinalityBranch[i])
	}

	return json.Marshal(&lightClientUpdateJSON{
		AttestedHeader:          l.AttestedHeader,
		NextSyncCommittee:       l.NextSyncCommittee,
		NextSyncCommitteeBranch: nb,
		FinalizedHeader:         l.FinalizedHeader,
		FinalityBranch:          fb,
		SyncAggregate:           l.SyncAggregate,
		SignatureSlot:           fmt.Sprintf("%d", l.SignatureSlot),
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (l *LightClientUpdate) UnmarshalJSON(input []byte) error {
	var err error

	var jso lightClientUpdateJSON
	if err = json.Unmarshal(input, &jso); err != nil {
		return errors.Wrap(err, "invalid JSON")
	}

	if jso.AttestedHeader == nil {
		return errors.New("attested_header missing")
	}
	l.AttestedHeader = jso.AttestedHeader

	if jso.NextSyncCommittee == nil {
		return errors.New("next_sync_committee missing")
	}
	l.NextSyncCommittee = jso.NextSyncCommittee

	if jso.NextSyncCommitteeBranch == nil {
		return errors.New("next_sync_committee_branch missing")
	}
	if len(jso.NextSyncCommitteeBranch) == 0 {
		return errors.New("next_sync_committee_branch length cannot be 0")
	}
	l.NextSyncCommitteeBranch = make([][]byte, len(jso.NextSyncCommitteeBranch))
	for i := range jso.NextSyncCommitteeBranch {
		if jso.NextSyncCommitteeBranch[i] == "" {
			return errors.Errorf("next_sync_committee_branch[%d] missing", i)
		}
		l.NextSyncCommitteeBranch[i], err = hex.DecodeString(strings.TrimPrefix(jso.NextSyncCommitteeBranch[i], "0x"))
		if err != nil {
			return errors.Wrapf(err, "invalid value for next_sync_committee_branch[%d]", i)
		}
		if len(l.NextSyncCommitteeBranch[i]) != 32 {
			return errors.Errorf("invalid length of next_sync_committee_branch[%d]", i)
		}
	}

	if jso.FinalizedHeader == nil {
		return errors.New("finalized_header missing")
	}
	l.FinalizedHeader = jso.FinalizedHeader

	if jso.FinalityBranch == nil {
		return errors.New("finality_branch missing")
	}
	if len(jso.FinalityBranch) == 0 {
		return errors.New("finality_branch length cannot be 0")
	}
	l.FinalityBranch = make([][]byte, len(jso.FinalityBranch))
	for i := range jso.FinalityBranch {
		if jso.FinalityBranch[i] == "" {
			return errors.Errorf("finality_branch[%d] missing", i)
		}
		l.FinalityBranch[i], err = hex.DecodeString(strings.TrimPrefix(jso.FinalityBranch[i], "0x"))
		if err != nil {
			return errors.Wrapf(err, "invalid value for finality_branch[%d]", i)
		}
		if len(l.FinalityBranch[i]) != 32 {
			return errors.Errorf("invalid length of finality_branch[%d]", i)
		}
	}

	if jso.SyncAggregate == nil {
		return errors.New("sync_aggregate missing")
	}
	l.SyncAggregate = jso.SyncAggregate

	if jso.SignatureSlot == "" {
		return errors.New("signature_slot missing")
	}
	slot, err := strconv.ParseUint(jso.SignatureSlot, 10, 64)
	if err != nil {
		return errors.Wrap(err, "invalid value for slot")
	}
	l.SignatureSlot = phase0.Slot(slot)
	return nil
}

// String returns a string version of the structure.
func (l *LightClientUpdate) String() string {
	data, err := json.Marshal(l)
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}
	return string(data)
}
