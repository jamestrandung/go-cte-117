package server

import (
	"github.com/jamestrandung/go-cte-117/sample/dependencies/configsfetcher"
	"github.com/jamestrandung/go-cte-117/sample/dependencies/mapservice"
	"github.com/jamestrandung/go-cte-117/sample/service-no-framework/endpoint"
)

type dependencies struct {
	configsFetcher configsfetcher.Fetcher
	mapService     mapservice.Service
}

func (d dependencies) GetConfigsFetcher() configsfetcher.Fetcher {
	return d.configsFetcher
}

func (d dependencies) GetMapService() mapservice.Service {
	return d.mapService
}

var Dependencies dependencies
var Handler *endpoint.Handler

func Serve() {
	Dependencies = dependencies{
		configsFetcher: configsfetcher.Fetcher{},
		mapService:     mapservice.Service{},
	}

	Handler = endpoint.NewHandler(configsfetcher.Fetcher{}, mapservice.Service{})
}
