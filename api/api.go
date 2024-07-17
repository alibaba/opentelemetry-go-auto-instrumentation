package api

// -----------------------------------------------------------------------------
// Call Context
//
// The CallContext struct is used to pass information between the OnEnter and
// OnExit callbacks. The SkipCall field is used to skip the function call if set
// to true. Params and ReturnVals holds address of parameters and return values
// of the original function call. Modification of the Params and ReturnVals will
// affect the original function call thus should be used with caution.

type CallContext struct {
	Params     []interface{} // Address of parameters of original function
	ReturnVals []interface{} // Address of return values of original function
	Data       interface{}   // User defined data
	SkipCall   bool          // Skip the original function call if set to true
}

func (ctx *CallContext) SetSkipCall(skip bool) {
	ctx.SkipCall = skip
}

func (ctx *CallContext) SetData(data interface{}) {
	ctx.Data = data
}

func (ctx *CallContext) GetData() interface{} {
	return ctx.Data
}

func (ctx *CallContext) SetKeyData(key, val string) {
	if ctx.Data == nil {
		ctx.Data = make(map[string]string)
	}
	ctx.Data.(map[string]string)[key] = val
}

func (ctx *CallContext) GetKeyData(key string) string {
	if ctx.Data == nil {
		return ""
	}
	return ctx.Data.(map[string]string)[key]
}

func (ctx *CallContext) HasKeyData(key string) bool {
	if ctx.Data == nil {
		return false
	}
	_, ok := ctx.Data.(map[string]string)[key]
	return ok
}
