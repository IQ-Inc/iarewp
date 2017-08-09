/*Package iarewp defines a library for interacting with IAR's EWP files.

Notes:

- The library does not contain routines for interacting with configurations
or file groups. The structs maintain the data as raw, unmarshalled XML.
The library provides no means of interacting with the data.

- The library has no affiliation with IAR.

Ian McIntyre
*/
package iarewp

import (
	"encoding/xml"
	"sort"
	"strings"
)

const (
	// ProjectDir is the project directory prefix for file paths
	ProjectDir string = "$PROJ_DIR$"
)

// Ewp the top-level type of the EWP
type Ewp struct {
	XMLName xml.Name `xml:"project"`

	// FileVersion is the EWP file version
	FileVersion int `xml:"fileVersion"`

	// Configuration is the configuration tags.
	// They are unused, so we consume all the inner structure
	Configuration []struct {
		Unused string `xml:",innerxml"`
	} `xml:"configuration"`

	// Group is the file groups. Also unused, so we
	// consume all the inner XML here.
	Group []struct {
		Unused string `xml:",innerxml"`
	} `xml:"group"`

	// Files is a list of files
	Files []File `xml:"file"`
}

// File defines the file
type File struct {
	// Name is the file name
	Name string `xml:"name"`

	// Excluded is an optional struct that defines
	// the configurations from which the file is
	// excluded.
	Exclusions *Excluded `xml:"excluded,omitempty"`
}

// Excluded defines the configurations that
// omit the file
type Excluded struct {
	Configurations []string `xml:"configuration"`
}

// ByFileName wraps a slice of files for sorting
// by file name
type ByFileName []File

// Len is a method of the sort interface
func (fs ByFileName) Len() int {
	return len(fs)
}

// Less is a method of the sort interface
func (fs ByFileName) Less(i, j int) bool {
	return fs[i].Name < fs[j].Name
}

// Swap is a method of the sort interface
func (fs ByFileName) Swap(i, j int) {
	fs[i], fs[j] = fs[j], fs[i]
}

// FileName gets the name of the file
func (f *File) FileName() string {
	words := strings.Split(f.Name, "\\")
	return words[len(words)-1]
}

// MakeFile creates a file with zero (nil) or more exclusions
func MakeFile(name string, exclusions ...string) File {
	path := ProjectDir + "\\" + name
	var ex *Excluded

	if exclusions != nil {
		ex = &Excluded{exclusions}
	}

	return File{Name: path, Exclusions: ex}
}

// InsertFile will insert a file into an EWP
// while maintaining file order
func (ewp *Ewp) InsertFile(f File) {
	fs := append(ewp.Files, f)
	sort.Sort(ByFileName(fs))
	ewp.Files = fs
}

// Contains checks if the provided file is in
// the EWP
func (ewp *Ewp) Contains(f File) bool {
	for _, file := range ewp.Files {
		if file.Name == f.Name {
			return true
		}
	}
	return false
}
