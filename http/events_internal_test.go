// Copyright © 2020, 2021 Attestant Limited.
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

package http

import (
	"bytes"
	"context"
	"os"
	"testing"
	"time"

	"github.com/r3labs/sse/v2"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"

	client "github.com/attestantio/go-eth2-client"
	api "github.com/attestantio/go-eth2-client/api/v1"
)

// timeout for tests.
var timeout = 60 * time.Second

func TestEventHandler(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	handled := false
	handler := func(*api.Event) {
		handled = true
	}

	tests := []struct {
		name    string
		message *sse.Event
		handler client.EventHandlerFunc
		handled bool
	}{
		{
			name:    "MessageNil",
			handler: handler,
			handled: false,
		},
		{
			name:    "MessageEmpty",
			message: &sse.Event{},
			handler: handler,
			handled: false,
		},
		{
			name: "EventUnknown",
			message: &sse.Event{
				Event: []byte("unknown"),
			},
			handler: handler,
			handled: false,
		},
		{
			name: "HandlerNil",
			message: &sse.Event{
				Event: []byte("head"),
			},
			handled: false,
		},
		{
			name: "HeadGood",
			message: &sse.Event{
				Event: []byte("head"),
				Data:  []byte(`{"slot":"4095940","block":"0x73d83c5f925716c9bd2d1e9c339fb99b0ec4addef3e93f6f35d4c5f1de7ae092","state":"0xead0e6eb4004576546864f10cfa4aeac31afbf96abc405a86c00cbda8f3e8ed0","epoch_transition":false,"previous_duty_dependent_root":"0xeca94cc9180212a2cff2659289cc7e6f2df08a645120e35e25d09c2ddc7db5f1","current_duty_dependent_root":"0xdda286c4a096fc8ec0d6ba9e14e688cbb046bfb33462fdf94953e75d0cea0074","execution_optimistic":false}`),
			},
			handler: handler,
			handled: true,
		},
		{
			name: "BlockGood",
			message: &sse.Event{
				Event: []byte("block"),
				Data:  []byte(`{"slot":"4095943","block":"0x1c3981b7439cd2dc53dca1a99122e1cacb36a13796d426d4c8a03ba745cb0c8b","execution_optimistic":false}`),
			},
			handler: handler,
			handled: true,
		},
		{
			name: "AttestationGood",
			message: &sse.Event{
				Event: []byte("attestation"),
				Data:  []byte(`{"aggregation_bits":"0x00002840403040000000020008800040008042800000020220","data":{"slot":"4095945","index":"12","beacon_block_root":"0xff27c7551bf4cfe4dc4cce00920e7a5c5074860d1dbd8aa8b3b5f888523f51ff","source":{"epoch":"127997","root":"0x38758fb180459583bd5e8e1a31711eb09e63eb92be974485397e9a2c57de2783"},"target":{"epoch":"127998","root":"0x46d4629861bd81cfc94007501b4edb1b3ca9444b41d7a98681b6c2f4bdb978bd"}},"signature":"0xacb9f562a28c4ef5b60b88678068ea51573a3237d3331dda3b2d845a0d03bc56ab2994d2deb90d9f074a8bdab59945150d0a7717e74b1bf2627f8971c81091f724c211dfce8fa16fb839c6a1bfd341ddec5e7eb88472682fd1a170e373660534"}`),
			},
			handler: handler,
			handled: true,
		},
		{
			name: "VoluntaryExitGood",
			message: &sse.Event{
				Event: []byte("voluntary_exit"),
				Data:  []byte(`{"message":{"epoch":"1", "validator_index":"1"}, "signature":"0x1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505cc411d61252fb6cb3fa0017b679f8bb2305b26a285fa2737f175668d0dff91cc1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505"}`),
			},
			handler: handler,
			handled: true,
		},
		{
			name: "FinalizedCheckpointGood",
			message: &sse.Event{
				Event: []byte("finalized_checkpoint"),
				Data:  []byte(`{"block":"0x38758fb180459583bd5e8e1a31711eb09e63eb92be974485397e9a2c57de2783","state":"0x9c237b2a66df8636f816e6b2c8860ba287fc5b817d882b1be8b7111486fb4ddc","epoch":"127997","execution_optimistic":false}`),
			},
			handler: handler,
			handled: true,
		},
		{
			name: "ChainReorgGood",
			message: &sse.Event{
				Event: []byte("chain_reorg"),
				Data:  []byte(`{"slot":"4100237","depth":"2","old_head_block":"0x5c988c12b7d8638c06e6b9511e09e5e28511a16a33c153413416d3fd5d95353a","new_head_block":"0x783ccb4310c0ddab3ea500f8d2b88c5ad8d6b2d601513f4ebf491066cda1d180","old_head_state":"0xdbad017808a1c5a77866fccc3e15f14a67585fff25b3a4d947ae9cc6d937b4ab","new_head_state":"0x67f0302bb939f64b1feaaa907b30c9631c0588480f2305578377bfd87df68b95","epoch":"128132","execution_optimistic":false}`),
			},
			handler: handler,
			handled: true,
		},
		{
			name: "ContributionAndProofGood",
			message: &sse.Event{
				Event: []byte("contribution_and_proof"),
				Data:  []byte(`{"message":{"aggregator_index":"355177","contribution":{"slot":"4095970","beacon_block_root":"0xfe17cc29bb937740eead4b84716751b00e361891dae7c8f0e98f4deb5f753cbc","subcommittee_index":"0","aggregation_bits":"0xffffffffffffffffffffdfffffffffff","signature":"0x926e8fbf2b8599f76a42e2dd02b954853d3841577a0c68303fb9a5690f7973e95454bc1df03118d53c80e3cf13dc33490f2516aefb4cb7766a724a10dd0536811ec43b7e5f08442ef7dbc072b4b484ea1acde78e5aae8d636b06dd18677c535b"},"selection_proof":"0x84c805b21a40315dc19ea89e6d64d8f5e913d5e003a813d74a23f16add72250791a950b7508e2ea69adab8169c2800b70d26fea4cdb51f1fb96d939baaa44567eb9e96d2021799478fd3f557326a62060215be95465fcc9c49f67dd8685a20bc"},"signature":"0x8093efce898e36cab5ab2b198a48046d029b36909a29ec33ca7075f389133288c4d7e13cf3e20396612050d4aebe9212154fd5a2be4bf356e6191600d65906d5c404bd46c95ae20fe4bc5e18c6e2808c97a4572f995bf90db8aaf3fd84fb87ac"}`),
			},
			handler: handler,
			handled: true,
		},
		{
			name: "LightClientFinalityUpdateGood",
			message: &sse.Event{
				Event: []byte("light_client_finality_update"),
				Data:  []byte(`{"attested_header":{"slot":"4957824","proposer_index":"102015","parent_root":"0x8b8c9f91b6d14a06098db07040dfc46c4ab21f55d213559f6f02636ff15977d5","state_root":"0x6f135bf53593ce342bbfaff7278613b7a8c3709b065b612e8b47784f76ae4da0","body_root":"0x61595997dfea35dce29332455b260377ab1d6a788d49269372ad008e0b78215b"},"finalized_header":{"slot":"4957760","proposer_index":"346845","parent_root":"0xcbde4e99a232e6df7a6d8ef0ec495558e5d282f926ab4e5048c9af321d5fda72","state_root":"0xd357a7b09ead054ac22786ff45a3459f067ac1795d4a149f6d8e5edece964501","body_root":"0xd69dc5cde1f986f74dd8330bb055154ce6d098857e17e927362c47a343e18f3f"},"finality_branch":["0x325d020000000000000000000000000000000000000000000000000000000000","0x4ea8e197f4d8cb99074425bdfe07eae17c31a2117cfd7af9257c371edbeb0dbf","0xf40e5559be4761afc95bfaec1b87116a0fcd09e3c01f1d5fbc8befc0704ac1f2","0x91e25054f15600fd48b054361e570ff35425e143392ea3ad566d3cb0f8ee6a2f","0x929ee3ff1bce3066825611a03e03b89e94046d8f158bbf24969074d3fe26143f","0x365a7a573241a9a5b278f33d46242f8177a77c4467cb750d50a9943773200a09"],"sync_aggregate":{"sync_committee_bits":"0xffffefffeffffbfffff9efffffffbffffffdf7ff7ffffffbdffffdffffffffffffff7ffffeffffffffffffffbffffffefffffbffff7f7fffffffffffffffeeff","sync_committee_signature":"0xa63cd88b564d6db06fb3a3d953ddddd9d0e8d6501e2bc9f8df12b27b490de6a98fc537f5551289fc6cc954d80afa63d10bf5f14a42aa6c94a5ee1df6480c48c89a8bd1c8d785f38a09a45e16feb584e5b8417b72202d1c520f8231229f38cbc3"},"signature_slot": "4957825"}`),
			},
			handler: handler,
			handled: true,
		},
		{
			name: "LightClientOptimisticUpdateGood",
			message: &sse.Event{
				Event: []byte("light_client_optimistic_update"),
				Data:  []byte(`{"attested_header":{"slot":"4957842","proposer_index":"11256","parent_root":"0xfc8390a8789b7fefcf950e009d493779f53e422b72a6816b98c35b1024d09c40","state_root":"0x6504e27e01f50f9b86702abe6e06fe656b722c95204d1f97e2204de63a672d18","body_root":"0x31b30a8f4d053fe117f3c9571d9989ef597cf04f178f0034e0bdd9139365f90d"},"sync_aggregate":{"sync_committee_bits":"0xffffffffffffffffffffffffffffffffffffffff7ffffffffffffffefffffffffbfffffffeffffefffffffffffffffffffffffffffffffffffffdfffffffffff","sync_committee_signature":"0x9300e6e2359d1fd4a96623635d0dc1315623cdf6d5b3406efadee3e55a3907ec27acec609eb18658270fcc528e1eb8ae0fa2e3440ce3dabc814fae612af90ddef6f0a97f471e0eeba47d3d0c807873d5d9efed0f6d84c6f2f0a7cba299d88de6"},"signature_slot":"4957843"}`),
			},
			handler: handler,
			handled: true,
		},
	}

	s, err := New(ctx,
		WithTimeout(timeout),
		WithAddress(os.Getenv("HTTP_ADDRESS")),
	)
	require.NoError(t, err)

	h, isHTTPService := s.(*Service)
	require.True(t, isHTTPService)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			handled = false
			log := zerolog.New(&bytes.Buffer{})
			ctx = log.WithContext(ctx)
			h.handleEvent(ctx, test.message, test.handler)
			require.Equal(t, test.handled, handled)
		})
	}
}
