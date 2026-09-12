package goseg_test

import (
	"fmt"

	"goseg"
)

func ExampleSegment() {
	text := "Hello world. My name is Mr. Smith. I work for the U.S. Government and I live in the U.S. I live in New York."
	for _, s := range goseg.Segment(text) {
		fmt.Println(s)
	}
	// Output:
	// Hello world.
	// My name is Mr. Smith.
	// I work for the U.S. Government and I live in the U.S.
	// I live in New York.
}

func ExampleSegment_language() {
	text := "Այսօր երկուշաբթի է: Ես գնում եմ աշխատանքի:"
	for _, s := range goseg.Segment(text, goseg.Language("hy")) {
		fmt.Println(s)
	}
	// Output:
	// Այսօր երկուշաբթի է:
	// Ես գնում եմ աշխատանքի:
}
