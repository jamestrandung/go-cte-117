package costconfigs

import (
	"github.com/jamestrandung/go-cte-117/cte"
	"github.com/jamestrandung/go-cte-117/sample/dependencies/configsfetcher"
)

type inout interface {
	Dependencies
}

type Dependencies interface {
	GetConfigsFetcher() configsfetcher.Fetcher
}

type CostConfigs cte.Result

func (c CostConfigs) CTEMetadata() interface{} {
	return struct {
		computer computer
		inout    inout
	}{}
}

func (c CostConfigs) GetBaseCost() float64 {
	outcome, err := c.Task.Outcome()
	if err != nil {
		return 0
	}

	return outcome.(configsfetcher.MergedCostConfigs).BaseCost
}

func (c CostConfigs) GetCostPerKilometer() float64 {
	outcome, err := c.Task.Outcome()
	if err != nil {
		return 0
	}

	return outcome.(configsfetcher.MergedCostConfigs).CostPerKilometer
}

func (c CostConfigs) GetCostPerMinute() float64 {
	outcome, err := c.Task.Outcome()
	if err != nil {
		return 0
	}

	return outcome.(configsfetcher.MergedCostConfigs).CostPerMinute
}

func (c CostConfigs) GetPlatformFee() float64 {
	outcome, err := c.Task.Outcome()
	if err != nil {
		return 0
	}

	return outcome.(configsfetcher.MergedCostConfigs).PlatformFee
}

func (c CostConfigs) GetVATPercent() float64 {
	outcome, err := c.Task.Outcome()
	if err != nil {
		return 0
	}

	return outcome.(configsfetcher.MergedCostConfigs).VATPercent
}

func (c CostConfigs) GetIsFixedCostEnabled() bool {
	outcome, err := c.Task.Outcome()
	if err != nil {
		return false
	}

	return outcome.(configsfetcher.MergedCostConfigs).IsFixedCostEnabled
}

func (c CostConfigs) GetFixedCost() float64 {
	outcome, err := c.Task.Outcome()
	if err != nil {
		return 0
	}

	return outcome.(configsfetcher.MergedCostConfigs).FixedCost
}
