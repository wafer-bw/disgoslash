// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package disgoslash

import (
	http "net/http"

	mock "github.com/stretchr/testify/mock"
)

// MockAuth is an autogenerated mock type for the Auth type
type MockAuth struct {
	mock.Mock
}

// Verify provides a mock function with given fields: rawBody, headers
func (_m *MockAuth) Verify(rawBody []byte, headers http.Header) bool {
	ret := _m.Called(rawBody, headers)

	var r0 bool
	if rf, ok := ret.Get(0).(func([]byte, http.Header) bool); ok {
		r0 = rf(rawBody, headers)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}