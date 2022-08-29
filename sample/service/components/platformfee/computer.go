package platformfee

import (
	"context"

	"github.com/jamestrandung/go-cte-117/cte"
)

type computer struct{}

func (c computer) Compute(ctx context.Context, p cte.MasterPlan) error {
	casted := p.(inout)

	c.addPlatformFee(casted)

	return nil
}

func (computer) addPlatformFee(p inout) {
	p.SetTotalCost(p.GetTotalCost() + p.GetPlatformFee())
}
