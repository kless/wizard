// Copyright 2010, The 'gowizard' Authors.  All rights reserved.
// Use of this source code is governed by the Simplified BSD License
// that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"
)


// Exit status code if there is any error
const ERROR = 2

// Permissions
const _PERM_DIRECTORY = 0755

// Licenses available
var license = map[string]string{
	"apache": "Apache (version 2.0)",
	"bsd-2":  "Simplified BSD",
	"bsd-3":  "New BSD",
	"cc0":    "Creative Commons CC0 1.0 Universal",
}

// Flags for the command line
var (
	fDebug   = flag.Bool("d", false, "debug mode")
	fList    = flag.Bool("l", false, "show the list of licenses for the flag `license`")
	fWeb     = flag.Bool("w", false, "web application")
	fProject = flag.String("project", "", "name of the project (e.g. 'goweb-foo')")
	fPkg     = flag.String("pkg", "", "name of the package (e.g. 'foo')")
	fLicense = flag.String("license", "bsd-2", "kind of license")
)

// Headers for source code files
const (
	t_HEADER     = `// Copyright {year}, The '{project}' Authors.  All rights reserved.
// Use of this source code is governed by the {license} License
// that can be found in the LICENSE file.
`
	t_HEADER_CC0 = `// To the extent possible under law, Authors have waived all copyright and
// related or neighboring rights to '{project}'.
`
)

// === Template and data to build the file

const t_PAGE = "{header}\n{content}"

type page struct {
	header  string
	content string
}


func checkFlags() {
	usage := func() {
		fmt.Fprintf(os.Stderr, "Usage: gowizard -project [-license]\n\n")
		flag.PrintDefaults()
		os.Exit(ERROR)
	}
	flag.Usage = usage
	flag.Parse()

	reGo := regexp.MustCompile(`^go`) // To remove it from the project name

	if *fList {
		fmt.Printf("Licenses\n\n")
		for k, v := range license {
			fmt.Printf("  %s: %s\n", k, v)
		}
		os.Exit(0)
	}

	*fLicense = strings.ToLower(*fLicense)
	if _, present := license[*fLicense]; !present {
		log.Exitf("license unavailable %s", *fLicense)
	}

	if *fProject == "" {
		usage()
	}
	*fProject = strings.TrimSpace(*fProject)

	if *fPkg == "" {
		// The package name is created:
		// getting the last string after of the dash ('-'), if any,
		// and removing 'go'. Finally, it's lower cased.
		pkg := strings.Split(*fProject, "-", -1)
		*fPkg = reGo.ReplaceAllString(strings.ToLower(pkg[len(pkg)-1]), "")
	} else {
		*fPkg = strings.ToLower(strings.TrimSpace(*fPkg))
	}

	return
}

func addApp() {
	
}


// Main program execution
func main() {
	var renderedHeader string

	checkFlags()

	// Tags to pass to the templates
	tag := map[string]string{
		"license": license[*fLicense],
		"pkg":     *fPkg,
		"project": *fProject,
	}

	// === Renders the header

	if *fLicense == "cc0" {
		renderedHeader = parse(t_HEADER_CC0, tag)
	} else {
		tag["year"] = strconv.Itoa64(time.LocalTime().Year)
		renderedHeader = parse(t_HEADER, tag)
	}

	if *fDebug {
		fmt.Printf("Debug\n\n")
		for k, v := range tag {
			fmt.Printf("  %s: %s\n", k, v)
		}
		os.Exit(0)
	}

	// === Creates directories
	os.MkdirAll(path.Join(strings.ToLower(*fProject), *fPkg), _PERM_DIRECTORY)

	// === Renders files for normal project

	if !*fWeb {

	} else {

	}

	renderedContent := parseFile("web-setup", tag)

	tagPage := &page{
		header: renderedHeader,
		content: renderedContent,
	}

	end := parse(t_PAGE, tagPage)
	fmt.Println(end)
}

