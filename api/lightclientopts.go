// Copyright © 2023 Attestant Limited.
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

package api

// LightClientBootstrapOpts are the options for obtaining light client bootstrap.
type LightClientBootstrapOpts struct {
	Common CommonOpts

	// Block is the ID of the block which the data is obtained.
	Block string
}

// LightClientUpdatesOpts are the options for obtaining light client update instances.
type LightClientUpdatesOpts struct {
	Common CommonOpts

	StartPeriod, Count uint64
}
