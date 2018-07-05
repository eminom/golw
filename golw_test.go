package lwm2m

import (
	"math"
	"math/rand"
	"sort"
	"testing"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func getAttrs() []DataItem {
	return []DataItem{
		NewString(0, "Manufacture-XB"),
		NewString(1, "Model-XB"),
		NewString(16, "Binding-U"),
		NewInteger(0x2018, 0x20180624),
		NewFloat(0x2019, math.Pi),
	}
}

func genChunk(uri *UriT) ([]byte, error) {
	return EncodeTlv(uri, getAttrs())
}

func TestLwm2mSort(t *testing.T) {
	var arr []DataItem
	for i := 0; i < 300000; i++ {
		arr = append(arr, NewString(rand.Intn(90000), "X"))
	}
	sort.Sort(DataItemArray(arr))
	for i := 0; i < len(arr)-1; i++ {
		if arr[i].ID > arr[i+1].ID {
			t.Fatalf("error sorting")
		}
	}
	t.Logf("sort tested one")
}

func TestLwm2mObjects(t *testing.T) {

	chunk, err := EncodeTlv(nil, []DataItem{
		NewArray(0, getAttrs()),
		NewArray(1, getAttrs()),
		NewArray(2, getAttrs()),
	})
	if err != nil {
		t.Fatalf("error encode object array: %v", err)
	}

	items, left := ParseTlv(chunk)
	if left != 0 {
		t.Fatalf("error parsing: left %v", left)
	}
	for i, item := range items {
		if item.Type != TypeObjectInstance {
			t.Fatalf("error type: %v", item.Type)
		}
		t.Logf("<%v>", i)
		for _, attr := range item.AsArray() {
			t.Logf("  <%v> %v", attr.ID, attr)
		}
	}
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

	t.Logf("0:<id=%v> %v", items[0].ID, items[0].AsString())
	t.Logf("1:<id=%v> %v", items[1].ID, items[1].AsString())
	t.Logf("2:<id=%v> %v", items[2].ID, items[2].AsString())
	t.Logf("3:<id=%v> 0x%x", items[3].ID, items[3].AsInteger())
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

func TestLwm2mEnc3(t *testing.T) {
	inner := []DataItem{
		NewString(0, "Manufacture-XB"),
		NewString(1, "Model-XB"),
		NewString(16, "Binding-U"),
		NewInteger(0x2018, 0x20180624),
		NewFloat(0x2019, math.Pi),
	}

	chunk, err := EncodeTlv(nil, []DataItem{
		NewArray(20, inner),
		NewString(18, "tail"),
	})
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	items, left := ParseTlv(chunk)
	if left != 0 {
		t.Fatalf("error parsing tlv: %v byte(s) left", left)
	}

	t.Logf("%v", items)

	t.Logf("[0] %v", items[0])
	t.Logf("%v", items[1].AsString())
}
