package vat

import (
	"github.com/jamestrandung/go-cte-117/sample/dependencies/configsfetcher"
	"github.com/jamestrandung/go-cte-117/sample/dto"
)

func AddVAT(quote *dto.Quote, costConfigs configsfetcher.MergedCostConfigs) {
	vatAmount := quote.TotalCost * costConfigs.VATPercent / 100

	quote.VATAmount = vatAmount
	quote.TotalCost = quote.TotalCost + vatAmount
}
