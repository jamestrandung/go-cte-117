package cte

import (
	"context"
)

type Plan interface {
	IsSequentialCTEPlan() bool
}

type MasterPlan interface {
	Plan
	Execute(ctx context.Context) error
}

type Pre interface {
	PreExecute(p Plan) error
}

type Post interface {
	PostExecute(p Plan) error
}
