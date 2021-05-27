package msgp

import (
	"bytes"
	"math"
	"testing"
)

func TestNumber(t *testing.T) {

	n := Number{}

	if n.Type() != IntType {
		t.Errorf("expected zero-value type to be %s; got %s", IntType, n.Type())
	}

	if n.String() != "0" {
		t.Errorf("expected Number{}.String() to be \"0\" but got %q", n.String())
	}

	n.AsInt(248)
	i, ok := n.Int()
	if !ok || i != 248 || n.Type() != IntType || n.String() != "248" {
		t.Errorf("%d in; %d out!", 248, i)
	}

	n.AsFloat64(3.141)
	f, ok := n.Float()
	if !ok || f != 3.141 || n.Type() != Float64Type || n.String() != "3.141" {
		t.Errorf("%f in; %f out!", 3.141, f)
	}

	n.AsUint(40000)
	u, ok := n.Uint()
	if !ok || u != 40000 || n.Type() != UintType || n.String() != "40000" {
		t.Errorf("%d in; %d out!", 40000, u)
	}

	nums := []interface{}{
		float64(3.14159),
		int64(-29081),
		uint64(90821983),
		float32(3.141),
	}

	var dat []byte
	var buf bytes.Buffer
	wr := NewWriter(&buf)
	for _, n := range nums {
		dat, _ = AppendIntf(dat, n)
		wr.WriteIntf(n)
	}
	wr.Flush()

	mout := make([]Number, len(nums))
	dout := make([]Number, len(nums))

	rd := NewReader(&buf)
	unm := dat
	for i := range nums {
		var err error
		unm, err = mout[i].UnmarshalMsg(unm)
		if err != nil {
			t.Fatal("unmarshal error:", err)
		}
		err = dout[i].DecodeMsg(rd)
		if err != nil {
			t.Fatal("decode error:", err)
		}
		if mout[i] != dout[i] {
			t.Errorf("for %#v, got %#v from unmarshal and %#v from decode", nums[i], mout[i], dout[i])
		}
	}

	buf.Reset()
	var odat []byte
	for i := range nums {
		var err error
		odat, err = mout[i].MarshalMsg(odat)
		if err != nil {
			t.Fatal("marshal error:", err)
		}
		err = dout[i].EncodeMsg(wr)
		if err != nil {
			t.Fatal("encode error", err)
		}
	}
	wr.Flush()

	if !bytes.Equal(dat, odat) {
		t.Errorf("marshal: expected output %#v; got %#v", dat, odat)
	}

	if !bytes.Equal(dat, buf.Bytes()) {
		t.Errorf("encode: expected output %#v; got %#v", dat, buf.Bytes())
	}

}

func TestNumber_CastInt64(t *testing.T) {
	var n Number
	n.AsUint(math.MaxUint64)

	_, ok := n.CastInt64()
	if ok {
		t.Error("CastInt64() failed: MaxUint64 > MaxInt64")
	}
}

func TestNumber_CastInt32(t *testing.T) {
	cases := []struct {
		in   int64
		want bool
	}{
		{math.MinInt32, true},
		{math.MinInt32 - 1, false},
		{math.MaxInt32, true},
		{math.MaxInt32 + 1, false},
		{0, true},
	}

	var n Number
	for _, c := range cases {
		n.AsInt(c.in)

		_, ok := n.CastInt32()
		if ok != c.want {
			t.Errorf("cast %v to int32 invalid: %v, got %v", c.in, c.want, ok)
		}
	}
}

func TestNumber_CastUint64(t *testing.T) {
	var n Number
	n.AsInt(math.MinInt64)

	_, ok := n.CastUint64()
	if ok {
		t.Error("CastUint64() failed: MinInt64 < 0")
	}
}

func TestNumber_CastUint32(t *testing.T) {
	cases := []struct {
		in   int64
		want bool
	}{
		{math.MinInt64, false},
		{math.MaxInt32, true},
		{math.MaxInt32 + 1, false},
		{0, true},
	}

	var n Number
	for _, c := range cases {
		n.AsInt(c.in)

		_, ok := n.CastUint32()
		if ok != c.want {
			t.Errorf("cast %v to uint32 invalid: %v, got %v", c.in, c.want, ok)
		}
	}
}
