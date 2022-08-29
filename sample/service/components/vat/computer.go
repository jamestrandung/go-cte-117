package vat

import (
	"context"

	"github.com/jamestrandung/go-cte-117/cte"
)

type computer struct{}

func (c computer) Compute(ctx context.Context, p cte.MasterPlan) (interface{}, error) {
	casted := p.(inout)

	vatAmount := casted.GetTotalCost() * casted.GetVATPercent() / 100
	casted.SetTotalCost(casted.GetTotalCost() + vatAmount)

	return vatAmount, nil
}
