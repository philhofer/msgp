package gen

import (
	"fmt"
	"io"
)

func omitempies(w io.Writer) *omitemptyGen {
	return &omitemptyGen{
		p: printer{w: w},
	}
}

type omitemptyGen struct {
	passes
	p    printer
	expr string
}

func (s *omitemptyGen) Method() Method { return Size }

func (s *omitemptyGen) Apply(dirs []string) error {
	return nil
}

func (s *omitemptyGen) Execute(p Elem) error {
	if !s.p.ok() {
		return s.p.err
	}
	p = s.applyall(p)
	if p == nil {
		return nil
	}
	if !IsPrintable(p) {
		return nil
	}

	s.p.comment("ShouldOmitAsEmpty returns true if all struct fields have default value")
	s.p.printf("\nfunc (%s %s) ShouldOmitAsEmpty() (bool) {", p.Varname(), imutMethodReceiver(p))
	s.p.printf("\nreturn false")
	s.p.closeblock()
	return s.p.err
}

func (s *omitemptyGen) gStruct(st *Struct) {
	// Always serialize if at least one fields always serialized.
	empty := omitemptyGen{}
	for i := range st.Fields {
		empty.expr = ""
		next(&empty, st.Fields[i].FieldElem)
		if empty.expr == "" {
			s.expr = ""
			return
		}
	}

	// Check it object empty code.
	vname := ""
	for i := range st.Fields {
		s.expr = ""
		next(s, st.Fields[i].FieldElem)
		if s.expr != "true" {
			if vname == "" {
				vname = randIdent()
				if s.p.w != nil {
					s.p.printf("\n%s := %s", vname, s.expr)
				}
			} else if s.p.w != nil {
				s.p.printf("\n%s = %s && %s", vname, vname, s.expr)
			}
		}
	}
	if vname == "" {
		vname = "true"
	}
	s.expr = vname
}

func (s *omitemptyGen) gPtr(p *Ptr) {
	s.expr = fmt.Sprintf("(%s == nil)", p.Varname())
}

func (s *omitemptyGen) gSlice(sl *Slice) {
	s.expr = fmt.Sprintf("(len(%s) == 0)", sl.Varname())
}

func (s *omitemptyGen) gArray(a *Array) {
	s.expr = ""
}

func (s *omitemptyGen) gMap(m *Map) {
	s.expr = fmt.Sprintf("((%s == nil) || (len(%s) == 0))", m.Varname(), m.Varname())
}

func (s *omitemptyGen) gBase(b *BaseElem) {
	vname := b.Varname()
	var verr string
	if b.Convert {
		if b.ShimMode == Cast {
			vname = tobaseConvert(b)
		} else {
			vname = randIdent()
			verr = randIdent()
			if s.p.w != nil {
				s.p.printf("\nvar %s %s", vname, b.BaseType())
				s.p.printf("\nvar %s error", verr)
				s.p.printf("\n%s, %s = %s", vname, verr, tobaseConvert(b))
			}
		}
	}

	if b.Value == IDENT { // unknown identity
		s.expr = fmt.Sprintf("%s.ShouldOmitAsEmpty()", vname)
	} else { // typical case
		s.expr = fmt.Sprintf("msgp.ShouldOmitAsEmpty%s(%s)", b.BaseName(), vname)
	}
	if verr != "" {
		s.expr = fmt.Sprintf("((%s == nil) && %s)", verr, s.expr)
	}
}
