package endpoint

import (
	"context"
	"github.com/jamestrandung/go-cte-117/sample/dependencies/configsfetcher"
	"github.com/jamestrandung/go-cte-117/sample/dependencies/mapservice"
	"github.com/jamestrandung/go-cte-117/sample/dto"
	"github.com/jamestrandung/go-cte-117/sample/service-no-framework/helpers/future"
	"github.com/jamestrandung/go-cte-117/sample/service-no-framework/modules/costconfigs"
	"github.com/jamestrandung/go-cte-117/sample/service-no-framework/modules/platformfee"
	"github.com/jamestrandung/go-cte-117/sample/service-no-framework/modules/streaming"
	"github.com/jamestrandung/go-cte-117/sample/service-no-framework/modules/travelcost"
	"github.com/jamestrandung/go-cte-117/sample/service-no-framework/modules/travelplan"
	"github.com/jamestrandung/go-cte-117/sample/service-no-framework/modules/vat"
)

type Handler struct {
	costConfigsDAL configsfetcher.Fetcher
	routeDAL       mapservice.Service
}

func NewHandler(costConfigsDAL configsfetcher.Fetcher, routeDAL mapservice.Service) *Handler {
	return &Handler{
		costConfigsDAL: costConfigsDAL,
		routeDAL:       routeDAL,
	}
}

func (h *Handler) Handle(ctx context.Context, request dto.CostRequest) (*dto.Quote, error) {
	costConfigsLoader, costConfigsGetter := costconfigs.GenerateCostConfigsLoader(ctx, h.costConfigsDAL)
	travelRouteLoader, travelRouteGetter := travelplan.GenerateRouteLoader(ctx, h.routeDAL, request.PointA, request.PointB)
	travelCostLoader, travelCostGetter := travelcost.GenerateTravelCostLoader(ctx, costConfigsGetter, travelRouteGetter)

	future.Load(costConfigsLoader, travelRouteLoader, travelCostLoader)

	costConfigs, err := costConfigsGetter()
	if err != nil {
		return nil, err
	}

	quote, err := func() (*dto.Quote, error) {
		if costConfigs.IsFixedCostEnabled {
			return h.calculateQuoteWithFixedCost(costConfigs), nil
		}

		return h.calculateQuoteUsingStandardRoute(costConfigs, travelCostGetter)
	}()

	if err != nil {
		return nil, err
	}

	streaming.StreamQuote(quote)

	return quote, nil
}

func (h *Handler) calculateQuoteWithFixedCost(costConfigs configsfetcher.MergedCostConfigs) *dto.Quote {
	quote := &dto.Quote{
		TotalCost: costConfigs.FixedCost,
	}

	vat.AddVAT(quote, costConfigs)

	return quote
}

func (h *Handler) calculateQuoteUsingStandardRoute(costConfigs configsfetcher.MergedCostConfigs, travelCostGetter travelcost.GetterFunc) (*dto.Quote, error) {
	travelCost, err := travelCostGetter()
	if err != nil {
		return nil, err
	}

	quote := &dto.Quote{
		TotalCost: travelCost,
	}

	platformfee.AddPlatformFee(quote, costConfigs)
	vat.AddVAT(quote, costConfigs)

	return quote, nil
}
