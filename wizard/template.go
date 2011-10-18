// Copyright 2010  The "GoWizard" Authors
//
// Use of this source code is governed by the BSD-2 Clause license
// that can be found in the LICENSE file.
//
// This software is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES
// OR CONDITIONS OF ANY KIND, either express or implied. See the License
// for more details.

package wizard

import (
	"log"
	"path/filepath"
	"strings"
	"template"
	"time"
)

// Copyright and licenses
const (
	tmplCopyright = `Copyright {{.Year}}  The "{{.ProjectName}}" Authors`
	tmplCopyleft  = `Written in {{.Year}} by the "{{.ProjectName}}" Authors`

	tmplBSD = `{{.Comment}} {{template "Copyright" .}}
{{.Comment}}
{{.Comment}} Use of this source code is governed by the {{.FullLicense}}
{{.Comment}} that can be found in the LICENSE file.
{{.Comment}}
{{.Comment}} This software is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES
{{.Comment}} OR CONDITIONS OF ANY KIND, either express or implied. See the License
{{.Comment}} for more details.
`

	tmplApache = `{{.Comment}} {{template "Copyright" .}}
{{.Comment}}
{{.Comment}} Licensed under the Apache License, Version 2.0 (the "License");
{{.Comment}} you may not use this file except in compliance with the License.
{{.Comment}} You may obtain a copy of the License at
{{.Comment}}
{{.Comment}}     http://www.apache.org/licenses/LICENSE-2.0
{{.Comment}}
{{.Comment}} Unless required by applicable law or agreed to in writing, software
{{.Comment}} distributed under the License is distributed on an "AS IS" BASIS,
{{.Comment}} WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
{{.Comment}} See the License for the specific language governing permissions and
{{.Comment}} limitations under the License.
`

	tmplGNU = `{{.Comment}} {{template "Copyright" .}}
{{.Comment}}
{{.Comment}} This program is free software: you can redistribute it and/or modify
{{.Comment}} it under the terms of the GNU {{with .GNUextra}}{{.}} {{end}}General Public License as published by
{{.Comment}} the Free Software Foundation, either version 3 of the License, or
{{.Comment}} (at your option) any later version.
{{.Comment}}
{{.Comment}} This program is distributed in the hope that it will be useful,
{{.Comment}} but WITHOUT ANY WARRANTY; without even the implied warranty of
{{.Comment}} MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
{{.Comment}} GNU {{with .GNUextra}}{{.}} {{end}}General Public License for more details.
{{.Comment}}
{{.Comment}} You should have received a copy of the GNU {{with .GNUextra}}{{.}} {{end}}General Public License
{{.Comment}} along with this program.  If not, see <http://www.gnu.org/licenses/>.
`

	tmplNone = `{{.Comment}} {{template "Copyright" .}}
`

	tmplCC0 = `{{.Comment}} {{template "Copyright" .}}
{{.Comment}}
{{.Comment}} To the extent possible under law, the author(s) have waived all copyright
{{.Comment}} and related or neighboring rights to this software to the public domain worldwide.
{{.Comment}} This software is distributed without any warranty.
{{.Comment}}
{{.Comment}} You should have received a copy of the CC0 Public Domain Dedication along
{{.Comment}} with this software. If not, see <http://creativecommons.org/publicdomain/zero/1.0/>.
`
)

// For source code files
const (
	tmplCmd      = `{{template "Header" .}}
package main

import (

)


`
	tmplPkg      = `{{template "Header" .}}
package {{.PackageName}}
{{if .IsCgoProject}}
import "C"{{end}}
import (

)


`
	tmplTest     = `{{template "Header" .}}
package {{.PackageName}}

import (
	"testing"
)

func Test(t *testing.T) {

}

`
	tmplMakefile = `include $(GOROOT)/src/Make.inc

TARG={{if .IsCmdProject}}{{else}}<< IMPORT PATH >>/{{end}}{{.PackageName}}
GOFILES=\
	{{.PackageName}}.go\

include $(GOROOT)/src/Make.{{if .IsCmdProject}}cmd{{else}}pkg{{end}}

`
)

// User configuration
const tmplUserConfig = `[DEFAULT]
author: {{.Author}}
email: {{.Email}}
license: {{.license}}
vcs: {{.vcs}}
`

// === File ignore for VCS
const hgIgnoreTop = "syntax: glob\n"

var tmplIgnore = `# Generic
*~
[._]*

# Go
*.[ao]
*.[568vq]
[568vq].out
main

# Cgo
*.cgo*
*.so
`

// Renders the template "src", creating a file in "dst".
func (p *project) parseFromFile(dst, src string, useNest bool) {
	file := createFile(dst)

	tmpl, err := template.ParseFile(src)
	if err != nil {
		log.Fatal("parsing error:", err)
	}

	if !useNest {
		if err = tmpl.Execute(file, p.cfg); err != nil {
			log.Fatal("execution failed:", err)
		}
	} else {
		p.set.Add(tmpl)

		if err = p.set.Execute(file, filepath.Base(src), p.cfg); err != nil {
			log.Fatal("execution failed:", err)
		}
	}
}

// Renders the template "tmplName" in "set" to the file "dst".
func (p *project) parseFromVar(dst string, tmplName string) {
	file := createFile(dst)

	if err := p.set.Execute(file, tmplName, p.cfg); err != nil {
		log.Fatal("execution failed:", err)
	}
}

// Parses the templates.
// "charComment" is the character used to comment in code files.
// If "year" is nil then gets the actual year.
func (p *project) parseTemplates(charComment string, year int) {
	var tmplHeader string

	licenseName := strings.Split(p.cfg.license, "-")[0]
	p.cfg.Comment = charComment

	if year == 0 {
		p.cfg.Year = int(time.LocalTime().Year)
	} else {
		p.cfg.Year = year
	}

	switch licenseName {
	case "apache":
		tmplHeader = tmplApache
	case "bsd":
		tmplHeader = tmplBSD
	case "cc0":
		tmplHeader = tmplCC0
	case "gpl", "lgpl", "agpl":
		if licenseName == "agpl" {
			p.cfg.GNUextra = "Affero"
		} else if licenseName == "lgpl" {
			p.cfg.GNUextra = "Lesser"
		}

		tmplHeader = tmplGNU
	case "none":
		tmplHeader = tmplNone
	}

	// === Add all templates
	if licenseName != "cc0" {
		p.set.Add(template.Must(template.New("Copyright").
			Parse(tmplCopyright)))
	} else {
		p.set.Add(template.Must(template.New("Copyright").
			Parse(tmplCopyleft)))
	}

	if p.cfg.projecType == "cmd" {
		p.set.Add(template.Must(template.New("Cmd").Parse(tmplCmd)))
	} else {
		tPkg := template.Must(template.New("Pkg").Parse(tmplPkg))
		tTest := template.Must(template.New("Test").Parse(tmplTest))
		p.set.Add(tPkg, tTest)
	}

	if p.cfg.addUserConf {
		p.set.Add(template.Must(template.New("Config").Parse(tmplUserConfig)))
	}

	tHeader := template.Must(template.New("Header").Parse(tmplHeader))
	tMakefile := template.Must(template.New("Makefile").Parse(tmplMakefile))
	p.set.Add(tHeader, tMakefile)
}
