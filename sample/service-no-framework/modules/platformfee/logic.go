package platformfee

import (
	"github.com/jamestrandung/go-cte-117/sample/dependencies/configsfetcher"
	"github.com/jamestrandung/go-cte-117/sample/dto"
)

func AddPlatformFee(quote *dto.Quote, costConfigs configsfetcher.MergedCostConfigs) {
	quote.TotalCost = quote.TotalCost + costConfigs.PlatformFee
}
