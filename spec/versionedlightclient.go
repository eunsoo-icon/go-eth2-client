// Copyright © 2021 - 2023 Attestant Limited.
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

package spec

import (
	"github.com/attestantio/go-eth2-client/spec/altair"
	"github.com/attestantio/go-eth2-client/spec/capella"
	"github.com/attestantio/go-eth2-client/spec/deneb"
)

// VersionedLCHeader contains a versioned light client header.
type VersionedLCHeader struct {
	Version   DataVersion
	Altair    *altair.LightClientHeader
	Bellatrix *altair.LightClientHeader
	Capella   *capella.LightClientHeader
	Deneb     *deneb.LightClientHeader
}

// String returns a string version of the structure.
func (v *VersionedLCHeader) String() string {
	switch v.Version {
	case DataVersionPhase0:
		v.Version.NotSupport(v)
		return ""
	case DataVersionAltair:
		if v.Altair == nil {
			return ""
		}

		return v.Altair.String()
	case DataVersionBellatrix:
		if v.Bellatrix == nil {
			return ""
		}

		return v.Bellatrix.String()
	case DataVersionCapella:
		if v.Capella == nil {
			return ""
		}

		return v.Capella.String()
	case DataVersionDeneb:
		if v.Deneb == nil {
			return ""
		}

		return v.Deneb.String()
	default:
		return "unknown version"
	}
}

// VersionedLCBootstrap contains a versioned light client bootstrap.
type VersionedLCBootstrap struct {
	Version   DataVersion
	Altair    *altair.LightClientBootstrap
	Bellatrix *altair.LightClientBootstrap
	Capella   *capella.LightClientBootstrap
	Deneb     *deneb.LightClientBootstrap
}

// String returns a string version of the structure.
func (v *VersionedLCBootstrap) String() string {
	switch v.Version {
	case DataVersionPhase0:
		v.Version.NotSupport(v)
		return ""
	case DataVersionAltair:
		if v.Altair == nil {
			return ""
		}

		return v.Altair.String()
	case DataVersionBellatrix:
		if v.Bellatrix == nil {
			return ""
		}

		return v.Bellatrix.String()
	case DataVersionCapella:
		if v.Capella == nil {
			return ""
		}

		return v.Capella.String()
	case DataVersionDeneb:
		if v.Deneb == nil {
			return ""
		}

		return v.Deneb.String()
	default:
		return "unknown version"
	}
}

// VersionedLCFinalityUpdate contains a versioned light client finality update.
type VersionedLCFinalityUpdate struct {
	Version   DataVersion
	Altair    *altair.LightClientFinalityUpdate
	Bellatrix *altair.LightClientFinalityUpdate
	Capella   *capella.LightClientFinalityUpdate
	Deneb     *deneb.LightClientFinalityUpdate
}

// String returns a string version of the structure.
func (v *VersionedLCFinalityUpdate) String() string {
	switch v.Version {
	case DataVersionPhase0:
		v.Version.NotSupport(v)
		return ""
	case DataVersionAltair:
		if v.Altair == nil {
			return ""
		}

		return v.Altair.String()
	case DataVersionBellatrix:
		if v.Bellatrix == nil {
			return ""
		}

		return v.Bellatrix.String()
	case DataVersionCapella:
		if v.Capella == nil {
			return ""
		}

		return v.Capella.String()
	case DataVersionDeneb:
		if v.Deneb == nil {
			return ""
		}

		return v.Deneb.String()
	default:
		return "unknown version"
	}
}

// VersionedLCOptimisticUpdate contains a versioned light client optimistic update.
type VersionedLCOptimisticUpdate struct {
	Version   DataVersion
	Altair    *altair.LightClientOptimisticUpdate
	Bellatrix *altair.LightClientOptimisticUpdate
	Capella   *capella.LightClientOptimisticUpdate
	Deneb     *deneb.LightClientOptimisticUpdate
}

// String returns a string version of the structure.
func (v *VersionedLCOptimisticUpdate) String() string {
	switch v.Version {
	case DataVersionPhase0:
		v.Version.NotSupport(v)
		return ""
	case DataVersionAltair:
		if v.Altair == nil {
			return ""
		}

		return v.Altair.String()
	case DataVersionBellatrix:
		if v.Bellatrix == nil {
			return ""
		}

		return v.Bellatrix.String()
	case DataVersionCapella:
		if v.Capella == nil {
			return ""
		}

		return v.Capella.String()
	case DataVersionDeneb:
		if v.Deneb == nil {
			return ""
		}

		return v.Deneb.String()
	default:
		return "unknown version"
	}
}

// VersionedLCUpdate contains a versioned light client update.
type VersionedLCUpdate struct {
	Version   DataVersion
	Altair    *altair.LightClientUpdate
	Bellatrix *altair.LightClientUpdate
	Capella   *capella.LightClientUpdate
	Deneb     *deneb.LightClientUpdate
}

// String returns a string version of the structure.
func (v *VersionedLCUpdate) String() string {
	switch v.Version {
	case DataVersionPhase0:
		v.Version.NotSupport(v)
		return ""
	case DataVersionAltair:
		if v.Altair == nil {
			return ""
		}

		return v.Altair.String()
	case DataVersionBellatrix:
		if v.Bellatrix == nil {
			return ""
		}

		return v.Bellatrix.String()
	case DataVersionCapella:
		if v.Capella == nil {
			return ""
		}

		return v.Capella.String()
	case DataVersionDeneb:
		if v.Deneb == nil {
			return ""
		}

		return v.Deneb.String()
	default:
		return "unknown version"
	}
}