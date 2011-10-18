// Copyright 2010  The "GoWizard" Authors
//
// Use of this source code is governed by the BSD-2 Clause license
// that can be found in the LICENSE file.
//
// This software is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES
// OR CONDITIONS OF ANY KIND, either express or implied. See the License
// for more details.

package main

import (
	"os"

	"github.com/kless/GoWizard/wizard"
)

func main() {
	p := wizard.NewProject(true)
	p.Create()

	os.Exit(0)
}
