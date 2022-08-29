package main

import (
	"context"
	"fmt"
	"github.com/jamestrandung/go-cte-117/sample/dto"
	"github.com/jamestrandung/go-cte-117/sample/service/scaffolding/endpoint"

	"github.com/jamestrandung/go-cte-117/sample/config"
	"github.com/jamestrandung/go-cte-117/sample/server"
)

func main() {
	err := config.Engine.VerifyConfigurations()
	fmt.Println("Engine configuration error:", err)

	testEngine()
}

func testEngine() {
	server.Serve()

	p := endpoint.NewPlan(
		dto.CostRequest{
			PointA: "Clementi",
			PointB: "Changi Airport",
		},
		server.Dependencies,
	)

	if err := p.Execute(context.Background()); err != nil {
		fmt.Println(err)
	}

	//config.Print(p.GetTravelCost())
	config.Print(p.GetTotalCost())
	config.Print(p.GetVATAmount())
}
