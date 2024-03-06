package config

import (
	"github.com/rocket-pool/node-manager-core/config/ids"
)

// Common configuration for all validator clients
type ValidatorClientCommonConfig struct {
	// Custom proposal graffiti
	Graffiti Parameter[string]

	// Toggle for enabling doppelganger detection
	DoppelgangerDetection Parameter[bool]

	// The port to expose VC metrics on
	MetricsPort Parameter[uint16]
}

// Generates a new common VC configuration
func NewValidatorClientCommonConfig() *ValidatorClientCommonConfig {
	return &ValidatorClientCommonConfig{
		Graffiti: Parameter[string]{
			ParameterCommon: &ParameterCommon{
				ID:          ids.GraffitiID,
				Name:        "Custom Graffiti",
				Description: "Add a short message to any blocks you propose, so the world can see what you have to say!\nIt has a 16 character limit.",
				MaxLength:   16,
				AffectsContainers: []ContainerID{
					ContainerID_ValidatorClient,
				},
				CanBeBlank:         true,
				OverwriteOnUpgrade: false,
			},
			Default: map[Network]string{
				Network_All: "",
			},
		},

		DoppelgangerDetection: Parameter[bool]{
			ParameterCommon: &ParameterCommon{
				ID:                 ids.DoppelgangerDetectionID,
				Name:               "Enable Doppelgänger Detection",
				Description:        "If enabled, your client will *intentionally* miss 1 or 2 attestations on startup to check if validator keys are already running elsewhere. If they are, it will disable validation duties for them to prevent you from being slashed.",
				AffectsContainers:  []ContainerID{ContainerID_ValidatorClient},
				CanBeBlank:         false,
				OverwriteOnUpgrade: false,
			},
			Default: map[Network]bool{
				Network_All: true,
			},
		},

		MetricsPort: Parameter[uint16]{
			ParameterCommon: &ParameterCommon{
				ID:                 ids.MetricsPortID,
				Name:               "Validator Client Metrics Port",
				Description:        "The port your Validator Client should expose its metrics on, if metrics collection is enabled.",
				AffectsContainers:  []ContainerID{ContainerID_ValidatorClient, ContainerID_Prometheus},
				CanBeBlank:         false,
				OverwriteOnUpgrade: false,
			},
			Default: map[Network]uint16{
				Network_All: 9101,
			},
		},
	}
}

// Get the title for the config
func (cfg *ValidatorClientCommonConfig) GetTitle() string {
	return "Common Validator Client"
}

// Get the parameters for this config
func (cfg *ValidatorClientCommonConfig) GetParameters() []IParameter {
	return []IParameter{
		&cfg.Graffiti,
		&cfg.DoppelgangerDetection,
		&cfg.MetricsPort,
	}
}

// Get the sections underneath this one
func (cfg *ValidatorClientCommonConfig) GetSubconfigs() map[string]IConfigSection {
	return map[string]IConfigSection{}
}
