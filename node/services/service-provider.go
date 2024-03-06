package services

import (
	"fmt"
	"path/filepath"
	"runtime"
	"time"

	"github.com/docker/docker/client"
	"github.com/fatih/color"
	"github.com/rocket-pool/node-manager-core/config"
	"github.com/rocket-pool/node-manager-core/eth"
	"github.com/rocket-pool/node-manager-core/node/wallet"
	"github.com/rocket-pool/node-manager-core/utils/log"
)

const (
	apiLogColor    color.Attribute = color.FgHiCyan
	walletLogColor color.Attribute = color.FgYellow
)

// A container for all of the various services used by the node service
type ServiceProvider struct {
	// Services
	cfg        config.IConfig
	resources  *config.NetworkResources
	nodeWallet *wallet.Wallet
	ecManager  *ExecutionClientManager
	bcManager  *BeaconClientManager
	docker     *client.Client
	txMgr      *eth.TransactionManager
	queryMgr   *eth.QueryManager
	debugMode  bool

	// TODO: find a better place for this than the common service provider
	apiLogger    *log.ColorLogger
	walletLogger *log.ColorLogger
}

// Creates a new ServiceProvider instance
func NewServiceProvider(cfg config.IConfig, clientTimeout time.Duration, debugMode bool) (*ServiceProvider, error) {
	// Loggers
	apiLogger := log.NewColorLogger(apiLogColor)
	walletLogger := log.NewColorLogger(walletLogColor)

	// Wallet
	resources := cfg.GetNetworkResources()
	nodeAddressPath := filepath.Join(cfg.GetNodeAddressFilePath())
	walletDataPath := filepath.Join(cfg.GetWalletFilePath())
	passwordPath := filepath.Join(cfg.GetPasswordFilePath())
	nodeWallet, err := wallet.NewWallet(&walletLogger, walletDataPath, nodeAddressPath, passwordPath, resources.ChainID)
	if err != nil {
		return nil, fmt.Errorf("error creating node wallet: %w", err)
	}

	// EC Manager
	primaryEcUrl, fallbackEcUrl := cfg.GetExecutionClientUrls()
	ecManager, err := NewExecutionClientManager(primaryEcUrl, fallbackEcUrl, resources.ChainID, clientTimeout)
	if err != nil {
		return nil, fmt.Errorf("error creating executon client manager: %w", err)
	}

	// Beacon manager
	primaryBnUrl, fallbackBnUrl := cfg.GetBeaconNodeUrls()
	bcManager, err := NewBeaconClientManager(primaryBnUrl, fallbackBnUrl, clientTimeout)
	if err != nil {
		return nil, fmt.Errorf("error creating Beacon client manager: %w", err)
	}

	// Docker client
	dockerClient, err := client.NewClientWithOpts(client.WithVersion(config.DockerApiVersion))
	if err != nil {
		return nil, fmt.Errorf("error creating Docker client: %w", err)
	}

	// TX Manager
	txMgr, err := eth.NewTransactionManager(ecManager, eth.DefaultSafeGasBuffer, eth.DefaultSafeGasMultiplier)
	if err != nil {
		return nil, fmt.Errorf("error creating transaction manager: %w", err)
	}

	// Query Manager - set the default concurrent run limit to half the CPUs so the EC doesn't get overwhelmed
	concurrentCallLimit := runtime.NumCPU()
	if concurrentCallLimit < 1 {
		concurrentCallLimit = 1
	}
	queryMgr := eth.NewQueryManager(ecManager, resources.MulticallAddress, concurrentCallLimit)

	// Create the provider
	provider := &ServiceProvider{
		cfg:        cfg,
		resources:  resources,
		nodeWallet: nodeWallet,
		ecManager:  ecManager,
		bcManager:  bcManager,
		docker:     dockerClient,
		txMgr:      txMgr,
		queryMgr:   queryMgr,
		apiLogger:  &apiLogger,
		debugMode:  debugMode,
	}
	return provider, nil
}

// ===============
// === Getters ===
// ===============

func (p *ServiceProvider) GetConfig() config.IConfig {
	return p.cfg
}

func (p *ServiceProvider) GetNetworkResources() *config.NetworkResources {
	return p.resources
}

func (p *ServiceProvider) GetWallet() *wallet.Wallet {
	return p.nodeWallet
}

func (p *ServiceProvider) GetEthClient() *ExecutionClientManager {
	return p.ecManager
}

func (p *ServiceProvider) GetBeaconClient() *BeaconClientManager {
	return p.bcManager
}

func (p *ServiceProvider) GetDocker() *client.Client {
	return p.docker
}

func (p *ServiceProvider) GetTransactionManager() *eth.TransactionManager {
	return p.txMgr
}

func (p *ServiceProvider) GetQueryManager() *eth.QueryManager {
	return p.queryMgr
}

func (p *ServiceProvider) GetApiLogger() *log.ColorLogger {
	return p.apiLogger
}

func (p *ServiceProvider) IsDebugMode() bool {
	return p.debugMode
}
