package travelplan

import (
	"github.com/jamestrandung/go-cte-117/cte"
	"github.com/jamestrandung/go-cte-117/sample/dependencies/mapservice"
)

type inout interface {
	Input
}

type Dependencies interface {
	GetMapService() mapservice.Service
}

type Input interface {
	Dependencies
	GetPointA() string
	GetPointB() string
}

type TravelPlan cte.Result

func (p TravelPlan) CTEMetadata() interface{} {
	return struct {
		computer computer
		inout    inout
	}{}
}

func (p TravelPlan) GetTravelDistance() float64 {
	outcome, err := p.Task.Outcome()
	if err != nil {
		return 0
	}

	return outcome.(mapservice.Route).Distance
}

func (p TravelPlan) GetTravelDuration() float64 {
	outcome, err := p.Task.Outcome()
	if err != nil {
		return 0
	}

	return outcome.(mapservice.Route).Duration
}
