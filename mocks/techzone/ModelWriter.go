// Code generated by mockery v2.25.1. DO NOT EDIT.

package mocks

import (
	io "io"

	mock "github.com/stretchr/testify/mock"
)

// ModelWriter is an autogenerated mock type for the ModelWriter type
type ModelWriter struct {
	mock.Mock
}

// WriteMany provides a mock function with given fields: w, val
func (_m *ModelWriter) WriteMany(w io.Writer, val interface{}) error {
	ret := _m.Called(w, val)

	var r0 error
	if rf, ok := ret.Get(0).(func(io.Writer, interface{}) error); ok {
		r0 = rf(w, val)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// WriteOne provides a mock function with given fields: w, val
func (_m *ModelWriter) WriteOne(w io.Writer, val interface{}) error {
	ret := _m.Called(w, val)

	var r0 error
	if rf, ok := ret.Get(0).(func(io.Writer, interface{}) error); ok {
		r0 = rf(w, val)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewModelWriter interface {
	mock.TestingT
	Cleanup(func())
}

// NewModelWriter creates a new instance of ModelWriter. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewModelWriter(t mockConstructorTestingTNewModelWriter) *ModelWriter {
	mock := &ModelWriter{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}