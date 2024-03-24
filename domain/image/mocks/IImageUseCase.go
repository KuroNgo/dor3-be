// Code generated by mockery v2.40.1. DO NOT EDIT.

package mocks

import (
	image_domain "clean-architecture/domain/image"
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// IImageUseCase is an autogenerated mock type for the IImageUseCase type
type IImageUseCase struct {
	mock.Mock
}

// CreateMany provides a mock function with given fields: ctx, image
func (_m *IImageUseCase) CreateMany(ctx context.Context, image []*image_domain.Image) error {
	ret := _m.Called(ctx, image)

	if len(ret) == 0 {
		panic("no return value specified for CreateMany")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, []*image_domain.Image) error); ok {
		r0 = rf(ctx, image)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CreateOne provides a mock function with given fields: ctx, image
func (_m *IImageUseCase) CreateOne(ctx context.Context, image *image_domain.Image) error {
	ret := _m.Called(ctx, image)

	if len(ret) == 0 {
		panic("no return value specified for CreateOne")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *image_domain.Image) error); ok {
		r0 = rf(ctx, image)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteMany provides a mock function with given fields: ctx, imageID
func (_m *IImageUseCase) DeleteMany(ctx context.Context, imageID ...string) error {
	_va := make([]interface{}, len(imageID))
	for _i := range imageID {
		_va[_i] = imageID[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for DeleteMany")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, ...string) error); ok {
		r0 = rf(ctx, imageID...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteOne provides a mock function with given fields: ctx, imageID
func (_m *IImageUseCase) DeleteOne(ctx context.Context, imageID string) error {
	ret := _m.Called(ctx, imageID)

	if len(ret) == 0 {
		panic("no return value specified for DeleteOne")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, imageID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// FetchMany provides a mock function with given fields: ctx
func (_m *IImageUseCase) FetchMany(ctx context.Context) ([]image_domain.Image, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for FetchMany")
	}

	var r0 []image_domain.Image
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) ([]image_domain.Image, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) []image_domain.Image); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]image_domain.Image)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetURLByName provides a mock function with given fields: ctx, name
func (_m *IImageUseCase) GetURLByName(ctx context.Context, name string) (image_domain.Image, error) {
	ret := _m.Called(ctx, name)

	if len(ret) == 0 {
		panic("no return value specified for GetURLByName")
	}

	var r0 image_domain.Image
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (image_domain.Image, error)); ok {
		return rf(ctx, name)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) image_domain.Image); ok {
		r0 = rf(ctx, name)
	} else {
		r0 = ret.Get(0).(image_domain.Image)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateOne provides a mock function with given fields: ctx, imageID, image
func (_m *IImageUseCase) UpdateOne(ctx context.Context, imageID string, image image_domain.Image) error {
	ret := _m.Called(ctx, imageID, image)

	if len(ret) == 0 {
		panic("no return value specified for UpdateOne")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, image_domain.Image) error); ok {
		r0 = rf(ctx, imageID, image)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewIImageUseCase creates a new instance of IImageUseCase. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewIImageUseCase(t interface {
	mock.TestingT
	Cleanup(func())
}) *IImageUseCase {
	mock := &IImageUseCase{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}