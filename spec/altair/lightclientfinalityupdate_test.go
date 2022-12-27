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

	"github.com/stretchr/testify/require"
	"gotest.tools/assert"

	"github.com/attestantio/go-eth2-client/spec/altair"
)

func TestLightClientFinalityUpdateJSON(t *testing.T) {
	const (
		format        = `{"attested_header":%s,"finalized_header":%s,"finality_branch":%s,"sync_aggregate":%s,"signature_slot":%s}`
		header        = `{"slot":"4943744","proposer_index":"222870","parent_root":"0x0c594acb2c7ec3564590fd2feb6724cfcf786faf51fe2a284154516c2903c153","state_root":"0x237962d02698b2f5f37f3a7c43dfae0e2fe28e103225237bc7f09938c8527eaa","body_root":"0xff42d5726526628ce27c4ca89172ccf5c562adbfec64c22d494b6f8bd03dc034"}`
		branch        = `["0x65af40980fe7dfc2f8587fd1d75044f8adcf8e0e8b142363f5bf3bce21e66bb5","0x26648104944ae0085548cea356ebdd0c5c4b73aa440bcaf0c7b2821325b28f66","0xf97cbc51dd5b8ffffb73783e6938e3eee934448eaa08c9f50e136cb00635cf9f","0xf2adfbbfc2a4e45f01f90752b069b5fcd136b89dfa473dacbea52f6fefc3936c","0x8b32360624a233863dc23f40fa08618cf49651e1b556e21b2f992d12a6cd84c2","0x42ad9040275048b9209b72096dd0d3551e08962a04dfc0602dd1459c4a165367"]`
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
			err:   "invalid JSON: json: cannot unmarshal array into Go value of type altair.lightClientFinalityUpdateJSON",
		},
		{
			name:  "AttestedHeaderMissing",
			input: []byte(fmt.Sprintf(`{"finalized_header":%s,"finality_branch":%s,"sync_aggregate":%s,"signature_slot":%s}`, header, branch, syncAggregate, slot)),
			err:   "attested_header missing",
		},
		{
			name:  "FinalizedHeaderMissing",
			input: []byte(fmt.Sprintf(`{"attested_header":%s,"finality_branch":%s,"sync_aggregate":%s,"signature_slot":%s}`, header, branch, syncAggregate, slot)),
			err:   "finalized_header missing",
		},
		{
			name:  "FinalityBranchMissing",
			input: []byte(fmt.Sprintf(`{"attested_header":%s,"finalized_header":%s,"sync_aggregate":%s,"signature_slot":%s}`, header, header, syncAggregate, slot)),
			err:   "finality_branch missing",
		},
		{
			name:  "FinalityBranchEmpty",
			input: []byte(fmt.Sprintf(format, header, header, "[]", syncAggregate, slot)),
			err:   "finality_branch length cannot be 0",
		},
		{
			name:  "FinalityBranchWrongType",
			input: []byte(fmt.Sprintf(format, header, header, "true", syncAggregate, slot)),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field lightClientFinalityUpdateJSON.finality_branch of type []string",
		},
		{
			name:  "FinalityBranchWrongValueType",
			input: []byte(fmt.Sprintf(format, header, header, "[123]", syncAggregate, slot)),
			err:   "invalid JSON: json: cannot unmarshal number into Go struct field lightClientFinalityUpdateJSON.finality_branch of type string",
		},
		{
			name:  "FinalityBranchWrongValueFormat",
			input: []byte(fmt.Sprintf(format, header, header, `["invalid"]`, syncAggregate, slot)),
			err:   "invalid value for finality_branch[0]: encoding/hex: invalid byte: U+0069 'i'",
		},
		{
			name:  "FinalityBranchWrongValueFormat",
			input: []byte(fmt.Sprintf(format, header, header, `["0x12acde"]`, syncAggregate, slot)),
			err:   "invalid length of finality_branch[0]",
		},
		{
			name:  "SyncAggregateMissing",
			input: []byte(fmt.Sprintf(`{"attested_header":%s,"finalized_header":%s,"finality_branch":%s,"signature_slot":%s}`, header, header, branch, slot)),
			err:   "sync_aggregate missing",
		},
		{
			name:  "SignatureSlotMissing",
			input: []byte(fmt.Sprintf(`{"attested_header":%s,"finalized_header":%s,"finality_branch":%s,"sync_aggregate":%s}`, header, header, branch, syncAggregate)),
			err:   "signature_slot missing",
		},
		{
			name:  "SignatureSlotWrongType",
			input: []byte(fmt.Sprintf(format, header, header, branch, syncAggregate, "true")),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field lightClientFinalityUpdateJSON.signature_slot of type string",
		},
		{
			name:  "SignatureSlotInvalid",
			input: []byte(fmt.Sprintf(format, header, header, branch, syncAggregate, `"-1"`)),
			err:   "invalid value for slot: strconv.ParseUint: parsing \"-1\": invalid syntax",
		},
		{
			name:  "Good",
			input: []byte(fmt.Sprintf(format, header, header, branch, syncAggregate, slot)),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var res altair.LightClientFinalityUpdate
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
