package costconfigs

import (
	"context"

	"github.com/jamestrandung/go-cte-117/cte"

	"github.com/jamestrandung/go-cte-117/sample/dependencies/configsfetcher"
)

type computer struct{}

func (c computer) Compute(ctx context.Context, p cte.MasterPlan) (interface{}, error) {
	casted := p.(inout)

	return c.doFetch(casted), nil
}

func (c computer) doFetch(p inout) configsfetcher.MergedCostConfigs {
	return p.GetConfigsFetcher().Fetch()
}
