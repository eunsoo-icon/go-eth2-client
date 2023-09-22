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

package capella

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/pkg/errors"

	"github.com/attestantio/go-eth2-client/spec/altair"
	"github.com/attestantio/go-eth2-client/spec/phase0"
)

// LightClientHeader is the data providing light client header
type LightClientHeader struct {
	Beacon          *phase0.BeaconBlockHeader
	Execution       *ExecutionPayloadHeader
	ExecutionBranch [][]byte `ssz-size:"4,32"`
}

// lightClientHeaderJSON is the spec representation of the struct.
type lightClientHeaderJSON struct {
	Beacon          *phase0.BeaconBlockHeader `json:"beacon"`
	Execution       *ExecutionPayloadHeader   `json:"execution"`
	ExecutionBranch []string                  `json:"execution_branch"`
}

// MarshalJSON implements json.Marshaler.
func (l *LightClientHeader) MarshalJSON() ([]byte, error) {
	branch := make([]string, len(l.ExecutionBranch))
	for i := range l.ExecutionBranch {
		branch[i] = fmt.Sprintf("%#x", l.ExecutionBranch[i])
	}

	return json.Marshal(&lightClientHeaderJSON{
		Beacon:          l.Beacon,
		Execution:       l.Execution,
		ExecutionBranch: branch,
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (l *LightClientHeader) UnmarshalJSON(input []byte) error {
	var err error

	var jso lightClientHeaderJSON
	if err = json.Unmarshal(input, &jso); err != nil {
		return errors.Wrap(err, "invalid JSON")
	}

	if jso.Beacon == nil {
		return errors.New("beacon missing")
	}
	l.Beacon = jso.Beacon

	if jso.Execution == nil {
		return errors.New("execution missing")
	}
	l.Execution = jso.Execution

	if jso.ExecutionBranch == nil {
		return errors.New("execution_branch missing")
	}
	if len(jso.ExecutionBranch) == 0 {
		return errors.New("execution_branch length cannot be 0")
	}
	l.ExecutionBranch = make([][]byte, len(jso.ExecutionBranch))
	for i := range jso.ExecutionBranch {
		if jso.ExecutionBranch[i] == "" {
			return errors.Errorf("execution_branch[%d] missing", i)
		}
		l.ExecutionBranch[i], err = hex.DecodeString(strings.TrimPrefix(jso.ExecutionBranch[i], "0x"))
		if err != nil {
			return errors.Wrapf(err, "invalid value for execution_branch[%d]", i)
		}
		if len(l.ExecutionBranch[i]) != 32 {
			return errors.Errorf("invalid length of execution_branch[%d]", i)
		}
	}

	return nil
}

// String returns a string version of the structure.
func (l *LightClientHeader) String() string {
	data, err := json.Marshal(l)
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}
	return string(data)
}

func (l *LightClientHeader) ToAltair() *altair.LightClientHeader {
	return &altair.LightClientHeader{Beacon: l.Beacon}
}
