package endpoint

import (
	"github.com/jamestrandung/go-cte-117/sample/service/scaffolding/loading"
)

type Dependencies interface {
	loading.Dependencies
}

type Input interface {
	Request
}

type Request interface {
	GetPointA() string
	GetPointB() string
}
