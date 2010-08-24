// Copyright 2010  The "gowizard" Authors
//
// Use of this source code is governed by the Simplified BSD License
// that can be found in the LICENSE file.
//
// This software is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES
// OR CONDITIONS OF ANY KIND, either express or implied. See the License
// for more details.

/* Based on Metadata for Python Software Packages. */

package main

import (
	"log"
	"path"
	"reflect"

	"github.com/kless/goconfig/config"
)


const _FILE_NAME = "Metadata"

// Available application types
var listProject = map[string]string{
	"cmd":    "command line",
	"pkg":    "package",
	"web.go": "web environment",
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

The field 'Name' has been substituted by 'Project-name' and 'Application-name'.
The field 'License' needs a value from the map 'license'.

It has been added 'Project-type'.

For 'Classifier' see on http://pypi.python.org/pypi?%3Aaction=list_classifiers
*/
type metadata struct {
	MetadataVersion string "Metadata-Version" // Version of the file format
	ProjectType     string "Project-type"
	ProjectName     string "Project-name"
	ApplicationName string "Application-name"
	Version         string
	Summary         string
	DownloadURL     string "Download-URL"
	Author          string
	AuthorEmail     string "Author-email"
	License         string

	// === Optional
	Platform string
	//Description string
	Keywords string
	HomePage string "Home-page"
	//Classifier  []string

	// Config file
	file *config.File
}

/* Creates a new metadata with the basic fields to build the project. */
func NewMetadata(ProjectType, ProjectName, ApplicationName, Author,
AuthorEmail, License string, file *config.File) *metadata {
	metadata := new(metadata)

	metadata.MetadataVersion = "1.1"
	metadata.ProjectType = ProjectType
	metadata.ProjectName = ProjectName
	metadata.ApplicationName = ApplicationName
	metadata.Author = Author
	metadata.AuthorEmail = AuthorEmail
	metadata.License = License

	metadata.file = file

	return metadata
}

/* Serializes to INI format and write it to file `_FILE_NAME` in directory `dir`.
 */
func (self *metadata) WriteINI(dir string) {
	header := "Generated by gowizard"
	reflectMetadata := self.getStruct()

	default_ := []string{
		"MetadataVersion",
		"ProjectName",
		"ApplicationName",
		"License",
	}

	base := []string{
		"ProjectType",
		"Version",
		"Summary",
		"DownloadURL",
		"Author",
		"AuthorEmail",
	}

	optional := []string{
		"Platform",
		//"Description",
		"Keywords",
		"HomePage",
		//"Classifier",
	}

	for i := 0; i < len(default_); i++ {
		name, value := reflectMetadata.name_value(default_[i])
		self.file.AddOption("", name, value)
	}

	for i := 0; i < len(base); i++ {
		name, value := reflectMetadata.name_value(base[i])
		self.file.AddOption("base", name, value)
	}

	for i := 0; i < len(optional); i++ {
		name, value := reflectMetadata.name_value(optional[i])
		self.file.AddOption("optional", name, value)
	}

	filePath := path.Join(dir, _FILE_NAME)
	if err := self.file.WriteFile(filePath, PERM_FILE, header); err != nil {
		log.Exit(err)
	}
}

func ReadMetadata() {
	file, err := config.ReadFile(_FILE_NAME)
	if err != nil {
		log.Exit(err)
	}

	s := file.Sections()
	println(s)
}


// === Reflection
// ===

// To handle the reflection of a struct
type reflectStruct struct {
	strType  *reflect.StructType
	strValue *reflect.StructValue
}

/* Gets the structs that represent the type 'metadata'. */
func (self *metadata) getStruct() *reflectStruct {
	v := reflect.NewValue(self).(*reflect.PtrValue)

	strType := v.Elem().Type().(*reflect.StructType)
	strValue := v.Elem().(*reflect.StructValue)

	return &reflectStruct{strType, strValue}
}

/* Gets the tag or field name and its value, given the field name. */
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

