package iarewp

import (
	"encoding/xml"
	"sort"
	"testing"
)

func TestMakeFile(t *testing.T) {
	exclusions := []string{"configA", "configB"}
	expected := ProjectDir + "\\abc.cpp"

	f := MakeFile("abc.cpp", exclusions...)

	if f.Name != expected {
		t.Fatal("Expected file name", f.Name, "to equal", expected)
	}

	for i := range exclusions {
		if f.Exclusions.Configurations[i] != exclusions[i] {
			t.Fatal("Expected exclusions to be equal")
		}
	}
}
func TestFileName(t *testing.T) {
	f := MakeFile("abc.cpp")
	n := f.FileName()
	if n != "abc.cpp" {
		t.Fatal("Expected file name", n, "to equal abc.cpp")
	}
}

func TestSortByFileName(t *testing.T) {
	main := MakeFile("main.cpp")
	math := MakeFile("math.h")
	boostcpp := MakeFile("boost.cpp")
	boosth := MakeFile("boost.h")
	gmock := MakeFile("gmock\\gmock.h")
	gtest := MakeFile("gtest\\gtest.h")

	fs := []File{
		main,
		math,
		boostcpp,
		boosth,
		gmock,
		gtest,
	}

	sort.Sort(ByFileName(fs))

	expected := []File{
		boostcpp,
		boosth,
		gmock,
		gtest,
		main,
		math,
	}

	for i := range expected {
		if expected[i].FileName() != fs[i].FileName() {
			t.Fatal("Failed to sort collection of files by name")
		}
	}
}

// Heads up: this test is touchy because
// it's whitespace sensitive...
func TestParseUnused(t *testing.T) {
	data := []byte(`
	<project>
		<fileVersion>3</fileVersion>
		<file>
			<name>$PROJ_DIR$\\main.cpp</name>
		</file>
		<file>
			<name>$PROJ_DIR$\\math.h</name>
			<excluded>
				<configuration>MyConfig</configuration>
				<configuration>AnotherConfig</configuration>
			</excluded>
		</file>
		<configuration>
			<content>
				This is additional content
			</content>
		</configuration>
		<configuration>
			<hello>
				hello world
			</hello>
		</configuration>
		<group>
			<world>
				hello world
			</world>
		</group>
	</project>
	`)

	var proj Ewp
	err := xml.Unmarshal(data, &proj)

	if err != nil {
		t.Fatal(err)
	}

	exclusions := [2]string{"MyConfig", "AnotherConfig"}

	files := []File{
		MakeFile("main.cpp"),
		MakeFile("math.h", exclusions[:]...),
	}

	config := []string{
		`
			<content>
				This is additional content
			</content>
		`,
		`
			<hello>
				hello world
			</hello>
		`,
	}

	group := []string{
		`
			<world>
				hello world
			</world>
		`,
	}

	// Check for all files
	for i := range files {
		if proj.Files[i].FileName() != files[i].FileName() {
			t.Fatal("Failed to parse list of files")
		}
	}

	// Check for exclusions to math.h...
	var math *File
	for i := range proj.Files {
		if proj.Files[i].FileName() == "math.h" {
			math = &proj.Files[i]
			break
		}
	}

	if math == nil {
		t.Fatal("Expected to find math.h file")
	}

	for i := range math.Exclusions.Configurations {
		if math.Exclusions.Configurations[i] != exclusions[i] {
			t.Fatal("Failed to maintain file exclusion")
		}
	}

	// Check for file version
	if proj.FileVersion != 3 {
		t.Fatal("Failed to parse file version")
	}

	// Check for configurations
	for i := range config {
		if proj.Configuration[i].Unused != config[i] {
			t.Fatal("Failed to maintain unused configuration")
		}
	}

	// Check for groups
	for i := range group {
		if proj.Group[i].Unused != group[i] {
			t.Fatal("Failed to maintain unused configuration")
		}
	}
}

func TestInsertFile(t *testing.T) {
	main := MakeFile("main.cpp")
	math := MakeFile("math.h")
	boostcpp := MakeFile("boost.cpp")
	boosth := MakeFile("boost.h")
	gmock := MakeFile("gmock\\gmock.h")
	gtest := MakeFile("gtest\\gtest.h")

	fs := []File{
		main,
		math,
		boostcpp,
		boosth,
		gmock,
		gtest,
	}

	addition := MakeFile("foo.cpp")

	expected := append(fs, addition)
	sort.Sort(ByFileName(expected))

	ewp := Ewp{Files: fs}

	ewp.InsertFile(addition)

	for i := range ewp.Files {
		if ewp.Files[i].Name != expected[i].Name {
			t.Fatal("Failed to insert new file into EWP while maintaining order")
		}
	}
}

func TestContains(t *testing.T) {
	main := MakeFile("main.cpp")
	math := MakeFile("math.h")
	boostcpp := MakeFile("boost.cpp")
	boosth := MakeFile("boost.h")
	gmock := MakeFile("gmock\\gmock.h")
	gtest := MakeFile("gtest\\gtest.h")

	fs := []File{
		main,
		math,
		boostcpp,
		boosth,
		gmock,
		gtest,
	}

	ewp := Ewp{Files: fs}

	foo := MakeFile("foo.cpp")

	if ewp.Contains(foo) {
		t.Fatal("Expected EWP to not contain foo.cpp")
	}

	if !ewp.Contains(boosth) {
		t.Fatal("Expected EWP to contain boost.h")
	}
}
