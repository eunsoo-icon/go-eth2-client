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

	"github.com/pkg/errors"

	"github.com/attestantio/go-eth2-client/spec/phase0"
)

// LightClientHeader is the data providing light client header
type LightClientHeader struct {
	Beacon *phase0.BeaconBlockHeader
}

// lightClientHeaderJSON is the spec representation of the struct.
type lightClientHeaderJSON struct {
	Beacon *phase0.BeaconBlockHeader `json:"beacon"`
}

// MarshalJSON implements json.Marshaler.
func (l *LightClientHeader) MarshalJSON() ([]byte, error) {
	return json.Marshal(&lightClientHeaderJSON{
		Beacon: l.Beacon,
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
