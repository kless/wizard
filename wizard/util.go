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
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"
)

// Creates the user configuration file.
func AddConfig(cfg *Conf) error {
	tmpl := template.Must(template.New("Config").Parse(tmplUserConfig))

	envHome := os.Getenv("HOME")
	if envHome == "" {
		return errors.New("could not add user configuration file because $HOME is not set")
	}

	file, err := createFile(filepath.Join(envHome, _USER_CONFIG))
	if err != nil {
		return err
	}

	if err := tmpl.Execute(file, cfg); err != nil {
		return fmt.Errorf("execution failed: %s", err)
	}
	return nil
}

// Creates a license file.
func AddLicense(p *project, isNewProject bool) error {
	licenseLower := p.cfg.License
	dirProject := p.cfg.ProjectName

	if !isNewProject {
		licenseLower = strings.ToLower(licenseLower)
		if err := CheckLicense(licenseLower); err != nil {
			return err
		}

		dirProject = "." // actual directory
	}

	dirData := filepath.Join(p.dirData, "license")
	license := ListLowerLicense[licenseLower]

	filename := func(name string) string {
		if strings.HasPrefix(name, "BSD") {
			name = strings.TrimRight(name, "-23")
		}

		if name == "unlicense" {
			return "UNLICENSE.txt"
		}
		if isNewProject {
			return "LICENSE.txt"
		}
		return "LICENSE-" + name + ".txt"
	}

	switch licenseLower {
	case "none":
		break
	case "bsd-2", "bsd-3":
		p.parseFromFile(filepath.Join(dirProject, filename(license)),
			filepath.Join(dirData, license+".txt"))
	default:
		copyFile(filepath.Join(dirProject, filename(license)),
			filepath.Join(dirData, license+".txt"), _PERM_FILE)

		// License LGPL must also add the GPL license text.
		if licenseLower == "lgpl" {
			isNewProject = false
			copyFile(filepath.Join(dirProject, filename("GPL")),
				filepath.Join(dirData, "GPL.txt"), _PERM_FILE)
		}
	}

	return nil
}

// Finds the first line that matches the copyright header to return the year.
func ProjectYear(filename string) (int, error) {
	file, err := os.Open(filename)
	if err != nil {
		return 0, fmt.Errorf("no project directory: %s", err)
	}
	defer file.Close()

	fileBuf := bufio.NewReader(file)

	for {
		line, err := fileBuf.ReadString('\n')
		if err == io.EOF {
			break
		}

		if reCopyright.MatchString(line) || reCopyleft.MatchString(line) {
			//strYear := strings.Split(line, " ")[1]
			return strconv.Atoi(strings.Split(line, " ")[1])
		}
	}
	panic("unreached")
}

// * * *

// Copies a file from source to destination.
func copyFile(destination, source string, perm uint32) error {
	src, err := ioutil.ReadFile(source)
	if err != nil {
		return fmt.Errorf("copy error reading: %s", err)
	}

	err = ioutil.WriteFile(destination, src, perm)
	if err != nil {
		return fmt.Errorf("copy error writing: %s", err)
	}

	return nil
}

// Creates a file.
func createFile(dst string) (*os.File, error) {
	file, err := os.Create(dst)
	if err != nil {
		return nil, fmt.Errorf("file error: %s", err)
	}
	if err = file.Chmod(_PERM_FILE); err != nil {
		return nil, fmt.Errorf("file error: %s", err)
	}

	return file, nil
}

// Gets the path of the templates directory.
func dirData() (string, error) {
	goEnv := os.Getenv("GOPATH")

	if goEnv != "" {
		goto _Found
	}
	if goEnv = os.Getenv("GOROOT"); goEnv != "" {
		goto _Found
	}
	if goEnv = os.Getenv("GOROOT_FINAL"); goEnv != "" {
		goto _Found
	}

_Found:
	if goEnv == "" {
		return "", errors.New("environment variable GOROOT neither" +
			" GOROOT_FINAL has been set")
	}

	return filepath.Join(goEnv, _DIR_DATA), nil
}
