// Copyright 2010  The "gowizard" Authors
//
// Use of this source code is governed by the Simplified BSD License
// that can be found in the LICENSE file.
//
// This software is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES
// OR CONDITIONS OF ANY KIND, either express or implied. See the License
// for more details.

/* Based on Metadata for Python Software Packages.

Description of fields that are not set via 'flag':

* Version: A string containing the package's version number.

* Summary: A one-line summary of what the package does.

* Download-URL: A string containing the URL from which this version of the
	package can be downloaded.

* Platform: A comma-separated list of platform specifications, summarizing
	the operating systems supported by the package which are not listed
	in the "Operating System" Trove classifiers.

* Description: A longer description of the package that can run to several
	paragraphs.

* Keywords: A list of additional keywords to be used to assist searching for
	the package in a larger catalog.

* Home-page: A string containing the URL for the package's home page.

* Classifier: Each entry is a string giving a single classification value
	for the package.

*/

package main

import (
	"os"
	"path"
	"reflect"

	"github.com/kless/goconfig/config"
)


const _META_FILE = "Metadata"
const _VERSION = "1.1"

// Project types
var listProject = map[string]string{
	"tool": "Development tool",
	"app":  "Program",
	"cgo":  "Package that calls C code",
	"lib":  "Library",
}

// Available licenses
var listLicense = map[string]string{
	"apache-2": "Apache License, version 2.0",
	"bsd-2":    "Simplified BSD License",
	"bsd-3":    "New BSD License",
	"cc0-1":    "Creative Commons CC0, version 1.0 Universal",
	"gpl-3":    "GNU General Public License, version 3 or later",
	"agpl-3":   "GNU Affero General Public License, version 3 or later",
	"none":     "Proprietary License",
}

// Version control systems (VCS)
var listVCS = map[string]string{
	"bzr":   "Bazaar",
	"git":   "Git",
	"hg":    "Mercurial",
	"other": "other VCS",
	"none":  "none",
}


// === Errors
type MetadataFieldError string

func (self MetadataFieldError) String() string {
	return "metadata: default section has not field '" + string(self) + "'"
}


/* v1.1 http://www.python.org/dev/peps/pep-0314/

The next fields have not been taken:

	Supported-Platform
	Requires
	Provides
	Obsoletes

Neither the next ones because they are only useful on Python since they are
added to pages on packages index:

	Description
	Classifier

The field 'Name' has been substituted by 'Project-name' and 'Package-name'.
The field 'License' needs a value from the map 'license'.

It has been added 'Project-type', and 'VCS'.

For 'Classifier' see on http://pypi.python.org/pypi?%3Aaction=list_classifiers
*/
type Metadata struct {
	MetadataVersion string "metadata-version" // Version of the file format
	ProjectType     string "project-type"
	ProjectName     string "project-name"
	PackageName     string "package-name"
	Version         string "version"
	Summary         string "summary"
	DownloadURL     string "download-url"
	Author          string "author"
	AuthorEmail     string "author-email"
	License         string "license"
	VCS             string "vcs"

	// === Optional
	//Platform string "platform"
	//Description string "description"
	Keywords string "keywords"
	HomePage string "home-page"
	//Classifier  []string "classifier"

	// Configuration file
	cfg *config.Config
}

/* Creates a new metadata with the basic fields to build the project. */
func NewMetadata(ProjectType, ProjectName, PackageName, License, Author,
AuthorEmail, vcs string) *Metadata {
	_Metadata := new(Metadata)
	_Metadata.cfg = config.NewDefault()

	_Metadata.MetadataVersion = _VERSION
	_Metadata.ProjectType = ProjectType
	_Metadata.ProjectName = ProjectName
	_Metadata.PackageName = PackageName
	_Metadata.License = License
	_Metadata.Author = Author
	_Metadata.AuthorEmail = AuthorEmail
	_Metadata.VCS = vcs

	return _Metadata
}

/* Reads metadata file. */
func ReadMetadata() (*Metadata, os.Error) {
	cfg, err := config.ReadDefault(_META_FILE)
	if err != nil {
		return nil, err
	}

	_Metadata := new(Metadata)

	_Metadata.MetadataVersion = _VERSION
	_Metadata.cfg = cfg

	// === Section 'CORE' has required fields.
	section := "CORE"

	field := "project-type"
	if s, err := cfg.String(section, field); err == nil {
		_Metadata.ProjectType = s
	} else {
		return nil, MetadataFieldError(field)
	}
	field = "project-name"
	if s, err := cfg.String(section, field); err == nil {
		_Metadata.ProjectName = s
	} else {
		return nil, MetadataFieldError(field)
	}
	field = "package-name"
	if s, err := cfg.String(section, field); err == nil {
		_Metadata.PackageName = s
	} else {
		return nil, MetadataFieldError(field)
	}
	field = "license"
	if s, err := cfg.String(section, field); err == nil {
		_Metadata.License = s
	} else {
		return nil, MetadataFieldError(field)
	}
	field = "vcs"
	if s, err := cfg.String(section, field); err == nil {
		_Metadata.VCS = s
	} else {
		return nil, MetadataFieldError(field)
	}

	section = "Main"
	// ===
	if s, err := cfg.String(section, "author"); err == nil {
		_Metadata.Author = s
	}
	if s, err := cfg.String(section, "author-email"); err == nil {
		_Metadata.AuthorEmail = s
	}
	if s, err := cfg.String(section, "version"); err == nil {
		_Metadata.Version = s
	}
	if s, err := cfg.String(section, "summary"); err == nil {
		_Metadata.Summary = s
	}
	if s, err := cfg.String(section, "download-url"); err == nil {
		_Metadata.DownloadURL = s
	}

	section = "Optional"
	// ===
	//if s, err := cfg.String(section, "platform"); err == nil {
		//_Metadata.Platform = s
	//}
	if s, err := cfg.String(section, "keywords"); err == nil {
		_Metadata.Keywords = s
	}
	if s, err := cfg.String(section, "home-page"); err == nil {
		_Metadata.HomePage = s
	}

	return _Metadata, nil
}

/* Serializes to INI format and write it to file `_META_FILE` in directory `dir`.
 */
func (self *Metadata) WriteINI(dir string) os.Error {
	header := "Generated by gowizard"
	reflectMetadata := self.getStruct()

	core := []string{
		"MetadataVersion",
		"ProjectType",
		"ProjectName",
		"PackageName",
		"License",
		"VCS",
	}

	main := []string{
		"Version",
		"Summary",
		"DownloadURL",
		"Author",
		"AuthorEmail",
	}

	optional := []string{
		//"Platform",
		//"Description",
		"HomePage",
		"Keywords",
		//"Classifier",
	}

	for i := 0; i < len(main); i++ {
		name, value := reflectMetadata.name_value(main[i])
		self.cfg.AddOption("Main", name, value)
	}

	for i := 0; i < len(optional); i++ {
		name, value := reflectMetadata.name_value(optional[i])
		self.cfg.AddOption("Optional", name, value)
	}

	for i := 0; i < len(core); i++ {
		name, value := reflectMetadata.name_value(core[i])
		self.cfg.AddOption("CORE", name, value)
	}

	filePath := path.Join(dir, _META_FILE)
	if err := self.cfg.WriteFile(filePath, PERM_FILE, header); err != nil {
		return err
	}

	return nil
}


// === Reflection
// ===

// To handle the reflection of a struct
type reflectStruct struct {
	strType  *reflect.StructType
	strValue *reflect.StructValue
}

/* Gets structs that represent the type 'Metadata'. */
func (self *Metadata) getStruct() *reflectStruct {
	v := reflect.NewValue(self).(*reflect.PtrValue)

	strType := v.Elem().Type().(*reflect.StructType)
	strValue := v.Elem().(*reflect.StructValue)

	return &reflectStruct{strType, strValue}
}

/* Gets tag or field name and its value, given the field name. */
func (self *reflectStruct) name_value(fieldName string) (name, value string) {
	field, _ := self.strType.FieldByName(fieldName)
	value_ := self.strValue.FieldByName(fieldName)

	value = value_.(*reflect.StringValue).Get()

	if tag := field.Tag; tag != "" {
		name = tag
	} else {
		name = field.Name
	}

	return
}

