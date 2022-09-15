package cte

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"reflect"
	"strings"
	"testing"
)

func TestVerifyComponentCompleteness(t *testing.T) {
	scenarios := []struct {
		desc string
		test func(t *testing.T)
	}{
		{
			desc: "ErrInoutMetaMissing",
			test: func(t *testing.T) {
				pm := parsedMetadata{}
				sd := newStructDisassembler()
				cs := componentStack{}
				rootPlanName := "rootPlanName"
				componentID := "componentID"
				pType := "planType"

				v := &completenessValidator{
					rootPlanName: rootPlanName,
					sd:           sd,
				}

				err := v.verifyComponentCompleteness(pm, cs, componentID, pType)
				_, ok := pm.getInoutInterface()
				assert.False(t, ok)
				assert.Equal(t, ErrInoutMetaMissing.Err(componentID), err)
			},
		},
		{
			desc: "isInterfaceSatisfied return non-nil",
			test: func(t *testing.T) {
				pm := parsedMetadata{}
				pm[metaTypeInout] = reflect.TypeOf("dummy")

				sd := newStructDisassembler()
				cs := componentStack{}
				rootPlanName := "rootPlanName"
				componentID := "componentID"
				pType := "planType"
				vMock := &mockICompletenessValidator{}

				v := &completenessValidator{
					itself:       vMock,
					rootPlanName: rootPlanName,
					sd:           sd,
				}

				vMock.On("isInterfaceSatisfied", reflect.TypeOf("dummy")).
					Return(assert.AnError).
					Once()

				err := v.verifyComponentCompleteness(pm, cs, componentID, pType)
				expectedInout, ok := pm.getInoutInterface()
				assert.True(t, ok)
				assert.Equal(t, ErrPlanNotMeetingInoutRequirements.Err(pType, expectedInout, assert.AnError.Error(), cs.push("componentID")), err)
				mock.AssertExpectationsForObjects(t, vMock)
			},
		},
		{
			desc: "isInterfaceSatisfied return nil",
			test: func(t *testing.T) {
				pm := parsedMetadata{}
				pm[metaTypeInout] = reflect.TypeOf("dummy")

				sd := newStructDisassembler()
				cs := componentStack{}
				rootPlanName := "rootPlanName"
				componentID := "componentID"
				pType := "planType"
				vMock := &mockICompletenessValidator{}

				v := &completenessValidator{
					itself:       vMock,
					rootPlanName: rootPlanName,
					sd:           sd,
				}

				vMock.On("isInterfaceSatisfied", reflect.TypeOf("dummy")).
					Return(nil).
					Once()

				err := v.verifyComponentCompleteness(pm, cs, componentID, pType)
				_, ok := pm.getInoutInterface()
				assert.True(t, ok)
				assert.Equal(t, nil, err)
				mock.AssertExpectationsForObjects(t, vMock)
			},
		},
	}

	for _, scenario := range scenarios {
		s := scenario

		t.Run(s.desc, s.test)
	}
}

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

				v := &completenessValidator{
					rootPlanName: rootPlanName,
					sd:           sd,
				}

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

				err := v.isInterfaceSatisfied(expectedInterfaceType)
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

				v := &completenessValidator{
					rootPlanName: rootPlanName,
					sd:           sd,
				}

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

				err := v.isInterfaceSatisfied(expectedInterfaceType)
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

				v := &completenessValidator{
					rootPlanName: rootPlanName,
					sd:           sd,
				}

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

				err := v.isInterfaceSatisfied(expectedInterfaceType)
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

				v := &completenessValidator{
					rootPlanName: rootPlanName,
					sd:           sd,
				}

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

				err := v.isInterfaceSatisfied(expectedInterfaceType)
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

				v := &completenessValidator{
					rootPlanName: rootPlanName,
					sd:           sd,
				}

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

				err := v.isInterfaceSatisfied(expectedInterfaceType)
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
