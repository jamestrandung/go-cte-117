package calculation

import (
	"github.com/jamestrandung/go-cte-117/sample/service/components/platformfee"
	"github.com/jamestrandung/go-cte-117/sample/service/components/vat"
)

type Input interface {
	preIn
	vat.Input
	platformfee.Input
}
