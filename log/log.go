// Copyright 2023 FishGoddess. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package log

import stdlog "log"

var (
	// DebugFunc logs a message as debug.
	// Set to nil if you want debug logs are ignored.
	DebugFunc = func(format string, v ...any) {
		format = "[DEBUG] vex: " + format
		stdlog.Printf(format, v...)
	}

	// InfoFunc logs a message as info.
	// Set to nil if you want info logs are ignored.
	InfoFunc = func(format string, v ...any) {
		format = "[INFO] vex: " + format
		stdlog.Printf(format, v...)
	}

	// ErrorFunc logs a message as error.
	// Set to nil if you want error logs are ignored.
	ErrorFunc = func(err error, format string, v ...any) {
		format = "[ERROR] vex: " + format + ": %+v"
		stdlog.Printf(format, append(v, err)...)
	}
)

// Debug logs a debug message.
func Debug(format string, v ...interface{}) {
	if DebugFunc != nil {
		DebugFunc(format, v...)
	}
}

// Info logs an info message.
func Info(format string, v ...interface{}) {
	if InfoFunc != nil {
		InfoFunc(format, v...)
	}
}

// Error logs an error message.
func Error(err error, format string, v ...interface{}) {
	if ErrorFunc != nil {
		ErrorFunc(err, format, v...)
	}
}
