package shared

// All tool's flags are listed here.

// InToolexec true means this tool is being invoked in the go build process.
// This flag should not be set manually by users.
var InToolexec bool
var NameOfInToolexec = "in-toolexec"
var UsageOfIntoolexec = "true means this tool is being invoked in the go build process"
