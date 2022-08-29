package config

import (
	"fmt"

	"github.com/jamestrandung/go-cte-117/cte"
)

var Engine = &CostEngine{
	Engine: cte.NewEngine(),
}

var printDebugLog = false

func Print(values ...any) {
	if printDebugLog {
		fmt.Println(values...)
	}
}

func Printf(format string, values ...any) {
	if printDebugLog {
		fmt.Printf(format, values...)
	}
}

type CostEngine struct {
	cte.Engine
	// Add common utilities like logger, statsD, UCM client, etc.
	// for all component codes to share.
}
