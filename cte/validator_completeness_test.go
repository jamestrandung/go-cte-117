package cte

import (
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCompletenessValidator_Validate(t *testing.T) {
	vMock := &mockICompletenessValidator{}

	v := completenessValidator{
		itself:       vMock,
		planValue:    reflect.ValueOf("dummy"),
		rootPlanName: "rootPlanName",
	}

	t.Run(
		"doValidate return no errors", func(t *testing.T) {
			vMock.On("doValidate", mock.Anything, mock.Anything, mock.Anything).
				Return(
					func(planName string, cs componentStack, curPlanValue reflect.Value) error {
						assert.Equal(t, v.rootPlanName, planName)
						assert.Equal(t, v.planValue, curPlanValue)
						assert.Equal(t, 0, len(cs))

						return nil
					},
				).
				Once()

			err := v.validate()
			assert.Nil(t, err)
			mock.AssertExpectationsForObjects(t, vMock)
		},
	)

	t.Run(
		"doValidate return an error", func(t *testing.T) {
			vMock.On("doValidate", mock.Anything, mock.Anything, mock.Anything).
				Return(
					func(planName string, cs componentStack, curPlanValue reflect.Value) error {
						assert.Equal(t, v.rootPlanName, planName)
						assert.Equal(t, v.planValue, curPlanValue)
						assert.Equal(t, 0, len(cs))

						return assert.AnError
					},
				).
				Once()

			err := v.validate()
			assert.Equal(t, assert.AnError, err)
			mock.AssertExpectationsForObjects(t, vMock)
		},
	)
}

type doValidate_Hook struct{}

func (doValidate_Hook) PreExecute(p Plan) error {
	return nil
}

func (doValidate_Hook) PostExecute(p Plan) error {
	return nil
}

type doValidate_NestedPlan struct{}

type doValidate_ParentPlan struct {
	doValidate_NestedPlan
}

func TestCompletenessValidator_DoValidate(t *testing.T) {
	scenarios := []struct {
		desc string
		test func(t *testing.T)
	}{
		{
			desc: "pre hook should be processed, returning no errors",
			test: func(t *testing.T) {
				planName := "planName"
				cs := componentStack{}
				pValue := reflect.ValueOf(dummy{})

				eMock := &mockIEngine{}
				vMock := &mockICompletenessValidator{}

				v := completenessValidator{
					itself:    vMock,
					engine:    eMock,
					planValue: reflect.ValueOf("dummy"),
				}

				dummyHook := preHook{
					hook:     doValidate_Hook{},
					metadata: parsedMetadata{},
				}

				eMock.On("findAnalyzedPlan", planName, pValue).
					Return(
						analyzedPlan{
							preHooks: []preHook{dummyHook},
						},
					).
					Once()

				vMock.On("verifyComponentCompleteness", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					Return(
						func(pm parsedMetadata, cs componentStack, componentID string, planType string) error {
							assert.Equal(t, dummyHook.metadata, pm)
							assert.Equal(t, reflect.TypeOf(dummyHook.hook).String(), componentID)
							assert.Equal(t, v.planValue.Type().String(), planType)
							if assert.Equal(t, 1, len(cs)) {
								assert.Equal(t, planName, cs[0])
							}

							return nil
						},
					).
					Once()

				err := v.doValidate(planName, cs, pValue)
				assert.Nil(t, err)
				mock.AssertExpectationsForObjects(t, eMock, vMock)
			},
		},
		{
			desc: "pre hook should be processed, returning an error",
			test: func(t *testing.T) {
				planName := "planName"
				cs := componentStack{}
				pValue := reflect.ValueOf(dummy{})

				eMock := &mockIEngine{}
				vMock := &mockICompletenessValidator{}

				v := completenessValidator{
					itself:    vMock,
					engine:    eMock,
					planValue: reflect.ValueOf("dummy"),
				}

				dummyHook := preHook{
					hook:     doValidate_Hook{},
					metadata: parsedMetadata{},
				}

				eMock.On("findAnalyzedPlan", planName, pValue).
					Return(
						analyzedPlan{
							preHooks: []preHook{dummyHook},
						},
					).
					Once()

				vMock.On("verifyComponentCompleteness", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					Return(
						func(pm parsedMetadata, cs componentStack, componentID string, planType string) error {
							assert.Equal(t, dummyHook.metadata, pm)
							assert.Equal(t, reflect.TypeOf(dummyHook.hook).String(), componentID)
							assert.Equal(t, v.planValue.Type().String(), planType)
							if assert.Equal(t, 1, len(cs)) {
								assert.Equal(t, planName, cs[0])
							}

							return assert.AnError
						},
					).
					Once()

				err := v.doValidate(planName, cs, pValue)
				assert.Equal(t, assert.AnError, err)
				mock.AssertExpectationsForObjects(t, eMock, vMock)
			},
		},
		{
			desc: "post hook should be processed, returning no errors",
			test: func(t *testing.T) {
				planName := "planName"
				cs := componentStack{}
				pValue := reflect.ValueOf(dummy{})

				eMock := &mockIEngine{}
				vMock := &mockICompletenessValidator{}

				v := completenessValidator{
					itself:    vMock,
					engine:    eMock,
					planValue: reflect.ValueOf("dummy"),
				}

				dummyHook := postHook{
					hook:     doValidate_Hook{},
					metadata: parsedMetadata{},
				}

				eMock.On("findAnalyzedPlan", planName, pValue).
					Return(
						analyzedPlan{
							postHooks: []postHook{dummyHook},
						},
					).
					Once()

				vMock.On("verifyComponentCompleteness", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					Return(
						func(pm parsedMetadata, cs componentStack, componentID string, planType string) error {
							assert.Equal(t, dummyHook.metadata, pm)
							assert.Equal(t, reflect.TypeOf(dummyHook.hook).String(), componentID)
							assert.Equal(t, v.planValue.Type().String(), planType)
							if assert.Equal(t, 1, len(cs)) {
								assert.Equal(t, planName, cs[0])
							}

							return nil
						},
					).
					Once()

				err := v.doValidate(planName, cs, pValue)
				assert.Nil(t, err)
				mock.AssertExpectationsForObjects(t, eMock, vMock)
			},
		},
		{
			desc: "post hook should be processed, returning an error",
			test: func(t *testing.T) {
				planName := "planName"
				cs := componentStack{}
				pValue := reflect.ValueOf(dummy{})

				eMock := &mockIEngine{}
				vMock := &mockICompletenessValidator{}

				v := completenessValidator{
					itself:    vMock,
					engine:    eMock,
					planValue: reflect.ValueOf("dummy"),
				}

				dummyHook := postHook{
					hook:     doValidate_Hook{},
					metadata: parsedMetadata{},
				}

				eMock.On("findAnalyzedPlan", planName, pValue).
					Return(
						analyzedPlan{
							postHooks: []postHook{dummyHook},
						},
					).
					Once()

				vMock.On("verifyComponentCompleteness", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					Return(
						func(pm parsedMetadata, cs componentStack, componentID string, planType string) error {
							assert.Equal(t, dummyHook.metadata, pm)
							assert.Equal(t, reflect.TypeOf(dummyHook.hook).String(), componentID)
							assert.Equal(t, v.planValue.Type().String(), planType)
							if assert.Equal(t, 1, len(cs)) {
								assert.Equal(t, planName, cs[0])
							}

							return assert.AnError
						},
					).
					Once()

				err := v.doValidate(planName, cs, pValue)
				assert.Equal(t, assert.AnError, err)
				mock.AssertExpectationsForObjects(t, eMock, vMock)
			},
		},
		{
			desc: "computer should be processed, returning no errors",
			test: func(t *testing.T) {
				planName := "planName"
				cs := componentStack{}
				pValue := reflect.ValueOf(dummy{})

				eMock := &mockIEngine{}
				vMock := &mockICompletenessValidator{}

				v := completenessValidator{
					itself:    vMock,
					engine:    eMock,
					planValue: reflect.ValueOf("dummy"),
				}

				dummyComponent := parsedComponent{
					id: "dummyComputer",
				}

				eMock.On("findAnalyzedPlan", planName, pValue).
					Return(
						analyzedPlan{
							components: []parsedComponent{dummyComponent},
						},
					).
					Once()

				dummyComputer := registeredComputer{
					metadata: parsedMetadata{},
				}

				eMock.On("getComputer", dummyComponent.id).
					Return(dummyComputer, true).
					Once()

				vMock.On("verifyComponentCompleteness", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					Return(
						func(pm parsedMetadata, cs componentStack, componentID string, planType string) error {
							assert.Equal(t, dummyComputer.metadata, pm)
							assert.Equal(t, dummyComponent.id, componentID)
							assert.Equal(t, v.planValue.Type().String(), planType)
							if assert.Equal(t, 1, len(cs)) {
								assert.Equal(t, planName, cs[0])
							}

							return nil
						},
					).
					Once()

				err := v.doValidate(planName, cs, pValue)
				assert.Nil(t, err)
				mock.AssertExpectationsForObjects(t, eMock, vMock)
			},
		},
		{
			desc: "computer should be processed, returning an error",
			test: func(t *testing.T) {
				planName := "planName"
				cs := componentStack{}
				pValue := reflect.ValueOf(dummy{})

				eMock := &mockIEngine{}
				vMock := &mockICompletenessValidator{}

				v := completenessValidator{
					itself:    vMock,
					engine:    eMock,
					planValue: reflect.ValueOf("dummy"),
				}

				dummyComponent := parsedComponent{
					id: "dummyComputer",
				}

				eMock.On("findAnalyzedPlan", planName, pValue).
					Return(
						analyzedPlan{
							components: []parsedComponent{dummyComponent},
						},
					).
					Once()

				dummyComputer := registeredComputer{
					metadata: parsedMetadata{},
				}

				eMock.On("getComputer", dummyComponent.id).
					Return(dummyComputer, true).
					Once()

				vMock.On("verifyComponentCompleteness", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					Return(
						func(pm parsedMetadata, cs componentStack, componentID string, planType string) error {
							assert.Equal(t, dummyComputer.metadata, pm)
							assert.Equal(t, dummyComponent.id, componentID)
							assert.Equal(t, v.planValue.Type().String(), planType)
							if assert.Equal(t, 1, len(cs)) {
								assert.Equal(t, planName, cs[0])
							}

							return assert.AnError
						},
					).
					Once()

				err := v.doValidate(planName, cs, pValue)
				assert.Equal(t, assert.AnError, err)
				mock.AssertExpectationsForObjects(t, eMock, vMock)
			},
		},
		{
			desc: "nested plan should be processed, returning no errors",
			test: func(t *testing.T) {
				planName := "planName"
				cs := componentStack{}

				eMock := &mockIEngine{}
				vMock := &mockICompletenessValidator{}

				v := completenessValidator{
					itself: vMock,
					engine: eMock,
				}

				t.Run(
					"pointer plan value", func(t *testing.T) {
						pValue := reflect.ValueOf(&doValidate_ParentPlan{})

						dummyComponent := parsedComponent{
							id:       "nestedPlan",
							fieldIdx: 0,
						}

						eMock.On("findAnalyzedPlan", planName, pValue).
							Return(
								analyzedPlan{
									components: []parsedComponent{dummyComponent},
								},
							).
							Once()

						eMock.On("getComputer", dummyComponent.id).
							Return(registeredComputer{}, false).
							Once()

						eMock.On("getPlan", dummyComponent.id).
							Return(analyzedPlan{}, true).
							Once()

						vMock.On("doValidate", mock.Anything, mock.Anything, mock.Anything).
							Return(
								func(pName string, cs componentStack, curPlanValue reflect.Value) error {
									assert.Equal(t, dummyComponent.id, pName)
									assert.Equal(t, pValue.Elem().Field(dummyComponent.fieldIdx), curPlanValue)
									if assert.Equal(t, 1, len(cs)) {
										assert.Equal(t, planName, cs[0])
									}

									return nil
								},
							).
							Once()

						err := v.doValidate(planName, cs, pValue)
						assert.Nil(t, err)
						mock.AssertExpectationsForObjects(t, eMock, vMock)
					},
				)

				t.Run(
					"non-pointer plan value", func(t *testing.T) {
						pValue := reflect.ValueOf(doValidate_ParentPlan{})

						dummyComponent := parsedComponent{
							id:       "nestedPlan",
							fieldIdx: 0,
						}

						eMock.On("findAnalyzedPlan", planName, pValue).
							Return(
								analyzedPlan{
									components: []parsedComponent{dummyComponent},
								},
							).
							Once()

						eMock.On("getComputer", dummyComponent.id).
							Return(registeredComputer{}, false).
							Once()

						eMock.On("getPlan", dummyComponent.id).
							Return(analyzedPlan{}, true).
							Once()

						vMock.On("doValidate", mock.Anything, mock.Anything, mock.Anything).
							Return(
								func(pName string, cs componentStack, curPlanValue reflect.Value) error {
									assert.Equal(t, dummyComponent.id, pName)
									assert.Equal(t, pValue.Field(dummyComponent.fieldIdx), curPlanValue)
									if assert.Equal(t, 1, len(cs)) {
										assert.Equal(t, planName, cs[0])
									}

									return nil
								},
							).
							Once()

						err := v.doValidate(planName, cs, pValue)
						assert.Nil(t, err)
						mock.AssertExpectationsForObjects(t, eMock, vMock)
					},
				)
			},
		},
		{
			desc: "nested plan should be processed, returning an error",
			test: func(t *testing.T) {
				planName := "planName"
				cs := componentStack{}

				eMock := &mockIEngine{}
				vMock := &mockICompletenessValidator{}

				v := completenessValidator{
					itself: vMock,
					engine: eMock,
				}

				t.Run(
					"pointer plan value", func(t *testing.T) {
						pValue := reflect.ValueOf(&doValidate_ParentPlan{})

						dummyComponent := parsedComponent{
							id:       "nestedPlan",
							fieldIdx: 0,
						}

						eMock.On("findAnalyzedPlan", planName, pValue).
							Return(
								analyzedPlan{
									components: []parsedComponent{dummyComponent},
								},
							).
							Once()

						eMock.On("getComputer", dummyComponent.id).
							Return(registeredComputer{}, false).
							Once()

						eMock.On("getPlan", dummyComponent.id).
							Return(analyzedPlan{}, true).
							Once()

						vMock.On("doValidate", mock.Anything, mock.Anything, mock.Anything).
							Return(
								func(pName string, cs componentStack, curPlanValue reflect.Value) error {
									assert.Equal(t, dummyComponent.id, pName)
									assert.Equal(t, pValue.Elem().Field(dummyComponent.fieldIdx), curPlanValue)
									if assert.Equal(t, 1, len(cs)) {
										assert.Equal(t, planName, cs[0])
									}

									return assert.AnError
								},
							).
							Once()

						err := v.doValidate(planName, cs, pValue)
						assert.Equal(t, assert.AnError, err)
						mock.AssertExpectationsForObjects(t, eMock, vMock)
					},
				)

				t.Run(
					"non-pointer plan value", func(t *testing.T) {
						pValue := reflect.ValueOf(doValidate_ParentPlan{})

						dummyComponent := parsedComponent{
							id:       "nestedPlan",
							fieldIdx: 0,
						}

						eMock.On("findAnalyzedPlan", planName, pValue).
							Return(
								analyzedPlan{
									components: []parsedComponent{dummyComponent},
								},
							).
							Once()

						eMock.On("getComputer", dummyComponent.id).
							Return(registeredComputer{}, false).
							Once()

						eMock.On("getPlan", dummyComponent.id).
							Return(analyzedPlan{}, true).
							Once()

						vMock.On("doValidate", mock.Anything, mock.Anything, mock.Anything).
							Return(
								func(pName string, cs componentStack, curPlanValue reflect.Value) error {
									assert.Equal(t, dummyComponent.id, pName)
									assert.Equal(t, pValue.Field(dummyComponent.fieldIdx), curPlanValue)
									if assert.Equal(t, 1, len(cs)) {
										assert.Equal(t, planName, cs[0])
									}

									return assert.AnError
								},
							).
							Once()

						err := v.doValidate(planName, cs, pValue)
						assert.Equal(t, assert.AnError, err)
						mock.AssertExpectationsForObjects(t, eMock, vMock)
					},
				)
			},
		},
	}

	for _, scenario := range scenarios {
		s := scenario

		t.Run(s.desc, s.test)
	}
}

func TestCompletenessValidator_VerifyComponentCompleteness(t *testing.T) {
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

func TestCompletenessValidator_IsInterfaceSatisfied(t *testing.T) {
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
