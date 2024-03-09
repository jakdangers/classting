// Code generated by mockery v2.32.0. DO NOT EDIT.

package mocks

import (
	domain "classting/domain"
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// SubscriptionRepository is an autogenerated mock type for the SubscriptionRepository type
type SubscriptionRepository struct {
	mock.Mock
}

type SubscriptionRepository_Expecter struct {
	mock *mock.Mock
}

func (_m *SubscriptionRepository) EXPECT() *SubscriptionRepository_Expecter {
	return &SubscriptionRepository_Expecter{mock: &_m.Mock}
}

// CreateSubscription provides a mock function with given fields: ctx, subscription
func (_m *SubscriptionRepository) CreateSubscription(ctx context.Context, subscription domain.Subscription) (int, error) {
	ret := _m.Called(ctx, subscription)

	var r0 int
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, domain.Subscription) (int, error)); ok {
		return rf(ctx, subscription)
	}
	if rf, ok := ret.Get(0).(func(context.Context, domain.Subscription) int); ok {
		r0 = rf(ctx, subscription)
	} else {
		r0 = ret.Get(0).(int)
	}

	if rf, ok := ret.Get(1).(func(context.Context, domain.Subscription) error); ok {
		r1 = rf(ctx, subscription)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SubscriptionRepository_CreateSubscription_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreateSubscription'
type SubscriptionRepository_CreateSubscription_Call struct {
	*mock.Call
}

// CreateSubscription is a helper method to define mock.On call
//   - ctx context.Context
//   - subscription domain.Subscription
func (_e *SubscriptionRepository_Expecter) CreateSubscription(ctx interface{}, subscription interface{}) *SubscriptionRepository_CreateSubscription_Call {
	return &SubscriptionRepository_CreateSubscription_Call{Call: _e.mock.On("CreateSubscription", ctx, subscription)}
}

func (_c *SubscriptionRepository_CreateSubscription_Call) Run(run func(ctx context.Context, subscription domain.Subscription)) *SubscriptionRepository_CreateSubscription_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(domain.Subscription))
	})
	return _c
}

func (_c *SubscriptionRepository_CreateSubscription_Call) Return(_a0 int, _a1 error) *SubscriptionRepository_CreateSubscription_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *SubscriptionRepository_CreateSubscription_Call) RunAndReturn(run func(context.Context, domain.Subscription) (int, error)) *SubscriptionRepository_CreateSubscription_Call {
	_c.Call.Return(run)
	return _c
}

// DeleteSubscription provides a mock function with given fields: ctx, subscriptionID
func (_m *SubscriptionRepository) DeleteSubscription(ctx context.Context, subscriptionID int) error {
	ret := _m.Called(ctx, subscriptionID)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int) error); ok {
		r0 = rf(ctx, subscriptionID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SubscriptionRepository_DeleteSubscription_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DeleteSubscription'
type SubscriptionRepository_DeleteSubscription_Call struct {
	*mock.Call
}

// DeleteSubscription is a helper method to define mock.On call
//   - ctx context.Context
//   - subscriptionID int
func (_e *SubscriptionRepository_Expecter) DeleteSubscription(ctx interface{}, subscriptionID interface{}) *SubscriptionRepository_DeleteSubscription_Call {
	return &SubscriptionRepository_DeleteSubscription_Call{Call: _e.mock.On("DeleteSubscription", ctx, subscriptionID)}
}

func (_c *SubscriptionRepository_DeleteSubscription_Call) Run(run func(ctx context.Context, subscriptionID int)) *SubscriptionRepository_DeleteSubscription_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(int))
	})
	return _c
}

func (_c *SubscriptionRepository_DeleteSubscription_Call) Return(_a0 error) *SubscriptionRepository_DeleteSubscription_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *SubscriptionRepository_DeleteSubscription_Call) RunAndReturn(run func(context.Context, int) error) *SubscriptionRepository_DeleteSubscription_Call {
	_c.Call.Return(run)
	return _c
}

// FindSubscriptionByUserIDAndSchoolID provides a mock function with given fields: ctx, params
func (_m *SubscriptionRepository) FindSubscriptionByUserIDAndSchoolID(ctx context.Context, params domain.FindSubscriptionByUserIDAndSchoolIDParams) (*domain.Subscription, error) {
	ret := _m.Called(ctx, params)

	var r0 *domain.Subscription
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, domain.FindSubscriptionByUserIDAndSchoolIDParams) (*domain.Subscription, error)); ok {
		return rf(ctx, params)
	}
	if rf, ok := ret.Get(0).(func(context.Context, domain.FindSubscriptionByUserIDAndSchoolIDParams) *domain.Subscription); ok {
		r0 = rf(ctx, params)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.Subscription)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, domain.FindSubscriptionByUserIDAndSchoolIDParams) error); ok {
		r1 = rf(ctx, params)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SubscriptionRepository_FindSubscriptionByUserIDAndSchoolID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'FindSubscriptionByUserIDAndSchoolID'
type SubscriptionRepository_FindSubscriptionByUserIDAndSchoolID_Call struct {
	*mock.Call
}

// FindSubscriptionByUserIDAndSchoolID is a helper method to define mock.On call
//   - ctx context.Context
//   - params domain.FindSubscriptionByUserIDAndSchoolIDParams
func (_e *SubscriptionRepository_Expecter) FindSubscriptionByUserIDAndSchoolID(ctx interface{}, params interface{}) *SubscriptionRepository_FindSubscriptionByUserIDAndSchoolID_Call {
	return &SubscriptionRepository_FindSubscriptionByUserIDAndSchoolID_Call{Call: _e.mock.On("FindSubscriptionByUserIDAndSchoolID", ctx, params)}
}

func (_c *SubscriptionRepository_FindSubscriptionByUserIDAndSchoolID_Call) Run(run func(ctx context.Context, params domain.FindSubscriptionByUserIDAndSchoolIDParams)) *SubscriptionRepository_FindSubscriptionByUserIDAndSchoolID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(domain.FindSubscriptionByUserIDAndSchoolIDParams))
	})
	return _c
}

func (_c *SubscriptionRepository_FindSubscriptionByUserIDAndSchoolID_Call) Return(_a0 *domain.Subscription, _a1 error) *SubscriptionRepository_FindSubscriptionByUserIDAndSchoolID_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *SubscriptionRepository_FindSubscriptionByUserIDAndSchoolID_Call) RunAndReturn(run func(context.Context, domain.FindSubscriptionByUserIDAndSchoolIDParams) (*domain.Subscription, error)) *SubscriptionRepository_FindSubscriptionByUserIDAndSchoolID_Call {
	_c.Call.Return(run)
	return _c
}

// ListSubscriptionSchools provides a mock function with given fields: ctx, params
func (_m *SubscriptionRepository) ListSubscriptionSchools(ctx context.Context, params domain.ListSubscriptionSchoolsParams) ([]domain.SubscriptionSchool, error) {
	ret := _m.Called(ctx, params)

	var r0 []domain.SubscriptionSchool
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, domain.ListSubscriptionSchoolsParams) ([]domain.SubscriptionSchool, error)); ok {
		return rf(ctx, params)
	}
	if rf, ok := ret.Get(0).(func(context.Context, domain.ListSubscriptionSchoolsParams) []domain.SubscriptionSchool); ok {
		r0 = rf(ctx, params)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]domain.SubscriptionSchool)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, domain.ListSubscriptionSchoolsParams) error); ok {
		r1 = rf(ctx, params)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SubscriptionRepository_ListSubscriptionSchools_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ListSubscriptionSchools'
type SubscriptionRepository_ListSubscriptionSchools_Call struct {
	*mock.Call
}

// ListSubscriptionSchools is a helper method to define mock.On call
//   - ctx context.Context
//   - params domain.ListSubscriptionSchoolsParams
func (_e *SubscriptionRepository_Expecter) ListSubscriptionSchools(ctx interface{}, params interface{}) *SubscriptionRepository_ListSubscriptionSchools_Call {
	return &SubscriptionRepository_ListSubscriptionSchools_Call{Call: _e.mock.On("ListSubscriptionSchools", ctx, params)}
}

func (_c *SubscriptionRepository_ListSubscriptionSchools_Call) Run(run func(ctx context.Context, params domain.ListSubscriptionSchoolsParams)) *SubscriptionRepository_ListSubscriptionSchools_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(domain.ListSubscriptionSchoolsParams))
	})
	return _c
}

func (_c *SubscriptionRepository_ListSubscriptionSchools_Call) Return(_a0 []domain.SubscriptionSchool, _a1 error) *SubscriptionRepository_ListSubscriptionSchools_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *SubscriptionRepository_ListSubscriptionSchools_Call) RunAndReturn(run func(context.Context, domain.ListSubscriptionSchoolsParams) ([]domain.SubscriptionSchool, error)) *SubscriptionRepository_ListSubscriptionSchools_Call {
	_c.Call.Return(run)
	return _c
}

// NewSubscriptionRepository creates a new instance of SubscriptionRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewSubscriptionRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *SubscriptionRepository {
	mock := &SubscriptionRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}