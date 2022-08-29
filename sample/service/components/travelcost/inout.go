package travelcost

import (
	"github.com/jamestrandung/go-cte-117/cte"
)

type inout interface {
	Input
}

type Input interface {
	GetBaseCost() float64
	GetTravelDistance() float64
	GetTravelDuration() float64
	GetCostPerKilometer() float64
	GetCostPerMinute() float64
}

type TravelCost cte.Result

func (r TravelCost) CTEMetadata() interface{} {
	return struct {
		computer Computer
		inout    inout
	}{}
}

func (r TravelCost) GetTravelCost() float64 {
	outcome, err := r.Task.Outcome()
	if err != nil {
		return 0
	}

	return outcome.(float64)
}
