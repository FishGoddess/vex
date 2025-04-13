// Copyright 2025 FishGoddess. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"io"

	"github.com/FishGoddess/vex/log"
)

func main() {
	// All logs outputted by client and server can be ignored by some functions:
	log.DisableDebug()
	log.DisableInfo()
	log.DisableError()

	// If you want to ignore all logs, try it:
	log.DisableAll()

	// Also, you might want to customize the logging functions, like using your own logging library.
	// Just overwrite the variables of log:
	log.DebugFunc = func(format string, v ...any) {
		format = "[DEBUG] " + format + "\n"
		fmt.Printf(format, v...)
	}

	log.InfoFunc = func(format string, v ...any) {
		format = "[INFO] " + format + "\n"
		fmt.Printf(format, v...)
	}

	log.ErrorFunc = func(err error, format string, v ...any) {
		format = "[ERROR] " + format + " | err=%+v\n"
		v = append(v, err)
		fmt.Printf(format, v...)
	}

	// Then, try to log:
	log.Debug("wow oooooooooh %d", 123)
	log.Info("calm down boy %.2f", 3.14)
	log.Error(nil, "something happening? %s", "No.")
	log.Error(io.EOF, "oh no %+v", true)
}
