package meshutil

import "testing"

func TestWildcardPatternMatch(t *testing.T) {
	if WildcardPatternMatch("", "") != true {
		t.Error("err")
	}
	if WildcardPatternMatch("a", "") != false {
		t.Error("err")
	}
	if WildcardPatternMatch("baaabab", "*****ba*****ab") != true {
		t.Error("err")
	}
	if WildcardPatternMatch("baaabab", "baaa?ab") != true {
		t.Error("err")
	}
	if WildcardPatternMatch("baaabab", "ba*a?") != true {
		t.Error("err")
	}
	if WildcardPatternMatch("baaabab", "a*ab") != false {
		t.Error("err")
	}
	if WildcardPatternMatch("aa", "a") != false {
		t.Error("err")
	}
	if WildcardPatternMatch("aa", "aa") != true {
		t.Error("err")
	}
	if WildcardPatternMatch("aaa", "aa") != false {
		t.Error("err")
	}
	if WildcardPatternMatch("aa", "*") != true {
		t.Error("err")
	}
	if WildcardPatternMatch("aa", "a*") != true {
		t.Error("err")
	}
	if WildcardPatternMatch("ab", "?*") != true {
		t.Error("err")
	}
	if WildcardPatternMatch("aab", "c*a*b") != false {
		t.Error("err")
	}
}
