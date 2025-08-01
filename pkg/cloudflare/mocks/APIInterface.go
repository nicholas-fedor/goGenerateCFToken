// Code generated by mockery v2.53.3. DO NOT EDIT.

package mocks

import (
	context "context"

	pagination "github.com/cloudflare/cloudflare-go/v5/packages/pagination"
	mock "github.com/stretchr/testify/mock"

	user "github.com/cloudflare/cloudflare-go/v5/user"

	zones "github.com/cloudflare/cloudflare-go/v5/zones"
)

// MockAPIInterface is an autogenerated mock type for the APIInterface type
type MockAPIInterface struct {
	mock.Mock
}

type MockAPIInterface_Expecter struct {
	mock *mock.Mock
}

func (_m *MockAPIInterface) EXPECT() *MockAPIInterface_Expecter {
	return &MockAPIInterface_Expecter{mock: &_m.Mock}
}

// CreateAPIToken provides a mock function with given fields: ctx, params
func (_m *MockAPIInterface) CreateAPIToken(ctx context.Context, params user.TokenNewParams) (*user.TokenNewResponse, error) {
	ret := _m.Called(ctx, params)

	if len(ret) == 0 {
		panic("no return value specified for CreateAPIToken")
	}

	var r0 *user.TokenNewResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, user.TokenNewParams) (*user.TokenNewResponse, error)); ok {
		return rf(ctx, params)
	}
	if rf, ok := ret.Get(0).(func(context.Context, user.TokenNewParams) *user.TokenNewResponse); ok {
		r0 = rf(ctx, params)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*user.TokenNewResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, user.TokenNewParams) error); ok {
		r1 = rf(ctx, params)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockAPIInterface_CreateAPIToken_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreateAPIToken'
type MockAPIInterface_CreateAPIToken_Call struct {
	*mock.Call
}

// CreateAPIToken is a helper method to define mock.On call
//   - ctx context.Context
//   - params user.TokenNewParams
func (_e *MockAPIInterface_Expecter) CreateAPIToken(ctx interface{}, params interface{}) *MockAPIInterface_CreateAPIToken_Call {
	return &MockAPIInterface_CreateAPIToken_Call{Call: _e.mock.On("CreateAPIToken", ctx, params)}
}

func (_c *MockAPIInterface_CreateAPIToken_Call) Run(run func(ctx context.Context, params user.TokenNewParams)) *MockAPIInterface_CreateAPIToken_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(user.TokenNewParams))
	})
	return _c
}

func (_c *MockAPIInterface_CreateAPIToken_Call) Return(_a0 *user.TokenNewResponse, _a1 error) *MockAPIInterface_CreateAPIToken_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockAPIInterface_CreateAPIToken_Call) RunAndReturn(run func(context.Context, user.TokenNewParams) (*user.TokenNewResponse, error)) *MockAPIInterface_CreateAPIToken_Call {
	_c.Call.Return(run)
	return _c
}

// ListZones provides a mock function with given fields: ctx, params
func (_m *MockAPIInterface) ListZones(ctx context.Context, params zones.ZoneListParams) (*pagination.V4PagePaginationArray[zones.Zone], error) {
	ret := _m.Called(ctx, params)

	if len(ret) == 0 {
		panic("no return value specified for ListZones")
	}

	var r0 *pagination.V4PagePaginationArray[zones.Zone]
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, zones.ZoneListParams) (*pagination.V4PagePaginationArray[zones.Zone], error)); ok {
		return rf(ctx, params)
	}
	if rf, ok := ret.Get(0).(func(context.Context, zones.ZoneListParams) *pagination.V4PagePaginationArray[zones.Zone]); ok {
		r0 = rf(ctx, params)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*pagination.V4PagePaginationArray[zones.Zone])
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, zones.ZoneListParams) error); ok {
		r1 = rf(ctx, params)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockAPIInterface_ListZones_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ListZones'
type MockAPIInterface_ListZones_Call struct {
	*mock.Call
}

// ListZones is a helper method to define mock.On call
//   - ctx context.Context
//   - params zones.ZoneListParams
func (_e *MockAPIInterface_Expecter) ListZones(ctx interface{}, params interface{}) *MockAPIInterface_ListZones_Call {
	return &MockAPIInterface_ListZones_Call{Call: _e.mock.On("ListZones", ctx, params)}
}

func (_c *MockAPIInterface_ListZones_Call) Run(run func(ctx context.Context, params zones.ZoneListParams)) *MockAPIInterface_ListZones_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(zones.ZoneListParams))
	})
	return _c
}

func (_c *MockAPIInterface_ListZones_Call) Return(_a0 *pagination.V4PagePaginationArray[zones.Zone], _a1 error) *MockAPIInterface_ListZones_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockAPIInterface_ListZones_Call) RunAndReturn(run func(context.Context, zones.ZoneListParams) (*pagination.V4PagePaginationArray[zones.Zone], error)) *MockAPIInterface_ListZones_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockAPIInterface creates a new instance of MockAPIInterface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockAPIInterface(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockAPIInterface {
	mock := &MockAPIInterface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
