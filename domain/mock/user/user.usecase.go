package mocks

//type UserUseCase struct {
//	mock.Mock
//}
//
//// CreateAccessToken provides a mock function with given fields: user, secret, expiry
//func (_m *UserUseCase) CreateAccessToken(user *user.User, secret string, expiry int) (string, error) {
//	ret := _m.Called(user, secret, expiry)
//
//	var r0 string
//	if rf, ok := ret.Get(0).(func(*user.User, string, int) string); ok {
//		r0 = rf(user, secret, expiry)
//	} else {
//		r0 = ret.Get(0).(string)
//	}
//
//	var r1 error
//	if rf, ok := ret.Get(1).(func(*user.User, string, int) error); ok {
//		r1 = rf(user, secret, expiry)
//	} else {
//		r1 = ret.Error(1)
//	}
//
//	return r0, r1
//}
//
//// CreateRefreshToken provides a mock function with given fields: user, secret, expiry
//func (_m *UserUseCase) CreateRefreshToken(user *user.User, secret string, expiry int) (string, error) {
//	ret := _m.Called(user, secret, expiry)
//
//	var r0 string
//	if rf, ok := ret.Get(0).(func(*user.User, string, int) string); ok {
//		r0 = rf(user, secret, expiry)
//	} else {
//		r0 = ret.Get(0).(string)
//	}
//
//	var r1 error
//	if rf, ok := ret.Get(1).(func(*user.User, string, int) error); ok {
//		r1 = rf(user, secret, expiry)
//	} else {
//		r1 = ret.Error(1)
//	}
//
//	return r0, r1
//}
//
//// GetUserByEmail provides a mock function with given fields: c, email
//func (_m *UserUseCase) GetUserByEmail(c context.Context, email string) (user.User, error) {
//	ret := _m.Called(c, email)
//
//	var r0 user.User
//	if rf, ok := ret.Get(0).(func(context.Context, string) user.User); ok {
//		r0 = rf(c, email)
//	} else {
//		r0 = ret.Get(0).(user.User)
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
//type mockConstructorTestingTNewLoginUsecase interface {
//	mock.TestingT
//	Cleanup(func())
//}
//
//// NewUserUseCase creates a new instance of LoginUsecase. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
//func NewUserUseCase(t mockConstructorTestingTNewLoginUsecase) *UserUseCase {
//	mock := &UserUseCase{}
//	mock.Mock.Test(t)
//
//	t.Cleanup(func() { mock.AssertExpectations(t) })
//
//	return mock
//}
