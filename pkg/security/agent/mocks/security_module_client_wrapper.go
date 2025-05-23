// Code generated by mockery v2.49.2. DO NOT EDIT.

package mocks

import (
	api "github.com/DataDog/datadog-agent/pkg/security/proto/api"
	mock "github.com/stretchr/testify/mock"
)

// SecurityModuleClientWrapper is an autogenerated mock type for the SecurityModuleClientWrapper type
type SecurityModuleClientWrapper struct {
	mock.Mock
}

type SecurityModuleClientWrapper_Expecter struct {
	mock *mock.Mock
}

func (_m *SecurityModuleClientWrapper) EXPECT() *SecurityModuleClientWrapper_Expecter {
	return &SecurityModuleClientWrapper_Expecter{mock: &_m.Mock}
}

// Close provides a mock function with no fields
func (_m *SecurityModuleClientWrapper) Close() {
	_m.Called()
}

// SecurityModuleClientWrapper_Close_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Close'
type SecurityModuleClientWrapper_Close_Call struct {
	*mock.Call
}

// Close is a helper method to define mock.On call
func (_e *SecurityModuleClientWrapper_Expecter) Close() *SecurityModuleClientWrapper_Close_Call {
	return &SecurityModuleClientWrapper_Close_Call{Call: _e.mock.On("Close")}
}

func (_c *SecurityModuleClientWrapper_Close_Call) Run(run func()) *SecurityModuleClientWrapper_Close_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *SecurityModuleClientWrapper_Close_Call) Return() *SecurityModuleClientWrapper_Close_Call {
	_c.Call.Return()
	return _c
}

func (_c *SecurityModuleClientWrapper_Close_Call) RunAndReturn(run func()) *SecurityModuleClientWrapper_Close_Call {
	_c.Run(run)
	return _c
}

// DumpDiscarders provides a mock function with no fields
func (_m *SecurityModuleClientWrapper) DumpDiscarders() (string, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for DumpDiscarders")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func() (string, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SecurityModuleClientWrapper_DumpDiscarders_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DumpDiscarders'
type SecurityModuleClientWrapper_DumpDiscarders_Call struct {
	*mock.Call
}

// DumpDiscarders is a helper method to define mock.On call
func (_e *SecurityModuleClientWrapper_Expecter) DumpDiscarders() *SecurityModuleClientWrapper_DumpDiscarders_Call {
	return &SecurityModuleClientWrapper_DumpDiscarders_Call{Call: _e.mock.On("DumpDiscarders")}
}

func (_c *SecurityModuleClientWrapper_DumpDiscarders_Call) Run(run func()) *SecurityModuleClientWrapper_DumpDiscarders_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *SecurityModuleClientWrapper_DumpDiscarders_Call) Return(_a0 string, _a1 error) *SecurityModuleClientWrapper_DumpDiscarders_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *SecurityModuleClientWrapper_DumpDiscarders_Call) RunAndReturn(run func() (string, error)) *SecurityModuleClientWrapper_DumpDiscarders_Call {
	_c.Call.Return(run)
	return _c
}

// DumpNetworkNamespace provides a mock function with given fields: snapshotInterfaces
func (_m *SecurityModuleClientWrapper) DumpNetworkNamespace(snapshotInterfaces bool) (*api.DumpNetworkNamespaceMessage, error) {
	ret := _m.Called(snapshotInterfaces)

	if len(ret) == 0 {
		panic("no return value specified for DumpNetworkNamespace")
	}

	var r0 *api.DumpNetworkNamespaceMessage
	var r1 error
	if rf, ok := ret.Get(0).(func(bool) (*api.DumpNetworkNamespaceMessage, error)); ok {
		return rf(snapshotInterfaces)
	}
	if rf, ok := ret.Get(0).(func(bool) *api.DumpNetworkNamespaceMessage); ok {
		r0 = rf(snapshotInterfaces)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*api.DumpNetworkNamespaceMessage)
		}
	}

	if rf, ok := ret.Get(1).(func(bool) error); ok {
		r1 = rf(snapshotInterfaces)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SecurityModuleClientWrapper_DumpNetworkNamespace_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DumpNetworkNamespace'
type SecurityModuleClientWrapper_DumpNetworkNamespace_Call struct {
	*mock.Call
}

// DumpNetworkNamespace is a helper method to define mock.On call
//   - snapshotInterfaces bool
func (_e *SecurityModuleClientWrapper_Expecter) DumpNetworkNamespace(snapshotInterfaces interface{}) *SecurityModuleClientWrapper_DumpNetworkNamespace_Call {
	return &SecurityModuleClientWrapper_DumpNetworkNamespace_Call{Call: _e.mock.On("DumpNetworkNamespace", snapshotInterfaces)}
}

func (_c *SecurityModuleClientWrapper_DumpNetworkNamespace_Call) Run(run func(snapshotInterfaces bool)) *SecurityModuleClientWrapper_DumpNetworkNamespace_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(bool))
	})
	return _c
}

func (_c *SecurityModuleClientWrapper_DumpNetworkNamespace_Call) Return(_a0 *api.DumpNetworkNamespaceMessage, _a1 error) *SecurityModuleClientWrapper_DumpNetworkNamespace_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *SecurityModuleClientWrapper_DumpNetworkNamespace_Call) RunAndReturn(run func(bool) (*api.DumpNetworkNamespaceMessage, error)) *SecurityModuleClientWrapper_DumpNetworkNamespace_Call {
	_c.Call.Return(run)
	return _c
}

// DumpProcessCache provides a mock function with given fields: withArgs, format
func (_m *SecurityModuleClientWrapper) DumpProcessCache(withArgs bool, format string) (string, error) {
	ret := _m.Called(withArgs, format)

	if len(ret) == 0 {
		panic("no return value specified for DumpProcessCache")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(bool, string) (string, error)); ok {
		return rf(withArgs, format)
	}
	if rf, ok := ret.Get(0).(func(bool, string) string); ok {
		r0 = rf(withArgs, format)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(bool, string) error); ok {
		r1 = rf(withArgs, format)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SecurityModuleClientWrapper_DumpProcessCache_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DumpProcessCache'
type SecurityModuleClientWrapper_DumpProcessCache_Call struct {
	*mock.Call
}

// DumpProcessCache is a helper method to define mock.On call
//   - withArgs bool
//   - format string
func (_e *SecurityModuleClientWrapper_Expecter) DumpProcessCache(withArgs interface{}, format interface{}) *SecurityModuleClientWrapper_DumpProcessCache_Call {
	return &SecurityModuleClientWrapper_DumpProcessCache_Call{Call: _e.mock.On("DumpProcessCache", withArgs, format)}
}

func (_c *SecurityModuleClientWrapper_DumpProcessCache_Call) Run(run func(withArgs bool, format string)) *SecurityModuleClientWrapper_DumpProcessCache_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(bool), args[1].(string))
	})
	return _c
}

func (_c *SecurityModuleClientWrapper_DumpProcessCache_Call) Return(_a0 string, _a1 error) *SecurityModuleClientWrapper_DumpProcessCache_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *SecurityModuleClientWrapper_DumpProcessCache_Call) RunAndReturn(run func(bool, string) (string, error)) *SecurityModuleClientWrapper_DumpProcessCache_Call {
	_c.Call.Return(run)
	return _c
}

// GenerateActivityDump provides a mock function with given fields: request
func (_m *SecurityModuleClientWrapper) GenerateActivityDump(request *api.ActivityDumpParams) (*api.ActivityDumpMessage, error) {
	ret := _m.Called(request)

	if len(ret) == 0 {
		panic("no return value specified for GenerateActivityDump")
	}

	var r0 *api.ActivityDumpMessage
	var r1 error
	if rf, ok := ret.Get(0).(func(*api.ActivityDumpParams) (*api.ActivityDumpMessage, error)); ok {
		return rf(request)
	}
	if rf, ok := ret.Get(0).(func(*api.ActivityDumpParams) *api.ActivityDumpMessage); ok {
		r0 = rf(request)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*api.ActivityDumpMessage)
		}
	}

	if rf, ok := ret.Get(1).(func(*api.ActivityDumpParams) error); ok {
		r1 = rf(request)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SecurityModuleClientWrapper_GenerateActivityDump_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GenerateActivityDump'
type SecurityModuleClientWrapper_GenerateActivityDump_Call struct {
	*mock.Call
}

// GenerateActivityDump is a helper method to define mock.On call
//   - request *api.ActivityDumpParams
func (_e *SecurityModuleClientWrapper_Expecter) GenerateActivityDump(request interface{}) *SecurityModuleClientWrapper_GenerateActivityDump_Call {
	return &SecurityModuleClientWrapper_GenerateActivityDump_Call{Call: _e.mock.On("GenerateActivityDump", request)}
}

func (_c *SecurityModuleClientWrapper_GenerateActivityDump_Call) Run(run func(request *api.ActivityDumpParams)) *SecurityModuleClientWrapper_GenerateActivityDump_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*api.ActivityDumpParams))
	})
	return _c
}

func (_c *SecurityModuleClientWrapper_GenerateActivityDump_Call) Return(_a0 *api.ActivityDumpMessage, _a1 error) *SecurityModuleClientWrapper_GenerateActivityDump_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *SecurityModuleClientWrapper_GenerateActivityDump_Call) RunAndReturn(run func(*api.ActivityDumpParams) (*api.ActivityDumpMessage, error)) *SecurityModuleClientWrapper_GenerateActivityDump_Call {
	_c.Call.Return(run)
	return _c
}

// GenerateEncoding provides a mock function with given fields: request
func (_m *SecurityModuleClientWrapper) GenerateEncoding(request *api.TranscodingRequestParams) (*api.TranscodingRequestMessage, error) {
	ret := _m.Called(request)

	if len(ret) == 0 {
		panic("no return value specified for GenerateEncoding")
	}

	var r0 *api.TranscodingRequestMessage
	var r1 error
	if rf, ok := ret.Get(0).(func(*api.TranscodingRequestParams) (*api.TranscodingRequestMessage, error)); ok {
		return rf(request)
	}
	if rf, ok := ret.Get(0).(func(*api.TranscodingRequestParams) *api.TranscodingRequestMessage); ok {
		r0 = rf(request)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*api.TranscodingRequestMessage)
		}
	}

	if rf, ok := ret.Get(1).(func(*api.TranscodingRequestParams) error); ok {
		r1 = rf(request)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SecurityModuleClientWrapper_GenerateEncoding_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GenerateEncoding'
type SecurityModuleClientWrapper_GenerateEncoding_Call struct {
	*mock.Call
}

// GenerateEncoding is a helper method to define mock.On call
//   - request *api.TranscodingRequestParams
func (_e *SecurityModuleClientWrapper_Expecter) GenerateEncoding(request interface{}) *SecurityModuleClientWrapper_GenerateEncoding_Call {
	return &SecurityModuleClientWrapper_GenerateEncoding_Call{Call: _e.mock.On("GenerateEncoding", request)}
}

func (_c *SecurityModuleClientWrapper_GenerateEncoding_Call) Run(run func(request *api.TranscodingRequestParams)) *SecurityModuleClientWrapper_GenerateEncoding_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*api.TranscodingRequestParams))
	})
	return _c
}

func (_c *SecurityModuleClientWrapper_GenerateEncoding_Call) Return(_a0 *api.TranscodingRequestMessage, _a1 error) *SecurityModuleClientWrapper_GenerateEncoding_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *SecurityModuleClientWrapper_GenerateEncoding_Call) RunAndReturn(run func(*api.TranscodingRequestParams) (*api.TranscodingRequestMessage, error)) *SecurityModuleClientWrapper_GenerateEncoding_Call {
	_c.Call.Return(run)
	return _c
}

// GetActivityDumpStream provides a mock function with no fields
func (_m *SecurityModuleClientWrapper) GetActivityDumpStream() (api.SecurityModule_GetActivityDumpStreamClient, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetActivityDumpStream")
	}

	var r0 api.SecurityModule_GetActivityDumpStreamClient
	var r1 error
	if rf, ok := ret.Get(0).(func() (api.SecurityModule_GetActivityDumpStreamClient, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() api.SecurityModule_GetActivityDumpStreamClient); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(api.SecurityModule_GetActivityDumpStreamClient)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SecurityModuleClientWrapper_GetActivityDumpStream_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetActivityDumpStream'
type SecurityModuleClientWrapper_GetActivityDumpStream_Call struct {
	*mock.Call
}

// GetActivityDumpStream is a helper method to define mock.On call
func (_e *SecurityModuleClientWrapper_Expecter) GetActivityDumpStream() *SecurityModuleClientWrapper_GetActivityDumpStream_Call {
	return &SecurityModuleClientWrapper_GetActivityDumpStream_Call{Call: _e.mock.On("GetActivityDumpStream")}
}

func (_c *SecurityModuleClientWrapper_GetActivityDumpStream_Call) Run(run func()) *SecurityModuleClientWrapper_GetActivityDumpStream_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *SecurityModuleClientWrapper_GetActivityDumpStream_Call) Return(_a0 api.SecurityModule_GetActivityDumpStreamClient, _a1 error) *SecurityModuleClientWrapper_GetActivityDumpStream_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *SecurityModuleClientWrapper_GetActivityDumpStream_Call) RunAndReturn(run func() (api.SecurityModule_GetActivityDumpStreamClient, error)) *SecurityModuleClientWrapper_GetActivityDumpStream_Call {
	_c.Call.Return(run)
	return _c
}

// GetConfig provides a mock function with no fields
func (_m *SecurityModuleClientWrapper) GetConfig() (*api.SecurityConfigMessage, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetConfig")
	}

	var r0 *api.SecurityConfigMessage
	var r1 error
	if rf, ok := ret.Get(0).(func() (*api.SecurityConfigMessage, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() *api.SecurityConfigMessage); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*api.SecurityConfigMessage)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SecurityModuleClientWrapper_GetConfig_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetConfig'
type SecurityModuleClientWrapper_GetConfig_Call struct {
	*mock.Call
}

// GetConfig is a helper method to define mock.On call
func (_e *SecurityModuleClientWrapper_Expecter) GetConfig() *SecurityModuleClientWrapper_GetConfig_Call {
	return &SecurityModuleClientWrapper_GetConfig_Call{Call: _e.mock.On("GetConfig")}
}

func (_c *SecurityModuleClientWrapper_GetConfig_Call) Run(run func()) *SecurityModuleClientWrapper_GetConfig_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *SecurityModuleClientWrapper_GetConfig_Call) Return(_a0 *api.SecurityConfigMessage, _a1 error) *SecurityModuleClientWrapper_GetConfig_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *SecurityModuleClientWrapper_GetConfig_Call) RunAndReturn(run func() (*api.SecurityConfigMessage, error)) *SecurityModuleClientWrapper_GetConfig_Call {
	_c.Call.Return(run)
	return _c
}

// GetEvents provides a mock function with no fields
func (_m *SecurityModuleClientWrapper) GetEvents() (api.SecurityModule_GetEventsClient, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetEvents")
	}

	var r0 api.SecurityModule_GetEventsClient
	var r1 error
	if rf, ok := ret.Get(0).(func() (api.SecurityModule_GetEventsClient, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() api.SecurityModule_GetEventsClient); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(api.SecurityModule_GetEventsClient)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SecurityModuleClientWrapper_GetEvents_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetEvents'
type SecurityModuleClientWrapper_GetEvents_Call struct {
	*mock.Call
}

// GetEvents is a helper method to define mock.On call
func (_e *SecurityModuleClientWrapper_Expecter) GetEvents() *SecurityModuleClientWrapper_GetEvents_Call {
	return &SecurityModuleClientWrapper_GetEvents_Call{Call: _e.mock.On("GetEvents")}
}

func (_c *SecurityModuleClientWrapper_GetEvents_Call) Run(run func()) *SecurityModuleClientWrapper_GetEvents_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *SecurityModuleClientWrapper_GetEvents_Call) Return(_a0 api.SecurityModule_GetEventsClient, _a1 error) *SecurityModuleClientWrapper_GetEvents_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *SecurityModuleClientWrapper_GetEvents_Call) RunAndReturn(run func() (api.SecurityModule_GetEventsClient, error)) *SecurityModuleClientWrapper_GetEvents_Call {
	_c.Call.Return(run)
	return _c
}

// GetRuleSetReport provides a mock function with no fields
func (_m *SecurityModuleClientWrapper) GetRuleSetReport() (*api.GetRuleSetReportMessage, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetRuleSetReport")
	}

	var r0 *api.GetRuleSetReportMessage
	var r1 error
	if rf, ok := ret.Get(0).(func() (*api.GetRuleSetReportMessage, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() *api.GetRuleSetReportMessage); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*api.GetRuleSetReportMessage)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SecurityModuleClientWrapper_GetRuleSetReport_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetRuleSetReport'
type SecurityModuleClientWrapper_GetRuleSetReport_Call struct {
	*mock.Call
}

// GetRuleSetReport is a helper method to define mock.On call
func (_e *SecurityModuleClientWrapper_Expecter) GetRuleSetReport() *SecurityModuleClientWrapper_GetRuleSetReport_Call {
	return &SecurityModuleClientWrapper_GetRuleSetReport_Call{Call: _e.mock.On("GetRuleSetReport")}
}

func (_c *SecurityModuleClientWrapper_GetRuleSetReport_Call) Run(run func()) *SecurityModuleClientWrapper_GetRuleSetReport_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *SecurityModuleClientWrapper_GetRuleSetReport_Call) Return(_a0 *api.GetRuleSetReportMessage, _a1 error) *SecurityModuleClientWrapper_GetRuleSetReport_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *SecurityModuleClientWrapper_GetRuleSetReport_Call) RunAndReturn(run func() (*api.GetRuleSetReportMessage, error)) *SecurityModuleClientWrapper_GetRuleSetReport_Call {
	_c.Call.Return(run)
	return _c
}

// GetStatus provides a mock function with no fields
func (_m *SecurityModuleClientWrapper) GetStatus() (*api.Status, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetStatus")
	}

	var r0 *api.Status
	var r1 error
	if rf, ok := ret.Get(0).(func() (*api.Status, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() *api.Status); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*api.Status)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SecurityModuleClientWrapper_GetStatus_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetStatus'
type SecurityModuleClientWrapper_GetStatus_Call struct {
	*mock.Call
}

// GetStatus is a helper method to define mock.On call
func (_e *SecurityModuleClientWrapper_Expecter) GetStatus() *SecurityModuleClientWrapper_GetStatus_Call {
	return &SecurityModuleClientWrapper_GetStatus_Call{Call: _e.mock.On("GetStatus")}
}

func (_c *SecurityModuleClientWrapper_GetStatus_Call) Run(run func()) *SecurityModuleClientWrapper_GetStatus_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *SecurityModuleClientWrapper_GetStatus_Call) Return(_a0 *api.Status, _a1 error) *SecurityModuleClientWrapper_GetStatus_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *SecurityModuleClientWrapper_GetStatus_Call) RunAndReturn(run func() (*api.Status, error)) *SecurityModuleClientWrapper_GetStatus_Call {
	_c.Call.Return(run)
	return _c
}

// ListActivityDumps provides a mock function with no fields
func (_m *SecurityModuleClientWrapper) ListActivityDumps() (*api.ActivityDumpListMessage, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for ListActivityDumps")
	}

	var r0 *api.ActivityDumpListMessage
	var r1 error
	if rf, ok := ret.Get(0).(func() (*api.ActivityDumpListMessage, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() *api.ActivityDumpListMessage); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*api.ActivityDumpListMessage)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SecurityModuleClientWrapper_ListActivityDumps_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ListActivityDumps'
type SecurityModuleClientWrapper_ListActivityDumps_Call struct {
	*mock.Call
}

// ListActivityDumps is a helper method to define mock.On call
func (_e *SecurityModuleClientWrapper_Expecter) ListActivityDumps() *SecurityModuleClientWrapper_ListActivityDumps_Call {
	return &SecurityModuleClientWrapper_ListActivityDumps_Call{Call: _e.mock.On("ListActivityDumps")}
}

func (_c *SecurityModuleClientWrapper_ListActivityDumps_Call) Run(run func()) *SecurityModuleClientWrapper_ListActivityDumps_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *SecurityModuleClientWrapper_ListActivityDumps_Call) Return(_a0 *api.ActivityDumpListMessage, _a1 error) *SecurityModuleClientWrapper_ListActivityDumps_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *SecurityModuleClientWrapper_ListActivityDumps_Call) RunAndReturn(run func() (*api.ActivityDumpListMessage, error)) *SecurityModuleClientWrapper_ListActivityDumps_Call {
	_c.Call.Return(run)
	return _c
}

// ListSecurityProfiles provides a mock function with given fields: includeCache
func (_m *SecurityModuleClientWrapper) ListSecurityProfiles(includeCache bool) (*api.SecurityProfileListMessage, error) {
	ret := _m.Called(includeCache)

	if len(ret) == 0 {
		panic("no return value specified for ListSecurityProfiles")
	}

	var r0 *api.SecurityProfileListMessage
	var r1 error
	if rf, ok := ret.Get(0).(func(bool) (*api.SecurityProfileListMessage, error)); ok {
		return rf(includeCache)
	}
	if rf, ok := ret.Get(0).(func(bool) *api.SecurityProfileListMessage); ok {
		r0 = rf(includeCache)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*api.SecurityProfileListMessage)
		}
	}

	if rf, ok := ret.Get(1).(func(bool) error); ok {
		r1 = rf(includeCache)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SecurityModuleClientWrapper_ListSecurityProfiles_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ListSecurityProfiles'
type SecurityModuleClientWrapper_ListSecurityProfiles_Call struct {
	*mock.Call
}

// ListSecurityProfiles is a helper method to define mock.On call
//   - includeCache bool
func (_e *SecurityModuleClientWrapper_Expecter) ListSecurityProfiles(includeCache interface{}) *SecurityModuleClientWrapper_ListSecurityProfiles_Call {
	return &SecurityModuleClientWrapper_ListSecurityProfiles_Call{Call: _e.mock.On("ListSecurityProfiles", includeCache)}
}

func (_c *SecurityModuleClientWrapper_ListSecurityProfiles_Call) Run(run func(includeCache bool)) *SecurityModuleClientWrapper_ListSecurityProfiles_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(bool))
	})
	return _c
}

func (_c *SecurityModuleClientWrapper_ListSecurityProfiles_Call) Return(_a0 *api.SecurityProfileListMessage, _a1 error) *SecurityModuleClientWrapper_ListSecurityProfiles_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *SecurityModuleClientWrapper_ListSecurityProfiles_Call) RunAndReturn(run func(bool) (*api.SecurityProfileListMessage, error)) *SecurityModuleClientWrapper_ListSecurityProfiles_Call {
	_c.Call.Return(run)
	return _c
}

// ReloadPolicies provides a mock function with no fields
func (_m *SecurityModuleClientWrapper) ReloadPolicies() (*api.ReloadPoliciesResultMessage, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for ReloadPolicies")
	}

	var r0 *api.ReloadPoliciesResultMessage
	var r1 error
	if rf, ok := ret.Get(0).(func() (*api.ReloadPoliciesResultMessage, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() *api.ReloadPoliciesResultMessage); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*api.ReloadPoliciesResultMessage)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SecurityModuleClientWrapper_ReloadPolicies_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ReloadPolicies'
type SecurityModuleClientWrapper_ReloadPolicies_Call struct {
	*mock.Call
}

// ReloadPolicies is a helper method to define mock.On call
func (_e *SecurityModuleClientWrapper_Expecter) ReloadPolicies() *SecurityModuleClientWrapper_ReloadPolicies_Call {
	return &SecurityModuleClientWrapper_ReloadPolicies_Call{Call: _e.mock.On("ReloadPolicies")}
}

func (_c *SecurityModuleClientWrapper_ReloadPolicies_Call) Run(run func()) *SecurityModuleClientWrapper_ReloadPolicies_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *SecurityModuleClientWrapper_ReloadPolicies_Call) Return(_a0 *api.ReloadPoliciesResultMessage, _a1 error) *SecurityModuleClientWrapper_ReloadPolicies_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *SecurityModuleClientWrapper_ReloadPolicies_Call) RunAndReturn(run func() (*api.ReloadPoliciesResultMessage, error)) *SecurityModuleClientWrapper_ReloadPolicies_Call {
	_c.Call.Return(run)
	return _c
}

// RunSelfTest provides a mock function with no fields
func (_m *SecurityModuleClientWrapper) RunSelfTest() (*api.SecuritySelfTestResultMessage, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for RunSelfTest")
	}

	var r0 *api.SecuritySelfTestResultMessage
	var r1 error
	if rf, ok := ret.Get(0).(func() (*api.SecuritySelfTestResultMessage, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() *api.SecuritySelfTestResultMessage); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*api.SecuritySelfTestResultMessage)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SecurityModuleClientWrapper_RunSelfTest_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'RunSelfTest'
type SecurityModuleClientWrapper_RunSelfTest_Call struct {
	*mock.Call
}

// RunSelfTest is a helper method to define mock.On call
func (_e *SecurityModuleClientWrapper_Expecter) RunSelfTest() *SecurityModuleClientWrapper_RunSelfTest_Call {
	return &SecurityModuleClientWrapper_RunSelfTest_Call{Call: _e.mock.On("RunSelfTest")}
}

func (_c *SecurityModuleClientWrapper_RunSelfTest_Call) Run(run func()) *SecurityModuleClientWrapper_RunSelfTest_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *SecurityModuleClientWrapper_RunSelfTest_Call) Return(_a0 *api.SecuritySelfTestResultMessage, _a1 error) *SecurityModuleClientWrapper_RunSelfTest_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *SecurityModuleClientWrapper_RunSelfTest_Call) RunAndReturn(run func() (*api.SecuritySelfTestResultMessage, error)) *SecurityModuleClientWrapper_RunSelfTest_Call {
	_c.Call.Return(run)
	return _c
}

// SaveSecurityProfile provides a mock function with given fields: name, tag
func (_m *SecurityModuleClientWrapper) SaveSecurityProfile(name string, tag string) (*api.SecurityProfileSaveMessage, error) {
	ret := _m.Called(name, tag)

	if len(ret) == 0 {
		panic("no return value specified for SaveSecurityProfile")
	}

	var r0 *api.SecurityProfileSaveMessage
	var r1 error
	if rf, ok := ret.Get(0).(func(string, string) (*api.SecurityProfileSaveMessage, error)); ok {
		return rf(name, tag)
	}
	if rf, ok := ret.Get(0).(func(string, string) *api.SecurityProfileSaveMessage); ok {
		r0 = rf(name, tag)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*api.SecurityProfileSaveMessage)
		}
	}

	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(name, tag)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SecurityModuleClientWrapper_SaveSecurityProfile_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SaveSecurityProfile'
type SecurityModuleClientWrapper_SaveSecurityProfile_Call struct {
	*mock.Call
}

// SaveSecurityProfile is a helper method to define mock.On call
//   - name string
//   - tag string
func (_e *SecurityModuleClientWrapper_Expecter) SaveSecurityProfile(name interface{}, tag interface{}) *SecurityModuleClientWrapper_SaveSecurityProfile_Call {
	return &SecurityModuleClientWrapper_SaveSecurityProfile_Call{Call: _e.mock.On("SaveSecurityProfile", name, tag)}
}

func (_c *SecurityModuleClientWrapper_SaveSecurityProfile_Call) Run(run func(name string, tag string)) *SecurityModuleClientWrapper_SaveSecurityProfile_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(string))
	})
	return _c
}

func (_c *SecurityModuleClientWrapper_SaveSecurityProfile_Call) Return(_a0 *api.SecurityProfileSaveMessage, _a1 error) *SecurityModuleClientWrapper_SaveSecurityProfile_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *SecurityModuleClientWrapper_SaveSecurityProfile_Call) RunAndReturn(run func(string, string) (*api.SecurityProfileSaveMessage, error)) *SecurityModuleClientWrapper_SaveSecurityProfile_Call {
	_c.Call.Return(run)
	return _c
}

// StopActivityDump provides a mock function with given fields: name, container, cgroup
func (_m *SecurityModuleClientWrapper) StopActivityDump(name string, container string, cgroup string) (*api.ActivityDumpStopMessage, error) {
	ret := _m.Called(name, container, cgroup)

	if len(ret) == 0 {
		panic("no return value specified for StopActivityDump")
	}

	var r0 *api.ActivityDumpStopMessage
	var r1 error
	if rf, ok := ret.Get(0).(func(string, string, string) (*api.ActivityDumpStopMessage, error)); ok {
		return rf(name, container, cgroup)
	}
	if rf, ok := ret.Get(0).(func(string, string, string) *api.ActivityDumpStopMessage); ok {
		r0 = rf(name, container, cgroup)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*api.ActivityDumpStopMessage)
		}
	}

	if rf, ok := ret.Get(1).(func(string, string, string) error); ok {
		r1 = rf(name, container, cgroup)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SecurityModuleClientWrapper_StopActivityDump_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'StopActivityDump'
type SecurityModuleClientWrapper_StopActivityDump_Call struct {
	*mock.Call
}

// StopActivityDump is a helper method to define mock.On call
//   - name string
//   - container string
//   - cgroup string
func (_e *SecurityModuleClientWrapper_Expecter) StopActivityDump(name interface{}, container interface{}, cgroup interface{}) *SecurityModuleClientWrapper_StopActivityDump_Call {
	return &SecurityModuleClientWrapper_StopActivityDump_Call{Call: _e.mock.On("StopActivityDump", name, container, cgroup)}
}

func (_c *SecurityModuleClientWrapper_StopActivityDump_Call) Run(run func(name string, container string, cgroup string)) *SecurityModuleClientWrapper_StopActivityDump_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(string), args[2].(string))
	})
	return _c
}

func (_c *SecurityModuleClientWrapper_StopActivityDump_Call) Return(_a0 *api.ActivityDumpStopMessage, _a1 error) *SecurityModuleClientWrapper_StopActivityDump_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *SecurityModuleClientWrapper_StopActivityDump_Call) RunAndReturn(run func(string, string, string) (*api.ActivityDumpStopMessage, error)) *SecurityModuleClientWrapper_StopActivityDump_Call {
	_c.Call.Return(run)
	return _c
}

// NewSecurityModuleClientWrapper creates a new instance of SecurityModuleClientWrapper. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewSecurityModuleClientWrapper(t interface {
	mock.TestingT
	Cleanup(func())
}) *SecurityModuleClientWrapper {
	mock := &SecurityModuleClientWrapper{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
