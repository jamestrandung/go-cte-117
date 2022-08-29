package vat

import (
	"github.com/jamestrandung/go-cte-117/cte"
)

type inout interface {
	Input
	Output
}

type Input interface {
	GetVATPercent() float64
	GetTotalCost() float64
}

type Output interface {
	SetTotalCost(float64)
}

type VATAmount cte.SyncResult

func (a VATAmount) CTEMetadata() interface{} {
	return struct {
		computer computer
		inout    inout
	}{}
}

func (a VATAmount) GetVATAmount() float64 {
	if a.Outcome == nil {
		return 0
	}

	return a.Outcome.(float64)
}
