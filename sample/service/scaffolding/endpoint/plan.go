package endpoint

import (
	"context"
	"github.com/jamestrandung/go-cte-117/sample/config"
	"github.com/jamestrandung/go-cte-117/sample/service/components/quote"
	"github.com/jamestrandung/go-cte-117/sample/service/components/streaming"
	"github.com/jamestrandung/go-cte-117/sample/service/scaffolding/loading"
)

type SequentialPlan struct {
	Input
	Dependencies
	loading.ParallelPlan
	quote.FixedCostBranch
	streaming.CostStreaming
}

func NewPlan(r Input, d Dependencies) *SequentialPlan {
	return &SequentialPlan{
		Input:        r,
		Dependencies: d,
	}
}

func (p *SequentialPlan) IsSequentialCTEPlan() bool {
	return true
}

func (p *SequentialPlan) Execute(ctx context.Context) error {
	return config.Engine.ExecuteMasterPlan(ctx, p)
}
