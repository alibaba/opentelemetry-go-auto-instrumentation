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

package ex

import (
	"errors"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
)

// Error represents an error with stack trace information
type stackfulError struct {
	message string
	frame   string
	wrapped error
}

func (e *stackfulError) Error() string {
	return e.message
}

// currentFrame returns the "current frame" whose caller is the function that
// called Errorf.
func currentFrame(skip int) string {
	pc := make([]uintptr, 1)
	n := runtime.Callers(skip, pc)
	if n == 0 {
		return ""
	}
	pc = pc[:n]
	frames := runtime.CallersFrames(pc)
	frame, _ := frames.Next()
	shortFunc := frame.Function
	const prefix = "github.com/alibaba/loongsuite-go-agent/"
	shortFunc = strings.TrimPrefix(shortFunc, prefix)
	return frame.File + ":" + strconv.Itoa(frame.Line) + " " + shortFunc
}

func fetchFrames(err error, cnt int) string {
	e := &stackfulError{}
	if errors.As(err, &e) {
		frame := fmt.Sprintf("[%d] %s\n", cnt, e.frame)
		return fetchFrames(e.wrapped, cnt+1) + frame
	}
	return ""
}

func Error(previousErr error) error {
	if previousErr == nil {
		previousErr = errors.New("unknown error")
	}
	e := &stackfulError{
		message: previousErr.Error(),
		frame:   currentFrame(3),
		wrapped: previousErr,
	}
	return e
}

func Errorf(previousErr error, format string, args ...any) error {
	if previousErr == nil {
		previousErr = errors.New("unknown error")
	}
	e := &stackfulError{
		message: fmt.Sprintf(format, args...),
		frame:   currentFrame(3),
		wrapped: previousErr,
	}
	return e
}

func Fatalf(format string, args ...any) {
	Fatal(Errorf(nil, format, args...))
}

func Fatal(err error) {
	if err == nil {
		panic("Fatal error: unknown")
	}
	err = &stackfulError{
		message: err.Error(),
		frame:   currentFrame(3), // skip the Fatal caller
		wrapped: err,
	}
	e := &stackfulError{}
	if errors.As(err, &e) {
		frames := fetchFrames(err, 0)
		msg := fmt.Sprintf("%s\n\nStack:\n%s", e.message, frames)
		_, _ = fmt.Fprint(os.Stderr, msg)
		os.Exit(1)
	}
	panic(err)
}
