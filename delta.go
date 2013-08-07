// Delta is a simple package for calculating differences between variants of
// text on a single word level. It can output HTML or plain text. It is built off
// of the pseudocode for the Longest Common Subsequence problem found on
// http://en.wikipedia.org/wiki/Longest_common_subsequence_problem
//
// Examples:
//
//  delta.Calculate("hello world", "hello earth", false)
//		// "hello <del>world</del> <ins>earth</ins>"
//
//	delta.Calculate("hello world", "hello earth", true)
//		// "hello ---world--- +++earth+++"
package delta

import (
	"html"
	"regexp"
	"strings"
)

var (
	regexpNewline = regexp.MustCompile(`\r\n?`)
	regexpDouble = regexp.MustCompile(`(\s*?)&__DOUBLE__;(\s*)`)
	regexpSingle = regexp.MustCompile(`(\s*?)&__SINGLE__;(\s*)`)
)

// Calculate accepts the two revisions of text, first one being the previous
// (older) and second being the current (newer) version. It returns the string
// representation of the diff, in either HTML or plain text.
func Calculate(prev, curr string, plaintext bool) string {
	p, c := preprocess(prev), preprocess(curr)

	return postprocess(print(sequence(p, c), p, c, len(p)-1, len(c)-1), plaintext)
}

// sequence builds the necessary matrix and computes the length of it. It
// also reads through the matrix and computes the longest common subsequence.
func sequence(prev, curr []string) map[int]map[int]int {
	c := make(map[int]map[int]int)

	for i := -1; i <= len(prev); i++ {
		// For some reason, when making multidimensional matrices, Go
		// only makes the outermost one, and all deeper ones need to
		// be created manually. An afternoon well spent.
		_, ok := c[i]
		if !ok {
			c[i] = make(map[int]int)
		}

		for j := -1; j <= len(curr); j++ {
			c[i][j] = 1
		}
	}

	for i := 0; i <= len(prev)-1; i++ {
		for j := 0; j <= len(curr)-1; j++ {
			if prev[i] == curr[j] {
				c[i][j] = c[i-1][j-1] + 1
			} else {
				var max int
				if c[i][j-1] > c[i-1][j] {
					max = c[i][j-1]
				} else {
					max = c[i-1][j]
				}
				c[i][j] = max
			}
		}
	}

	return c
}

// print backtracks over the input sequences and prints out the calculated
// differences between them. The output is HTML which is later processed and
// cleaned up.
func print(c map[int]map[int]int, prev, curr []string, i, j int) string {
	var output string

	if i >= 0 && j >= 0 && prev[i] == curr[j] {
		output += print(c, prev, curr, i-1, j-1)
		output += string(prev[i]) + " "
	} else {
		if j >= 0 && (i == -1 || c[i][j-1] >= c[i-1][j]) {
			output += print(c, prev, curr, i, j-1)
			output += "<ins>" + string(curr[j]) + "</ins> "
		} else if i >= 0 && (j == -1 || c[i][j-1] < c[i-1][j]) {
			output += print(c, prev, curr, i-1, j)
			output += "<del>" + string(prev[i]) + "</del> "
		}
	}

	return output
}

// preprocess normalizes new lines, escapes the input, and replacesv new
// lines so that changes that span across more lines get caught as such
// and treated accordingly.
func preprocess(input string) []string {
	input = regexpNewline.ReplaceAllString(input, "\n")
	input = strings.TrimSpace(html.EscapeString(input))
	input = strings.Replace(input, "\n\n", " &__DOUBLE__; ", -1)
	input = strings.Replace(input, "\n", " &__SINGLE__; ", -1)
	return strings.Split(input, " ")
}

// postprocess finalizes the output by merging adjacent HTML tags, returning
// the previously removed new lines, and converting between HTML and text.
func postprocess(input string, plaintext bool) string {
	input = strings.Replace(input, "</del> <del>", " ", -1)
	input = strings.Replace(input, "</ins> <ins>", " ", -1)

	// Plain text is different than HTML in way that HTML variant
	// uses the <ins> and <del> tags, while plain text variant
	// uses three + and - characters to wrap added and removed
	// pieces of text.
	if plaintext {
		input = strings.Replace(input, "<ins>", "+++", -1)
		input = strings.Replace(input, "</ins>", "+++", -1)
		input = strings.Replace(input, "<del>", "---", -1)
		input = strings.Replace(input, "</del>", "---", -1)
	}

	input = regexpDouble.ReplaceAllString(input, "\n\n")
	input = regexpSingle.ReplaceAllString(input, "\n")
	return strings.TrimSpace(input)
}
