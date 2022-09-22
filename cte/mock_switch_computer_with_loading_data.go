// Code generated by mockery v2.14.0. DO NOT EDIT.

package cte

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// MockSwitchComputerWithLoadingData is an autogenerated mock type for the SwitchComputerWithLoadingData type
type MockSwitchComputerWithLoadingData struct {
	mock.Mock
}

// Load provides a mock function with given fields: ctx, p
func (_m *MockSwitchComputerWithLoadingData) Load(ctx context.Context, p MasterPlan) (interface{}, error) {
	ret := _m.Called(ctx, p)

	var r0 interface{}
	if rf, ok := ret.Get(0).(func(context.Context, MasterPlan) interface{}); ok {
		r0 = rf(ctx, p)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interface{})
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, MasterPlan) error); ok {
		r1 = rf(ctx, p)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Switch provides a mock function with given fields: ctx, p, data
func (_m *MockSwitchComputerWithLoadingData) Switch(ctx context.Context, p MasterPlan, data LoadingData) (MasterPlan, error) {
	ret := _m.Called(ctx, p, data)

	var r0 MasterPlan
	if rf, ok := ret.Get(0).(func(context.Context, MasterPlan, LoadingData) MasterPlan); ok {
		r0 = rf(ctx, p, data)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(MasterPlan)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, MasterPlan, LoadingData) error); ok {
		r1 = rf(ctx, p, data)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewMockSwitchComputerWithLoadingData interface {
	mock.TestingT
	Cleanup(func())
}

// NewMockSwitchComputerWithLoadingData creates a new instance of MockSwitchComputerWithLoadingData. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockSwitchComputerWithLoadingData(t mockConstructorTestingTNewMockSwitchComputerWithLoadingData) *MockSwitchComputerWithLoadingData {
	mock := &MockSwitchComputerWithLoadingData{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}