package mocks

//type UserRepository struct {
//	mock.Mock
//}
//
//// Create provides a mock function with given fields: c, user
//func (_m *UserRepository) Create(c context.Context, user *user.User) error {
//	ret := _m.Called(c, user)
//
//	var r0 error
//	if rf, ok := ret.Get(0).(func(context.Context, user.) error); ok {
//		r0 = rf(c, user)
//	} else {
//		r0 = ret.Error(0)
//	}
//
//	return r0
//}
//
//// Fetch provides a mock function with given fields: c
//func (_m *UserRepository) Fetch(c context.Context) ([]user.User, error) {
//	ret := _m.Called(c)
//
//	var r0 []domain.domain
//	if rf, ok := ret.Get(0).(func(context.Context) []domain.domain); ok {
//		r0 = rf(c)
//	} else {
//		if ret.Get(0) != nil {
//			r0 = ret.Get(0).([]domain.domain)
//		}
//	}
//
//	var r1 error
//	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
//		r1 = rf(c)
//	} else {
//		r1 = ret.Error(1)
//	}
//
//	return r0, r1
//}
//
//// GetByEmail provides a mock function with given fields: c, email
//func (_m *UserRepository) GetByEmail(c context.Context, email string) (domain.domain, error) {
//	ret := _m.Called(c, email)
//
//	var r0 domain.domain
//	if rf, ok := ret.Get(0).(func(context.Context, string) domain.domain); ok {
//		r0 = rf(c, email)
//	} else {
//		r0 = ret.Get(0).(domain.domain)
//	}
//
//	var r1 error
//	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
//		r1 = rf(c, email)
//	} else {
//		r1 = ret.Error(1)
//	}
//
//	return r0, r1
//}
//
//// GetByID provides a mock function with given fields: c, id
//func (_m *UserRepository) GetByID(c context.Context, id string) (domain.domain, error) {
//	ret := _m.Called(c, id)
//
//	var r0 domain.domain
//	if rf, ok := ret.Get(0).(func(context.Context, string) domain.domain); ok {
//		r0 = rf(c, id)
//	} else {
//		r0 = ret.Get(0).(domain.domain)
//	}
//
//	var r1 error
//	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
//		r1 = rf(c, id)
//	} else {
//		r1 = ret.Error(1)
//	}
//
//	return r0, r1
//}
//
//type mockConstructorTestingTNewUserRepository interface {
//	mock.TestingT
//	Cleanup(func())
//}

//// NewUserRepository creates a new instance of UserRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
//func NewUserRepository(t mockConstructorTestingTNewUserRepository) *UserRepository {
//	mock := &UserRepository{}
//	mock.Mock.Test(t)
//
//	t.Cleanup(func() { mock.AssertExpectations(t) })
//
//	return mock
//}
