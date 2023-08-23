// Code generated by mockery v2.32.4. DO NOT EDIT.

package fake

import (
	context "context"
	http "net/http"

	mock "github.com/stretchr/testify/mock"
)

// GoryaServiceHandler is an autogenerated mock type for the GoryaServiceHandler type
type GoryaServiceHandler struct {
	mock.Mock
}

// AddPolicy provides a mock function with given fields: ctx
func (_m *GoryaServiceHandler) AddPolicy(ctx context.Context) http.Handler {
	ret := _m.Called(ctx)

	var r0 http.Handler
	if rf, ok := ret.Get(0).(func(context.Context) http.Handler); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(http.Handler)
		}
	}

	return r0
}

// AddSchedule provides a mock function with given fields: ctx
func (_m *GoryaServiceHandler) AddSchedule(ctx context.Context) http.Handler {
	ret := _m.Called(ctx)

	var r0 http.Handler
	if rf, ok := ret.Get(0).(func(context.Context) http.Handler); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(http.Handler)
		}
	}

	return r0
}

// ChangeState provides a mock function with given fields: ctx
func (_m *GoryaServiceHandler) ChangeState(ctx context.Context) http.Handler {
	ret := _m.Called(ctx)

	var r0 http.Handler
	if rf, ok := ret.Get(0).(func(context.Context) http.Handler); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(http.Handler)
		}
	}

	return r0
}

// DeletePolicy provides a mock function with given fields: ctx
func (_m *GoryaServiceHandler) DeletePolicy(ctx context.Context) http.Handler {
	ret := _m.Called(ctx)

	var r0 http.Handler
	if rf, ok := ret.Get(0).(func(context.Context) http.Handler); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(http.Handler)
		}
	}

	return r0
}

// DeleteSchedule provides a mock function with given fields: ctx
func (_m *GoryaServiceHandler) DeleteSchedule(ctx context.Context) http.Handler {
	ret := _m.Called(ctx)

	var r0 http.Handler
	if rf, ok := ret.Get(0).(func(context.Context) http.Handler); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(http.Handler)
		}
	}

	return r0
}

// GetPolicy provides a mock function with given fields: ctx
func (_m *GoryaServiceHandler) GetPolicy(ctx context.Context) http.Handler {
	ret := _m.Called(ctx)

	var r0 http.Handler
	if rf, ok := ret.Get(0).(func(context.Context) http.Handler); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(http.Handler)
		}
	}

	return r0
}

// GetSchedule provides a mock function with given fields: ctx
func (_m *GoryaServiceHandler) GetSchedule(ctx context.Context) http.Handler {
	ret := _m.Called(ctx)

	var r0 http.Handler
	if rf, ok := ret.Get(0).(func(context.Context) http.Handler); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(http.Handler)
		}
	}

	return r0
}

// GetTimeZone provides a mock function with given fields:
func (_m *GoryaServiceHandler) GetTimeZone() http.Handler {
	ret := _m.Called()

	var r0 http.Handler
	if rf, ok := ret.Get(0).(func() http.Handler); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(http.Handler)
		}
	}

	return r0
}

// GetVersionInfo provides a mock function with given fields:
func (_m *GoryaServiceHandler) GetVersionInfo() http.Handler {
	ret := _m.Called()

	var r0 http.Handler
	if rf, ok := ret.Get(0).(func() http.Handler); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(http.Handler)
		}
	}

	return r0
}

// ListPolicy provides a mock function with given fields: ctx
func (_m *GoryaServiceHandler) ListPolicy(ctx context.Context) http.Handler {
	ret := _m.Called(ctx)

	var r0 http.Handler
	if rf, ok := ret.Get(0).(func(context.Context) http.Handler); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(http.Handler)
		}
	}

	return r0
}

// ListSchedule provides a mock function with given fields: ctx
func (_m *GoryaServiceHandler) ListSchedule(ctx context.Context) http.Handler {
	ret := _m.Called(ctx)

	var r0 http.Handler
	if rf, ok := ret.Get(0).(func(context.Context) http.Handler); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(http.Handler)
		}
	}

	return r0
}

// ScheduleTask provides a mock function with given fields: ctx
func (_m *GoryaServiceHandler) ScheduleTask(ctx context.Context) http.Handler {
	ret := _m.Called(ctx)

	var r0 http.Handler
	if rf, ok := ret.Get(0).(func(context.Context) http.Handler); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(http.Handler)
		}
	}

	return r0
}

// NewGoryaServiceHandler creates a new instance of GoryaServiceHandler. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewGoryaServiceHandler(t interface {
	mock.TestingT
	Cleanup(func())
}) *GoryaServiceHandler {
	mock := &GoryaServiceHandler{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
