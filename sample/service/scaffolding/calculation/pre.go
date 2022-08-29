package calculation

import (
	"github.com/jamestrandung/go-cte-117/cte"
	"github.com/jamestrandung/go-cte-117/sample/config"
)

type pre interface {
	preIn
	preOut
}

type preIn interface {
	GetTravelCost() float64
}

type preOut interface {
	SetTotalCost(float64)
}

type preHook struct{}

func (preHook) CTEMetadata() interface{} {
	return struct {
		inout pre
	}{}
}

func (preHook) PreExecute(p cte.Plan) error {
	config.Print("Before executing sequential plan")
	casted := p.(pre)

	casted.SetTotalCost(casted.GetTravelCost())

	return nil
}
