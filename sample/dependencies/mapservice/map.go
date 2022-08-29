package mapservice

import (
	"github.com/jamestrandung/go-cte-117/sample/config"
	"github.com/stretchr/testify/assert"
)

type Route struct {
	Distance float64
	Duration float64
}

type Service struct{}

func (Service) GetRoute(pointA, pointB string) (Route, error) {
	config.Printf("Building route from %s to %s using real map\n", pointA, pointB)
	return Route{
		Distance: 2,
		Duration: 3,
	}, assert.AnError
}
