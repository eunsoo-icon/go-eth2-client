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
	lcUpdateFormat = `{"attested_header":%s,"next_sync_committee":%s,"next_sync_committee_branch":%s,"finalized_header":%s,"finality_branch":%s,"sync_aggregate":%s,"signature_slot":%s}`
	slot           = `"1234"`
)

func TestLightClientUpdateJSON(t *testing.T) {
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
			err:   "invalid JSON: json: cannot unmarshal array into Go value of type deneb.lightClientUpdateJSON",
		},
		{
			name:  "AttestedHeaderMissing",
			input: []byte(fmt.Sprintf(`{"next_sync_committee":%s,"next_sync_committee_branch":%s,"finalized_header":%s,"finality_branch":%s,"sync_aggregate":%s,"signature_slot":%s}`, syncCommittee, proofBranch, lcHeaderValue, proofBranch, syncAggregate, slot)),
			err:   "attested_header missing",
		},
		{
			name:  "NextSyncCommitteeMissing",
			input: []byte(fmt.Sprintf(`{"attested_header":%s,"next_sync_committee_branch":%s,"finalized_header":%s,"finality_branch":%s,"sync_aggregate":%s,"signature_slot":%s}`, lcHeaderValue, proofBranch, lcHeaderValue, proofBranch, syncAggregate, slot)),
			err:   "next_sync_committee missing",
		},
		{
			name:  "NextSyncCommitteeBranchMissing",
			input: []byte(fmt.Sprintf(`{"attested_header":%s,"next_sync_committee":%s,"finalized_header":%s,"finality_branch":%s,"sync_aggregate":%s,"signature_slot":%s}`, lcHeaderValue, syncCommittee, lcHeaderValue, proofBranch, syncAggregate, slot)),
			err:   "next_sync_committee_branch missing",
		},
		{
			name:  "NextSyncCommitteeBranchEmpty",
			input: []byte(fmt.Sprintf(lcUpdateFormat, lcHeaderValue, syncCommittee, "[]", lcHeaderValue, proofBranch, syncAggregate, slot)),
			err:   "next_sync_committee_branch length cannot be 0",
		},
		{
			name:  "NextSyncCommitteeBranchWrongType",
			input: []byte(fmt.Sprintf(lcUpdateFormat, lcHeaderValue, syncCommittee, "true", lcHeaderValue, proofBranch, syncAggregate, slot)),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field lightClientUpdateJSON.next_sync_committee_branch of type []string",
		},
		{
			name:  "NextSyncCommitteeBranchWrongValueType",
			input: []byte(fmt.Sprintf(lcUpdateFormat, lcHeaderValue, syncCommittee, "[123]", lcHeaderValue, proofBranch, syncAggregate, slot)),
			err:   "invalid JSON: json: cannot unmarshal number into Go struct field lightClientUpdateJSON.next_sync_committee_branch of type string",
		},
		{
			name:  "NextSyncCommitteeBranchWrongValueFormat",
			input: []byte(fmt.Sprintf(lcUpdateFormat, lcHeaderValue, syncCommittee, `["invalid"]`, lcHeaderValue, proofBranch, syncAggregate, slot)),
			err:   "invalid value for next_sync_committee_branch[0]: encoding/hex: invalid byte: U+0069 'i'",
		},
		{
			name:  "NextSyncCommitteeBranchWrongValueFormat",
			input: []byte(fmt.Sprintf(lcUpdateFormat, lcHeaderValue, syncCommittee, `["0x12acde"]`, lcHeaderValue, proofBranch, syncAggregate, slot)),
			err:   "invalid length of next_sync_committee_branch[0]",
		},
		{
			name:  "FinalizedHeaderMissing",
			input: []byte(fmt.Sprintf(`{"attested_header":%s,"next_sync_committee":%s,"next_sync_committee_branch":%s,"finality_branch":%s,"sync_aggregate":%s,"signature_slot":%s}`, lcHeaderValue, syncCommittee, proofBranch, proofBranch, syncAggregate, slot)),
			err:   "finalized_header missing",
		},
		{
			name:  "FinalityBranchMissing",
			input: []byte(fmt.Sprintf(`{"attested_header":%s,"next_sync_committee":%s,"next_sync_committee_branch":%s,"finalized_header":%s,"sync_aggregate":%s,"signature_slot":%s}`, lcHeaderValue, syncCommittee, proofBranch, lcHeaderValue, syncAggregate, slot)),
			err:   "finality_branch missing",
		},
		{
			name:  "FinalityBranchEmpty",
			input: []byte(fmt.Sprintf(lcUpdateFormat, lcHeaderValue, syncCommittee, proofBranch, lcHeaderValue, "[]", syncAggregate, slot)),
			err:   "finality_branch length cannot be 0",
		},
		{
			name:  "FinalityBranchWrongType",
			input: []byte(fmt.Sprintf(lcUpdateFormat, lcHeaderValue, syncCommittee, proofBranch, lcHeaderValue, "true", syncAggregate, slot)),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field lightClientUpdateJSON.finality_branch of type []string",
		},
		{
			name:  "FinalityBranchWrongValueType",
			input: []byte(fmt.Sprintf(lcUpdateFormat, lcHeaderValue, syncCommittee, proofBranch, lcHeaderValue, "[123]", syncAggregate, slot)),
			err:   "invalid JSON: json: cannot unmarshal number into Go struct field lightClientUpdateJSON.finality_branch of type string",
		},
		{
			name:  "FinalityBranchWrongValueFormat",
			input: []byte(fmt.Sprintf(lcUpdateFormat, lcHeaderValue, syncCommittee, proofBranch, lcHeaderValue, `["invalid"]`, syncAggregate, slot)),
			err:   "invalid value for finality_branch[0]: encoding/hex: invalid byte: U+0069 'i'",
		},
		{
			name:  "FinalityBranchWrongValueFormat",
			input: []byte(fmt.Sprintf(lcUpdateFormat, lcHeaderValue, syncCommittee, proofBranch, lcHeaderValue, `["0x12acde"]`, syncAggregate, slot)),
			err:   "invalid length of finality_branch[0]",
		},
		{
			name:  "SyncAggregateMissing",
			input: []byte(fmt.Sprintf(`{"attested_header":%s,"next_sync_committee":%s,"next_sync_committee_branch":%s,"finalized_header":%s,"finality_branch":%s,"signature_slot":%s}`, lcHeaderValue, syncCommittee, proofBranch, lcHeaderValue, proofBranch, slot)),
			err:   "sync_aggregate missing",
		},
		{
			name:  "SignatureSlotMissing",
			input: []byte(fmt.Sprintf(`{"attested_header":%s,"next_sync_committee":%s,"next_sync_committee_branch":%s,"finalized_header":%s,"finality_branch":%s,"sync_aggregate":%s}`, lcHeaderValue, syncCommittee, proofBranch, lcHeaderValue, proofBranch, syncAggregate)),
			err:   "signature_slot missing",
		},
		{
			name:  "SignatureSlotWrongType",
			input: []byte(fmt.Sprintf(lcUpdateFormat, lcHeaderValue, syncCommittee, proofBranch, lcHeaderValue, proofBranch, syncAggregate, "true")),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field lightClientUpdateJSON.signature_slot of type string",
		},
		{
			name:  "SignatureSlotInvalid",
			input: []byte(fmt.Sprintf(lcUpdateFormat, lcHeaderValue, syncCommittee, proofBranch, lcHeaderValue, proofBranch, syncAggregate, `"-1"`)),
			err:   "invalid value for slot: strconv.ParseUint: parsing \"-1\": invalid syntax",
		},
		{
			name:  "Good",
			input: []byte(fmt.Sprintf(lcUpdateFormat, lcHeaderValue, syncCommittee, proofBranch, lcHeaderValue, proofBranch, syncAggregate, slot)),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var res deneb.LightClientUpdate
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
