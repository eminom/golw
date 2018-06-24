package lwm2m

import (
	"math"
	"testing"
)

func genChunk(uri *UriT) ([]byte, error) {
	return EncodeTlv(uri, []DataItem{
		NewString(0, "Manufacture-XB"),
		NewString(1, "Model-XB"),
		NewString(16, "Binding-U"),
		NewInteger(0x2018, 0x20180624),
		NewFloat(0x2019, math.Pi),
	})
}

func TestLwm2m(t *testing.T) {
	chunk, err := genChunk(nil)
	if err != nil {
		t.Fatalf("error marshaling: %v", err)
	}
	t.Logf("Test One: %x", chunk)
	t.Logf("encoded length: %v", len(chunk))

	items, left := ParseTlv(chunk)
	if left != 0 {
		t.Fatalf("error parsing tlv: %v byte(s) left", left)
	}
	t.Logf("%v", items)

	t.Logf("0: %v", items[0].AsString())
	t.Logf("1: %v", items[1].AsString())
	t.Logf("2: %v", items[2].AsString())
	t.Logf("3: 0x%x", items[3].AsInteger())
}

func TestLwm2mEnc2(t *testing.T) {
	a := UriT{}
	a.SetResourceID(100)
	chunk, err := genChunk(&a)
	if err != nil {
		t.Fatalf("error encoding: %v", err)
	}
	t.Logf("Test Two: %x", chunk)
	t.Logf("encode length: %v", len(chunk))

	items, left := ParseTlv(chunk)
	if left != 0 {
		t.Fatalf("error parsing tlv: %v byte(s) left", left)
	}
	t.Logf("%v", items)

	t.Logf("0: %v", items[0].AsString())
	t.Logf("1: %v", items[1].AsString())
	t.Logf("2: %v", items[2].AsString())
	t.Logf("3: 0x%x", items[3].AsInteger())

}
