// Copyright 2010  The "Go-Wizard" Authors
//
// Use of this source code is governed by the BSD-2 Clause license
// that can be found in the LICENSE file.
//
// This software is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES
// OR CONDITIONS OF ANY KIND, either express or implied. See the License
// for more details.

package main

import (
	"fmt"
	"os"
)

const ERROR = 1 // Exit status code if there is any error.

func reportExit(err os.Error) {
	fmt.Fprintf(os.Stderr, "gowizard: %s\n", err.String())
	os.Exit(ERROR)
}

func fatalf(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, "gowizard: "+format, a...)
	os.Exit(ERROR)
}
