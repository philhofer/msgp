package _generated

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/tinylib/msgp/msgp"
	"time"
)

func TestRuneEncodeDecode(t *testing.T) {
	tt := &TestType{}
	r := 'r'
	rp := &r
	tt.Rune = r
	tt.RunePtr = &r
	tt.RunePtrPtr = &rp
	tt.RuneSlice = []rune{'a', 'b', '😳'}

	var buf bytes.Buffer
	wrt := msgp.NewWriter(&buf)
	if err := tt.EncodeMsg(wrt); err != nil {
		t.Errorf("%v", err)
	}
	wrt.Flush()

	var out TestType
	rdr := msgp.NewReader(&buf)
	if err := (&out).DecodeMsg(rdr); err != nil {
		t.Errorf("%v", err)
	}
	if r != out.Rune {
		t.Errorf("rune mismatch: expected %c found %c", r, out.Rune)
	}
	if r != *out.RunePtr {
		t.Errorf("rune ptr mismatch: expected %c found %c", r, *out.RunePtr)
	}
	if r != **out.RunePtrPtr {
		t.Errorf("rune ptr ptr mismatch: expected %c found %c", r, **out.RunePtrPtr)
	}
	if !reflect.DeepEqual(tt.RuneSlice, out.RuneSlice) {
		t.Errorf("rune slice mismatch")
	}
}

func TestRuneMarshalUnmarshal(t *testing.T) {
	tt := &TestType{}
	r := 'r'
	rp := &r
	tt.Rune = r
	tt.RunePtr = &r
	tt.RunePtrPtr = &rp
	tt.RuneSlice = []rune{'a', 'b', '😳'}

	bts, err := tt.MarshalMsg(nil)
	if err != nil {
		t.Errorf("%v", err)
	}

	var out TestType
	if _, err := (&out).UnmarshalMsg(bts); err != nil {
		t.Errorf("%v", err)
	}
	if r != out.Rune {
		t.Errorf("rune mismatch: expected %c found %c", r, out.Rune)
	}
	if r != *out.RunePtr {
		t.Errorf("rune ptr mismatch: expected %c found %c", r, *out.RunePtr)
	}
	if r != **out.RunePtrPtr {
		t.Errorf("rune ptr ptr mismatch: expected %c found %c", r, **out.RunePtrPtr)
	}
	if !reflect.DeepEqual(tt.RuneSlice, out.RuneSlice) {
		t.Errorf("rune slice mismatch")
	}
}

func TestOmitEncodeEmpty(t *testing.T) {
	tt := &TestOmit{}

	var buf bytes.Buffer
	wrt := msgp.NewWriter(&buf)
	if err := tt.EncodeMsg(wrt); err != nil {
		t.Errorf("%v", err)
	}
	wrt.Flush()

	bts := buf.Bytes()
	// Expected empty map
	if len(bts) != 1 || bts[0] != 0x80 {
		t.Errorf("data mismatch: expected empty map object, found %v", bts)
	}
}

func TestOmitEncodePart(t *testing.T) {
	now := time.Now()
	tt := &TestOmit{
		Time: now,
		Byte: 42,
		Slice1: []string{
			"foo", "bar",
		},
	}
	tp := &TestPart{
		Time: now,
		Byte: 42,
		Slice1: []string{
			"foo", "bar",
		},
	}

	var buf bytes.Buffer
	wrt := msgp.NewWriter(&buf)
	if err := tt.EncodeMsg(wrt); err != nil {
		t.Errorf("%v", err)
	}
	wrt.Flush()

	exp, err := tp.MarshalMsg(nil)
	if err != nil {
		t.Errorf("%v", err)
	}

	bts := buf.Bytes()
	// Expected same binary data
	if !bytes.Equal(exp, bts) {
		t.Errorf("data mismatch: expected %v, found %v", exp, bts)
	}
}

func TestOmitEncodeDecode(t *testing.T) {
	tt := &TestOmit{}
	r := 'r'
	rp := &r
	tt.Rune = r
	tt.RunePtr = &r
	tt.RunePtrPtr = &rp
	tt.RuneSlice = []rune{'a', 'b', '😳'}

	var buf bytes.Buffer
	wrt := msgp.NewWriter(&buf)
	if err := tt.EncodeMsg(wrt); err != nil {
		t.Errorf("%v", err)
	}
	wrt.Flush()

	var out TestOmit
	rdr := msgp.NewReader(&buf)
	if err := (&out).DecodeMsg(rdr); err != nil {
		t.Errorf("%v", err)
	}
	if r != out.Rune {
		t.Errorf("rune mismatch: expected %c found %c", r, out.Rune)
	}
	if r != *out.RunePtr {
		t.Errorf("rune ptr mismatch: expected %c found %c", r, *out.RunePtr)
	}
	if r != **out.RunePtrPtr {
		t.Errorf("rune ptr ptr mismatch: expected %c found %c", r, **out.RunePtrPtr)
	}
	if !reflect.DeepEqual(tt.RuneSlice, out.RuneSlice) {
		t.Errorf("rune slice mismatch")
	}
}

func TestOmitMarshalEmpty(t *testing.T) {
	tt := &TestOmit{}

	bts, err := tt.MarshalMsg(nil)
	if err != nil {
		t.Errorf("%v", err)
	}

	// Expected empty map
	if len(bts) != 1 || bts[0] != 0x80 {
		t.Errorf("data mismatch: expected empty map object, found %v", bts)
	}
}

func TestOmitMarshalPart(t *testing.T) {
	now := time.Now()
	tt := &TestOmit{
		Time: now,
		Byte: 42,
		Slice1: []string{
			"foo", "bar",
		},
	}
	tp := &TestPart{
		Time: now,
		Byte: 42,
		Slice1: []string{
			"foo", "bar",
		},
	}

	bts, err := tt.MarshalMsg(nil)
	if err != nil {
		t.Errorf("%v", err)
	}

	exp, err := tp.MarshalMsg(nil)
	if err != nil {
		t.Errorf("%v", err)
	}

	// Expected same bi
	if !bytes.Equal(exp, bts) {
		t.Errorf("data mismatch: expected %v, found %v", exp, bts)
	}
}

func TestOmitMarshalUnmarshal(t *testing.T) {
	tt := &TestOmit{}
	r := 'r'
	rp := &r
	tt.Rune = r
	tt.RunePtr = &r
	tt.RunePtrPtr = &rp
	tt.RuneSlice = []rune{'a', 'b', '😳'}

	bts, err := tt.MarshalMsg(nil)
	if err != nil {
		t.Errorf("%v", err)
	}

	var out TestOmit
	if _, err := (&out).UnmarshalMsg(bts); err != nil {
		t.Errorf("%v", err)
	}
	if r != out.Rune {
		t.Errorf("rune mismatch: expected %c found %c", r, out.Rune)
	}
	if r != *out.RunePtr {
		t.Errorf("rune ptr mismatch: expected %c found %c", r, *out.RunePtr)
	}
	if r != **out.RunePtrPtr {
		t.Errorf("rune ptr ptr mismatch: expected %c found %c", r, **out.RunePtrPtr)
	}
	if !reflect.DeepEqual(tt.RuneSlice, out.RuneSlice) {
		t.Errorf("rune slice mismatch")
	}
}
