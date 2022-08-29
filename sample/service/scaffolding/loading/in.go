package loading

import (
	"github.com/jamestrandung/go-cte-117/sample/service/components/costconfigs"
	"github.com/jamestrandung/go-cte-117/sample/service/components/travelplan"
)

type Dependencies interface {
	costconfigs.Dependencies
	travelplan.Dependencies
}
