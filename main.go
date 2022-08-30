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

	server.Serve()

	testEngine()
	testPlainGo()
}

func testEngine() {
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

func testPlainGo() {
	quote, err := server.Handler.Handle(context.Background(), dto.CostRequest{
		PointA: "Clementi",
		PointB: "Changi Airport",
	})

	if err != nil {
		fmt.Println("Plain Go error:", err)
	}

	config.Print(quote.TotalCost)
	config.Print(quote.VATAmount)
}
