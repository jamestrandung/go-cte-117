package streaming

import "github.com/jamestrandung/go-cte-117/cte"

type inout interface {
	Input
}

type Input interface {
	GetTotalCost() float64
}

type CostStreaming cte.SideEffect

func (c CostStreaming) CTEMetadata() interface{} {
	return struct {
		computer computer
		inout    inout
	}{}
}
