package streaming

import (
	"github.com/jamestrandung/go-cte-117/sample/config"
	"github.com/jamestrandung/go-cte-117/sample/dto"
)

func StreamQuote(quote *dto.Quote) {
	config.Print("Streaming calculated cost:", quote.TotalCost)
}
