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
	"strings"

	"github.com/pkg/errors"

	"github.com/attestantio/go-eth2-client/spec/altair"
	"github.com/attestantio/go-eth2-client/spec/phase0"
)

// LightClientBootstrap is the data providing light client bootstrap
type LightClientBootstrap struct {
	Header                     *phase0.BeaconBlockHeader
	CurrentSyncCommittee       *altair.SyncCommittee
	CurrentSyncCommitteeBranch [][32]byte
}

// lightClientBootstrapJSON is the spec representation of the struct.
type lightClientBootstrapJSON struct {
	Header                     *phase0.BeaconBlockHeader `json:"header"`
	CurrentSyncCommittee       *altair.SyncCommittee     `json:"current_sync_committee"`
	CurrentSyncCommitteeBranch []string                  `json:"current_sync_committee_branch"`
}

// MarshalJSON implements json.Marshaler.
func (l *LightClientBootstrap) MarshalJSON() ([]byte, error) {
	branch := make([]string, len(l.CurrentSyncCommitteeBranch))
	for i := range l.CurrentSyncCommitteeBranch {
		branch[i] = fmt.Sprintf("%#x", l.CurrentSyncCommitteeBranch[i])
	}

	return json.Marshal(&lightClientBootstrapJSON{
		Header:                     l.Header,
		CurrentSyncCommittee:       l.CurrentSyncCommittee,
		CurrentSyncCommitteeBranch: branch,
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (l *LightClientBootstrap) UnmarshalJSON(input []byte) error {
	var err error

	var jso lightClientBootstrapJSON
	if err = json.Unmarshal(input, &jso); err != nil {
		return errors.Wrap(err, "invalid JSON")
	}

	if jso.Header == nil {
		return errors.New("header missing")
	}
	l.Header = jso.Header

	if jso.CurrentSyncCommittee == nil {
		return errors.New("current_sync_committee missing")
	}
	l.CurrentSyncCommittee = jso.CurrentSyncCommittee

	if jso.CurrentSyncCommitteeBranch == nil {
		return errors.New("current_sync_committee_branch missing")
	}
	if len(jso.CurrentSyncCommitteeBranch) == 0 {
		return errors.New("current_sync_committee_branch length cannot be 0")
	}
	l.CurrentSyncCommitteeBranch = make([][32]byte, len(jso.CurrentSyncCommitteeBranch))
	for i := range jso.CurrentSyncCommitteeBranch {
		branch, err := hex.DecodeString(strings.TrimPrefix(jso.CurrentSyncCommitteeBranch[i], "0x"))
		if err != nil {
			return errors.Wrapf(err, "invalid value for current_sync_committee_branch[%d]", i)
		}
		if len(branch) != 32 {
			return errors.Errorf("invalid length of current_sync_committee_branch[%d]", i)
		}
		copy(l.CurrentSyncCommitteeBranch[i][:], branch)
	}

	return nil
}

// String returns a string version of the structure.
func (l *LightClientBootstrap) String() string {
	data, err := json.Marshal(l)
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}
	return string(data)
}
