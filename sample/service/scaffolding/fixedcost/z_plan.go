package fixedcost

import (
	"github.com/jamestrandung/go-cte-117/sample/config"
)

func init() {
	config.Engine.AnalyzePlan(&SequentialPlan{})
}

func (p *SequentialPlan) GetTotalCost() float64 {
	return p.totalCost
}

func (p *SequentialPlan) SetTotalCost(totalCost float64) {
	p.totalCost = totalCost
}
