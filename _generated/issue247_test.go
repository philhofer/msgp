package _generated

import (
	"bytes"
	"github.com/tinylib/msgp/msgp"
	"testing"
)

func TestIssue247(t *testing.T) {
	buf := bytes.Buffer{}
	test := T247{}

	msgp.Encode(&buf, &T247{S: []string{}})
	msgp.Decode(&buf, &test)

	if test.S == nil {
		t.Errorf("expected slice to not be nil")
	}
}
