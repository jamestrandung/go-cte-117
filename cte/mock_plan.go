// Code generated by mockery v2.14.0. DO NOT EDIT.

package cte

import mock "github.com/stretchr/testify/mock"

// MockPlan is an autogenerated mock type for the Plan type
type MockPlan struct {
	mock.Mock
}

// IsSequentialCTEPlan provides a mock function with given fields:
func (_m *MockPlan) IsSequentialCTEPlan() bool {
	ret := _m.Called()

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

type mockConstructorTestingTNewMockPlan interface {
	mock.TestingT
	Cleanup(func())
}

// NewMockPlan creates a new instance of MockPlan. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockPlan(t mockConstructorTestingTNewMockPlan) *MockPlan {
	mock := &MockPlan{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
