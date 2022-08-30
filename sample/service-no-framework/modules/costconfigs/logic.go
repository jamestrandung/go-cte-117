package costconfigs

import (
	"context"
	"github.com/jamestrandung/go-cte-117/sample/dependencies/configsfetcher"
	"github.com/jamestrandung/go-cte-117/sample/service-no-framework/helpers/future"
)

type GetterFunc func() (configsfetcher.MergedCostConfigs, error)

func GenerateCostConfigsLoader(ctx context.Context, dal configsfetcher.Fetcher) (future.LoaderFunc, GetterFunc) {
	resultChan := make(chan *future.Result, 1)
	return func() {
		defer close(resultChan)

		costConfigs := dal.Fetch()
		resultChan <- &future.Result{
			Value: costConfigs,
			Err:   nil,
		}
	}, newFuture(ctx, resultChan).Get
}

type Future struct {
	ftr *future.Future
}

func newFuture(parent context.Context, ch <-chan *future.Result) *Future {
	f := &Future{
		ftr: future.NewFuture(parent, ch),
	}

	return f
}

func (f *Future) Get() (configsfetcher.MergedCostConfigs, error) {
	result, err := f.ftr.Get()
	if result == nil {
		return configsfetcher.MergedCostConfigs{}, err
	}

	return result.(configsfetcher.MergedCostConfigs), err
}
