package endpoint

import "github.com/jamestrandung/go-cte-117/sample/config"

func init() {
	config.Engine.AnalyzePlan(&SequentialPlan{})
}
