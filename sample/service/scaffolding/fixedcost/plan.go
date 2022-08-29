package fixedcost

import (
	"context"

	"github.com/jamestrandung/go-cte-117/sample/config"
	"github.com/jamestrandung/go-cte-117/sample/service/components/vat"
)

type SequentialPlan struct {
	Input
	totalCost float64
	vat.VATAmount
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
	p.preExecute()

	return config.Engine.ExecuteMasterPlan(ctx, p)
}

func (p *SequentialPlan) preExecute() {
	p.totalCost = p.GetFixedCost()
}
