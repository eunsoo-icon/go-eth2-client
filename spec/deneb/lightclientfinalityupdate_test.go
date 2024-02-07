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
	"github.com/stretchr/testify/require"

	"github.com/attestantio/go-eth2-client/spec/deneb"
)

const (
	lcFinalityUpdateFormat = `{"attested_header":%s,"finalized_header":%s,"finality_branch":%s,"sync_aggregate":%s,"signature_slot":%s}`
	syncAggregate          = `{"sync_committee_bits":"0xfffffffffffbffffffffffffffffffffffffffffffffffff7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff","sync_committee_signature":"0x97e006ecbe9df2f082eb450e1c07ace045da0d4e367f453170bfb32911e72fc9f08237d348e99b3500531c8cba770fc119d844c22950c094d860cfa784ba237debe681e55875994a75c72689d9e289f72c8bb7559ae91b3788e5e769aee0705a"}`
)

func TestLightClientFinalityUpdateJSON(t *testing.T) {
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
			err:   "invalid JSON: json: cannot unmarshal array into Go value of type deneb.lightClientFinalityUpdateJSON",
		},
		{
			name:  "AttestedHeaderMissing",
			input: []byte(fmt.Sprintf(`{"finalized_header":%s,"finality_branch":%s,"sync_aggregate":%s,"signature_slot":%s}`, lcHeaderValue, proofBranch, syncAggregate, slot)),
			err:   "attested_header missing",
		},
		{
			name:  "FinalizedHeaderMissing",
			input: []byte(fmt.Sprintf(`{"attested_header":%s,"finality_branch":%s,"sync_aggregate":%s,"signature_slot":%s}`, lcHeaderValue, proofBranch, syncAggregate, slot)),
			err:   "finalized_header missing",
		},
		{
			name:  "FinalityBranchMissing",
			input: []byte(fmt.Sprintf(`{"attested_header":%s,"finalized_header":%s,"sync_aggregate":%s,"signature_slot":%s}`, lcHeaderValue, lcHeaderValue, syncAggregate, slot)),
			err:   "finality_branch missing",
		},
		{
			name:  "FinalityBranchEmpty",
			input: []byte(fmt.Sprintf(lcFinalityUpdateFormat, lcHeaderValue, lcHeaderValue, "[]", syncAggregate, slot)),
			err:   "finality_branch length cannot be 0",
		},
		{
			name:  "FinalityBranchWrongType",
			input: []byte(fmt.Sprintf(lcFinalityUpdateFormat, lcHeaderValue, lcHeaderValue, "true", syncAggregate, slot)),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field lightClientFinalityUpdateJSON.finality_branch of type []string",
		},
		{
			name:  "FinalityBranchWrongValueType",
			input: []byte(fmt.Sprintf(lcFinalityUpdateFormat, lcHeaderValue, lcHeaderValue, "[123]", syncAggregate, slot)),
			err:   "invalid JSON: json: cannot unmarshal number into Go struct field lightClientFinalityUpdateJSON.finality_branch of type string",
		},
		{
			name:  "FinalityBranchWrongValueFormat",
			input: []byte(fmt.Sprintf(lcFinalityUpdateFormat, lcHeaderValue, lcHeaderValue, `["invalid"]`, syncAggregate, slot)),
			err:   "invalid value for finality_branch[0]: encoding/hex: invalid byte: U+0069 'i'",
		},
		{
			name:  "FinalityBranchWrongValueFormat",
			input: []byte(fmt.Sprintf(lcFinalityUpdateFormat, lcHeaderValue, lcHeaderValue, `["0x12acde"]`, syncAggregate, slot)),
			err:   "invalid length of finality_branch[0]",
		},
		{
			name:  "SyncAggregateMissing",
			input: []byte(fmt.Sprintf(`{"attested_header":%s,"finalized_header":%s,"finality_branch":%s,"signature_slot":%s}`, lcHeaderValue, lcHeaderValue, proofBranch, slot)),
			err:   "sync_aggregate missing",
		},
		{
			name:  "SignatureSlotMissing",
			input: []byte(fmt.Sprintf(`{"attested_header":%s,"finalized_header":%s,"finality_branch":%s,"sync_aggregate":%s}`, lcHeaderValue, lcHeaderValue, proofBranch, syncAggregate)),
			err:   "signature_slot missing",
		},
		{
			name:  "SignatureSlotWrongType",
			input: []byte(fmt.Sprintf(lcFinalityUpdateFormat, lcHeaderValue, lcHeaderValue, proofBranch, syncAggregate, "true")),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field lightClientFinalityUpdateJSON.signature_slot of type string",
		},
		{
			name:  "SignatureSlotInvalid",
			input: []byte(fmt.Sprintf(lcFinalityUpdateFormat, lcHeaderValue, lcHeaderValue, proofBranch, syncAggregate, `"-1"`)),
			err:   "invalid value for slot: strconv.ParseUint: parsing \"-1\": invalid syntax",
		},
		{
			name:  "Good",
			input: []byte(fmt.Sprintf(lcFinalityUpdateFormat, lcHeaderValue, lcHeaderValue, proofBranch, syncAggregate, slot)),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var res deneb.LightClientFinalityUpdate
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
