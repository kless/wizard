// Copyright 2010  The "GoWizard" Authors
//
// Use of this source code is governed by the BSD 2-Clause License
// that can be found in the LICENSE file.
//
// This software is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES
// OR CONDITIONS OF ANY KIND, either express or implied. See the License
// for more details.

package wizard

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

const (
	// Permissions
	_PERM_DIRECTORY = 0755
	_PERM_FILE      = 0644

	CHAR_COMMENT = "//" // For comments in source code files
	_CHAR_HEADER = "="  // Header under the project name

	// Configuration file per user
	_USER_CONFIG = ".gowizard"

	// Subdirectory where is installed through "goinstall"
	_DIR_DATA = "src/pkg/github.com/kless/GoWizard/data"
)

/*// VCS configuration files to push to a server.
var configVCS = map[string]string{
	"bzr": ".bzr/branch/branch.conf",
	"git": ".git/config",
	"hg":  ".hg/hgrc",
}*/

// Project types
var ListProject = map[string]string{
	"cmd": "Command line program",
	"pkg": "Package",
	"cgo": "Package that calls C code",
}

// Available licenses
// The first field in []string is the file name without extension.
var ListLicense = map[string][]string{
	"apache": {"Apache", "Apache License, version 2.0"},
	"bsd-2":  {"BSD-2", "BSD 2-Clause License"},
	"bsd-3":  {"BSD-3", "BSD 3-Clause License"},
	"cc0":    {"CC0", "Creative Commons CC0, version 1.0 Universal"},
	"gpl":    {"GPL", "GNU General Public License, version 3 or later"},
	"lgpl":   {"LGPL", "GNU Lesser General Public License, version 3 or later"},
	"agpl":   {"AGPL", "GNU Affero General Public License, version 3 or later"},
	"none":   {"none", "Proprietary license"},
}

// Version control systems (VCS)
var ListVCS = map[string]string{
	"bzr":   "Bazaar",
	"git":   "Git",
	"hg":    "Mercurial",
	"other": "other VCS",
	"none":  "none",
}

// * * *

// Represents all information to create a project
type project struct {
	dirData    string // directory with templates
	dirProject string // directory of project created

	cfg *Conf
	set *template.Set // set of templates
}

// Creates information for the project.
// "isFirstRun" indicates if it is the first time in be called.
func NewProject(cfg *Conf) (*project, error) {
	var err error
	p := new(project)

	if p.dirData, err = dirData(); err != nil {
		return nil, err
	}

	p.dirProject = filepath.Join(cfg.ProjectName, cfg.PackageName)
	p.set = new(template.Set)
	p.cfg = cfg

	return p, nil
}

// Creates a new project.
func (p *project) Create() error {
	if err := os.MkdirAll(p.dirProject, _PERM_DIRECTORY); err != nil {
		return fmt.Errorf("directory error: %s", err)
	}

	p.ParseLicense(CHAR_COMMENT, 0)
	p.parseProject()

	// === Render project files
	if p.cfg.ProjecType != "cmd" {
		p.parseFromVar(filepath.Join(p.dirProject, p.cfg.PackageName)+".go",
			"Pkg")
		p.parseFromVar(filepath.Join(p.dirProject, p.cfg.PackageName)+"_test.go",
			"Test")
	} else {
		p.parseFromVar(filepath.Join(p.dirProject, p.cfg.PackageName)+".go",
			"Cmd")
	}
	p.parseFromVar(filepath.Join(p.dirProject, "Makefile"), "Makefile")

	// === Render common files
	dirTmpl := filepath.Join(p.dirData, "templ") // Base directory of templates

	p.parseFromFile(filepath.Join(p.cfg.ProjectName, "CONTRIBUTORS.md"),
		filepath.Join(dirTmpl, "CONTRIBUTORS.md"), false)
	p.parseFromFile(filepath.Join(p.cfg.ProjectName, "NEWS.md"),
		filepath.Join(dirTmpl, "NEWS.md"), false)
	p.parseFromFile(filepath.Join(p.cfg.ProjectName, "README.md"),
		filepath.Join(dirTmpl, "README.md"), true)

	// The file AUTHORS is for copyright holders.
	if !strings.HasPrefix(p.cfg.License, "cc0") {
		p.parseFromFile(filepath.Join(p.cfg.ProjectName, "AUTHORS.md"),
			filepath.Join(dirTmpl, "AUTHORS.md"), false)
	}

	// === Add file related to VCS
	switch p.cfg.VCS {
	case "other", "none":
		break
	default:
		ignoreFile := "." + p.cfg.VCS + "ignore"

		if p.cfg.VCS == "hg" {
			tmplIgnore = hgIgnoreTop + tmplIgnore
		}

		if err := ioutil.WriteFile(filepath.Join(p.cfg.ProjectName, ignoreFile),
			[]byte(tmplIgnore), _PERM_FILE); err != nil {
			return fmt.Errorf("write error: %s", err)
		}
	}

	// === License file
	AddLicense(p, true)

	return nil
}
