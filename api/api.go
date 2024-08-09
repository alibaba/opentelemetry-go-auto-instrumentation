package api

// -----------------------------------------------------------------------------
// Call Context
//
// The CallContext struct is used to pass information between the OnEnter and
// OnExit callbacks. The SkipCall field is used to skip the function call if set
// to true. Params and ReturnVals holds address of parameters and return values
// of the original function call. Modification of the Params and ReturnVals will
// affect the original function call thus should be used with caution.

type CallContext interface {
	// Skip the original function call
	SetSkipCall(bool)
	// Check if the original function call should be skipped
	IsSkipCall() bool
	// Set the data field, can be used to pass information between OnEnter&OnExit
	SetData(interface{})
	// Get the data field, can be used to pass information between OnEnter&OnExit
	GetData() interface{}
	// Get the map data field by key
	GetKeyData(key string) interface{}
	// Set the map data field by key
	SetKeyData(key string, val interface{})
	// Has the map data field by key
	HasKeyData(key string) bool
	// Get the original function parameter at index idx
	GetParam(idx int) interface{}
	// Change the original function parameter at index idx
	SetParam(idx int, val interface{})
	// Get the original function return value at index idx
	GetReturnVal(idx int) interface{}
	// Change the original function return value at index idx
	SetReturnVal(idx int, val interface{})
}
