// Code generated by mockery. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// PasswordService is an autogenerated mock type for the PasswordService type
type PasswordService struct {
	mock.Mock
}

type PasswordService_Expecter struct {
	mock *mock.Mock
}

func (_m *PasswordService) EXPECT() *PasswordService_Expecter {
	return &PasswordService_Expecter{mock: &_m.Mock}
}

// Compare provides a mock function with given fields: currPassword, hashedPassword, salt
func (_m *PasswordService) Compare(currPassword string, hashedPassword string, salt string) (bool, error) {
	ret := _m.Called(currPassword, hashedPassword, salt)

	if len(ret) == 0 {
		panic("no return value specified for Compare")
	}

	var r0 bool
	var r1 error
	if rf, ok := ret.Get(0).(func(string, string, string) (bool, error)); ok {
		return rf(currPassword, hashedPassword, salt)
	}
	if rf, ok := ret.Get(0).(func(string, string, string) bool); ok {
		r0 = rf(currPassword, hashedPassword, salt)
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func(string, string, string) error); ok {
		r1 = rf(currPassword, hashedPassword, salt)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PasswordService_Compare_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Compare'
type PasswordService_Compare_Call struct {
	*mock.Call
}

// Compare is a helper method to define mock.On call
//   - currPassword string
//   - hashedPassword string
//   - salt string
func (_e *PasswordService_Expecter) Compare(currPassword interface{}, hashedPassword interface{}, salt interface{}) *PasswordService_Compare_Call {
	return &PasswordService_Compare_Call{Call: _e.mock.On("Compare", currPassword, hashedPassword, salt)}
}

func (_c *PasswordService_Compare_Call) Run(run func(currPassword string, hashedPassword string, salt string)) *PasswordService_Compare_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(string), args[2].(string))
	})
	return _c
}

func (_c *PasswordService_Compare_Call) Return(_a0 bool, _a1 error) *PasswordService_Compare_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *PasswordService_Compare_Call) RunAndReturn(run func(string, string, string) (bool, error)) *PasswordService_Compare_Call {
	_c.Call.Return(run)
	return _c
}

// GenerateRandomSalt provides a mock function with given fields: saltSize
func (_m *PasswordService) GenerateRandomSalt(saltSize int) (string, error) {
	ret := _m.Called(saltSize)

	if len(ret) == 0 {
		panic("no return value specified for GenerateRandomSalt")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(int) (string, error)); ok {
		return rf(saltSize)
	}
	if rf, ok := ret.Get(0).(func(int) string); ok {
		r0 = rf(saltSize)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(int) error); ok {
		r1 = rf(saltSize)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PasswordService_GenerateRandomSalt_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GenerateRandomSalt'
type PasswordService_GenerateRandomSalt_Call struct {
	*mock.Call
}

// GenerateRandomSalt is a helper method to define mock.On call
//   - saltSize int
func (_e *PasswordService_Expecter) GenerateRandomSalt(saltSize interface{}) *PasswordService_GenerateRandomSalt_Call {
	return &PasswordService_GenerateRandomSalt_Call{Call: _e.mock.On("GenerateRandomSalt", saltSize)}
}

func (_c *PasswordService_GenerateRandomSalt_Call) Run(run func(saltSize int)) *PasswordService_GenerateRandomSalt_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(int))
	})
	return _c
}

func (_c *PasswordService_GenerateRandomSalt_Call) Return(_a0 string, _a1 error) *PasswordService_GenerateRandomSalt_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *PasswordService_GenerateRandomSalt_Call) RunAndReturn(run func(int) (string, error)) *PasswordService_GenerateRandomSalt_Call {
	_c.Call.Return(run)
	return _c
}

// Hash provides a mock function with given fields: password, salt
func (_m *PasswordService) Hash(password string, salt string) (string, error) {
	ret := _m.Called(password, salt)

	if len(ret) == 0 {
		panic("no return value specified for Hash")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(string, string) (string, error)); ok {
		return rf(password, salt)
	}
	if rf, ok := ret.Get(0).(func(string, string) string); ok {
		r0 = rf(password, salt)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(password, salt)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PasswordService_Hash_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Hash'
type PasswordService_Hash_Call struct {
	*mock.Call
}

// Hash is a helper method to define mock.On call
//   - password string
//   - salt string
func (_e *PasswordService_Expecter) Hash(password interface{}, salt interface{}) *PasswordService_Hash_Call {
	return &PasswordService_Hash_Call{Call: _e.mock.On("Hash", password, salt)}
}

func (_c *PasswordService_Hash_Call) Run(run func(password string, salt string)) *PasswordService_Hash_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(string))
	})
	return _c
}

func (_c *PasswordService_Hash_Call) Return(_a0 string, _a1 error) *PasswordService_Hash_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *PasswordService_Hash_Call) RunAndReturn(run func(string, string) (string, error)) *PasswordService_Hash_Call {
	_c.Call.Return(run)
	return _c
}

// Validate provides a mock function with given fields: password
func (_m *PasswordService) Validate(password string) bool {
	ret := _m.Called(password)

	if len(ret) == 0 {
		panic("no return value specified for Validate")
	}

	var r0 bool
	if rf, ok := ret.Get(0).(func(string) bool); ok {
		r0 = rf(password)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// PasswordService_Validate_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Validate'
type PasswordService_Validate_Call struct {
	*mock.Call
}

// Validate is a helper method to define mock.On call
//   - password string
func (_e *PasswordService_Expecter) Validate(password interface{}) *PasswordService_Validate_Call {
	return &PasswordService_Validate_Call{Call: _e.mock.On("Validate", password)}
}

func (_c *PasswordService_Validate_Call) Run(run func(password string)) *PasswordService_Validate_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *PasswordService_Validate_Call) Return(_a0 bool) *PasswordService_Validate_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *PasswordService_Validate_Call) RunAndReturn(run func(string) bool) *PasswordService_Validate_Call {
	_c.Call.Return(run)
	return _c
}

// NewPasswordService creates a new instance of PasswordService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewPasswordService(t interface {
	mock.TestingT
	Cleanup(func())
}) *PasswordService {
	mock := &PasswordService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}