// Code generated by mockery v2.14.0. DO NOT EDIT.

package cte

import mock "github.com/stretchr/testify/mock"

// mockIFieldAnalyzer is an autogenerated mock type for the iFieldAnalyzer type
type mockIFieldAnalyzer struct {
	mock.Mock
}

// analyze provides a mock function with given fields:
func (_m *mockIFieldAnalyzer) analyze() (*parsedComponent, *preHook, *postHook) {
	ret := _m.Called()

	var r0 *parsedComponent
	if rf, ok := ret.Get(0).(func() *parsedComponent); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*parsedComponent)
		}
	}

	var r1 *preHook
	if rf, ok := ret.Get(1).(func() *preHook); ok {
		r1 = rf()
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*preHook)
		}
	}

	var r2 *postHook
	if rf, ok := ret.Get(2).(func() *postHook); ok {
		r2 = rf()
	} else {
		if ret.Get(2) != nil {
			r2 = ret.Get(2).(*postHook)
		}
	}

	return r0, r1, r2
}

// createComputerComponent provides a mock function with given fields: componentID
func (_m *mockIFieldAnalyzer) createComputerComponent(componentID string) *parsedComponent {
	ret := _m.Called(componentID)

	var r0 *parsedComponent
	if rf, ok := ret.Get(0).(func(string) *parsedComponent); ok {
		r0 = rf(componentID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*parsedComponent)
		}
	}

	return r0
}

// handleComputer provides a mock function with given fields:
func (_m *mockIFieldAnalyzer) handleComputer() *parsedComponent {
	ret := _m.Called()

	var r0 *parsedComponent
	if rf, ok := ret.Get(0).(func() *parsedComponent); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*parsedComponent)
		}
	}

	return r0
}

// handleHooks provides a mock function with given fields:
func (_m *mockIFieldAnalyzer) handleHooks() (*preHook, *postHook) {
	ret := _m.Called()

	var r0 *preHook
	if rf, ok := ret.Get(0).(func() *preHook); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*preHook)
		}
	}

	var r1 *postHook
	if rf, ok := ret.Get(1).(func() *postHook); ok {
		r1 = rf()
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*postHook)
		}
	}

	return r0, r1
}

// handleNestedPlan provides a mock function with given fields:
func (_m *mockIFieldAnalyzer) handleNestedPlan() *parsedComponent {
	ret := _m.Called()

	var r0 *parsedComponent
	if rf, ok := ret.Get(0).(func() *parsedComponent); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*parsedComponent)
		}
	}

	return r0
}

type mockConstructorTestingTnewMockIFieldAnalyzer interface {
	mock.TestingT
	Cleanup(func())
}

// newMockIFieldAnalyzer creates a new instance of mockIFieldAnalyzer. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func newMockIFieldAnalyzer(t mockConstructorTestingTnewMockIFieldAnalyzer) *mockIFieldAnalyzer {
	mock := &mockIFieldAnalyzer{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
