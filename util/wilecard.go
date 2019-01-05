package meshutil

// https://siongui.github.io/2017/04/11/go-wildcard-pattern-matching/
func stringToRuneSlice(s string) []rune {
	var r []rune
	for _, runeValue := range s {
		r = append(r, runeValue)
	}
	return r
}

func initLookupTable(row, column int) [][]bool {
	lookup := make([][]bool, row)
	for i := range lookup {
		lookup[i] = make([]bool, column)
	}
	return lookup
}

// Function that matches input str with given wildcard pattern
func WildcardPatternMatch(str, pattern string) bool {
	s := stringToRuneSlice(str)
	p := stringToRuneSlice(pattern)

	// empty pattern can only match with empty string
	if len(p) == 0 {
		return len(s) == 0
	}

	// lookup table for storing results of subproblems
	// zero value of lookup is false
	lookup := initLookupTable(len(s)+1, len(p)+1)

	// empty pattern can match with empty string
	lookup[0][0] = true

	// Only '*' can match with empty string
	for j := 1; j < len(p)+1; j++ {
		if p[j-1] == '*' {
			lookup[0][j] = lookup[0][j-1]
		}
	}

	// fill the table in bottom-up fashion
	for i := 1; i < len(s)+1; i++ {
		for j := 1; j < len(p)+1; j++ {
			if p[j-1] == '*' {
				// Two cases if we see a '*'
				// a) We ignore ‘*’ character and move
				//    to next  character in the pattern,
				//     i.e., ‘*’ indicates an empty sequence.
				// b) '*' character matches with ith
				//     character in input
				lookup[i][j] = lookup[i][j-1] || lookup[i-1][j]

			} else if p[j-1] == '?' || s[i-1] == p[j-1] {
				// Current characters are considered as
				// matching in two cases
				// (a) current character of pattern is '?'
				// (b) characters actually match
				lookup[i][j] = lookup[i-1][j-1]

			} else {
				// If characters don't match
				lookup[i][j] = false
			}
		}
	}

	return lookup[len(s)][len(p)]
}
