// Copyright 2010, The "gowizard" Authors.  All rights reserved.
// Use of this source code is governed by the Simplified BSD License
// that can be found in the LICENSE file.

/* Based on Metadata for Python Software Packages. */

package main

import (
	"log"
	"path"
	"reflect"

	conf "goconf.googlecode.com/hg"
)


const _FILE_NAME = "Metadata"

// Available application types
var listApp = map[string]string{
	"cmd":    "command line",
	"pkg":    "package",
	"web.go": "web environment",
}

// Available licenses
var listLicense = map[string]string{
	"apache-2": "Apache License (version 2.0)",
	"bsd-2":    "Simplified BSD License",
	"bsd-3":    "New BSD License",
	"cc0-1":    "Creative Commons CC0 1.0 Universal",
	"gpl-3":    "GNU General Public License",
	"agpl-3":   "GNU Affero General Public License",
	"none":     "Proprietary License",
}


/* v1.1 http://www.python.org/dev/peps/pep-0314/

The next fields have not been taken:

	Supported-Platform
	Requires
	Provides
	Obsoletes

The field 'Name' has been substituted by 'Project-name' and 'Application-name'.
The field 'License' needs a value from the map 'license'.

It has been added 'Application-type'.

For 'Classifier' see on http://pypi.python.org/pypi?%3Aaction=list_classifiers
*/
type metadata struct {
	MetadataVersion string "Metadata-Version" // Version of the file format
	ProjectName     string "Project-name"
	ApplicationName string "Application-name"
	ApplicationType string "Application-type"
	Version         string
	Summary         string
	DownloadURL     string "Download-URL"
	Author          string
	AuthorEmail     string "Author-email"
	License         string

	// === Optional
	Platform    string
	Description string
	Keywords    string
	HomePage    string "Home-page"
	Classifier  []string

	// Config file
	file *conf.ConfigFile
}

/* Creates a new metadata with the basic fields to build the project. */
func NewMetadata(ProjectName, ApplicationName, ApplicationType, Author,
AuthorEmail, License string, file *conf.ConfigFile) *metadata {
	metadata := new(metadata)

	metadata.MetadataVersion = "1.1"
	metadata.ProjectName = ProjectName
	metadata.ApplicationName = ApplicationName
	metadata.ApplicationType = ApplicationType
	metadata.Author = Author
	metadata.AuthorEmail = AuthorEmail
	metadata.License = License

	metadata.file = file

	return metadata
}

/* Serializes to INI format and write it to file `_FILE_NAME` in directory `dir`.
 */
func (self *metadata) WriteINI(dir string) {
	var name, value string

	header := "Created by gowizard\n"
	reflectMetadata := self.getStruct()

	default_ := []string{
		"MetadataVersion",
		"ProjectName",
		"ApplicationName",
		"Summary",
		"License",
	}

	base := []string{
		"ApplicationType",
		"Version",
		"DownloadURL",
		"Author",
		"AuthorEmail",
	}

	optional := []string{
		"Platform",
		"Description",
		"Keywords",
		"HomePage",
		//"Classifier",
	}

	for i := 0; i < len(default_); i++ {
		name, value = reflectMetadata.getName_Value(default_[i])
		self.file.AddOption(conf.DefaultSection, name, value)
	}

	for i := 0; i < len(base); i++ {
		name, value = reflectMetadata.getName_Value(base[i])
		self.file.AddOption("base", name, value)
	}

	for i := 0; i < len(optional); i++ {
		name, value = reflectMetadata.getName_Value(optional[i])
		self.file.AddOption("optional", name, value)
	}

	filePath := path.Join(dir, _FILE_NAME)
	if err := self.file.WriteConfigFile(filePath, PERM_FILE, header); err != nil {
		log.Exit(err)
	}
}

func ReadMetadata() {
	file, err := conf.ReadConfigFile(_FILE_NAME)
	if err != nil {
		log.Exit(err)
	}

	s := file.GetSections()
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
func (self *reflectStruct) getName_Value(fieldName string) (name, value string) {
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

