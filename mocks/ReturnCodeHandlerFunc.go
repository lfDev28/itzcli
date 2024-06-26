// Code generated by mockery v2.42.3. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// ReturnCodeHandlerFunc is an autogenerated mock type for the ReturnCodeHandlerFunc type
type ReturnCodeHandlerFunc struct {
	mock.Mock
}

// Execute provides a mock function with given fields: code
func (_m *ReturnCodeHandlerFunc) Execute(code int) error {
	ret := _m.Called(code)

	if len(ret) == 0 {
		panic("no return value specified for Execute")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(int) error); ok {
		r0 = rf(code)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewReturnCodeHandlerFunc creates a new instance of ReturnCodeHandlerFunc. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewReturnCodeHandlerFunc(t interface {
	mock.TestingT
	Cleanup(func())
}) *ReturnCodeHandlerFunc {
	mock := &ReturnCodeHandlerFunc{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
