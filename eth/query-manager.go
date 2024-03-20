package eth

import (
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	batch "github.com/rocket-pool/batch-query"
	"golang.org/x/sync/errgroup"
)

// Manages multicall-capable queries to the Execution layer.
type QueryManager struct {
	// The client to use when querying the chain.
	client IExecutionClient

	// Address of the multicall contract to use.
	multicallAddress common.Address

	// The maximum number of batches to query in parallel
	concurrentCallLimit int
}

// Creates a new query manager.
// concurrentCallLimit should be the maximum number of batches to query in parallel for batch calls. Negative values mean no limit.
func NewQueryManager(client IExecutionClient, multicallAddress common.Address, concurrentCallLimit int) *QueryManager {
	return &QueryManager{
		client:              client,
		multicallAddress:    multicallAddress,
		concurrentCallLimit: concurrentCallLimit,
	}
}

// Run a multicall query that doesn't perform any return type allocation.
// The 'query' function is an optional general-purpose function you can use to add whatever you want to the multicall
// before running it. The 'queryables' can be used to simply list a collection of IQueryable objects, each of which will
// run 'AddToQuery()' on the multicall for convenience.
func (q *QueryManager) Query(query func(*batch.MultiCaller) error, opts *bind.CallOpts, queryables ...IQueryable) error {
	// Create the multicaller
	mc, err := batch.NewMultiCaller(q.client, q.multicallAddress)
	if err != nil {
		return fmt.Errorf("error creating multicaller: %w", err)
	}

	// Add the query function
	if query != nil {
		err = query(mc)
		if err != nil {
			return fmt.Errorf("error running multicall query: %w", err)
		}
	}

	// Add the queryables
	AddQueryablesToMulticall(mc, queryables...)

	// Execute the multicall
	_, err = mc.FlexibleCall(true, opts)
	if err != nil {
		return fmt.Errorf("error executing multicall: %w", err)
	}

	return nil
}

// Run a multicall query that doesn't perform any return type allocation
// Use this if one of the calls is allowed to fail without interrupting the others; the returned result array provides information about the success of each call.
// The 'query' function is an optional general-purpose function you can use to add whatever you want to the multicall
// before running it. The 'queryables' can be used to simply list a collection of IQueryable objects, each of which will
// run 'AddToQuery()' on the multicall for convenience.
func (q *QueryManager) FlexQuery(query func(*batch.MultiCaller) error, opts *bind.CallOpts, queryables ...IQueryable) ([]bool, error) {
	// Create the multicaller
	mc, err := batch.NewMultiCaller(q.client, q.multicallAddress)
	if err != nil {
		return nil, fmt.Errorf("error creating multicaller: %w", err)
	}

	// Run the query
	if query != nil {
		err = query(mc)
		if err != nil {
			return nil, fmt.Errorf("error running multicall query: %w", err)
		}
	}

	// Add the queryables
	AddQueryablesToMulticall(mc, queryables...)

	// Execute the multicall
	return mc.FlexibleCall(false, opts)
}

// Create and execute a multicall query that is too big for one call and must be run in batches
func (q *QueryManager) BatchQuery(count int, batchSize int, query func(*batch.MultiCaller, int) error, opts *bind.CallOpts) error {
	// Sync
	var wg errgroup.Group
	wg.SetLimit(q.concurrentCallLimit)

	// Run getters in batches
	for i := 0; i < count; i += batchSize {
		i := i
		max := i + batchSize
		if max > count {
			max = count
		}

		// Load details
		wg.Go(func() error {
			mc, err := batch.NewMultiCaller(q.client, q.multicallAddress)
			if err != nil {
				return err
			}
			for j := i; j < max; j++ {
				err := query(mc, j)
				if err != nil {
					return fmt.Errorf("error running query adder: %w", err)
				}
			}
			_, err = mc.FlexibleCall(true, opts)
			if err != nil {
				return fmt.Errorf("error executing multicall: %w", err)
			}
			return nil
		})
	}

	// Wait for them all to complete
	if err := wg.Wait(); err != nil {
		return fmt.Errorf("error during multicall query: %w", err)
	}

	return nil
}

// Create and execute a multicall query that is too big for one call and must be run in batches.
// Use this if one of the calls is allowed to fail without interrupting the others; the returned result array provides information about the success of each call.
func (q *QueryManager) FlexBatchQuery(count int, batchSize int, query func(*batch.MultiCaller, int) error, handleResult func(bool, int) error, opts *bind.CallOpts) error {
	// Sync
	var wg errgroup.Group
	wg.SetLimit(q.concurrentCallLimit)

	// Run getters in batches
	for i := 0; i < count; i += batchSize {
		i := i
		max := i + batchSize
		if max > count {
			max = count
		}

		// Load details
		wg.Go(func() error {
			mc, err := batch.NewMultiCaller(q.client, q.multicallAddress)
			if err != nil {
				return err
			}
			for j := i; j < max; j++ {
				err := query(mc, j)
				if err != nil {
					return fmt.Errorf("error running query adder: %w", err)
				}
			}
			results, err := mc.FlexibleCall(false, opts)
			if err != nil {
				return fmt.Errorf("error executing multicall: %w", err)
			}
			for j, result := range results {
				err = handleResult(result, j+i)
				if err != nil {
					return fmt.Errorf("error running query result handler: %w", err)
				}
			}

			return nil
		})
	}

	// Wait for them all to complete
	if err := wg.Wait(); err != nil {
		return fmt.Errorf("error during multicall query: %w", err)
	}

	// Return
	return nil
}
