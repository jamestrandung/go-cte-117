package cte

import (
	"context"
)

//go:generate mockery --name Plan --case=underscore --inpackage
type Plan interface {
	IsSequentialCTEPlan() bool
}

//go:generate mockery --name MasterPlan --case=underscore --inpackage
type MasterPlan interface {
	Plan
	Execute(ctx context.Context) error
}

//go:generate mockery --name Pre --case=underscore --inpackage
type Pre interface {
	PreExecute(p Plan) error
}

//go:generate mockery --name Post --case=underscore --inpackage
type Post interface {
	PostExecute(p Plan) error
}
