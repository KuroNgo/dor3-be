// Code generated by mockery v2.40.1. DO NOT EDIT.

package mocks

import (
	user_domain "clean-architecture/domain/request/user"
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// IProfileUseCase is an autogenerated mock type for the IProfileUseCase type
type IProfileUseCase struct {
	mock.Mock
}

// GetProfileByID provides a mock function with given fields: c, userID
func (_m *IProfileUseCase) GetProfileByID(c context.Context, userID string) (*user_domain.Profile, error) {
	ret := _m.Called(c, userID)

	if len(ret) == 0 {
		panic("no return value specified for GetProfileByID")
	}

	var r0 *user_domain.Profile
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*user_domain.Profile, error)); ok {
		return rf(c, userID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *user_domain.Profile); ok {
		r0 = rf(c, userID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*user_domain.Profile)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(c, userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewIProfileUseCase creates a new instance of IProfileUseCase. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewIProfileUseCase(t interface {
	mock.TestingT
	Cleanup(func())
}) *IProfileUseCase {
	mock := &IProfileUseCase{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
