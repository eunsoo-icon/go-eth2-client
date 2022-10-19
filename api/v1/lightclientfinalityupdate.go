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

// LightClientFinalityUpdate is the data providing light client finality update
type LightClientFinalityUpdate struct {
	AttestedHeader  *phase0.BeaconBlockHeader // The beacon block header that is attested to by the sync committee
	FinalizedHeader *phase0.BeaconBlockHeader // The finalized beacon block header attested to by Merkle branch
	FinalityBranch  [][32]byte
	SyncAggregate   *altair.SyncAggregate // Sync committee aggregate signature
	SignatureSlot   phase0.Slot           // Slot at which the aggregate signature was created (untrusted)
}

// lightClientFinalityUpdateJSON is the spec representation of the struct.
type lightClientFinalityUpdateJSON struct {
	AttestedHeader  *phase0.BeaconBlockHeader `json:"attested_header"`
	FinalizedHeader *phase0.BeaconBlockHeader `json:"finalized_header"`
	FinalityBranch  []string                  `json:"finality_branch"`
	SyncAggregate   *altair.SyncAggregate     `json:"sync_aggregate"`
	SignatureSlot   string                    `json:"signature_slot"`
}

// MarshalJSON implements json.Marshaler.
func (l *LightClientFinalityUpdate) MarshalJSON() ([]byte, error) {
	branch := make([]string, len(l.FinalityBranch))
	for i := range l.FinalityBranch {
		branch[i] = fmt.Sprintf("%#x", l.FinalityBranch[i])
	}

	return json.Marshal(&lightClientFinalityUpdateJSON{
		AttestedHeader:  l.AttestedHeader,
		FinalizedHeader: l.FinalizedHeader,
		FinalityBranch:  branch,
		SyncAggregate:   l.SyncAggregate,
		SignatureSlot:   fmt.Sprintf("%d", l.SignatureSlot),
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (l *LightClientFinalityUpdate) UnmarshalJSON(input []byte) error {
	var err error

	var jso lightClientFinalityUpdateJSON
	if err = json.Unmarshal(input, &jso); err != nil {
		return errors.Wrap(err, "invalid JSON")
	}

	if jso.AttestedHeader == nil {
		return errors.New("attested_header missing")
	}
	l.AttestedHeader = jso.AttestedHeader

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
	l.FinalityBranch = make([][32]byte, len(jso.FinalityBranch))
	for i := range jso.FinalityBranch {
		branch, err := hex.DecodeString(strings.TrimPrefix(jso.FinalityBranch[i], "0x"))
		if err != nil {
			return errors.Wrapf(err, "invalid value for finality_branch[%d]", i)
		}
		if len(branch) != 32 {
			return errors.Errorf("invalid length of finality_branch[%d]", i)
		}
		copy(l.FinalityBranch[i][:], branch)
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
func (l *LightClientFinalityUpdate) String() string {
	data, err := json.Marshal(l)
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}
	return string(data)
}
