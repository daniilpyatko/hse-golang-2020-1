package post

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// MockCursor is an autogenerated mock type for the MockCursor type
type MockCursor struct {
	mock.Mock
}

// Close provides a mock function with given fields: _a0
func (_m *MockCursor) Close(_a0 context.Context) {
	_m.Called(_a0)
}

// Decode provides a mock function with given fields: _a0
func (_m *MockCursor) Decode(_a0 interface{}) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(interface{}) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Next provides a mock function with given fields: _a0
func (_m *MockCursor) Next(_a0 context.Context) bool {
	ret := _m.Called(_a0)

	var r0 bool
	if rf, ok := ret.Get(0).(func(context.Context) bool); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}
