package calculation

import (
	"context"
	"github.com/jamestrandung/go-cte-117/sample/service/components/vat"

	"github.com/jamestrandung/go-cte-117/sample/config"
	"github.com/jamestrandung/go-cte-117/sample/service/components/platformfee"
)

type SequentialPlan struct {
	Input
	preHook
	totalCost float64
	platformfee.PlatformFee
	vat.VATAmount
	postHook
}

func NewPlan(in Input) *SequentialPlan {
	return &SequentialPlan{
		Input: in,
	}
}

func (p *SequentialPlan) IsSequentialCTEPlan() bool {
	return true
}

func (p *SequentialPlan) Execute(ctx context.Context) error {
	return config.Engine.ExecuteMasterPlan(ctx, p)
}
