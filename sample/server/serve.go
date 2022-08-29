package server

import (
	"github.com/jamestrandung/go-cte-117/sample/dependencies/configsfetcher"
	"github.com/jamestrandung/go-cte-117/sample/dependencies/mapservice"
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

func Serve() {
	Dependencies = dependencies{
		configsFetcher: configsfetcher.Fetcher{},
		mapService:     mapservice.Service{},
	}
}
