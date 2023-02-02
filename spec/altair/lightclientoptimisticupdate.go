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

package altair

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/pkg/errors"

	"github.com/attestantio/go-eth2-client/spec/phase0"
)

// LightClientOptimisticUpdate is the data providing light client finality update
type LightClientOptimisticUpdate struct {
	AttestedHeader *LightClientHeader // The beacon block header that is attested to by the sync committee
	SyncAggregate  *SyncAggregate     // Sync committee aggregate signature
	SignatureSlot  phase0.Slot        // Slot at which the aggregate signature was created (untrusted)
}

// lightClientOptimisticUpdateJSON is the spec representation of the struct.
type lightClientOptimisticUpdateJSON struct {
	AttestedHeader *LightClientHeader `json:"attested_header"`
	SyncAggregate  *SyncAggregate     `json:"sync_aggregate"`
	SignatureSlot  string             `json:"signature_slot"`
}

// MarshalJSON implements json.Marshaler.
func (l *LightClientOptimisticUpdate) MarshalJSON() ([]byte, error) {
	return json.Marshal(&lightClientOptimisticUpdateJSON{
		AttestedHeader: l.AttestedHeader,
		SyncAggregate:  l.SyncAggregate,
		SignatureSlot:  fmt.Sprintf("%d", l.SignatureSlot),
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (l *LightClientOptimisticUpdate) UnmarshalJSON(input []byte) error {
	var err error

	var jso lightClientOptimisticUpdateJSON
	if err = json.Unmarshal(input, &jso); err != nil {
		return errors.Wrap(err, "invalid JSON")
	}

	if jso.AttestedHeader == nil {
		return errors.New("attested_header missing")
	}
	l.AttestedHeader = jso.AttestedHeader

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
func (l *LightClientOptimisticUpdate) String() string {
	data, err := json.Marshal(l)
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}
	return string(data)
}
