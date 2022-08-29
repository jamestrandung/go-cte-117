package quote

import (
    "context"
    "github.com/jamestrandung/go-cte-117/cte"
    "github.com/jamestrandung/go-cte-117/sample/service/scaffolding/calculation"
    "github.com/jamestrandung/go-cte-117/sample/service/scaffolding/fixedcost"
)

type computer struct{}

// TODO: Due to pre execution can return nil, clients must take care of handling nil plan in getters
func (c computer) Switch(ctx context.Context, p cte.MasterPlan) (cte.MasterPlan, error) {
    casted := p.(inout)

    if casted.GetIsFixedCostEnabled() {
        return fixedcost.NewPlan(casted), nil
    }

    return calculation.NewPlan(casted), nil
}
