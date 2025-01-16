// Copyright (c) 2025 Alibaba Group Holding Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package errc

import (
	"fmt"
	"runtime/debug"
)

const (
	ErrOpenFile = 1000 + iota
	ErrCreateFile
	ErrCloseFile
	ErrRemoveAll
	ErrReadDir
	ErrCopyFile
	ErrWriteFile
	ErrWalkDir
	ErrStat
	ErrMkdirAll
	ErrNotExist
	ErrInvalidRule
	ErrMatchRule
	ErrInternal
	ErrRunCmd
	ErrInvalidJSON
	ErrGetwd
	ErrSetupRule
	ErrParseCode
	ErrAbsPath
	ErrNotModularized
	ErrGetExecutable
	ErrInstrument
	ErrPreprocess
)

var errMessages = map[int]string{
	ErrOpenFile:       "Failed to open file",
	ErrCreateFile:     "Failed to create file",
	ErrCloseFile:      "Failed to close file",
	ErrRemoveAll:      "Failed to remove all files",
	ErrReadDir:        "Failed to read directory",
	ErrCopyFile:       "Failed to copy file",
	ErrWriteFile:      "Failed to write file",
	ErrWalkDir:        "Failed to walk directory",
	ErrStat:           "Failed to get file info",
	ErrMkdirAll:       "Failed to create directory",
	ErrNotExist:       "File does not exist",
	ErrInvalidRule:    "Invalid rule",
	ErrMatchRule:      "Failed to match rule",
	ErrInternal:       "Internal error",
	ErrRunCmd:         "Failed to run command",
	ErrInvalidJSON:    "Invalid JSON",
	ErrGetwd:          "Failed to get working directory",
	ErrSetupRule:      "Failed to setup rule",
	ErrParseCode:      "Failed to parse Go source code",
	ErrAbsPath:        "Failed to get absolute path",
	ErrNotModularized: "Not a modularized project",
	ErrGetExecutable:  "Failed to get executable",
	ErrInstrument:     "Failed to instrument",
}

type PlentifulError struct {
	ErrorMsg string
	Reason   string
	Cause    string
	Details  map[string]string
}

func (e *PlentifulError) Error() string {
	str := ""
	str += fmt.Sprintf("%-11s: %v\n", "Error", e.ErrorMsg)
	str += fmt.Sprintf("%-11s: %v\n", "Reason", e.Reason)
	str += fmt.Sprintf("%-11s: %v\n", "Cause", e.Cause)
	for k, v := range e.Details {
		str += fmt.Sprintf("%-11s: %v\n", "Detail."+k, v)
	}
	return str
}

func New(code int, message string) *PlentifulError {
	e := &PlentifulError{
		ErrorMsg: errMessages[code],
		Reason:   message,
		Details:  make(map[string]string),
	}
	stackTrace := debug.Stack()
	e.Cause = string(stackTrace)
	return e
}

func (pe *PlentifulError) With(key, value string) *PlentifulError {
	pe.Details[key] = value
	return pe
}

func Adhere(err error, key, value string) error {
	if perr, ok := err.(*PlentifulError); ok {
		perr.Details[key] = value
		return perr
	}
	return err
}
