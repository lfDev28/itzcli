// Code generated by mockery v2.42.3. DO NOT EDIT.

package mocks

import (
	techzone "github.com/cloud-native-toolkit/itzcli/pkg/techzone"
	mock "github.com/stretchr/testify/mock"
)

// Filter is an autogenerated mock type for the Filter type
type Filter struct {
	mock.Mock
}

// Execute provides a mock function with given fields: _a0
func (_m *Filter) Execute(_a0 techzone.Reservation) bool {
	ret := _m.Called(_a0)

	if len(ret) == 0 {
		panic("no return value specified for Execute")
	}

	var r0 bool
	if rf, ok := ret.Get(0).(func(techzone.Reservation) bool); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// NewFilter creates a new instance of Filter. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewFilter(t interface {
	mock.TestingT
	Cleanup(func())
}) *Filter {
	mock := &Filter{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
