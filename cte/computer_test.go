package cte

import (
	"context"
	"reflect"
	"testing"

	"github.com/jamestrandung/go-concurrency-117/async"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewDelegatingComputer(t *testing.T) {
	scenarios := []struct {
		desc string
		test func(t *testing.T)
	}{
		{
			desc: "ImpureComputerWithLoadingData",
			test: func(t *testing.T) {
				mpMock := &MockMasterPlan{}
				dummyData := LoadingData{
					Data: 2,
					Err:  assert.AnError,
				}

				cMock := &MockImpureComputerWithLoadingData{}
				cMock.On("Load", context.Background(), mpMock).
					Return(1, assert.AnError).
					Once()
				cMock.On("Compute", context.Background(), mpMock, dummyData).
					Return(2, assert.AnError).
					Once()

				dc := newDelegatingComputer(cMock)
				assert.NotNil(t, dc.loadFn)

				assert.NotPanics(
					t, func() {
						result, err := dc.Load(context.Background(), mpMock)
						assert.Equal(t, 1, result)
						assert.Equal(t, assert.AnError, err)
					},
				)

				result, err := dc.Compute(context.Background(), mpMock, dummyData)
				assert.Equal(t, 2, result)
				assert.Equal(t, assert.AnError, err)

				mock.AssertExpectationsForObjects(t, cMock)
			},
		},
		{
			desc: "ImpureComputer",
			test: func(t *testing.T) {
				mpMock := &MockMasterPlan{}

				cMock := &MockImpureComputer{}
				cMock.On("Compute", context.Background(), mpMock).
					Return(2, assert.AnError).
					Once()

				dc := newDelegatingComputer(cMock)
				assert.Nil(t, dc.loadFn)

				assert.NotPanics(
					t, func() {
						result, err := dc.Load(context.Background(), mpMock)
						assert.Nil(t, result)
						assert.Nil(t, err)
					},
				)

				result, err := dc.Compute(context.Background(), mpMock, LoadingData{})
				assert.Equal(t, 2, result)
				assert.Equal(t, assert.AnError, err)

				mock.AssertExpectationsForObjects(t, cMock)
			},
		},
		{
			desc: "SideEffectComputerWithLoadingData",
			test: func(t *testing.T) {
				mpMock := &MockMasterPlan{}
				dummyData := LoadingData{
					Data: 2,
					Err:  assert.AnError,
				}

				cMock := &MockSideEffectComputerWithLoadingData{}
				cMock.On("Load", context.Background(), mpMock).
					Return(1, assert.AnError).
					Once()
				cMock.On("Compute", context.Background(), mpMock, dummyData).
					Return(assert.AnError).
					Once()

				dc := newDelegatingComputer(cMock)
				assert.NotNil(t, dc.loadFn)

				assert.NotPanics(
					t, func() {
						result, err := dc.Load(context.Background(), mpMock)
						assert.Equal(t, 1, result)
						assert.Equal(t, assert.AnError, err)
					},
				)

				_, err := dc.Compute(context.Background(), mpMock, dummyData)
				assert.Equal(t, assert.AnError, err)

				mock.AssertExpectationsForObjects(t, cMock)
			},
		},
		{
			desc: "SideEffectComputer",
			test: func(t *testing.T) {
				mpMock := &MockMasterPlan{}

				cMock := &MockSideEffectComputer{}
				cMock.On("Compute", context.Background(), mpMock).
					Return(assert.AnError).
					Once()

				dc := newDelegatingComputer(cMock)
				assert.Nil(t, dc.loadFn)

				assert.NotPanics(
					t, func() {
						result, err := dc.Load(context.Background(), mpMock)
						assert.Nil(t, result)
						assert.Nil(t, err)
					},
				)

				_, err := dc.Compute(context.Background(), mpMock, LoadingData{})
				assert.Equal(t, assert.AnError, err)

				mock.AssertExpectationsForObjects(t, cMock)
			},
		},
		{
			desc: "SwitchComputerWithLoadingData",
			test: func(t *testing.T) {
				mpMock1 := &MockMasterPlan{}
				mpMock2 := &MockMasterPlan{}
				dummyData := LoadingData{
					Data: 2,
					Err:  assert.AnError,
				}

				cMock := &MockSwitchComputerWithLoadingData{}
				cMock.On("Load", context.Background(), mpMock1).
					Return(1, assert.AnError).
					Once()
				cMock.On("Switch", context.Background(), mpMock1, dummyData).
					Return(mpMock2, assert.AnError).
					Once()

				dc := newDelegatingComputer(cMock)
				assert.NotNil(t, dc.loadFn)

				assert.NotPanics(
					t, func() {
						result, err := dc.Load(context.Background(), mpMock1)
						assert.Equal(t, 1, result)
						assert.Equal(t, assert.AnError, err)
					},
				)

				result, err := dc.Compute(context.Background(), mpMock1, dummyData)
				assert.Equal(t, toExecutePlan{mp: mpMock2}, result)
				assert.Equal(t, assert.AnError, err)

				mock.AssertExpectationsForObjects(t, cMock)
			},
		},
		{
			desc: "SwitchComputer",
			test: func(t *testing.T) {
				mpMock1 := &MockMasterPlan{}
				mpMock2 := &MockMasterPlan{}

				cMock := &MockSwitchComputer{}
				cMock.On("Switch", context.Background(), mpMock1).
					Return(mpMock2, assert.AnError).
					Once()

				dc := newDelegatingComputer(cMock)
				assert.Nil(t, dc.loadFn)

				assert.NotPanics(
					t, func() {
						result, err := dc.Load(context.Background(), mpMock1)
						assert.Nil(t, result)
						assert.Nil(t, err)
					},
				)

				result, err := dc.Compute(context.Background(), mpMock1, LoadingData{})
				assert.Equal(t, toExecutePlan{mp: mpMock2}, result)
				assert.Equal(t, assert.AnError, err)

				mock.AssertExpectationsForObjects(t, cMock)
			},
		},
		{
			desc: "invalid",
			test: func(t *testing.T) {
				invalid := "string"

				assert.PanicsWithError(
					t, ErrInvalidComputerType.Err(reflect.TypeOf(invalid)).Error(), func() {
						newDelegatingComputer(invalid)
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

func TestDelegatingComputer_Load(t *testing.T) {
	scenarios := []struct {
		desc string
		test func(t *testing.T)
	}{
		{
			desc: "nil loadFn",
			test: func(t *testing.T) {
				dc := delegatingComputer{
					loadFn: nil,
				}

				assert.NotPanics(
					t, func() {
						result, err := dc.Load(context.Background(), &MockMasterPlan{})
						assert.Nil(t, result)
						assert.Nil(t, err)
					},
				)
			},
		},
		{
			desc: "non-nil loadFn",
			test: func(t *testing.T) {
				mpMock := &MockMasterPlan{}

				dc := delegatingComputer{
					loadFn: func(ctx context.Context, p MasterPlan) (interface{}, error) {
						assert.Equal(t, context.Background(), ctx)
						assert.Equal(t, mpMock, p)

						return 1, assert.AnError
					},
				}

				result, err := dc.Load(context.Background(), mpMock)
				assert.Equal(t, 1, result)
				assert.Equal(t, assert.AnError, err)
			},
		},
	}

	for _, scenario := range scenarios {
		s := scenario

		t.Run(s.desc, s.test)
	}
}

func TestDelegatingComputer_Compute(t *testing.T) {
	mpMock := &MockMasterPlan{}
	dummyData := LoadingData{
		Data: 2,
		Err:  assert.AnError,
	}

	dc := delegatingComputer{
		computeFn: func(ctx context.Context, p MasterPlan, data LoadingData) (interface{}, error) {
			assert.Equal(t, context.Background(), ctx)
			assert.Equal(t, mpMock, p)
			assert.Equal(t, dummyData, data)

			return 1, assert.AnError
		},
	}

	result, err := dc.Compute(context.Background(), mpMock, dummyData)
	assert.Equal(t, 1, result)
	assert.Equal(t, assert.AnError, err)
}

func TestNewResult(t *testing.T) {
	task := async.Completed("test", assert.AnError)

	actual := newResult(task)
	assert.Equal(t, Result{task}, actual)
}

func TestNewSyncResult(t *testing.T) {
	outcome := "string"

	actual := newSyncResult(outcome)
	assert.Equal(t, SyncResult{outcome}, actual)
}
