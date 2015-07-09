/*
Tests
*/
package main

import (
	"fmt"
	// "strings"
	"testing"
)

func TestIndentConnects(t *testing.T) {
	source := readFileToSlice("./tests_data/source_indentConnects.cpp")
	result := readFileToSlice("./tests_data/result_indentConnects.cpp")

	indentConnects(source)

	for i, line := range source {
		if line != result[i] {
			fmt.Println("|" + line + "|")
			fmt.Println("|" + result[i] + "|")
			t.FailNow()
		}
	}
}
