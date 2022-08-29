package configsfetcher

type MergedCostConfigs struct {
	BaseCost           float64
	CostPerKilometer   float64
	CostPerMinute      float64
	PlatformFee        float64
	VATPercent         float64
	IsFixedCostEnabled bool
	FixedCost          float64
}

type Fetcher struct{}

func (Fetcher) Fetch() MergedCostConfigs {
	return MergedCostConfigs{
		BaseCost:           1,
		CostPerKilometer:   4,
		CostPerMinute:      5,
		PlatformFee:        8,
		VATPercent:         10,
		IsFixedCostEnabled: false,
		FixedCost:          10,
	}
}
