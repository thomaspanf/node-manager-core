package client

import "context"

type IBeaconApiProvider interface {
	Beacon_Attestations(ctx context.Context, blockId string) (AttestationsResponse, bool, error)
	Beacon_Block(ctx context.Context, blockId string) (BeaconBlockResponse, bool, error)
	Beacon_BlsToExecutionChanges_Post(ctx context.Context, request BLSToExecutionChangeRequest) error
	Beacon_Committees(ctx context.Context, stateId string, epoch *uint64) (CommitteesResponse, error)
	Beacon_FinalityCheckpoints(ctx context.Context, stateId string) (FinalityCheckpointsResponse, error)
	Beacon_Genesis(ctx context.Context) (GenesisResponse, error)
	Beacon_Header(ctx context.Context, blockId string) (BeaconBlockHeaderResponse, bool, error)
	Beacon_Validators(ctx context.Context, stateId string, ids []string) (ValidatorsResponse, error)
	Beacon_VoluntaryExits_Post(ctx context.Context, request VoluntaryExitRequest) error
	Config_DepositContract(ctx context.Context) (Eth2DepositContractResponse, error)
	Config_Spec(ctx context.Context) (Eth2ConfigResponse, error)
	Node_Syncing(ctx context.Context) (SyncStatusResponse, error)
	Validator_DutiesProposer(ctx context.Context, indices []string, epoch uint64) (ProposerDutiesResponse, error)
	Validator_DutiesSync_Post(ctx context.Context, indices []string, epoch uint64) (SyncDutiesResponse, error)
}
