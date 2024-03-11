package ids

const (
	// Shared
	MaxPeersID              string = "maxPeers"
	ContainerTagID          string = "containerTag"
	AdditionalFlagsID       string = "additionalFlags"
	HttpPortID              string = "httpPort"
	OpenHttpPortsID         string = "openHttpPort"
	P2pPortID               string = "p2pPort"
	PortID                  string = "port"
	OpenPortID              string = "openPort"
	HttpUrlID               string = "httpUrl"
	EcID                    string = "executionClient"
	BnID                    string = "beaconNode"
	GraffitiID              string = "graffiti"
	DoppelgangerDetectionID string = "doppelgangerDetection"
	MetricsPortID           string = "metricsPort"

	// Besu
	BesuJvmHeapSizeID   string = "jvmHeapSize"
	BesuMaxBackLayersID string = "maxBackLayers"

	// Bitfly
	BitflySecretID      string = "bitflySecret"
	BitflyEndpointID    string = "bitflyEndpoint"
	BitflyMachineNameID string = "bitflyMachineName"

	// Exporter
	ExporterEnableRootFsID string = "enableRootFs"

	// External Execution
	ExternalEcWebsocketUrlID string = "wsUrl"

	// Fallback
	FallbackUseFallbackClientsID string = "useFallbackClients"
	FallbackEcHttpUrlID          string = "ecHttpUrl"
	FallbackBnHttpUrlID          string = "bnHttpUrl"

	// Geth
	GethEnablePbssID string = "enablePbss"

	// Lighthouse
	LighthouseQuicPortID string = "p2pQuicPort"

	// Local Beacon Node
	LocalBnCheckpointSyncUrlID string = "checkpointSyncUrl"
	LocalBnLighthouseID        string = "lighthouse"
	LocalBnLodestarID          string = "lodestar"
	LocalBnNimbusID            string = "nimbus"
	LocalBnPrysmID             string = "prysm"
	LocalBnTekuID              string = "teku"

	// Local Execution Client
	LocalEcWebsocketPortID string = "wsPort"
	LocalEcEnginePortID    string = "enginePort"
	LocalEcOpenApiPortsID  string = "openApiPorts"
	LocalEcBesuID          string = "besu"
	LocalEcGethID          string = "geth"
	LocalEcNethermindID    string = "nethermind"

	// Metrics
	MetricsEnableID       string = "enableMetrics"
	MetricsEnableBitflyID string = "enableBitflyNodeMetrics"
	MetricsEcPortID       string = "ecMetricsPort"
	MetricsBnPortID       string = "bnMetricsPort"
	MetricsDaemonPortID   string = "daemonMetricsPort"
	MetricsExporterPortID string = "exporterMetricsPort"
	MetricsGrafanaID      string = "grafana"
	MetricsPrometheusID   string = "prometheus"
	MetricsExporterID     string = "exporter"
	MetricsBitflyID       string = "bitfly"

	// Nethermind
	NethermindCacheSizeID         string = "cacheSize"
	NethermindPruneMemSizeID      string = "pruneMemSize"
	NethermindAdditionalModulesID string = "additionalModules"
	NethermindAdditionalUrlsID    string = "additionalUrls"

	// Nimbus
	NimbusPruningModeID string = "pruningMode"

	// Prysm
	PrysmRpcPortID     string = "rpcPort"
	PrysmOpenRpcPortID string = "openRpcPort"
	PrysmRpcUrlID      string = "prysmRpcUrl"

	// Teku
	TekuJvmHeapSizeID string = "jvmHeapSize"
	TekuArchiveModeID string = "archiveMode"
)
