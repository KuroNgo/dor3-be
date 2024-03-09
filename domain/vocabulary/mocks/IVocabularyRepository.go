// Code generated by mockery v2.40.1. DO NOT EDIT.

package mocks

import (
	vocabulary_domain "clean-architecture/domain/vocabulary"
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// IVocabularyRepository is an autogenerated mock type for the IVocabularyRepository type
type IVocabularyRepository struct {
	mock.Mock
}

// CreateOne provides a mock function with given fields: ctx, vocabulary
func (_m *IVocabularyRepository) CreateOne(ctx context.Context, vocabulary *vocabulary_domain.Vocabulary) error {
	ret := _m.Called(ctx, vocabulary)

	if len(ret) == 0 {
		panic("no return value specified for CreateOne")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *vocabulary_domain.Vocabulary) error); ok {
		r0 = rf(ctx, vocabulary)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteOne provides a mock function with given fields: ctx, vocabularyID
func (_m *IVocabularyRepository) DeleteOne(ctx context.Context, vocabularyID string) error {
	ret := _m.Called(ctx, vocabularyID)

	if len(ret) == 0 {
		panic("no return value specified for DeleteOne")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, vocabularyID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// FetchByID provides a mock function with given fields: ctx, vocabularyID
func (_m *IVocabularyRepository) FetchByID(ctx context.Context, vocabularyID string) (*vocabulary_domain.Vocabulary, error) {
	ret := _m.Called(ctx, vocabularyID)

	if len(ret) == 0 {
		panic("no return value specified for FetchByID")
	}

	var r0 *vocabulary_domain.Vocabulary
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*vocabulary_domain.Vocabulary, error)); ok {
		return rf(ctx, vocabularyID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *vocabulary_domain.Vocabulary); ok {
		r0 = rf(ctx, vocabularyID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*vocabulary_domain.Vocabulary)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, vocabularyID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FetchMany provides a mock function with given fields: ctx
func (_m *IVocabularyRepository) FetchMany(ctx context.Context) ([]vocabulary_domain.Vocabulary, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for FetchMany")
	}

	var r0 []vocabulary_domain.Vocabulary
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) ([]vocabulary_domain.Vocabulary, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) []vocabulary_domain.Vocabulary); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]vocabulary_domain.Vocabulary)
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
func (_m *IVocabularyRepository) FetchToDeleteMany(ctx context.Context) (*[]vocabulary_domain.Vocabulary, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for FetchToDeleteMany")
	}

	var r0 *[]vocabulary_domain.Vocabulary
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) (*[]vocabulary_domain.Vocabulary, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) *[]vocabulary_domain.Vocabulary); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*[]vocabulary_domain.Vocabulary)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateOne provides a mock function with given fields: ctx, vocabularyID, vocabulary
func (_m *IVocabularyRepository) UpdateOne(ctx context.Context, vocabularyID string, vocabulary vocabulary_domain.Vocabulary) error {
	ret := _m.Called(ctx, vocabularyID, vocabulary)

	if len(ret) == 0 {
		panic("no return value specified for UpdateOne")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, vocabulary_domain.Vocabulary) error); ok {
		r0 = rf(ctx, vocabularyID, vocabulary)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpsertOne provides a mock function with given fields: c, id, vocabulary
func (_m *IVocabularyRepository) UpsertOne(c context.Context, id string, vocabulary *vocabulary_domain.Vocabulary) (*vocabulary_domain.Vocabulary, error) {
	ret := _m.Called(c, id, vocabulary)

	if len(ret) == 0 {
		panic("no return value specified for UpsertOne")
	}

	var r0 *vocabulary_domain.Vocabulary
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, *vocabulary_domain.Vocabulary) (*vocabulary_domain.Vocabulary, error)); ok {
		return rf(c, id, vocabulary)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, *vocabulary_domain.Vocabulary) *vocabulary_domain.Vocabulary); ok {
		r0 = rf(c, id, vocabulary)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*vocabulary_domain.Vocabulary)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, *vocabulary_domain.Vocabulary) error); ok {
		r1 = rf(c, id, vocabulary)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewIVocabularyRepository creates a new instance of IVocabularyRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewIVocabularyRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *IVocabularyRepository {
	mock := &IVocabularyRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
