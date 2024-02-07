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

package altair_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/attestantio/go-eth2-client/spec/altair"
)

func TestLightClientOptimisticUpdateJSON(t *testing.T) {
	const (
		format        = `{"attested_header":%s,"sync_aggregate":%s,"signature_slot":%s}`
		header        = `{"beacon":{"slot":"4943744","proposer_index":"222870","parent_root":"0x0c594acb2c7ec3564590fd2feb6724cfcf786faf51fe2a284154516c2903c153","state_root":"0x237962d02698b2f5f37f3a7c43dfae0e2fe28e103225237bc7f09938c8527eaa","body_root":"0xff42d5726526628ce27c4ca89172ccf5c562adbfec64c22d494b6f8bd03dc034"}}`
		syncAggregate = `{"sync_committee_bits":"0xfffffffffffbffffffffffffffffffffffffffffffffffff7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff","sync_committee_signature":"0x97e006ecbe9df2f082eb450e1c07ace045da0d4e367f453170bfb32911e72fc9f08237d348e99b3500531c8cba770fc119d844c22950c094d860cfa784ba237debe681e55875994a75c72689d9e289f72c8bb7559ae91b3788e5e769aee0705a"}`
		slot          = `"1234"`
	)

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
			err:   "invalid JSON: json: cannot unmarshal array into Go value of type altair.lightClientOptimisticUpdateJSON",
		},
		{
			name:  "AttestedHeaderMissing",
			input: []byte(fmt.Sprintf(`{"sync_aggregate":%s,"signature_slot":%s}`, syncAggregate, slot)),
			err:   "attested_header missing",
		},
		{
			name:  "SyncAggregateMissing",
			input: []byte(fmt.Sprintf(`{"attested_header":%s,"signature_slot":%s}`, header, slot)),
			err:   "sync_aggregate missing",
		},
		{
			name:  "SignatureSlotMissing",
			input: []byte(fmt.Sprintf(`{"attested_header":%s,"sync_aggregate":%s}`, header, syncAggregate)),
			err:   "signature_slot missing",
		},
		{
			name:  "SignatureSlotWrongType",
			input: []byte(fmt.Sprintf(format, header, syncAggregate, "true")),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field lightClientOptimisticUpdateJSON.signature_slot of type string",
		},
		{
			name:  "SignatureSlotInvalid",
			input: []byte(fmt.Sprintf(format, header, syncAggregate, `"-1"`)),
			err:   "invalid value for slot: strconv.ParseUint: parsing \"-1\": invalid syntax",
		},
		{
			name:  "Good",
			input: []byte(fmt.Sprintf(format, header, syncAggregate, slot)),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var res altair.LightClientOptimisticUpdate
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
