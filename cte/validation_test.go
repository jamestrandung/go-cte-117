package cte

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"reflect"
	"strings"
	"testing"
)

type isInterfaceSatisfied_interface interface {
	Do(int) string
}

func TestIsInterfaceSatisfied(t *testing.T) {
	defer func(original func(rm reflect.Method, ignoreFirstReceiverArgument bool) method) {
		extractMethodDetails = original
	}(extractMethodDetails)

	expectedInterfaceType := reflect.TypeOf((*isInterfaceSatisfied_interface)(nil)).Elem()

	scenarios := []struct {
		desc string
		test func(t *testing.T)
	}{
		{
			desc: "ErrPlanMissingMethod",
			test: func(t *testing.T) {
				rootPlanName := "rootPlanName"

				sd := newStructDisassembler()
				sdMock := &mockIStructDisassembler{}
				sd.itself = sdMock

				expectedMethod := method{
					name:    "Do",
					outputs: "string",
				}

				extractMethodDetails = func(rm reflect.Method, ignoreFirstReceiverArgument bool) method {
					assert.Equal(t, "Do", rm.Name)
					assert.Equal(t, "func(int) string", rm.Type.String())
					assert.Equal(t, false, ignoreFirstReceiverArgument)

					return expectedMethod
				}

				sdMock.On("findAvailableMethods", expectedMethod.name).
					Return(nil, false).
					Once()

				err := isInterfaceSatisfied(sd, expectedInterfaceType, rootPlanName)
				assert.Equal(t, ErrPlanMissingMethod.Err(expectedMethod), err)
				mock.AssertExpectationsForObjects(t, sdMock)
			},
		},
		{
			desc: "ErrPlanHavingAmbiguousMethods",
			test: func(t *testing.T) {
				rootPlanName := "rootPlanName"

				sd := newStructDisassembler()
				sdMock := &mockIStructDisassembler{}
				sd.itself = sdMock

				expectedMethod := method{
					name:      "Do",
					arguments: "int",
					outputs:   "string",
				}

				extractMethodDetails = func(rm reflect.Method, ignoreFirstReceiverArgument bool) method {
					assert.Equal(t, "Do", rm.Name)
					assert.Equal(t, "func(int) string", rm.Type.String())
					assert.Equal(t, false, ignoreFirstReceiverArgument)

					return expectedMethod
				}

				expectedMethodSet := methodSet{
					method{
						owningType: "owningType1",
						name:       "Do",
						arguments:  "int",
						outputs:    "string",
					}: struct{}{},
					method{
						owningType: "owningType2",
						name:       "Do",
						arguments:  "int",
						outputs:    "string",
					}: struct{}{},
				}

				sdMock.On("findAvailableMethods", expectedMethod.name).
					Return(expectedMethodSet, true).
					Once()

				expectedLocations := []string{"location1", "location2"}

				sdMock.On("findMethodLocations", mock.Anything, rootPlanName).
					Return(expectedLocations).
					Once()

				err := isInterfaceSatisfied(sd, expectedInterfaceType, rootPlanName)
				assert.Equal(t, ErrPlanHavingAmbiguousMethods.Err(expectedMethod, expectedMethodSet, strings.Join(expectedLocations, "; ")), err)
				mock.AssertExpectationsForObjects(t, sdMock)
			},
		},
		{
			desc: "ErrPlanHavingMethodButSignatureMismatched",
			test: func(t *testing.T) {
				rootPlanName := "rootPlanName"

				sd := newStructDisassembler()
				sdMock := &mockIStructDisassembler{}
				sd.itself = sdMock

				expectedMethod := method{
					name:    "Do",
					outputs: "string",
				}

				extractMethodDetails = func(rm reflect.Method, ignoreFirstReceiverArgument bool) method {
					assert.Equal(t, "Do", rm.Name)
					assert.Equal(t, "func(int) string", rm.Type.String())
					assert.Equal(t, false, ignoreFirstReceiverArgument)

					return expectedMethod
				}

				methodWithMismatchedSignature := method{
					owningType: "owningType",
					name:       "Do",
					arguments:  "float64",
					outputs:    "string",
				}

				sdMock.On("findAvailableMethods", expectedMethod.name).
					Return(methodSet{methodWithMismatchedSignature: struct{}{}}, true).
					Once()

				err := isInterfaceSatisfied(sd, expectedInterfaceType, rootPlanName)
				assert.Equal(t, ErrPlanHavingMethodButSignatureMismatched.Err(expectedMethod, methodWithMismatchedSignature), err)
				mock.AssertExpectationsForObjects(t, sdMock)
			},
		},
		{
			desc: "ErrPlanHavingSameMethodRegisteredMoreThanOnce",
			test: func(t *testing.T) {
				rootPlanName := "rootPlanName"

				sd := newStructDisassembler()
				sdMock := &mockIStructDisassembler{}
				sd.itself = sdMock

				expectedMethod := method{
					name:      "Do",
					arguments: "int",
					outputs:   "string",
				}

				extractMethodDetails = func(rm reflect.Method, ignoreFirstReceiverArgument bool) method {
					assert.Equal(t, "Do", rm.Name)
					assert.Equal(t, "func(int) string", rm.Type.String())
					assert.Equal(t, false, ignoreFirstReceiverArgument)

					return expectedMethod
				}

				duplicateMethod := method{
					owningType: "owningType",
					name:       "Do",
					arguments:  "int",
					outputs:    "string",
				}

				sdMock.On("findAvailableMethods", expectedMethod.name).
					Return(methodSet{duplicateMethod: struct{}{}}, true).
					Once()

				sdMock.On("isAvailableMoreThanOnce", duplicateMethod).
					Return(true).
					Once()

				expectedLocations := []string{"location1", "location2"}

				sdMock.On("findMethodLocations", mock.Anything, rootPlanName).
					Return(expectedLocations).
					Once()

				err := isInterfaceSatisfied(sd, expectedInterfaceType, rootPlanName)
				assert.Equal(t, ErrPlanHavingSameMethodRegisteredMoreThanOnce.Err(duplicateMethod, strings.Join(expectedLocations, "; ")), err)
				mock.AssertExpectationsForObjects(t, sdMock)
			},
		},
		{
			desc: "Happy",
			test: func(t *testing.T) {
				rootPlanName := "rootPlanName"

				sd := newStructDisassembler()
				sdMock := &mockIStructDisassembler{}
				sd.itself = sdMock

				expectedMethod := method{
					name:      "Do",
					arguments: "int",
					outputs:   "string",
				}

				extractMethodDetails = func(rm reflect.Method, ignoreFirstReceiverArgument bool) method {
					assert.Equal(t, "Do", rm.Name)
					assert.Equal(t, "func(int) string", rm.Type.String())
					assert.Equal(t, false, ignoreFirstReceiverArgument)

					return expectedMethod
				}

				matchingMethod := method{
					owningType: "owningType",
					name:       "Do",
					arguments:  "int",
					outputs:    "string",
				}

				sdMock.On("findAvailableMethods", expectedMethod.name).
					Return(methodSet{matchingMethod: struct{}{}}, true).
					Once()

				sdMock.On("isAvailableMoreThanOnce", matchingMethod).
					Return(false).
					Once()

				err := isInterfaceSatisfied(sd, expectedInterfaceType, rootPlanName)
				assert.Equal(t, nil, err)
				mock.AssertExpectationsForObjects(t, sdMock)
			},
		},
	}

	for _, scenario := range scenarios {
		s := scenario

		t.Run(s.desc, s.test)
	}
}

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
