package quote

import (
	"github.com/jamestrandung/go-cte-117/cte"
	"github.com/jamestrandung/go-cte-117/sample/service/scaffolding/calculation"
	"github.com/jamestrandung/go-cte-117/sample/service/scaffolding/fixedcost"
)

type inout interface {
	Input
}

type Input interface {
	fixedcost.Input
	calculation.Input
	GetIsFixedCostEnabled() bool
}

type result interface {
	GetTotalCost() float64
	GetVATAmount() float64
}

type FixedCostBranch cte.Result

func (c FixedCostBranch) CTEMetadata() interface{} {
	return struct {
		computer computer
		inout    inout
	}{}
}

func (c FixedCostBranch) GetTotalCost() float64 {
	outcome, err := c.Task.Outcome()
	if err != nil {
		return 0
	}

	return outcome.(result).GetTotalCost()
}

func (c FixedCostBranch) GetVATAmount() float64 {
	outcome, err := c.Task.Outcome()
	if err != nil {
		return 0
	}

	return outcome.(result).GetVATAmount()
}
