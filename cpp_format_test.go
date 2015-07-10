/*
Tests
*/
package main

import (
	"fmt"
	"github.com/fatih/color"
	"testing"
)

func TestIndentConnects(t *testing.T) {
	source := readFileToSlice("./tests_data/source_indentConnects.cpp")
	result := readFileToSlice("./tests_data/result_indentConnects.cpp")

	indentConnects(source)

	for i, line := range source {
		if line != result[i] {
			color.Set(color.FgYellow)
			fmt.Printf("\nTheses lines are different: \n%s\n%s\n\n", line, result[i])
			color.Set(color.FgRed)
			t.FailNow()
			color.Unset()
		}
	}
}
