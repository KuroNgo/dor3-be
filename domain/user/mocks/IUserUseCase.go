// Code generated by mockery v2.40.1. DO NOT EDIT.

package mocks

import (
	user_domain "clean-architecture/domain/user"
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// IUserUseCase is an autogenerated mock type for the IUserUseCase type
type IUserUseCase struct {
	mock.Mock
}

// Fetch provides a mock function with given fields: c
func (_m *IUserUseCase) Fetch(c context.Context) ([]user_domain.User, error) {
	ret := _m.Called(c)

	if len(ret) == 0 {
		panic("no return value specified for Fetch")
	}

	var r0 []user_domain.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) ([]user_domain.User, error)); ok {
		return rf(c)
	}
	if rf, ok := ret.Get(0).(func(context.Context) []user_domain.User); ok {
		r0 = rf(c)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]user_domain.User)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(c)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetByEmail provides a mock function with given fields: c, email
func (_m *IUserUseCase) GetByEmail(c context.Context, email string) (*user_domain.User, error) {
	ret := _m.Called(c, email)

	if len(ret) == 0 {
		panic("no return value specified for GetByEmail")
	}

	var r0 *user_domain.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*user_domain.User, error)); ok {
		return rf(c, email)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *user_domain.User); ok {
		r0 = rf(c, email)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*user_domain.User)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(c, email)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetByID provides a mock function with given fields: c, id
func (_m *IUserUseCase) GetByID(c context.Context, id string) (*user_domain.User, error) {
	ret := _m.Called(c, id)

	if len(ret) == 0 {
		panic("no return value specified for GetByID")
	}

	var r0 *user_domain.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*user_domain.User, error)); ok {
		return rf(c, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *user_domain.User); ok {
		r0 = rf(c, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*user_domain.User)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(c, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewIUserUseCase creates a new instance of IUserUseCase. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewIUserUseCase(t interface {
	mock.TestingT
	Cleanup(func())
}) *IUserUseCase {
	mock := &IUserUseCase{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
