package cte

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"reflect"
	"testing"
)

type handleHooks_BadPreHook struct{}

func (handleHooks_BadPreHook) PreExecute(p Plan) error {
	return nil
}

type handleHooks_BadPostHook struct{}

func (handleHooks_BadPostHook) PostExecute(p Plan) error {
	return nil
}

type handleHooks_PreHook struct{}

func (handleHooks_PreHook) CTEMetadata() interface{} {
	return struct{}{}
}

func (handleHooks_PreHook) PreExecute(p Plan) error {
	return nil
}

type handleHooks_PostHook struct{}

func (handleHooks_PostHook) CTEMetadata() interface{} {
	return struct{}{}
}

func (handleHooks_PostHook) PostExecute(p Plan) error {
	return nil
}

type handleHooks_Plan struct{}

func (handleHooks_Plan) IsSequentialCTEPlan() bool {
	return true
}

func TestFieldAnalyzer_HandleHooks(t *testing.T) {
	defer func(original func(mp MetadataProvider, isComputerKey bool) parsedMetadata) {
		extractMetadata = original
	}(extractMetadata)

	scenarios := []struct {
		desc string
		test func(test *testing.T)
	}{
		{
			desc: "field type is plan type",
			test: func(test *testing.T) {
				fa := fieldAnalyzer{
					fieldType:        reflect.TypeOf(handleHooks_Plan{}),
					fieldPointerType: reflect.TypeOf(&handleHooks_PreHook{}),
				}

				actualPre, actualPost := fa.handleHooks()
				assert.Nil(test, actualPre)
				assert.Nil(test, actualPost)
			},
		},
		{
			desc: "field pointer type is plan type",
			test: func(test *testing.T) {
				fa := fieldAnalyzer{
					fieldType:        reflect.TypeOf(handleHooks_PreHook{}),
					fieldPointerType: reflect.TypeOf(&handleHooks_Plan{}),
				}

				actualPre, actualPost := fa.handleHooks()
				assert.Nil(test, actualPre)
				assert.Nil(test, actualPost)
			},
		},
		{
			desc: "field is not pre or post hook type",
			test: func(test *testing.T) {
				fa := fieldAnalyzer{
					fieldType:        reflect.TypeOf(dummy{}),
					fieldPointerType: reflect.TypeOf(&dummy{}),
				}

				actualPre, actualPost := fa.handleHooks()
				assert.Nil(test, actualPre)
				assert.Nil(test, actualPost)
			},
		},
		{
			desc: "field is a pre hook but missing metadata",
			test: func(test *testing.T) {
				fa := fieldAnalyzer{
					fieldType:        reflect.TypeOf(handleHooks_BadPreHook{}),
					fieldPointerType: reflect.TypeOf(&handleHooks_BadPreHook{}),
				}

				assert.PanicsWithError(
					test, ErrMetadataMissing.Err(fa.fieldType).Error(), func() {
						fa.handleHooks()
					},
				)
			},
		},
		{
			desc: "field is a post hook but missing metadata",
			test: func(test *testing.T) {
				fa := fieldAnalyzer{
					fieldType:        reflect.TypeOf(handleHooks_BadPostHook{}),
					fieldPointerType: reflect.TypeOf(&handleHooks_BadPostHook{}),
				}

				assert.PanicsWithError(
					test, ErrMetadataMissing.Err(fa.fieldType).Error(), func() {
						fa.handleHooks()
					},
				)
			},
		},
		{
			desc: "field is a valid pre hook",
			test: func(test *testing.T) {
				fa := fieldAnalyzer{
					fieldType:        reflect.TypeOf(handleHooks_PreHook{}),
					fieldPointerType: reflect.TypeOf(&handleHooks_PreHook{}),
				}

				dummyMetadata := parsedMetadata{}

				extractMetadata = func(mp MetadataProvider, isComputerKey bool) parsedMetadata {
					assert.Equal(test, &handleHooks_PreHook{}, mp)
					assert.False(test, isComputerKey)

					return dummyMetadata
				}

				expectedPre := &preHook{
					hook:     &handleHooks_PreHook{},
					metadata: dummyMetadata,
				}

				actualPre, actualPost := fa.handleHooks()
				assert.Equal(test, expectedPre, actualPre)
				assert.Nil(test, actualPost)
			},
		},
		{
			desc: "field is a valid post hook",
			test: func(test *testing.T) {
				fa := fieldAnalyzer{
					fieldType:        reflect.TypeOf(handleHooks_PostHook{}),
					fieldPointerType: reflect.TypeOf(&handleHooks_PostHook{}),
				}

				dummyMetadata := parsedMetadata{}

				extractMetadata = func(mp MetadataProvider, isComputerKey bool) parsedMetadata {
					assert.Equal(test, &handleHooks_PostHook{}, mp)
					assert.False(test, isComputerKey)

					return dummyMetadata
				}

				expectedPost := &postHook{
					hook:     &handleHooks_PostHook{},
					metadata: dummyMetadata,
				}

				actualPre, actualPost := fa.handleHooks()
				assert.Nil(test, actualPre)
				assert.Equal(test, expectedPost, actualPost)
			},
		},
	}

	for _, scenario := range scenarios {
		s := scenario

		t.Run(s.desc, s.test)
	}
}

type handleNestedPlan_Plan struct{}

func (handleNestedPlan_Plan) IsSequentialCTEPlan() bool {
	return true
}

func TestFieldAnalyzer_HandleNestedPlan(t *testing.T) {
	defer func(original func(t reflect.Type) string) {
		extractFullNameFromType = original
	}(extractFullNameFromType)

	scenarios := []struct {
		desc string
		test func(test *testing.T)
	}{
		{
			desc: "fieldPointerType is not a plan",
			test: func(test *testing.T) {
				fa := fieldAnalyzer{
					fieldPointerType: reflect.TypeOf(&dummy{}),
				}

				actual := fa.handleNestedPlan()
				assert.Nil(test, actual)
			},
		},
		{
			desc: "fieldPointerType is a plan, field is a pointer",
			test: func(test *testing.T) {
				fa := fieldAnalyzer{
					pa: &planAnalyzer{
						planValue: reflect.ValueOf(dummy{}),
					},
					isPointerType:    true,
					fieldPointerType: reflect.TypeOf(&handleNestedPlan_Plan{}),
				}

				assert.PanicsWithError(
					test, ErrNestedPlanCannotBePointer.Err(fa.pa.planValue.Type(), fa.fieldType).Error(), func() {
						fa.handleNestedPlan()
					},
				)
			},
		},
		{
			desc: "fieldPointerType is a plan, field is not a pointer",
			test: func(test *testing.T) {
				eMock := &mockIEngine{}

				fa := fieldAnalyzer{
					pa: &planAnalyzer{
						engine:    eMock,
						planValue: reflect.ValueOf(dummy{}),
					},
					fieldIdx:         99,
					isPointerType:    false,
					fieldType:        reflect.TypeOf(handleNestedPlan_Plan{}),
					fieldPointerType: reflect.TypeOf(&handleNestedPlan_Plan{}),
				}

				eMock.On("AnalyzePlan", &handleNestedPlan_Plan{}).Once()

				dummyComponentName := "dummy"
				extractFullNameFromType = func(t reflect.Type) string {
					assert.Equal(test, fa.fieldType, t)

					return dummyComponentName
				}

				expected := &parsedComponent{
					id:       dummyComponentName,
					fieldIdx: fa.fieldIdx,
				}

				actual := fa.handleNestedPlan()
				assert.Equal(test, expected, actual)
				mock.AssertExpectationsForObjects(test, eMock)
			},
		},
	}

	for _, scenario := range scenarios {
		s := scenario

		t.Run(s.desc, s.test)
	}
}

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

				actual := fa.createComputerComponent(componentID)
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
						actual := fa.createComputerComponent(componentID)
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
						actual := fa.createComputerComponent(componentID)
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
						actual := fa.createComputerComponent(componentID)
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
						actual := fa.createComputerComponent(componentID)
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
						fa.createComputerComponent(componentID)
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
						actual := fa.createComputerComponent(componentID)
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
						fa.createComputerComponent(componentID)
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
						actual := fa.createComputerComponent(componentID)
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
						actual := fa.createComputerComponent(componentID)
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
						actual := fa.createComputerComponent(componentID)
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
						actual := fa.createComputerComponent(componentID)
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
						actual := fa.createComputerComponent(componentID)
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
						fa.createComputerComponent(componentID)
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
						actual := fa.createComputerComponent(componentID)
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
						fa.createComputerComponent(componentID)
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

				actual := fa.createComputerComponent(componentID)
				assert.Nil(t, actual)
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

				actual := fa.createComputerComponent(componentID)
				assert.Nil(t, actual)
			},
		},
	}

	for _, scenario := range scenarios {
		s := scenario

		t.Run(s.desc, s.test)
	}
}
