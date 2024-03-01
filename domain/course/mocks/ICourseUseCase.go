// Code generated by mockery v2.40.1. DO NOT EDIT.

package mocks

import (
	course_domain "clean-architecture/domain/course"
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// ICourseUseCase is an autogenerated mock type for the ICourseUseCase type
type ICourseUseCase struct {
	mock.Mock
}

// CreateOne provides a mock function with given fields: ctx, course
func (_m *ICourseUseCase) CreateOne(ctx context.Context, course *course_domain.Course) error {
	ret := _m.Called(ctx, course)

	if len(ret) == 0 {
		panic("no return value specified for CreateOne")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *course_domain.Course) error); ok {
		r0 = rf(ctx, course)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteOne provides a mock function with given fields: ctx, courseID
func (_m *ICourseUseCase) DeleteOne(ctx context.Context, courseID string) error {
	ret := _m.Called(ctx, courseID)

	if len(ret) == 0 {
		panic("no return value specified for DeleteOne")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, courseID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// FetchByID provides a mock function with given fields: ctx, courseID
func (_m *ICourseUseCase) FetchByID(ctx context.Context, courseID string) (*course_domain.Course, error) {
	ret := _m.Called(ctx, courseID)

	if len(ret) == 0 {
		panic("no return value specified for FetchByID")
	}

	var r0 *course_domain.Course
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*course_domain.Course, error)); ok {
		return rf(ctx, courseID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *course_domain.Course); ok {
		r0 = rf(ctx, courseID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*course_domain.Course)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, courseID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FetchMany provides a mock function with given fields: ctx
func (_m *ICourseUseCase) FetchMany(ctx context.Context) ([]course_domain.Course, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for FetchMany")
	}

	var r0 []course_domain.Course
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) ([]course_domain.Course, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) []course_domain.Course); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]course_domain.Course)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FetchToDeleteMany provides a mock function with given fields: ctx
func (_m *ICourseUseCase) FetchToDeleteMany(ctx context.Context) (*[]course_domain.Course, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for FetchToDeleteMany")
	}

	var r0 *[]course_domain.Course
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) (*[]course_domain.Course, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) *[]course_domain.Course); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*[]course_domain.Course)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateOne provides a mock function with given fields: ctx, courseID, course
func (_m *ICourseUseCase) UpdateOne(ctx context.Context, courseID string, course course_domain.Course) error {
	ret := _m.Called(ctx, courseID, course)

	if len(ret) == 0 {
		panic("no return value specified for UpdateOne")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, course_domain.Course) error); ok {
		r0 = rf(ctx, courseID, course)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpsertOne provides a mock function with given fields: ctx, id, course
func (_m *ICourseUseCase) UpsertOne(ctx context.Context, id string, course *course_domain.Course) (*course_domain.Response, error) {
	ret := _m.Called(ctx, id, course)

	if len(ret) == 0 {
		panic("no return value specified for UpsertOne")
	}

	var r0 *course_domain.Response
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, *course_domain.Course) (*course_domain.Response, error)); ok {
		return rf(ctx, id, course)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, *course_domain.Course) *course_domain.Response); ok {
		r0 = rf(ctx, id, course)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*course_domain.Response)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, *course_domain.Course) error); ok {
		r1 = rf(ctx, id, course)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewICourseUseCase creates a new instance of ICourseUseCase. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewICourseUseCase(t interface {
	mock.TestingT
	Cleanup(func())
}) *ICourseUseCase {
	mock := &ICourseUseCase{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
