package travelplan

import (
	"context"
	"github.com/jamestrandung/go-cte-117/sample/config"
	"github.com/jamestrandung/go-cte-117/sample/dependencies/mapservice"
	"github.com/jamestrandung/go-cte-117/sample/service-no-framework/helpers/future"
)

type GetterFunc func() (mapservice.Route, error)

func GenerateRouteLoader(ctx context.Context, dal mapservice.Service, pointA, pointB string) (future.LoaderFunc, GetterFunc) {
	resultChan := make(chan *future.Result, 1)
	return func() {
		defer close(resultChan)

		route, err := dal.GetRoute(pointA, pointB)
		if err != nil {
			config.Printf("Building route from %s to %s using straight-line distance\n", pointA, pointB)
			route = mapservice.Route{
				Distance: 4,
				Duration: 5,
			}
		}

		resultChan <- &future.Result{
			Value: route,
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

func (f *Future) Get() (mapservice.Route, error) {
	result, err := f.ftr.Get()
	if result == nil {
		return mapservice.Route{}, err
	}

	return result.(mapservice.Route), err
}
