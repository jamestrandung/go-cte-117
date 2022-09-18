package cte

import (
	"context"
)

type Plan interface {
	IsSequentialCTEPlan() bool
}

//go:generate mockery --name MasterPlan --case=underscore --inpackage
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
