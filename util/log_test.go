package meshutil

import "testing"

func TestSize(t *testing.T) {

	t.Log(ParseSizeString("10k"))
	t.Log(ParseSizeString("10M"))
	t.Log(ParseSizeString("10g"))

	t.Log(ParseSizeString("10"))

}
