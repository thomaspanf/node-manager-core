package services

import (
	"context"
	"fmt"
	"path/filepath"
	"runtime"
	"time"

	"github.com/docker/docker/client"
	"github.com/rocket-pool/node-manager-core/config"
	"github.com/rocket-pool/node-manager-core/eth"
	"github.com/rocket-pool/node-manager-core/log"
	"github.com/rocket-pool/node-manager-core/node/wallet"
)

const (
	DockerApiVersion string = "1.40"
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

	// Context for cancelling long operations
	ctx    context.Context
	cancel context.CancelFunc

	// Logging
	debugMode   bool
	apiLogger   *log.Logger
	tasksLogger *log.Logger
}

// Creates a new ServiceProvider instance
func NewServiceProvider(cfg config.IConfig, clientTimeout time.Duration, debugMode bool) (*ServiceProvider, error) {
	// Make the API logger
	apiLogger, err := log.NewLogger(cfg.GetApiLogFilePath(), debugMode, false)
	if err != nil {
		return nil, fmt.Errorf("error creating API logger: %w", err)
	}

	// Make the tasks logger
	tasksLogger, err := log.NewLogger(cfg.GetTasksLogFilePath(), debugMode, false)
	if err != nil {
		return nil, fmt.Errorf("error creating tasks logger: %w", err)
	}

	// Wallet
	resources := cfg.GetNetworkResources()
	nodeAddressPath := filepath.Join(cfg.GetNodeAddressFilePath())
	walletDataPath := filepath.Join(cfg.GetWalletFilePath())
	passwordPath := filepath.Join(cfg.GetPasswordFilePath())
	nodeWallet, err := wallet.NewWallet(tasksLogger, walletDataPath, nodeAddressPath, passwordPath, resources.ChainID)
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
	dockerClient, err := client.NewClientWithOpts(client.WithVersion(DockerApiVersion))
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

	// Context for handling task cancellation during shutdown
	ctx, cancel := context.WithCancel(context.Background())

	// Log startup
	apiLogger.Info("Starting API logger.")
	tasksLogger.Info("Starting Tasks logger.")

	// Create the provider
	provider := &ServiceProvider{
		cfg:         cfg,
		resources:   resources,
		nodeWallet:  nodeWallet,
		ecManager:   ecManager,
		bcManager:   bcManager,
		docker:      dockerClient,
		txMgr:       txMgr,
		queryMgr:    queryMgr,
		ctx:         ctx,
		cancel:      cancel,
		debugMode:   debugMode,
		apiLogger:   apiLogger,
		tasksLogger: tasksLogger,
	}
	return provider, nil
}

// Closes the service provider and its underlying services
func (p *ServiceProvider) Close() {
	p.apiLogger.Close()
	p.tasksLogger.Close()
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

func (p *ServiceProvider) GetApiLogger() *log.Logger {
	return p.apiLogger
}

func (p *ServiceProvider) GetTasksLogger() *log.Logger {
	return p.tasksLogger
}

func (p *ServiceProvider) IsDebugMode() bool {
	return p.debugMode
}

func (p *ServiceProvider) GetBaseContext() context.Context {
	return p.ctx
}

func (p *ServiceProvider) CancelContextOnShutdown() {
	p.cancel()
}
