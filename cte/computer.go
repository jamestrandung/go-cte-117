package cte

import (
	"context"

	"github.com/jamestrandung/go-concurrency-117/async"
)

type ImpureComputer interface {
	Compute(ctx context.Context, p MasterPlan) (interface{}, error)
}

type SideEffectComputer interface {
	Compute(ctx context.Context, p MasterPlan) error
}

type SwitchComputer interface {
	Switch(ctx context.Context, p MasterPlan) (MasterPlan, error)
}

type toExecutePlan struct {
	mp MasterPlan
}

type bridgeComputer struct {
	se SideEffectComputer
	sw SwitchComputer
}

func (bc bridgeComputer) Compute(ctx context.Context, p MasterPlan) (interface{}, error) {
	if bc.se != nil {
		return struct{}{}, bc.se.Compute(ctx, p)
	}

	mp, err := bc.sw.Switch(ctx, p)

	return toExecutePlan{
		mp: mp,
	}, err
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
