// Copyright 2023 FishGoddess. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package vex

import stdlog "log"

var (
	// LogDebug logs a message as debug.
	LogDebug = func(format string, v ...any) {
		format = "[DEBUG] vex: " + format
		stdlog.Printf(format, v...)
	}

	// LogInfo logs a message as info.
	LogInfo = func(format string, v ...any) {
		format = "[INFO] vex: " + format
		stdlog.Printf(format, v...)
	}

	// LogError logs a message as error.
	LogError = func(err error, format string, v ...any) {
		format = "[ERROR] vex: " + format + ": %+v"
		stdlog.Printf(format, append(v, err)...)
	}
)

func logDebug(format string, v ...interface{}) {
	if LogDebug != nil {
		LogDebug(format, v...)
	}
}

func logInfo(format string, v ...interface{}) {
	if LogInfo != nil {
		LogInfo(format, v...)
	}
}

func logError(err error, format string, v ...interface{}) {
	if LogError != nil {
		LogError(err, format, v...)
	}
}
