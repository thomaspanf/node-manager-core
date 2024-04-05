package config

import (
	"runtime"

	"github.com/pbnjay/memory"
	"github.com/rocket-pool/node-manager-core/config/ids"
)

// Constants
const (
	rethTagProd string = "ghcr.io/paradigmxyz/reth:v0.2.0-beta.5"
	rethTagTest string = "ghcr.io/paradigmxyz/reth:v0.2.0-beta.5"
)

// Configuration for Reth
type RethConfig struct {
	// Size of Reth's Cache
	CacheSize Parameter[uint64]

	// Max number of P2P peers to connect to
	MaxPeers Parameter[uint16]

	// The Docker Hub tag for Reth
	ContainerTag Parameter[string]

	// Custom command line flags
	AdditionalFlags Parameter[string]
}

// Generates a new Reth configuration
func NewRethConfig() *RethConfig {
	return &RethConfig{
		CacheSize: Parameter[uint64]{
			ParameterCommon: &ParameterCommon{
				ID:                 ids.CacheSizeID,
				Name:               "Cache Size",
				Description:        "The amount of RAM (in MB) you want Reth's cache to use. Larger values mean your disk space usage will increase slower, and you will have to prune less frequently. The default is based on how much total RAM your system has but you can adjust it manually.",
				AffectsContainers:  []ContainerID{ContainerID_ExecutionClient},
				CanBeBlank:         false,
				OverwriteOnUpgrade: false,
			},
			Default: map[Network]uint64{
				Network_All: calculateRethCache(),
			},
		},

		MaxPeers: Parameter[uint16]{
			ParameterCommon: &ParameterCommon{
				ID:                 ids.MaxPeersID,
				Name:               "Max Peers",
				Description:        "The maximum number of peers Reth should connect to. This can be lowered to improve performance on low-power systems or constrained networks. We recommend keeping it at 12 or higher.",
				AffectsContainers:  []ContainerID{ContainerID_ExecutionClient},
				CanBeBlank:         false,
				OverwriteOnUpgrade: false,
			},
			Default: map[Network]uint16{Network_All: calculateRethPeers()},
		},

		ContainerTag: Parameter[string]{
			ParameterCommon: &ParameterCommon{
				ID:                 ids.ContainerTagID,
				Name:               "Container Tag",
				Description:        "The tag name of the Reth container you want to use.",
				AffectsContainers:  []ContainerID{ContainerID_ExecutionClient},
				CanBeBlank:         false,
				OverwriteOnUpgrade: true,
			},
			Default: map[Network]string{
				Network_Mainnet: rethTagProd,
				Network_Holesky: rethTagTest,
			},
		},

		AdditionalFlags: Parameter[string]{
			ParameterCommon: &ParameterCommon{
				ID:                 ids.AdditionalFlagsID,
				Name:               "Additional Flags",
				Description:        "Additional custom command line flags you want to pass to Reth, to take advantage of other settings that aren't covered here.",
				AffectsContainers:  []ContainerID{ContainerID_ExecutionClient},
				CanBeBlank:         true,
				OverwriteOnUpgrade: false,
			},
			Default: map[Network]string{
				Network_All: "",
			},
		},
	}
}

// Get the title for the config
func (cfg *RethConfig) GetTitle() string {
	return "Reth"
}

// Get the config.Parameters for this config
func (cfg *RethConfig) GetParameters() []IParameter {
	return []IParameter{
		&cfg.CacheSize,
		&cfg.MaxPeers,
		&cfg.ContainerTag,
		&cfg.AdditionalFlags,
	}
}

// Get the sections underneath this one
func (cfg *RethConfig) GetSubconfigs() map[string]IConfigSection {
	return map[string]IConfigSection{}
}

// Calculate the recommended size for Reth's cache based on the amount of system RAM
func calculateRethCache() uint64 {
	totalMemoryGB := memory.TotalMemory() / 1024 / 1024 / 1024

	if totalMemoryGB == 0 {
		return 0
	} else if totalMemoryGB < 9 {
		return 256
	} else if totalMemoryGB < 13 {
		return 2048
	} else if totalMemoryGB < 17 {
		return 4096
	} else if totalMemoryGB < 25 {
		return 8192
	} else if totalMemoryGB < 33 {
		return 12288
	} else {
		return 16384
	}
}

// Calculate the default number of Reth peers
func calculateRethPeers() uint16 {
	if runtime.GOARCH == "arm64" {
		return 25
	}
	return 50
}
