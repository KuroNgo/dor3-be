// Code generated by mockery v2.40.1. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
	mongo "go.mongodb.org/mongo-driver/mongo"

	primitive "go.mongodb.org/mongo-driver/bson/primitive"

	unit_domain "clean-architecture/domain/unit"
)

// IUnitUseCase is an autogenerated mock type for the IUnitUseCase type
type IUnitUseCase struct {
	mock.Mock
}

// CreateOneByNameLessonInAdmin provides a mock function with given fields: ctx, unit
func (_m *IUnitUseCase) CreateOneByNameLessonInAdmin(ctx context.Context, unit *unit_domain.Unit) error {
	ret := _m.Called(ctx, unit)

	if len(ret) == 0 {
		panic("no return value specified for CreateOneByNameLessonInAdmin")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *unit_domain.Unit) error); ok {
		r0 = rf(ctx, unit)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CreateOneInAdmin provides a mock function with given fields: ctx, unit
func (_m *IUnitUseCase) CreateOneInAdmin(ctx context.Context, unit *unit_domain.Unit) error {
	ret := _m.Called(ctx, unit)

	if len(ret) == 0 {
		panic("no return value specified for CreateOneInAdmin")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *unit_domain.Unit) error); ok {
		r0 = rf(ctx, unit)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteOneInAdmin provides a mock function with given fields: ctx, unitID
func (_m *IUnitUseCase) DeleteOneInAdmin(ctx context.Context, unitID string) error {
	ret := _m.Called(ctx, unitID)

	if len(ret) == 0 {
		panic("no return value specified for DeleteOneInAdmin")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, unitID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// FetchByIdLessonInAdmin provides a mock function with given fields: ctx, idLesson, page
func (_m *IUnitUseCase) FetchByIdLessonInAdmin(ctx context.Context, idLesson string, page string) ([]unit_domain.UnitResponse, unit_domain.DetailResponse, error) {
	ret := _m.Called(ctx, idLesson, page)

	if len(ret) == 0 {
		panic("no return value specified for FetchByIdLessonInAdmin")
	}

	var r0 []unit_domain.UnitResponse
	var r1 unit_domain.DetailResponse
	var r2 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) ([]unit_domain.UnitResponse, unit_domain.DetailResponse, error)); ok {
		return rf(ctx, idLesson, page)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string) []unit_domain.UnitResponse); ok {
		r0 = rf(ctx, idLesson, page)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]unit_domain.UnitResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string) unit_domain.DetailResponse); ok {
		r1 = rf(ctx, idLesson, page)
	} else {
		r1 = ret.Get(1).(unit_domain.DetailResponse)
	}

	if rf, ok := ret.Get(2).(func(context.Context, string, string) error); ok {
		r2 = rf(ctx, idLesson, page)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// FetchByIdLessonInUser provides a mock function with given fields: ctx, user, idLesson, page
func (_m *IUnitUseCase) FetchByIdLessonInUser(ctx context.Context, user primitive.ObjectID, idLesson string, page string) ([]unit_domain.UnitProcessResponse, unit_domain.DetailResponse, error) {
	ret := _m.Called(ctx, user, idLesson, page)

	if len(ret) == 0 {
		panic("no return value specified for FetchByIdLessonInUser")
	}

	var r0 []unit_domain.UnitProcessResponse
	var r1 unit_domain.DetailResponse
	var r2 error
	if rf, ok := ret.Get(0).(func(context.Context, primitive.ObjectID, string, string) ([]unit_domain.UnitProcessResponse, unit_domain.DetailResponse, error)); ok {
		return rf(ctx, user, idLesson, page)
	}
	if rf, ok := ret.Get(0).(func(context.Context, primitive.ObjectID, string, string) []unit_domain.UnitProcessResponse); ok {
		r0 = rf(ctx, user, idLesson, page)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]unit_domain.UnitProcessResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, primitive.ObjectID, string, string) unit_domain.DetailResponse); ok {
		r1 = rf(ctx, user, idLesson, page)
	} else {
		r1 = ret.Get(1).(unit_domain.DetailResponse)
	}

	if rf, ok := ret.Get(2).(func(context.Context, primitive.ObjectID, string, string) error); ok {
		r2 = rf(ctx, user, idLesson, page)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// FetchManyInAdmin provides a mock function with given fields: ctx, page
func (_m *IUnitUseCase) FetchManyInAdmin(ctx context.Context, page string) ([]unit_domain.UnitResponse, unit_domain.DetailResponse, error) {
	ret := _m.Called(ctx, page)

	if len(ret) == 0 {
		panic("no return value specified for FetchManyInAdmin")
	}

	var r0 []unit_domain.UnitResponse
	var r1 unit_domain.DetailResponse
	var r2 error
	if rf, ok := ret.Get(0).(func(context.Context, string) ([]unit_domain.UnitResponse, unit_domain.DetailResponse, error)); ok {
		return rf(ctx, page)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) []unit_domain.UnitResponse); ok {
		r0 = rf(ctx, page)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]unit_domain.UnitResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) unit_domain.DetailResponse); ok {
		r1 = rf(ctx, page)
	} else {
		r1 = ret.Get(1).(unit_domain.DetailResponse)
	}

	if rf, ok := ret.Get(2).(func(context.Context, string) error); ok {
		r2 = rf(ctx, page)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// FetchManyInUser provides a mock function with given fields: ctx, user, page
func (_m *IUnitUseCase) FetchManyInUser(ctx context.Context, user primitive.ObjectID, page string) ([]unit_domain.UnitProcessResponse, unit_domain.DetailResponse, error) {
	ret := _m.Called(ctx, user, page)

	if len(ret) == 0 {
		panic("no return value specified for FetchManyInUser")
	}

	var r0 []unit_domain.UnitProcessResponse
	var r1 unit_domain.DetailResponse
	var r2 error
	if rf, ok := ret.Get(0).(func(context.Context, primitive.ObjectID, string) ([]unit_domain.UnitProcessResponse, unit_domain.DetailResponse, error)); ok {
		return rf(ctx, user, page)
	}
	if rf, ok := ret.Get(0).(func(context.Context, primitive.ObjectID, string) []unit_domain.UnitProcessResponse); ok {
		r0 = rf(ctx, user, page)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]unit_domain.UnitProcessResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, primitive.ObjectID, string) unit_domain.DetailResponse); ok {
		r1 = rf(ctx, user, page)
	} else {
		r1 = ret.Get(1).(unit_domain.DetailResponse)
	}

	if rf, ok := ret.Get(2).(func(context.Context, primitive.ObjectID, string) error); ok {
		r2 = rf(ctx, user, page)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// FetchManyNotPaginationInAdmin provides a mock function with given fields: ctx
func (_m *IUnitUseCase) FetchManyNotPaginationInAdmin(ctx context.Context) ([]unit_domain.UnitResponse, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for FetchManyNotPaginationInAdmin")
	}

	var r0 []unit_domain.UnitResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) ([]unit_domain.UnitResponse, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) []unit_domain.UnitResponse); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]unit_domain.UnitResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FetchManyNotPaginationInUser provides a mock function with given fields: ctx, user
func (_m *IUnitUseCase) FetchManyNotPaginationInUser(ctx context.Context, user primitive.ObjectID) ([]unit_domain.UnitProcessResponse, error) {
	ret := _m.Called(ctx, user)

	if len(ret) == 0 {
		panic("no return value specified for FetchManyNotPaginationInUser")
	}

	var r0 []unit_domain.UnitProcessResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, primitive.ObjectID) ([]unit_domain.UnitProcessResponse, error)); ok {
		return rf(ctx, user)
	}
	if rf, ok := ret.Get(0).(func(context.Context, primitive.ObjectID) []unit_domain.UnitProcessResponse); ok {
		r0 = rf(ctx, user)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]unit_domain.UnitProcessResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, primitive.ObjectID) error); ok {
		r1 = rf(ctx, user)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FetchOneByIDInAdmin provides a mock function with given fields: ctx, id
func (_m *IUnitUseCase) FetchOneByIDInAdmin(ctx context.Context, id string) (unit_domain.UnitResponse, error) {
	ret := _m.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for FetchOneByIDInAdmin")
	}

	var r0 unit_domain.UnitResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (unit_domain.UnitResponse, error)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) unit_domain.UnitResponse); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Get(0).(unit_domain.UnitResponse)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FetchOneByIDInUser provides a mock function with given fields: ctx, user, id
func (_m *IUnitUseCase) FetchOneByIDInUser(ctx context.Context, user primitive.ObjectID, id string) (unit_domain.UnitProcessResponse, error) {
	ret := _m.Called(ctx, user, id)

	if len(ret) == 0 {
		panic("no return value specified for FetchOneByIDInUser")
	}

	var r0 unit_domain.UnitProcessResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, primitive.ObjectID, string) (unit_domain.UnitProcessResponse, error)); ok {
		return rf(ctx, user, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, primitive.ObjectID, string) unit_domain.UnitProcessResponse); ok {
		r0 = rf(ctx, user, id)
	} else {
		r0 = ret.Get(0).(unit_domain.UnitProcessResponse)
	}

	if rf, ok := ret.Get(1).(func(context.Context, primitive.ObjectID, string) error); ok {
		r1 = rf(ctx, user, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindUnitIDByUnitLevelInAdmin provides a mock function with given fields: ctx, unitLevel, fieldOfIT
func (_m *IUnitUseCase) FindUnitIDByUnitLevelInAdmin(ctx context.Context, unitLevel int, fieldOfIT string) (primitive.ObjectID, error) {
	ret := _m.Called(ctx, unitLevel, fieldOfIT)

	if len(ret) == 0 {
		panic("no return value specified for FindUnitIDByUnitLevelInAdmin")
	}

	var r0 primitive.ObjectID
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int, string) (primitive.ObjectID, error)); ok {
		return rf(ctx, unitLevel, fieldOfIT)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int, string) primitive.ObjectID); ok {
		r0 = rf(ctx, unitLevel, fieldOfIT)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(primitive.ObjectID)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, int, string) error); ok {
		r1 = rf(ctx, unitLevel, fieldOfIT)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateCompleteInUser provides a mock function with given fields: ctx, user
func (_m *IUnitUseCase) UpdateCompleteInUser(ctx context.Context, user primitive.ObjectID) (*mongo.UpdateResult, error) {
	ret := _m.Called(ctx, user)

	if len(ret) == 0 {
		panic("no return value specified for UpdateCompleteInUser")
	}

	var r0 *mongo.UpdateResult
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, primitive.ObjectID) (*mongo.UpdateResult, error)); ok {
		return rf(ctx, user)
	}
	if rf, ok := ret.Get(0).(func(context.Context, primitive.ObjectID) *mongo.UpdateResult); ok {
		r0 = rf(ctx, user)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*mongo.UpdateResult)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, primitive.ObjectID) error); ok {
		r1 = rf(ctx, user)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateOneInAdmin provides a mock function with given fields: ctx, unit
func (_m *IUnitUseCase) UpdateOneInAdmin(ctx context.Context, unit *unit_domain.Unit) (*mongo.UpdateResult, error) {
	ret := _m.Called(ctx, unit)

	if len(ret) == 0 {
		panic("no return value specified for UpdateOneInAdmin")
	}

	var r0 *mongo.UpdateResult
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *unit_domain.Unit) (*mongo.UpdateResult, error)); ok {
		return rf(ctx, unit)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *unit_domain.Unit) *mongo.UpdateResult); ok {
		r0 = rf(ctx, unit)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*mongo.UpdateResult)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *unit_domain.Unit) error); ok {
		r1 = rf(ctx, unit)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewIUnitUseCase creates a new instance of IUnitUseCase. It also registers a testing interface on the mock and a cleanup function to assert the test expectations.
// The first argument is typically a *testing.T value.
func NewIUnitUseCase(t interface {
	mock.TestingT
	Cleanup(func())
}) *IUnitUseCase {
	mock := &IUnitUseCase{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}