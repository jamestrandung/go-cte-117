package travelcost

import (
	"context"
	"github.com/jamestrandung/go-cte-117/sample/service-no-framework/helpers/future"
	"github.com/jamestrandung/go-cte-117/sample/service-no-framework/modules/costconfigs"
	"github.com/jamestrandung/go-cte-117/sample/service-no-framework/modules/travelplan"
)

type GetterFunc func() (float64, error)

func GenerateTravelCostLoader(ctx context.Context, costConfigsGetter costconfigs.GetterFunc, travelRouteGetter travelplan.GetterFunc) (future.LoaderFunc, GetterFunc) {
	resultChan := make(chan *future.Result, 1)
	return func() {
		defer close(resultChan)

		costConfigs, err := costConfigsGetter()
		if err != nil {
			resultChan <- &future.Result{
				Value: nil,
				Err:   err,
			}

			return
		}

		route, err := travelRouteGetter()
		if err != nil {
			resultChan <- &future.Result{
				Value: nil,
				Err:   err,
			}

			return
		}

		travelCost := costConfigs.BaseCost + costConfigs.CostPerKilometer*route.Distance + costConfigs.CostPerMinute*route.Duration
		resultChan <- &future.Result{
			Value: travelCost,
			Err:   err,
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

func (f *Future) Get() (float64, error) {
	result, err := f.ftr.Get()
	if result == nil {
		return 0, err
	}

	return result.(float64), err
}
