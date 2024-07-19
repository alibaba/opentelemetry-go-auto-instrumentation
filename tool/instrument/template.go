//go:build ignore

package instrument

// @@ Modification on this trampoline template should be cautious, as it imposes
// many implicit constraints on generated code, known constraints are as follows:
// - It's performance critical, so it should be as simple as possible
// - It should not import any package because there is no guarantee that package
//   is existed in import config during the compilation, one practical approach
//   is to use function variables and setup these variables in preprocess stage
// - It should not panic as this affects user application
// - Function and variable names are coupled with the framework, any modification
//   on them should be synced with the framework

// Variable Declaration
var OtelGetStackImpl func() []byte = nil
var OtelPrintStackImpl func([]byte) = nil

// Function Declaration
func OtelOnEnterTrampoline() (*CallContext, bool) {
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
	callContext := &CallContext{
		Params:     nil,
		ReturnVals: nil,
		SkipCall:   false,
	}
	callContext.Params = []interface{}{}
	return callContext, callContext.SkipCall
}

func OtelOnExitTrampoline(callContext *CallContext) {
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
	callContext.ReturnVals = []interface{}{}
}
