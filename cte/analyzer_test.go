package cte

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"reflect"
	"testing"
)

type handleComputer_MetadataProvider struct{}

func (handleComputer_MetadataProvider) CTEMetadata() interface{} {
	return struct{}{}
}

func TestFieldAnalyzer_HandleComputer(t *testing.T) {
	defer func(original func(t reflect.Type) string) {
		extractFullNameFromType = original
	}(extractFullNameFromType)

	scenarios := []struct {
		desc string
		test func(test *testing.T)
	}{
		{
			desc: "createComputerComponent returns nil",
			test: func(test *testing.T) {
				eMock := &mockIEngine{}
				faMock := &mockIFieldAnalyzer{}

				fa := fieldAnalyzer{
					pa: &planAnalyzer{
						engine: eMock,
					},
					itself:    faMock,
					fieldType: reflect.TypeOf(dummy{}),
				}

				dummyComponentID := "dummy"
				extractFullNameFromType = func(t reflect.Type) string {
					assert.Equal(test, fa.fieldType, t)

					return dummyComponentID
				}

				faMock.On("createComputerComponent", dummyComponentID).
					Return(nil).
					Once()

				actual := fa.handleComputer()
				assert.Nil(test, actual)
				mock.AssertExpectationsForObjects(test, eMock, faMock)
			},
		},
		{
			desc: "field type is not MetadataProvider",
			test: func(test *testing.T) {
				eMock := &mockIEngine{}
				faMock := &mockIFieldAnalyzer{}

				fa := fieldAnalyzer{
					pa: &planAnalyzer{
						engine: eMock,
					},
					itself:    faMock,
					fieldType: reflect.TypeOf(dummy{}),
				}

				dummyComponentID := "dummy"
				extractFullNameFromType = func(t reflect.Type) string {
					assert.Equal(test, fa.fieldType, t)

					return dummyComponentID
				}

				faMock.On("createComputerComponent", dummyComponentID).
					Return(&parsedComponent{}).
					Once()

				assert.PanicsWithError(
					test, ErrMetadataMissing.Err(fa.fieldType).Error(), func() {
						fa.handleComputer()
					},
				)

				mock.AssertExpectationsForObjects(test, eMock, faMock)
			},
		},
		{
			desc: "field type is MetadataProvider",
			test: func(test *testing.T) {
				eMock := &mockIEngine{}
				faMock := &mockIFieldAnalyzer{}

				fa := fieldAnalyzer{
					pa: &planAnalyzer{
						engine: eMock,
					},
					itself:    faMock,
					fieldType: reflect.TypeOf(handleComputer_MetadataProvider{}),
				}

				dummyComponentID := "dummy"
				extractFullNameFromType = func(t reflect.Type) string {
					assert.Equal(test, fa.fieldType, t)

					return dummyComponentID
				}

				dummyComponent := &parsedComponent{}
				faMock.On("createComputerComponent", dummyComponentID).
					Return(dummyComponent).
					Once()

				eMock.On("registerComputer", &handleComputer_MetadataProvider{}).Once()

				actual := fa.handleComputer()
				assert.Equal(test, dummyComponent, actual)
				mock.AssertExpectationsForObjects(test, eMock, faMock)
			},
		},
	}

	for _, scenario := range scenarios {
		s := scenario

		t.Run(s.desc, s.test)
	}
}

type createComputerComponent_SequentialPlan struct{}

func (createComputerComponent_SequentialPlan) IsSequentialCTEPlan() bool {
	return true
}

type createComputerComponent_ParallelPlan struct{}

func (createComputerComponent_ParallelPlan) IsSequentialCTEPlan() bool {
	return false
}

func TestFieldAnalyzer_CreateComputerComponent(t *testing.T) {
	scenarios := []struct {
		desc string
		test func(t *testing.T)
	}{
		{
			desc: "Field is a *Result, Plan is sequential",
			test: func(t *testing.T) {
				type dummyResult Result
				componentID := "componentID"

				fa := fieldAnalyzer{
					pa: &planAnalyzer{
						plan:      createComputerComponent_SequentialPlan{},
						planValue: reflect.ValueOf(createComputerComponent_SequentialPlan{}),
					},
					fieldIdx:         99,
					isPointerType:    true,
					fieldType:        reflect.TypeOf(dummyResult{}),
					fieldPointerType: reflect.TypeOf(&dummyResult{}),
				}

				expected := &parsedComponent{
					id:            componentID,
					fieldIdx:      fa.fieldIdx,
					fieldType:     fa.fieldType,
					isSyncResult:  false,
					requireSet:    true,
					isPointerType: fa.isPointerType,
				}

				actual := fa.createComputerComponent("componentID")
				assert.Equal(t, expected, actual)
			},
		},
		{
			desc: "Field is a *Result, Plan is parallel",
			test: func(t *testing.T) {
				type dummyResult Result
				componentID := "componentID"

				fa := fieldAnalyzer{
					pa: &planAnalyzer{
						plan:      createComputerComponent_ParallelPlan{},
						planValue: reflect.ValueOf(createComputerComponent_ParallelPlan{}),
					},
					fieldIdx:         99,
					isPointerType:    true,
					fieldType:        reflect.TypeOf(dummyResult{}),
					fieldPointerType: reflect.TypeOf(&dummyResult{}),
				}

				expected := &parsedComponent{
					id:            componentID,
					fieldIdx:      fa.fieldIdx,
					fieldType:     fa.fieldType,
					isSyncResult:  false,
					requireSet:    true,
					isPointerType: fa.isPointerType,
				}

				assert.NotPanics(
					t, func() {
						actual := fa.createComputerComponent("componentID")
						assert.Equal(t, expected, actual)
					},
				)
			},
		},
		{
			desc: "Field is a Result, Plan is sequential",
			test: func(t *testing.T) {
				type dummyResult Result
				componentID := "componentID"

				fa := fieldAnalyzer{
					pa: &planAnalyzer{
						plan:      createComputerComponent_SequentialPlan{},
						planValue: reflect.ValueOf(createComputerComponent_SequentialPlan{}),
					},
					fieldIdx:         99,
					isPointerType:    false,
					fieldType:        reflect.TypeOf(dummyResult{}),
					fieldPointerType: reflect.TypeOf(&dummyResult{}),
				}

				expected := &parsedComponent{
					id:            componentID,
					fieldIdx:      fa.fieldIdx,
					fieldType:     fa.fieldType,
					isSyncResult:  false,
					requireSet:    true,
					isPointerType: fa.isPointerType,
				}

				assert.NotPanics(
					t, func() {
						actual := fa.createComputerComponent("componentID")
						assert.Equal(t, expected, actual)
					},
				)
			},
		},
		{
			desc: "Field is a Result, Plan is parallel",
			test: func(t *testing.T) {
				type dummyResult Result
				componentID := "componentID"

				fa := fieldAnalyzer{
					pa: &planAnalyzer{
						plan:      createComputerComponent_ParallelPlan{},
						planValue: reflect.ValueOf(createComputerComponent_ParallelPlan{}),
					},
					fieldIdx:         99,
					isPointerType:    false,
					fieldType:        reflect.TypeOf(dummyResult{}),
					fieldPointerType: reflect.TypeOf(&dummyResult{}),
				}

				expected := &parsedComponent{
					id:            componentID,
					fieldIdx:      fa.fieldIdx,
					fieldType:     fa.fieldType,
					isSyncResult:  false,
					requireSet:    true,
					isPointerType: fa.isPointerType,
				}

				assert.NotPanics(
					t, func() {
						actual := fa.createComputerComponent("componentID")
						assert.Equal(t, expected, actual)
					},
				)
			},
		},
		{
			desc: "Field is a *SyncResult, Plan is sequential",
			test: func(t *testing.T) {
				type dummyResult SyncResult
				componentID := "componentID"

				fa := fieldAnalyzer{
					pa: &planAnalyzer{
						plan:      createComputerComponent_SequentialPlan{},
						planValue: reflect.ValueOf(createComputerComponent_SequentialPlan{}),
					},
					fieldIdx:         99,
					isPointerType:    true,
					fieldType:        reflect.TypeOf(dummyResult{}),
					fieldPointerType: reflect.TypeOf(&dummyResult{}),
				}

				expected := &parsedComponent{
					id:            componentID,
					fieldIdx:      fa.fieldIdx,
					fieldType:     fa.fieldType,
					isSyncResult:  true,
					requireSet:    true,
					isPointerType: fa.isPointerType,
				}

				assert.NotPanics(
					t, func() {
						actual := fa.createComputerComponent("componentID")
						assert.Equal(t, expected, actual)
					},
				)
			},
		},
		{
			desc: "Field is a *SyncResult, Plan is parallel",
			test: func(t *testing.T) {
				type dummyResult SyncResult
				componentID := "componentID"

				fa := fieldAnalyzer{
					pa: &planAnalyzer{
						plan:      createComputerComponent_ParallelPlan{},
						planValue: reflect.ValueOf(createComputerComponent_ParallelPlan{}),
					},
					fieldIdx:         99,
					isPointerType:    true,
					fieldType:        reflect.TypeOf(dummyResult{}),
					fieldPointerType: reflect.TypeOf(&dummyResult{}),
				}

				assert.PanicsWithError(
					t, ErrParallelPlanCannotContainSyncResult.Err(fa.pa.planValue.Type(), extractShortName(componentID)).Error(), func() {
						fa.createComputerComponent("componentID")
					},
				)
			},
		},
		{
			desc: "Field is a SyncResult, Plan is sequential",
			test: func(t *testing.T) {
				type dummyResult SyncResult
				componentID := "componentID"

				fa := fieldAnalyzer{
					pa: &planAnalyzer{
						plan:      createComputerComponent_SequentialPlan{},
						planValue: reflect.ValueOf(createComputerComponent_SequentialPlan{}),
					},
					fieldIdx:         99,
					isPointerType:    false,
					fieldType:        reflect.TypeOf(dummyResult{}),
					fieldPointerType: reflect.TypeOf(&dummyResult{}),
				}

				expected := &parsedComponent{
					id:            componentID,
					fieldIdx:      fa.fieldIdx,
					fieldType:     fa.fieldType,
					isSyncResult:  true,
					requireSet:    true,
					isPointerType: fa.isPointerType,
				}

				assert.NotPanics(
					t, func() {
						actual := fa.createComputerComponent("componentID")
						assert.Equal(t, expected, actual)
					},
				)
			},
		},
		{
			desc: "Field is a SyncResult, Plan is parallel",
			test: func(t *testing.T) {
				type dummyResult SyncResult
				componentID := "componentID"

				fa := fieldAnalyzer{
					pa: &planAnalyzer{
						plan:      createComputerComponent_ParallelPlan{},
						planValue: reflect.ValueOf(createComputerComponent_ParallelPlan{}),
					},
					fieldIdx:         99,
					isPointerType:    false,
					fieldType:        reflect.TypeOf(dummyResult{}),
					fieldPointerType: reflect.TypeOf(&dummyResult{}),
				}

				assert.PanicsWithError(
					t, ErrParallelPlanCannotContainSyncResult.Err(fa.pa.planValue.Type(), extractShortName(componentID)).Error(), func() {
						fa.createComputerComponent("componentID")
					},
				)
			},
		},
		{
			desc: "Field is a *SideEffect, Plan is sequential",
			test: func(t *testing.T) {
				type dummyResult SideEffect
				componentID := "componentID"

				fa := fieldAnalyzer{
					pa: &planAnalyzer{
						plan:      createComputerComponent_SequentialPlan{},
						planValue: reflect.ValueOf(createComputerComponent_SequentialPlan{}),
					},
					fieldIdx:         99,
					isPointerType:    true,
					fieldType:        reflect.TypeOf(dummyResult{}),
					fieldPointerType: reflect.TypeOf(&dummyResult{}),
				}

				expected := &parsedComponent{
					id:       componentID,
					fieldIdx: fa.fieldIdx,
				}

				assert.NotPanics(
					t, func() {
						actual := fa.createComputerComponent("componentID")
						assert.Equal(t, expected, actual)
					},
				)
			},
		},
		{
			desc: "Field is a *SideEffect, Plan is parallel",
			test: func(t *testing.T) {
				type dummyResult SideEffect
				componentID := "componentID"

				fa := fieldAnalyzer{
					pa: &planAnalyzer{
						plan:      createComputerComponent_ParallelPlan{},
						planValue: reflect.ValueOf(createComputerComponent_ParallelPlan{}),
					},
					fieldIdx:         99,
					isPointerType:    true,
					fieldType:        reflect.TypeOf(dummyResult{}),
					fieldPointerType: reflect.TypeOf(&dummyResult{}),
				}

				expected := &parsedComponent{
					id:       componentID,
					fieldIdx: fa.fieldIdx,
				}

				assert.NotPanics(
					t, func() {
						actual := fa.createComputerComponent("componentID")
						assert.Equal(t, expected, actual)
					},
				)
			},
		},
		{
			desc: "Field is a SideEffect, Plan is sequential",
			test: func(t *testing.T) {
				type dummyResult SideEffect
				componentID := "componentID"

				fa := fieldAnalyzer{
					pa: &planAnalyzer{
						plan:      createComputerComponent_SequentialPlan{},
						planValue: reflect.ValueOf(createComputerComponent_SequentialPlan{}),
					},
					fieldIdx:         99,
					isPointerType:    false,
					fieldType:        reflect.TypeOf(dummyResult{}),
					fieldPointerType: reflect.TypeOf(&dummyResult{}),
				}

				expected := &parsedComponent{
					id:       componentID,
					fieldIdx: fa.fieldIdx,
				}

				assert.NotPanics(
					t, func() {
						actual := fa.createComputerComponent("componentID")
						assert.Equal(t, expected, actual)
					},
				)
			},
		},
		{
			desc: "Field is a SideEffect, Plan is parallel",
			test: func(t *testing.T) {
				type dummyResult SideEffect
				componentID := "componentID"

				fa := fieldAnalyzer{
					pa: &planAnalyzer{
						plan:      createComputerComponent_ParallelPlan{},
						planValue: reflect.ValueOf(createComputerComponent_ParallelPlan{}),
					},
					fieldIdx:         99,
					isPointerType:    false,
					fieldType:        reflect.TypeOf(dummyResult{}),
					fieldPointerType: reflect.TypeOf(&dummyResult{}),
				}

				expected := &parsedComponent{
					id:       componentID,
					fieldIdx: fa.fieldIdx,
				}

				assert.NotPanics(
					t, func() {
						actual := fa.createComputerComponent("componentID")
						assert.Equal(t, expected, actual)
					},
				)
			},
		},
		{
			desc: "Field is a *SyncSideEffect, Plan is sequential",
			test: func(t *testing.T) {
				type dummyResult SyncSideEffect
				componentID := "componentID"

				fa := fieldAnalyzer{
					pa: &planAnalyzer{
						plan:      createComputerComponent_SequentialPlan{},
						planValue: reflect.ValueOf(createComputerComponent_SequentialPlan{}),
					},
					fieldIdx:         99,
					isPointerType:    true,
					fieldType:        reflect.TypeOf(dummyResult{}),
					fieldPointerType: reflect.TypeOf(&dummyResult{}),
				}

				expected := &parsedComponent{
					id:       componentID,
					fieldIdx: fa.fieldIdx,
				}

				assert.NotPanics(
					t, func() {
						actual := fa.createComputerComponent("componentID")
						assert.Equal(t, expected, actual)
					},
				)
			},
		},
		{
			desc: "Field is a *SyncSideEffect, Plan is parallel",
			test: func(t *testing.T) {
				type dummyResult SyncSideEffect
				componentID := "componentID"

				fa := fieldAnalyzer{
					pa: &planAnalyzer{
						plan:      createComputerComponent_ParallelPlan{},
						planValue: reflect.ValueOf(createComputerComponent_ParallelPlan{}),
					},
					fieldIdx:         99,
					isPointerType:    true,
					fieldType:        reflect.TypeOf(dummyResult{}),
					fieldPointerType: reflect.TypeOf(&dummyResult{}),
				}

				assert.PanicsWithError(
					t, ErrParallelPlanCannotContainSyncSideEffect.Err(fa.pa.planValue.Type(), extractShortName(componentID)).Error(), func() {
						fa.createComputerComponent("componentID")
					},
				)
			},
		},
		{
			desc: "Field is a SyncSideEffect, Plan is sequential",
			test: func(t *testing.T) {
				type dummyResult SyncSideEffect
				componentID := "componentID"

				fa := fieldAnalyzer{
					pa: &planAnalyzer{
						plan:      createComputerComponent_SequentialPlan{},
						planValue: reflect.ValueOf(createComputerComponent_SequentialPlan{}),
					},
					fieldIdx:         99,
					isPointerType:    false,
					fieldType:        reflect.TypeOf(dummyResult{}),
					fieldPointerType: reflect.TypeOf(&dummyResult{}),
				}

				expected := &parsedComponent{
					id:       componentID,
					fieldIdx: fa.fieldIdx,
				}

				assert.NotPanics(
					t, func() {
						actual := fa.createComputerComponent("componentID")
						assert.Equal(t, expected, actual)
					},
				)
			},
		},
		{
			desc: "Field is a SyncSideEffect, Plan is parallel",
			test: func(t *testing.T) {
				type dummyResult SyncSideEffect
				componentID := "componentID"

				fa := fieldAnalyzer{
					pa: &planAnalyzer{
						plan:      createComputerComponent_ParallelPlan{},
						planValue: reflect.ValueOf(createComputerComponent_ParallelPlan{}),
					},
					fieldIdx:         99,
					isPointerType:    false,
					fieldType:        reflect.TypeOf(dummyResult{}),
					fieldPointerType: reflect.TypeOf(&dummyResult{}),
				}

				assert.PanicsWithError(
					t, ErrParallelPlanCannotContainSyncSideEffect.Err(fa.pa.planValue.Type(), extractShortName(componentID)).Error(), func() {
						fa.createComputerComponent("componentID")
					},
				)
			},
		},
		{
			desc: "Unknown computer key type, Plan is sequential",
			test: func(t *testing.T) {
				type dummyResult dummy
				componentID := "componentID"

				fa := fieldAnalyzer{
					pa: &planAnalyzer{
						plan:      createComputerComponent_SequentialPlan{},
						planValue: reflect.ValueOf(createComputerComponent_SequentialPlan{}),
					},
					fieldIdx:         99,
					isPointerType:    true,
					fieldType:        reflect.TypeOf(dummyResult{}),
					fieldPointerType: reflect.TypeOf(&dummyResult{}),
				}

				assert.PanicsWithError(
					t, ErrUnknownComputerKeyType.Err(fa.pa.planValue.Type(), extractShortName(componentID)).Error(), func() {
						fa.createComputerComponent("componentID")
					},
				)
			},
		},
		{
			desc: "Unknown computer key type, Plan is parallel",
			test: func(t *testing.T) {
				type dummyResult dummy
				componentID := "componentID"

				fa := fieldAnalyzer{
					pa: &planAnalyzer{
						plan:      createComputerComponent_ParallelPlan{},
						planValue: reflect.ValueOf(createComputerComponent_ParallelPlan{}),
					},
					fieldIdx:         99,
					isPointerType:    true,
					fieldType:        reflect.TypeOf(dummyResult{}),
					fieldPointerType: reflect.TypeOf(&dummyResult{}),
				}

				assert.PanicsWithError(
					t, ErrUnknownComputerKeyType.Err(fa.pa.planValue.Type(), extractShortName(componentID)).Error(), func() {
						fa.createComputerComponent("componentID")
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
