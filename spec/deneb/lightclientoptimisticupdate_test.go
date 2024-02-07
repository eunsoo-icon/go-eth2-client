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

package deneb_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	require "github.com/stretchr/testify/require"

	"github.com/attestantio/go-eth2-client/spec/deneb"
)

const (
	lcOptimisticUpdateFormat = `{"attested_header":%s,"sync_aggregate":%s,"signature_slot":%s}`
)

func TestLightClientOptimisticUpdateJSON(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
		err   string
	}{
		{
			name: "Empty",
			err:  "unexpected end of JSON input",
		},
		{
			name:  "JSONBad",
			input: []byte("[]"),
			err:   "invalid JSON: json: cannot unmarshal array into Go value of type deneb.lightClientOptimisticUpdateJSON",
		},
		{
			name:  "AttestedHeaderMissing",
			input: []byte(fmt.Sprintf(`{"sync_aggregate":%s,"signature_slot":%s}`, syncAggregate, slot)),
			err:   "attested_header missing",
		},
		{
			name:  "SyncAggregateMissing",
			input: []byte(fmt.Sprintf(`{"attested_header":%s,"signature_slot":%s}`, lcHeaderValue, slot)),
			err:   "sync_aggregate missing",
		},
		{
			name:  "SignatureSlotMissing",
			input: []byte(fmt.Sprintf(`{"attested_header":%s,"sync_aggregate":%s}`, lcHeaderValue, syncAggregate)),
			err:   "signature_slot missing",
		},
		{
			name:  "SignatureSlotWrongType",
			input: []byte(fmt.Sprintf(lcOptimisticUpdateFormat, lcHeaderValue, syncAggregate, "true")),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field lightClientOptimisticUpdateJSON.signature_slot of type string",
		},
		{
			name:  "SignatureSlotInvalid",
			input: []byte(fmt.Sprintf(lcOptimisticUpdateFormat, lcHeaderValue, syncAggregate, `"-1"`)),
			err:   "invalid value for slot: strconv.ParseUint: parsing \"-1\": invalid syntax",
		},
		{
			name:  "Good",
			input: []byte(fmt.Sprintf(lcOptimisticUpdateFormat, lcHeaderValue, syncAggregate, slot)),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var res deneb.LightClientOptimisticUpdate
			err := json.Unmarshal(test.input, &res)
			if test.err != "" {
				require.EqualError(t, err, test.err)
			} else {
				require.NoError(t, err)
				rt, err := json.Marshal(&res)
				require.NoError(t, err)
				assert.Equal(t, string(test.input), string(rt))
				assert.Equal(t, string(rt), res.String())
			}
		})
	}
}
