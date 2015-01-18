package bencode

import (
	"reflect"
	"strings"
	"testing"
)

func TestDecodeByteString(t *testing.T) {
	stringReader := strings.NewReader("4:spam")
	bencode, _ := Decode(stringReader)

	if bencode != "spam" {
		t.Errorf("for %q expected %q but got %q", "4:spam", "spam", bencode.(string))
	}
}

func TestInteger(t *testing.T) {
	intReader := strings.NewReader("i345e")
	bencode, _ := Decode(intReader)
	if bencode != 345 {
		t.Errorf("for %q expected %q but got %q", "i345e", "345", bencode)
	}

}

func TestList(t *testing.T) {
	listReader := strings.NewReader("l4:spami345ee")
	bencode, _ := Decode(listReader)
	bencodeSlice := reflect.ValueOf(bencode)
	if bencodeSlice.Len() < 2 {
		t.Errorf("for %q expected %q but got %d", "l4:spami34e", "2", bencodeSlice.Len())
	}

	stringVal := bencodeSlice.Index(0).Interface().(string)
	if stringVal != "spam" {
		t.Errorf("for %q expected %q but got %s", "l4:spami34e", "spam", stringVal)

	}

	intVal := bencodeSlice.Index(1).Interface().(int)
	if intVal != 345 {
		t.Errorf("for %q expected %q but got %d", "l4:spami34e", "345", intVal)

	}

}

func TestEmptyList(t *testing.T) {
	listReader := strings.NewReader("le")
	bencode, _ := Decode(listReader)
	bencodeSlice := reflect.ValueOf(bencode)
	if bencodeSlice.Len() != 0 {
		t.Errorf("for %q expected %q but got %s", "le", "0", bencodeSlice.Len())
	}
}

func TestError(t *testing.T) {
	listReader := strings.NewReader("l")
	_, err := Decode(listReader)
	if err == nil {
		t.Errorf("for %q expected %q but got %s", "l", "error", "nil")
	}
}

func TestDict(t *testing.T) {
	dictReader := strings.NewReader("d3:cow3:moo4:spam4:eggse")
	bencode, _ := Decode(dictReader)
	bencodeMap := reflect.ValueOf(bencode)

	value := bencodeMap.MapIndex(reflect.ValueOf("cow"))
	if value.Interface().(string) != "moo" {
		t.Errorf("for %q expected %q but got %s", "d3:cow3:moo4:spam4:eggse", "moo", value)

	}

	dictReader = strings.NewReader("d4:spaml1:a1:bee")
	bencode, _ = Decode(dictReader)
	bencodeMap = reflect.ValueOf(bencode)

	valueSlice := bencodeMap.MapIndex(reflect.ValueOf("spam"))
	stringVal := reflect.ValueOf(valueSlice.Interface()).Index(0).Interface().(string)
	if stringVal != "a" {
		t.Errorf("for %q expected %q but got %s", "d4:spaml1:a1:bee", "a", stringVal)

	}

	stringVal = reflect.ValueOf(valueSlice.Interface()).Index(1).Interface().(string)
	if stringVal != "b" {
		t.Errorf("for %q expected %q but got %s", "d4:spaml1:a1:bee", "b", stringVal)

	}
}

func TestEmptyDict(t *testing.T) {
	dictReader := strings.NewReader("de")
	bencode, _ := Decode(dictReader)
	bencodeMap := reflect.ValueOf(bencode)
	keys := bencodeMap.MapKeys()

	if len(keys) != 0 {
		t.Errorf("for %q expected %q but got %d keys", "de", "0 keys", len(keys))

	}
}
