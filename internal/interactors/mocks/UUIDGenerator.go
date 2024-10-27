// Code generated by mockery. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// UUIDGenerator is an autogenerated mock type for the UUIDGenerator type
type UUIDGenerator struct {
	mock.Mock
}

type UUIDGenerator_Expecter struct {
	mock *mock.Mock
}

func (_m *UUIDGenerator) EXPECT() *UUIDGenerator_Expecter {
	return &UUIDGenerator_Expecter{mock: &_m.Mock}
}

// Generate provides a mock function with given fields:
func (_m *UUIDGenerator) Generate() (string, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Generate")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func() (string, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UUIDGenerator_Generate_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Generate'
type UUIDGenerator_Generate_Call struct {
	*mock.Call
}

// Generate is a helper method to define mock.On call
func (_e *UUIDGenerator_Expecter) Generate() *UUIDGenerator_Generate_Call {
	return &UUIDGenerator_Generate_Call{Call: _e.mock.On("Generate")}
}

func (_c *UUIDGenerator_Generate_Call) Run(run func()) *UUIDGenerator_Generate_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *UUIDGenerator_Generate_Call) Return(_a0 string, _a1 error) *UUIDGenerator_Generate_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *UUIDGenerator_Generate_Call) RunAndReturn(run func() (string, error)) *UUIDGenerator_Generate_Call {
	_c.Call.Return(run)
	return _c
}

// NewUUIDGenerator creates a new instance of UUIDGenerator. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewUUIDGenerator(t interface {
	mock.TestingT
	Cleanup(func())
}) *UUIDGenerator {
	mock := &UUIDGenerator{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
