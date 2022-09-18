package cte

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatErr_Err(t *testing.T) {
	fe := formatErr{
		format: "%v and %v is an error",
	}

	assert.Equal(t, "THIS and THAT is an error", fe.Err("THIS", "THAT").Error())
}
