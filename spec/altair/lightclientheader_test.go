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

func TestLightClientHeaderJSON(t *testing.T) {
	const (
		format = `{"beacon":%s}`
		header = `{"slot":"4943744","proposer_index":"222870","parent_root":"0x0c594acb2c7ec3564590fd2feb6724cfcf786faf51fe2a284154516c2903c153","state_root":"0x237962d02698b2f5f37f3a7c43dfae0e2fe28e103225237bc7f09938c8527eaa","body_root":"0xff42d5726526628ce27c4ca89172ccf5c562adbfec64c22d494b6f8bd03dc034"}`
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
			err:   "invalid JSON: json: cannot unmarshal array into Go value of type altair.lightClientHeaderJSON",
		},
		{
			name:  "BeaconMissing",
			input: []byte(fmt.Sprintf(`{}`)),
			err:   "beacon missing",
		},
		{
			name:  "Good",
			input: []byte(fmt.Sprintf(format, header)),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var res altair.LightClientHeader
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
