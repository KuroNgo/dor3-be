// Code generated by mockery v2.40.1. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
	primitive "go.mongodb.org/mongo-driver/bson/primitive"

	user_domain "clean-architecture/domain/request/user"
)

// IUserRepository is an autogenerated mock type for the IUserRepository type
type IUserRepository struct {
	mock.Mock
}

// Create provides a mock function with given fields: c, user
func (_m *IUserRepository) Create(c context.Context, user *user_domain.User) error {
	ret := _m.Called(c, user)

	if len(ret) == 0 {
		panic("no return value specified for Create")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *user_domain.User) error); ok {
		r0 = rf(c, user)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CreateAsync provides a mock function with given fields: c, user
func (_m *IUserRepository) CreateAsync(c context.Context, user *user_domain.User) <-chan error {
	ret := _m.Called(c, user)

	if len(ret) == 0 {
		panic("no return value specified for CreateAsync")
	}

	var r0 <-chan error
	if rf, ok := ret.Get(0).(func(context.Context, *user_domain.User) <-chan error); ok {
		r0 = rf(c, user)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan error)
		}
	}

	return r0
}

// Delete provides a mock function with given fields: c, userID
func (_m *IUserRepository) Delete(c context.Context, userID primitive.ObjectID) error {
	ret := _m.Called(c, userID)

	if len(ret) == 0 {
		panic("no return value specified for Delete")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, primitive.ObjectID) error); ok {
		r0 = rf(c, userID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Fetch provides a mock function with given fields: c
func (_m *IUserRepository) Fetch(c context.Context) ([]user_domain.User, error) {
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
func (_m *IUserRepository) GetByEmail(c context.Context, email string) (user_domain.User, error) {
	ret := _m.Called(c, email)

	if len(ret) == 0 {
		panic("no return value specified for GetByEmail")
	}

	var r0 user_domain.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (user_domain.User, error)); ok {
		return rf(c, email)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) user_domain.User); ok {
		r0 = rf(c, email)
	} else {
		r0 = ret.Get(0).(user_domain.User)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(c, email)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetByID provides a mock function with given fields: c, id
func (_m *IUserRepository) GetByID(c context.Context, id primitive.ObjectID) (user_domain.User, error) {
	ret := _m.Called(c, id)

	if len(ret) == 0 {
		panic("no return value specified for GetByID")
	}

	var r0 user_domain.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, primitive.ObjectID) (user_domain.User, error)); ok {
		return rf(c, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, primitive.ObjectID) user_domain.User); ok {
		r0 = rf(c, id)
	} else {
		r0 = ret.Get(0).(user_domain.User)
	}

	if rf, ok := ret.Get(1).(func(context.Context, primitive.ObjectID) error); ok {
		r1 = rf(c, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetByUsername provides a mock function with given fields: c, username
func (_m *IUserRepository) GetByUsername(c context.Context, username string) (user_domain.User, error) {
	ret := _m.Called(c, username)

	if len(ret) == 0 {
		panic("no return value specified for GetByUsername")
	}

	var r0 user_domain.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (user_domain.User, error)); ok {
		return rf(c, username)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) user_domain.User); ok {
		r0 = rf(c, username)
	} else {
		r0 = ret.Get(0).(user_domain.User)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(c, username)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Update provides a mock function with given fields: c, userID, updatedUser
func (_m *IUserRepository) Update(c context.Context, userID primitive.ObjectID, updatedUser interface{}) error {
	ret := _m.Called(c, userID, updatedUser)

	if len(ret) == 0 {
		panic("no return value specified for Update")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, primitive.ObjectID, interface{}) error); ok {
		r0 = rf(c, userID, updatedUser)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpsertUser provides a mock function with given fields: email, user
func (_m *IUserRepository) UpsertUser(email string, user *user_domain.User) (*user_domain.Response, error) {
	ret := _m.Called(email, user)

	if len(ret) == 0 {
		panic("no return value specified for UpsertUser")
	}

	var r0 *user_domain.Response
	var r1 error
	if rf, ok := ret.Get(0).(func(string, *user_domain.User) (*user_domain.Response, error)); ok {
		return rf(email, user)
	}
	if rf, ok := ret.Get(0).(func(string, *user_domain.User) *user_domain.Response); ok {
		r0 = rf(email, user)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*user_domain.Response)
		}
	}

	if rf, ok := ret.Get(1).(func(string, *user_domain.User) error); ok {
		r1 = rf(email, user)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewIUserRepository creates a new instance of IUserRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewIUserRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *IUserRepository {
	mock := &IUserRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}