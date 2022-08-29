package loading

import (
	"github.com/jamestrandung/go-cte-117/sample/service/components/costconfigs"
	"github.com/jamestrandung/go-cte-117/sample/service/components/travelcost"
	"github.com/jamestrandung/go-cte-117/sample/service/components/travelplan"
)

type ParallelPlan struct {
	costconfigs.CostConfigs
	travelplan.TravelPlan
	travelcost.TravelCost
}

func (p *ParallelPlan) IsSequentialCTEPlan() bool {
	return false
}
