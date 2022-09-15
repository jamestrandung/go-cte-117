package cte

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestComponentStack_Push(t *testing.T) {
	cs := componentStack{}
	assert.Equal(t, 0, len(cs))

	cs = cs.push("dummy")
	assert.Equal(t, 1, len(cs))
	assert.Equal(t, "dummy", cs[0])
}

func TestComponentStack_Pop(t *testing.T) {
	cs := componentStack{}
	cs = cs.push("dummy1")
	cs = cs.push("dummy2")
	cs = cs.push("dummy3")
	assert.Equal(t, 3, len(cs))

	cs = cs.pop()
	assert.Equal(t, 2, len(cs))
	assert.Equal(t, "dummy1", cs[0])
	assert.Equal(t, "dummy2", cs[1])
}

func TestComponentStack_Clone(t *testing.T) {
	cs := componentStack{}
	cs = cs.push("dummy1")
	cs = cs.push("dummy2")
	assert.Equal(t, 2, len(cs))

	// csClone should be exactly the same as cs
	csClone := cs.clone()
	assert.Equal(t, 2, len(csClone))
	assert.Equal(t, "dummy1", csClone[0])
	assert.Equal(t, "dummy2", csClone[1])

	// Changes to cs must not affect csClone
	cs = cs.pop()
	cs = cs.push("dummy3")
	assert.Equal(t, 2, len(cs))
	assert.Equal(t, "dummy1", cs[0])
	assert.Equal(t, "dummy3", cs[1])
	assert.Equal(t, 2, len(csClone))
	assert.Equal(t, "dummy1", csClone[0])
	assert.Equal(t, "dummy2", csClone[1])
}

func TestComponentStack_String(t *testing.T) {
	cs := componentStack{}
	cs = cs.push("dummy1")
	assert.Equal(t, "dummy1", cs.String())

	cs = cs.push("dummy2")
	assert.Equal(t, "dummy1 >> dummy2", cs.String())
}
