// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package disgoslash

import mock "github.com/stretchr/testify/mock"

// MockSyncer is an autogenerated mock type for the Syncer type
type MockSyncer struct {
	mock.Mock
}

// Run provides a mock function with given fields: guildIDs, slashCommandMap
func (_m *MockSyncer) Run(guildIDs []string, slashCommandMap SlashCommandMap) []error {
	ret := _m.Called(guildIDs, slashCommandMap)

	var r0 []error
	if rf, ok := ret.Get(0).(func([]string, SlashCommandMap) []error); ok {
		r0 = rf(guildIDs, slashCommandMap)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]error)
		}
	}

	return r0
}
