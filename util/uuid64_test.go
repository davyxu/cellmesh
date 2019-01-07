package meshutil

import "testing"

func TestUUIDGenerator_Generate(t *testing.T) {
	gen := NewUUID64Generator()
	gen.AddSeqComponent(3)
	gen.AddTimeComponent(8)
	gen.AddComponent(&UUID64Component{func() uint64 { return 1 }, 2})
	gen.AddComponent(&UUID64Component{func() uint64 { return 2 }, 3})
	t.Logf("%x", gen.Generate())
}

func TestTimeSeqGen(t *testing.T) {
	gen := NewUUID64Generator()
	gen.AddTimeComponent(8)
	gen.AddSeqComponent(8)

	t.Logf("%x", gen.Generate())
}
