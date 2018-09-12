package _generated

import (
	"fmt"

	"github.com/tinylib/msgp/msgp"
)

//go:generate msgp

//msgp:ignore testStructProvider
type testStructProvider struct {
	Events []string
}

var prv = &testStructProvider{}

func resetStructProvider() {
	prv = &testStructProvider{}
}

func TestStructProvider() *testStructProvider {
	return prv
}

//msgp:intercept TestStructProvided using:TestStructProvider

type TestStructProvided struct {
	Foo string
}

type TestUsesStructProvided struct {
	Foo *TestStructProvided
}

func (p *testStructProvider) DecodeMsg(dc *msgp.Reader) (t *TestStructProvided, err error) {
	p.Events = append(p.Events, "decode")
	t = new(TestStructProvided)
	err = t.DecodeMsg(dc)
	return
}

func (p *testStructProvider) UnmarshalMsg(bts []byte) (t *TestStructProvided, o []byte, err error) {
	t = new(TestStructProvided)
	p.Events = append(p.Events, "unmarshal")
	o, err = t.UnmarshalMsg(bts)
	return
}

func (p *testStructProvider) EncodeMsg(t *TestStructProvided, en *msgp.Writer) (err error) {
	p.Events = append(p.Events, "encode")
	return t.EncodeMsg(en)
}

func (p *testStructProvider) MarshalMsg(t *TestStructProvided, b []byte) (o []byte, err error) {
	p.Events = append(p.Events, "marshal")
	return t.MarshalMsg(b)
}

func (p *testStructProvider) Msgsize(t *TestStructProvided) (s int) {
	p.Events = append(p.Events, "msgsize")
	return t.Msgsize()
}

//msgp:ignore testStringProvided
type testStringProvider struct {
	Events []string
}

var stringPrv = &testStringProvider{}

func resetStringProvider() {
	stringPrv = &testStringProvider{}
}

func TestStringProvider() *testStringProvider {
	return stringPrv
}

//msgp:intercept TestStringProvided using:TestStringProvider
type TestStringProvided string

type TestUsesStringProvided struct {
	Foo TestStringProvided
}

func (p *testStringProvider) DecodeMsg(dc *msgp.Reader) (t TestStringProvided, err error) {
	p.Events = append(p.Events, "decode")
	var s string
	s, err = dc.ReadString()
	if err != nil {
		return
	}
	t = TestStringProvided(s)
	return
}

func (p *testStringProvider) UnmarshalMsg(bts []byte) (t TestStringProvided, o []byte, err error) {
	p.Events = append(p.Events, "unmarshal")
	var s string
	s, o, err = msgp.ReadStringBytes(bts)
	if err != nil {
		return
	}
	t = TestStringProvided(s)
	return
}

func (p *testStringProvider) EncodeMsg(t TestStringProvided, en *msgp.Writer) (err error) {
	p.Events = append(p.Events, "encode")
	return en.WriteString(string(t))
}

func (p *testStringProvider) MarshalMsg(t TestStringProvided, b []byte) (o []byte, err error) {
	p.Events = append(p.Events, "marshal")
	o = msgp.AppendString(b, string(t))
	return
}

func (p *testStringProvider) Msgsize(t TestStringProvided) (s int) {
	return msgp.StringPrefixSize + len(t)
}

//msgp:ignore testIntfStructProvider
type testIntfStructProvider struct {
	Events []string
}

var intfStructPrv = &testIntfStructProvider{}

func resetIntfStructProvider() {
	intfStructPrv = &testIntfStructProvider{}
}

func TestIntfStructProvider() *testIntfStructProvider {
	return intfStructPrv
}

//msgp:intercept TestIntfStructProvided using:TestIntfStructProvider
type TestIntfStructProvided interface {
	msgp.Decodable
	msgp.Encodable
	msgp.MarshalSizer
	msgp.Unmarshaler
}

type TestUsesIntfStructProvided struct {
	Foo TestIntfStructProvided
}

type TestUsesIntfStructProvidedSlice struct {
	Foo []TestIntfStructProvided
}

type TestUsesIntfStructProvidedMap struct {
	Foo map[string]TestIntfStructProvided
}

type TestIntfA struct {
	Foo string
}

type TestIntfB struct {
	Bar string
}

func (p *testIntfStructProvider) DecodeMsg(dc *msgp.Reader) (t TestIntfStructProvided, err error) {
	p.Events = append(p.Events, "decode")

	if dc.IsNil() {
		err = dc.ReadNil()
	} else {
		var s string
		var sz uint32
		if sz, err = dc.ReadArrayHeader(); err != nil {
			return
		}
		if sz != 2 {
			err = fmt.Errorf("unexpected array length")
			return
		}
		s, err = dc.ReadString()
		if err != nil {
			return
		}
		switch s {
		case "a":
			t = new(TestIntfA)
		case "b":
			t = new(TestIntfB)
		default:
			err = fmt.Errorf("unexpected type")
			return
		}
		err = t.DecodeMsg(dc)
	}
	return
}

func (p *testIntfStructProvider) UnmarshalMsg(bts []byte) (t TestIntfStructProvided, o []byte, err error) {
	p.Events = append(p.Events, "unmarshal")

	o = bts
	if msgp.IsNil(bts) {
		o, err = msgp.ReadNilBytes(o)
	} else {
		var s string
		var sz uint32
		if sz, o, err = msgp.ReadArrayHeaderBytes(o); err != nil {
			return
		}
		if sz != 2 {
			err = fmt.Errorf("unexpected array length")
			return
		}
		s, o, err = msgp.ReadStringBytes(o)
		if err != nil {
			return
		}
		switch s {
		case "a":
			t = new(TestIntfA)
		case "b":
			t = new(TestIntfB)
		default:
			err = fmt.Errorf("unexpected type")
			return
		}
		o, err = t.UnmarshalMsg(o)
	}
	return
}

func (p *testIntfStructProvider) EncodeMsg(t TestIntfStructProvided, en *msgp.Writer) (err error) {
	p.Events = append(p.Events, "encode")
	if t == nil {
		return en.WriteNil()
	} else {
		if err = en.WriteArrayHeader(2); err != nil {
			return
		}
		var s string
		switch t.(type) {
		case *TestIntfA:
			s = "a"
		case *TestIntfB:
			s = "b"
		default:
			err = fmt.Errorf("unexpected type %T", t)
		}
		if err = en.WriteString(s); err != nil {
			return
		}
		return t.EncodeMsg(en)
	}
}

func (p *testIntfStructProvider) MarshalMsg(t TestIntfStructProvided, b []byte) (o []byte, err error) {
	p.Events = append(p.Events, "marshal")
	o = b
	if t == nil {
		o = msgp.AppendNil(o)
		return
	} else {
		o = msgp.AppendArrayHeader(o, 2)
		var s string
		switch t.(type) {
		case *TestIntfA:
			s = "a"
		case *TestIntfB:
			s = "b"
		default:
			err = fmt.Errorf("unexpected type %T", t)
		}
		o = msgp.AppendString(o, s)
		return t.MarshalMsg(o)
	}
}

func (p *testIntfStructProvider) Msgsize(t TestIntfStructProvided) (s int) {
	if t == nil {
		return msgp.NilSize
	}
	return t.Msgsize()
}
