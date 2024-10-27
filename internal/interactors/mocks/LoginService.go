// Code generated by mockery. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// LoginService is an autogenerated mock type for the LoginService type
type LoginService struct {
	mock.Mock
}

type LoginService_Expecter struct {
	mock *mock.Mock
}

func (_m *LoginService) EXPECT() *LoginService_Expecter {
	return &LoginService_Expecter{mock: &_m.Mock}
}

// Validate provides a mock function with given fields: login
func (_m *LoginService) Validate(login string) error {
	ret := _m.Called(login)

	if len(ret) == 0 {
		panic("no return value specified for Validate")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(login)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// LoginService_Validate_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Validate'
type LoginService_Validate_Call struct {
	*mock.Call
}

// Validate is a helper method to define mock.On call
//   - login string
func (_e *LoginService_Expecter) Validate(login interface{}) *LoginService_Validate_Call {
	return &LoginService_Validate_Call{Call: _e.mock.On("Validate", login)}
}

func (_c *LoginService_Validate_Call) Run(run func(login string)) *LoginService_Validate_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *LoginService_Validate_Call) Return(_a0 error) *LoginService_Validate_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *LoginService_Validate_Call) RunAndReturn(run func(string) error) *LoginService_Validate_Call {
	_c.Call.Return(run)
	return _c
}

// NewLoginService creates a new instance of LoginService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewLoginService(t interface {
	mock.TestingT
	Cleanup(func())
}) *LoginService {
	mock := &LoginService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
