package cte

import (
	"context"
	"reflect"

	"github.com/jamestrandung/go-concurrency-117/async"
)

//go:generate mockery --name ImpureComputer --case=underscore --inpackage
type ImpureComputer interface {
	Compute(ctx context.Context, p MasterPlan) (interface{}, error)
}

//go:generate mockery --name ImpureComputerWithLoadingData --case=underscore --inpackage
type ImpureComputerWithLoadingData interface {
	LoadingComputer
	Compute(ctx context.Context, p MasterPlan, data LoadingData) (interface{}, error)
}

//go:generate mockery --name SideEffectComputer --case=underscore --inpackage
type SideEffectComputer interface {
	Compute(ctx context.Context, p MasterPlan) error
}

//go:generate mockery --name SideEffectComputerWithLoadingData --case=underscore --inpackage
type SideEffectComputerWithLoadingData interface {
	LoadingComputer
	Compute(ctx context.Context, p MasterPlan, data LoadingData) error
}

//go:generate mockery --name SwitchComputer --case=underscore --inpackage
type SwitchComputer interface {
	Switch(ctx context.Context, p MasterPlan) (MasterPlan, error)
}

//go:generate mockery --name SwitchComputerWithLoadingData --case=underscore --inpackage
type SwitchComputerWithLoadingData interface {
	LoadingComputer
	Switch(ctx context.Context, p MasterPlan, data LoadingData) (MasterPlan, error)
}

type LoadingComputer interface {
	Load(ctx context.Context, p MasterPlan) (interface{}, error)
}

type LoadingData struct {
	Data interface{}
	Err  error
}

type toExecutePlan struct {
	mp MasterPlan
}

type loadFn func(ctx context.Context, p MasterPlan) (interface{}, error)
type computeFn func(ctx context.Context, p MasterPlan, data LoadingData) (interface{}, error)

type delegatingComputer struct {
	loadFn
	computeFn
}

func newDelegatingComputer(rawComputer interface{}) delegatingComputer {
	switch c := rawComputer.(type) {
	case ImpureComputerWithLoadingData:
		return delegatingComputer{
			loadFn: func(ctx context.Context, p MasterPlan) (interface{}, error) {
				return c.Load(ctx, p)
			},
			computeFn: func(ctx context.Context, p MasterPlan, data LoadingData) (interface{}, error) {
				return c.Compute(ctx, p, data)
			},
		}
	case ImpureComputer:
		return delegatingComputer{
			computeFn: func(ctx context.Context, p MasterPlan, data LoadingData) (interface{}, error) {
				return c.Compute(ctx, p)
			},
		}
	case SideEffectComputerWithLoadingData:
		return delegatingComputer{
			loadFn: func(ctx context.Context, p MasterPlan) (interface{}, error) {
				return c.Load(ctx, p)
			},
			computeFn: func(ctx context.Context, p MasterPlan, data LoadingData) (interface{}, error) {
				return struct{}{}, c.Compute(ctx, p, data)
			},
		}
	case SideEffectComputer:
		return delegatingComputer{
			computeFn: func(ctx context.Context, p MasterPlan, data LoadingData) (interface{}, error) {
				return struct{}{}, c.Compute(ctx, p)
			},
		}
	case SwitchComputerWithLoadingData:
		return delegatingComputer{
			loadFn: func(ctx context.Context, p MasterPlan) (interface{}, error) {
				return c.Load(ctx, p)
			},
			computeFn: func(ctx context.Context, p MasterPlan, data LoadingData) (interface{}, error) {
				mp, err := c.Switch(ctx, p, data)

				return toExecutePlan{
					mp: mp,
				}, err
			},
		}
	case SwitchComputer:
		return delegatingComputer{
			computeFn: func(ctx context.Context, p MasterPlan, data LoadingData) (interface{}, error) {
				mp, err := c.Switch(ctx, p)

				return toExecutePlan{
					mp: mp,
				}, err
			},
		}
	default:
		panic(ErrInvalidComputerType.Err(reflect.TypeOf(c)))
	}
}

func (dc delegatingComputer) Load(ctx context.Context, p MasterPlan) (interface{}, error) {
	if dc.loadFn == nil {
		return nil, nil
	}

	return dc.loadFn(ctx, p)
}

func (dc delegatingComputer) Compute(ctx context.Context, p MasterPlan, data LoadingData) (interface{}, error) {
	return dc.computeFn(ctx, p, data)
}

type SideEffect struct {
	isSync bool // This is a dummy field to prevent SideEffect & SyncSideEffect from being convertible to each other
	// to help fieldAnalyzer.createComputerComponent works
}

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
