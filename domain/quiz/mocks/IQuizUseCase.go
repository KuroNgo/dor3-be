// Code generated by mockery v2.40.1. DO NOT EDIT.

package mocks

import (
	quiz_domain "clean-architecture/domain/quiz"
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// IQuizUseCase is an autogenerated mock type for the IQuizUseCase type
type IQuizUseCase struct {
	mock.Mock
}

// Create provides a mock function with given fields: ctx, quiz
func (_m *IQuizUseCase) Create(ctx context.Context, quiz *quiz_domain.Input) error {
	ret := _m.Called(ctx, quiz)

	if len(ret) == 0 {
		panic("no return value specified for Create")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *quiz_domain.Input) error); ok {
		r0 = rf(ctx, quiz)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Delete provides a mock function with given fields: ctx, quizID
func (_m *IQuizUseCase) Delete(ctx context.Context, quizID string) error {
	ret := _m.Called(ctx, quizID)

	if len(ret) == 0 {
		panic("no return value specified for Delete")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, quizID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Fetch provides a mock function with given fields: ctx
func (_m *IQuizUseCase) Fetch(ctx context.Context) ([]quiz_domain.Quiz, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for Fetch")
	}

	var r0 []quiz_domain.Quiz
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) ([]quiz_domain.Quiz, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) []quiz_domain.Quiz); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]quiz_domain.Quiz)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FetchToDelete provides a mock function with given fields: ctx
func (_m *IQuizUseCase) FetchToDelete(ctx context.Context) (*[]quiz_domain.Quiz, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for FetchToDelete")
	}

	var r0 *[]quiz_domain.Quiz
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) (*[]quiz_domain.Quiz, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) *[]quiz_domain.Quiz); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*[]quiz_domain.Quiz)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Update provides a mock function with given fields: ctx, quizID, quiz
func (_m *IQuizUseCase) Update(ctx context.Context, quizID string, quiz quiz_domain.Quiz) error {
	ret := _m.Called(ctx, quizID, quiz)

	if len(ret) == 0 {
		panic("no return value specified for Update")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, quiz_domain.Quiz) error); ok {
		r0 = rf(ctx, quizID, quiz)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Upsert provides a mock function with given fields: c, question, quiz
func (_m *IQuizUseCase) Upsert(c context.Context, question string, quiz *quiz_domain.Quiz) (*quiz_domain.Response, error) {
	ret := _m.Called(c, question, quiz)

	if len(ret) == 0 {
		panic("no return value specified for Upsert")
	}

	var r0 *quiz_domain.Response
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, *quiz_domain.Quiz) (*quiz_domain.Response, error)); ok {
		return rf(c, question, quiz)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, *quiz_domain.Quiz) *quiz_domain.Response); ok {
		r0 = rf(c, question, quiz)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*quiz_domain.Response)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, *quiz_domain.Quiz) error); ok {
		r1 = rf(c, question, quiz)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewIQuizUseCase creates a new instance of IQuizUseCase. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewIQuizUseCase(t interface {
	mock.TestingT
	Cleanup(func())
}) *IQuizUseCase {
	mock := &IQuizUseCase{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}