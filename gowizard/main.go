// Copyright 2010  The "gowizard" Authors
//
// Use of this source code is governed by the Simplified BSD License
// that can be found in the LICENSE file.
//
// This software is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES
// OR CONDITIONS OF ANY KIND, either express or implied. See the License
// for more details.

package main

import (
	"container/vector"
	"fmt"
	"log"
	"os"
	"path"
	"strings"
)


// Permissions
const (
	PERM_DIRECTORY = 0755
	PERM_FILE      = 0644
)

// Characters
const (
	CHAR_COMMENT_CODE = "//" // For comments in source code files
	CHAR_COMMENT_MAKE = "#"  // For comments in file Makefile
	CHAR_HEADER       = '='  // Header under the project name
)

const ERROR = 2 // Exit status code if there is any error
const README = "README.mkd"

// Get data directory from `$(GOROOT)/lib/$(TARG)`
var dirData = path.Join(os.Getenv("GOROOT"), "lib", "gowizard")

var argv0 = os.Args[0] // Executable name


// === Main program execution
func main() {
	loadConfig()

	if !*fUpdate {
		createProject()
	} else {
		updateProject()
	}

	os.Exit(0)
}

/* Add license file in directory `dir`. */
func addLicense(dir string, tag map[string]string) {
	dirTmpl := dirData + "/license"

	switch *fLicense {
	case "none":
		break
	case "bsd-3":
		renderNewFile(dir+"/LICENSE", dirTmpl+"/bsd-3.txt", tag)
	default:
		if err := copyFile(dir+"/LICENSE",
			path.Join(dirTmpl, *fLicense+".txt")); err != nil {
			log.Exit(err)
		}
	}
}

/* Show data on 'tag'. */
func debug(tag map[string]string) {
	fmt.Println("  = Debug\n")

	for k, v := range tag {
		// Tags starting with '_' are not showed.
		if k[0] == '_' {
			continue
		}
		fmt.Printf("  %s: %s\n", k, v)
	}
	os.Exit(0)
}

// ===


/* Creates a new project. */
func createProject() {
	tag := tagsToCreate()
	if *fDebug {
		debug(tag)
	}

	headerCodeFile, headerMakefile := renderAllHeaders(tag, "")

	// === Create directories in lower case
	dirApp := path.Join(*fProjectName, *fPackageName)
	os.MkdirAll(dirApp, PERM_DIRECTORY)

	// === Render project files
	renderNesting(dirApp+"/Makefile", headerMakefile, tmplMakefile, tag)

	switch *fProjecType {
	case "lib", "cgo":
		renderNesting(dirApp+"/main.go", headerCodeFile, tmplPkgMain, tag)
		renderNesting(dirApp+"/main_test.go", headerCodeFile, tmplTest, tag)
	case "app", "tool":
		renderNesting(dirApp+"/main.go", headerCodeFile, tmplCmdMain, tag)
	}

	// === Render common files
	dirTmpl := dirData + "/tmpl" // Templates base directory

	renderFile(*fProjectName, dirTmpl+"/NEWS.mkd", tag)
	renderFile(*fProjectName, dirTmpl+"/README.mkd", tag)

	if strings.HasPrefix(*fLicense, "cc0") {
		renderNewFile(*fProjectName+"/AUTHORS.mkd",
			dirTmpl+"/AUTHORS-cc0.mkd", tag)
	} else {
		renderFile(*fProjectName, dirTmpl+"/AUTHORS.mkd", tag)
		renderFile(*fProjectName, dirTmpl+"/CONTRIBUTORS.mkd", tag)
	}

	// === Add file related to VCS
	switch *fVCS {
	case "other":
		break
	// File CHANGES is only necessary when is not used a VCS.
	case "none":
		renderFile(*fProjectName, dirTmpl+"/CHANGES.mkd", tag)
	default:
		fileIgnore := *fVCS + "ignore"

		if err := copyFile(path.Join(*fProjectName, "."+fileIgnore),
			path.Join(dirTmpl, fileIgnore)); err != nil {
			log.Exit(err)
		}
	}

	// === License file
	addLicense(*fProjectName, tag)

	// === Create file Metadata
	// tag["project_name"] has the original name (no in lower case).
	cfg := NewMetadata(*fProjecType, tag["project_name"], *fPackageName,
		*fLicense, *fAuthor, *fAuthorEmail)

	if err := cfg.WriteINI(*fProjectName); err != nil {
		log.Exit(err)
	}

	// === Print messages
	if tag["author_is_org"] != "" {
		fmt.Print(`
  * The organization has been added as author.
    Update `)

		if tag["license_is_cc0"] != "" {
			fmt.Print("AUTHORS")
		} else {
			fmt.Print("CONTRIBUTORS")
		}
		fmt.Print(" file to add people.\n")
	}
}

/* Updates some values from a project already created. */
func updateProject() {
	var filesUpdated vector.StringVector

	// 'cfg' has the old values.
	cfg, err := ReadMetadata()
	if err != nil {
		log.Exit(err)
	}

	// 'tag' and the flags have the new values.
	tag, update := tagsToUpdate(cfg)
	if *fDebug {
		debug(tag)
	}

	// === Update source code files
	if update["ProjectName"] || update["License"] || update["PackageInCode"] {
		bPackageName := []byte(tag["package_name"])
		files := finderGo(cfg.PackageName)

		for _, fname := range files {
			backup(fname)

			if err := replaceGoFile(fname, bPackageName, cfg, tag, update); err != nil {
				fmt.Fprintf(os.Stderr,
					"%s: file %q not updated: %s\n", argv0, fname, err)
			} else if *fVerbose {
				filesUpdated.Push(fname)
			}
		}

		// === Update Makefile
		fname := path.Join(cfg.PackageName, "Makefile")
		backup(fname)

		if err := replaceMakefile(fname, bPackageName, cfg, tag, update); err != nil {
			fmt.Fprintf(os.Stderr,
				"%s: file %q not updated: %s\n", argv0, fname, err)
		} else if *fVerbose {
			filesUpdated.Push(fname)
		}
	}

	// === Update text files with extension 'mkd'
	if update["ProjectName"] || update["License"] {
		bProjectName := []byte(tag["project_name"])
		files := finderMkd(".")

		for _, fname := range files {
			backup(fname)

			if err := replaceTextFile(fname, bProjectName, cfg, tag, update); err != nil {
				fmt.Fprintf(os.Stderr,
					"%s: file %q not updated: %s\n", argv0, fname, err)
			} else if *fVerbose {
				filesUpdated.Push(fname)
			}
		}
	}

	// === License file
	if update["License"] {
		addLicense(".", tag)

		if *fVerbose {
			filesUpdated.Push("LICENSE")
		}

		cfg.License = *fLicense // Metadata
	}

	// === Print messages
	if *fVerbose {
		fmt.Println("  = Files updated\n")

		for _, file := range filesUpdated {
			fmt.Printf(" * %s\n", file)
		}
	}

	if *fVerbose {
		fmt.Println("\n  = Directories renamed\n")
	}

	// === Rename directories
	if update["PackageName"] {
		if err := os.Rename(cfg.PackageName, *fPackageName); err != nil {
			log.Exit(err)
		} else if *fVerbose {
			fmt.Printf(" * Package: %q -> %q\n", cfg.PackageName, *fPackageName)
		}

		cfg.PackageName = *fPackageName // Metadata
	}

	if update["ProjectName"] {
		if err := os.Chdir(".."); err != nil {
			log.Exit(err)
		}

		cfgProjectName := strings.ToLower(cfg.ProjectName)

		if err := os.Rename(cfgProjectName, *fProjectName); err != nil {
			log.Exit(err)
		} else if *fVerbose {
			fmt.Printf(" * Project: %q -> %q\n", cfgProjectName, *fProjectName)
		}

		cfg.ProjectName = tag["project_name"] // Metadata
	}

	// === File Metadata
	backup(path.Join(*fProjectName, _FILE_NAME))

	if err := cfg.WriteINI(*fProjectName); err != nil {
		log.Exit(err)
	}
}

