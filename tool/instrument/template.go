//go:build ignore

package instrument

// Seeing is not always believing. The following template is a bit tricky, see
// trampoline.go for more details

// Struct Template
type CallContextImpl struct {
	Params     []interface{}
	ReturnVals []interface{}
	SkipCall   bool
	Data       interface{}
}

func (c *CallContextImpl) SetSkipCall(skip bool)    { c.SkipCall = skip }
func (c *CallContextImpl) IsSkipCall() bool         { return c.SkipCall }
func (c *CallContextImpl) SetData(data interface{}) { c.Data = data }
func (c *CallContextImpl) GetData() interface{}     { return c.Data }
func (c *CallContextImpl) GetParam(idx int) interface{} {
	switch idx {
	}
	return nil
}
func (c *CallContextImpl) SetParam(idx int, val interface{}) {
	switch idx {
	}
}
func (c *CallContextImpl) GetReturnVal(idx int) interface{} {
	switch idx {
	}
	return nil
}
func (c *CallContextImpl) SetReturnVal(idx int, val interface{}) {
	switch idx {
	}
}

// Variable Template
var OtelGetStackImpl func() []byte = nil
var OtelPrintStackImpl func([]byte) = nil

// Trampoline Template
func OtelOnEnterTrampoline() (CallContext, bool) {
	defer func() {
		if err := recover(); err != nil {
			println("failed to exec onEnter hook", "OtelOnEnterNamePlaceholder")
			if e, ok := err.(error); ok {
				println(e.Error())
			}
			fetchStack, printStack := OtelGetStackImpl, OtelPrintStackImpl
			if fetchStack != nil && printStack != nil {
				printStack(fetchStack())
			}
		}
	}()
	callContext := &CallContextImpl{}
	callContext.Params = []interface{}{}
	return callContext, callContext.SkipCall
}

func OtelOnExitTrampoline(callContext CallContext) {
	defer func() {
		if err := recover(); err != nil {
			println("failed to exec onExit hook", "OtelOnExitNamePlaceholder")
			if e, ok := err.(error); ok {
				println(e.Error())
			}
			fetchStack, printStack := OtelGetStackImpl, OtelPrintStackImpl
			if fetchStack != nil && printStack != nil {
				printStack(fetchStack())
			}
		}
	}()
	callContext.(*CallContextImpl).ReturnVals = []interface{}{}
}
