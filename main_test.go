package main

import (
	"context"
	"testing"

	"github.com/jamestrandung/go-cte-117/sample/service/scaffolding/endpoint"

	"github.com/jamestrandung/go-cte-117/sample/dto"

	"github.com/jamestrandung/go-cte-117/sample/config"
	"github.com/jamestrandung/go-cte-117/sample/server"
)

func BenchmarkEngine(b *testing.B) {
	server.Serve()

	p := endpoint.NewPlan(
		dto.CostRequest{
			PointA: "Clementi",
			PointB: "Changi Airport",
		},
		server.Dependencies,
	)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if err := p.Execute(context.Background()); err != nil {
			config.Print(err)
		}
	}
}

func BenchmarkPlainGo(b *testing.B) {
	server.Serve()

	request := dto.CostRequest{
		PointA: "Clementi",
		PointB: "Changi Airport",
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if _, err := server.Handler.Handle(context.Background(), request); err != nil {
			config.Print(err)
		}
	}
}
