package cte

import (
	"context"
	"github.com/jamestrandung/go-concurrency-117/async"
	"reflect"
)

type ImpureComputer interface {
	Compute(ctx context.Context, p MasterPlan) (interface{}, error)
}

type ImpureComputerWithLoadingData interface {
	LoadingComputer
	Compute(ctx context.Context, p MasterPlan, data LoadingData) (interface{}, error)
}

type SideEffectComputer interface {
	Compute(ctx context.Context, p MasterPlan) error
}

type SideEffectComputerWithLoadingData interface {
	LoadingComputer
	Compute(ctx context.Context, p MasterPlan, data LoadingData) error
}

type SwitchComputer interface {
	Switch(ctx context.Context, p MasterPlan) (MasterPlan, error)
}

type SwitchComputerWithLoadingData interface {
	LoadingComputer
	Switch(ctx context.Context, p MasterPlan, data LoadingData) (MasterPlan, error)
}

type LoadingComputer interface {
	Load(ctx context.Context, p MasterPlan) (interface{}, error)
}

var emptyLoadingData = LoadingData{}

type LoadingData struct {
	Data interface{}
	Err  error
}

type toExecutePlan struct {
	mp MasterPlan
}

type bridgeComputer struct {
	computeFn func(ctx context.Context, p MasterPlan, data LoadingData) (interface{}, error)
}

func (bc bridgeComputer) Compute(ctx context.Context, p MasterPlan, data LoadingData) (interface{}, error) {
	return bc.computeFn(ctx, p, data)
}

type computerWrapper struct {
	LoadingComputer
	bridgeComputer
}

func newComputerWrapper(rawComputer interface{}) computerWrapper {
	switch c := rawComputer.(type) {
	case ImpureComputerWithLoadingData:
		return computerWrapper{
			LoadingComputer: c,
			bridgeComputer: bridgeComputer{
				computeFn: func(ctx context.Context, p MasterPlan, data LoadingData) (interface{}, error) {
					return c.Compute(ctx, p, data)
				},
			},
		}
	case ImpureComputer:
		return computerWrapper{
			bridgeComputer: bridgeComputer{
				computeFn: func(ctx context.Context, p MasterPlan, data LoadingData) (interface{}, error) {
					return c.Compute(ctx, p)
				},
			},
		}
	case SideEffectComputerWithLoadingData:
		return computerWrapper{
			LoadingComputer: c,
			bridgeComputer: bridgeComputer{
				computeFn: func(ctx context.Context, p MasterPlan, data LoadingData) (interface{}, error) {
					return struct{}{}, c.Compute(ctx, p, data)
				},
			},
		}
	case SideEffectComputer:
		return computerWrapper{
			bridgeComputer: bridgeComputer{
				computeFn: func(ctx context.Context, p MasterPlan, data LoadingData) (interface{}, error) {
					return struct{}{}, c.Compute(ctx, p)
				},
			},
		}
	case SwitchComputerWithLoadingData:
		return computerWrapper{
			LoadingComputer: c,
			bridgeComputer: bridgeComputer{
				computeFn: func(ctx context.Context, p MasterPlan, data LoadingData) (interface{}, error) {
					mp, err := c.Switch(ctx, p, data)

					return toExecutePlan{
						mp: mp,
					}, err
				},
			},
		}
	case SwitchComputer:
		return computerWrapper{
			bridgeComputer: bridgeComputer{
				computeFn: func(ctx context.Context, p MasterPlan, data LoadingData) (interface{}, error) {
					mp, err := c.Switch(ctx, p)

					return toExecutePlan{
						mp: mp,
					}, err
				},
			},
		}
	default:
		panic(ErrInvalidComputerType.Err(reflect.TypeOf(c)))
	}
}

func (w computerWrapper) Compute(ctx context.Context, p MasterPlan, data LoadingData) (interface{}, error) {
	return w.bridgeComputer.Compute(ctx, p, data)
}

type SideEffect struct{}

type SyncSideEffect struct{}

type Result struct {
	Task async.Task
}

func newResult(t async.Task) Result {
	return Result{
		Task: t,
	}
}

type SyncResult struct {
	Outcome interface{}
}

func newSyncResult(o interface{}) SyncResult {
	return SyncResult{
		Outcome: o,
	}
}
