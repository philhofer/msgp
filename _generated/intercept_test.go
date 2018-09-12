package _generated

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/tinylib/msgp/msgp"
)

func TestInterceptEncodeDecodeStruct(t *testing.T) {
	resetStructProvider()

	in := TestUsesStructProvided{Foo: &TestStructProvided{Foo: "hi"}}

	var buf bytes.Buffer
	wrt := msgp.NewWriter(&buf)
	if err := in.EncodeMsg(wrt); err != nil {
		t.Errorf("%v", err)
	}
	wrt.Flush()

	var out TestUsesStructProvided
	rdr := msgp.NewReader(&buf)
	if err := (&out).DecodeMsg(rdr); err != nil {
		t.Errorf("%v", err)
	}

	if !reflect.DeepEqual(in, out) {
		t.Fatalf("provided encode decode failed")
	}

	if !reflect.DeepEqual([]string{"encode", "decode"}, TestStructProvider().Events) {
		t.Fatalf("unexpected events: %v", TestStructProvider().Events)
	}
}

func TestInterceptMarshalUnmarshalStruct(t *testing.T) {
	resetStructProvider()

	in := TestUsesStructProvided{Foo: &TestStructProvided{Foo: "hi"}}

	bts, err := in.MarshalMsg(nil)
	if err != nil {
		t.Fatalf("%v", err)
	}

	var out TestUsesStructProvided
	if _, err := (&out).UnmarshalMsg(bts); err != nil {
		t.Fatalf("%v", err)
	}

	if !reflect.DeepEqual(in, out) {
		t.Fatalf("provided unmarshal failed")
	}

	if !reflect.DeepEqual([]string{"msgsize", "marshal", "unmarshal"}, TestStructProvider().Events) {
		t.Fatalf("unexpected events: %v", TestStructProvider().Events)
	}
}

func TestInterceptEncodeDecodeString(t *testing.T) {
	resetStringProvider()

	in := TestUsesStringProvided{Foo: TestStringProvided("hi")}

	var buf bytes.Buffer
	wrt := msgp.NewWriter(&buf)
	if err := in.EncodeMsg(wrt); err != nil {
		t.Errorf("%v", err)
	}
	wrt.Flush()

	var out TestUsesStringProvided
	rdr := msgp.NewReader(&buf)
	if err := (&out).DecodeMsg(rdr); err != nil {
		t.Errorf("%v", err)
	}

	if !reflect.DeepEqual(in, out) {
		t.Fatalf("provided encode decode failed")
	}

	if !reflect.DeepEqual([]string{"encode", "decode"}, TestStringProvider().Events) {
		t.Fatalf("unexpected events: %v", TestStringProvider().Events)
	}
}

func TestInterceptMarshalUnmarshalString(t *testing.T) {
	resetStringProvider()

	in := TestUsesStringProvided{Foo: TestStringProvided("hi")}

	bts, err := in.MarshalMsg(nil)
	if err != nil {
		t.Fatalf("%v", err)
	}

	var out TestUsesStringProvided
	if _, err := (&out).UnmarshalMsg(bts); err != nil {
		t.Fatalf("%v", err)
	}

	if !reflect.DeepEqual(in, out) {
		t.Fatalf("provided unmarshal failed")
	}

	if !reflect.DeepEqual([]string{"marshal", "unmarshal"}, TestStringProvider().Events) {
		t.Fatalf("unexpected events, found: %v", TestStringProvider().Events)
	}
}

func TestInterceptInterfaceEncodeDecode(t *testing.T) {
	cases := []TestUsesIntfStructProvided{
		{Foo: &TestIntfA{Foo: "hello"}},
		{Foo: &TestIntfB{Bar: "world"}},
		{Foo: nil},
	}

	for _, in := range cases {
		resetIntfStructProvider()

		var buf bytes.Buffer
		wrt := msgp.NewWriter(&buf)
		if err := in.EncodeMsg(wrt); err != nil {
			t.Errorf("%v", err)
		}
		wrt.Flush()

		var out TestUsesIntfStructProvided
		rdr := msgp.NewReader(&buf)
		if err := (&out).DecodeMsg(rdr); err != nil {
			t.Errorf("%v", err)
		}

		if !reflect.DeepEqual(in, out) {
			t.Fatalf("provided encode decode failed")
		}

		if !reflect.DeepEqual([]string{"encode", "decode"}, TestIntfStructProvider().Events) {
			t.Fatalf("unexpected events: %v", TestIntfStructProvider().Events)
		}
	}
}

func TestInterceptInterfaceMarshalUnmarshal(t *testing.T) {
	cases := []TestUsesIntfStructProvided{
		{Foo: &TestIntfA{Foo: "hello"}},
		{Foo: &TestIntfB{Bar: "world"}},
		{Foo: nil},
	}

	for _, in := range cases {
		resetIntfStructProvider()

		bts, err := in.MarshalMsg(nil)
		if err != nil {
			t.Fatalf("%v", err)
		}

		var out TestUsesIntfStructProvided
		if _, err := (&out).UnmarshalMsg(bts); err != nil {
			t.Fatalf("%v", err)
		}

		if !reflect.DeepEqual(in, out) {
			t.Fatalf("provided marshal/unmarshal failed")
		}

		if !reflect.DeepEqual([]string{"marshal", "unmarshal"}, TestIntfStructProvider().Events) {
			t.Fatalf("unexpected events: %v", TestIntfStructProvider().Events)
		}
	}
}

func TestInterceptInterfaceSliceMarshalUnmarshal(t *testing.T) {
	cases := []TestUsesIntfStructProvidedSlice{
		{Foo: []TestIntfStructProvided{
			&TestIntfA{Foo: "hello"},
			&TestIntfB{Bar: "world"},
		}},

		// FIXME: empty slice unmarshals as nil, is this msgp?
		// {Foo: []TestIntfStructProvided{}},
	}

	for _, in := range cases {
		resetIntfStructProvider()

		bts, err := in.MarshalMsg(nil)
		if err != nil {
			t.Fatalf("%v", err)
		}

		var out TestUsesIntfStructProvidedSlice
		if _, err := (&out).UnmarshalMsg(bts); err != nil {
			t.Fatalf("%v", err)
		}

		if !reflect.DeepEqual(in, out) {
			t.Fatalf("provided marshal/unmarshal failed")
		}

		if !reflect.DeepEqual([]string{"marshal", "marshal", "unmarshal", "unmarshal"}, TestIntfStructProvider().Events) {
			t.Fatalf("unexpected events: %v", TestIntfStructProvider().Events)
		}
	}
}

func TestInterceptInterfaceMapMarshalUnmarshal(t *testing.T) {
	cases := []TestUsesIntfStructProvidedMap{
		{Foo: map[string]TestIntfStructProvided{
			"a": &TestIntfA{Foo: "hello"},
			"b": &TestIntfB{Bar: "world"},
		}},

		// FIXME: empty slice unmarshals as nil, is this msgp?
		// {Foo: []TestIntfStructProvided{}},
	}

	for _, in := range cases {
		resetIntfStructProvider()

		bts, err := in.MarshalMsg(nil)
		if err != nil {
			t.Fatalf("%v", err)
		}

		var out TestUsesIntfStructProvidedMap
		if _, err := (&out).UnmarshalMsg(bts); err != nil {
			t.Fatalf("%v", err)
		}

		if !reflect.DeepEqual(in, out) {
			t.Fatalf("provided marshal/unmarshal failed")
		}

		if !reflect.DeepEqual([]string{"marshal", "marshal", "unmarshal", "unmarshal"}, TestIntfStructProvider().Events) {
			t.Fatalf("unexpected events: %v", TestIntfStructProvider().Events)
		}
	}
}

func TestInterceptInterfaceUnmarshalAsJSON(t *testing.T) {
	cases := []struct {
		in  TestUsesIntfStructProvided
		out string
	}{
		{in: TestUsesIntfStructProvided{Foo: &TestIntfA{Foo: "hello"}}, out: `{"Foo":["a",{"Foo":"hello"}]}`},
		{in: TestUsesIntfStructProvided{Foo: &TestIntfB{Bar: "world"}}, out: `{"Foo":["b",{"Bar":"world"}]}`},
		{in: TestUsesIntfStructProvided{Foo: nil}, out: `{"Foo":null}`},
	}

	for idx, tcase := range cases {
		resetIntfStructProvider()

		bts, err := tcase.in.MarshalMsg(nil)
		if err != nil {
			t.Fatalf("%v", err)
		}

		var buf bytes.Buffer
		if _, err := msgp.UnmarshalAsJSON(&buf, bts); err != nil {
			t.Fatalf("%v", err)
		}

		if tcase.out != buf.String() {
			t.Fatalf("%d: unexpected JSON `%s`", idx, buf.String())
		}
	}
}
