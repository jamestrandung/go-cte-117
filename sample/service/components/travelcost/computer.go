package travelcost

import (
	"context"

	"github.com/jamestrandung/go-cte-117/cte"
)

type Computer struct{}

func (c Computer) Compute(ctx context.Context, p cte.MasterPlan) (interface{}, error) {
	casted := p.(inout)

	return c.calculateTravelCost(casted), nil
}

func (Computer) calculateTravelCost(in Input) float64 {
	return in.GetBaseCost() + in.GetTravelDuration()*in.GetCostPerMinute() + in.GetTravelDistance()*in.GetCostPerKilometer()
}
