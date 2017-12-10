package parse

import (
	"github.com/tinylib/msgp/gen"
)

// This file defines when and how we
// propagate type information from
// one type declaration to another.
// After the processing pass, every
// non-primitive type is marshalled/unmarshalled/etc.
// through a function call. Here, we propagate
// the type information into the caller's type
// tree *if* the child type is simple enough.
//
// For example, types like
//
//    type A [4]int
//
// will get pushed into parent methods,
// whereas types like
//
//    type B [3]map[string]struct{A, B [4]string}
//
// will not.

// this is an approximate measure
// of the number of children in a node
const maxComplex = 5

// begin recursive search for identities with the
// given name and replace them with be
func (fs *FileSet) findShim(id string, be *gen.BaseElem) {
	for name, el := range fs.Identities {
		pushstate(name)
		switch el := el.(type) {
		case *gen.Struct:
			for i := range el.Fields {
				fs.nextShim(&el.Fields[i].FieldElem, id, be)
			}
		case *gen.Array:
			fs.nextShim(&el.Els, id, be)
		case *gen.Slice:
			fs.nextShim(&el.Els, id, be)
		case *gen.Map:
			fs.nextShim(&el.Value, id, be)
		case *gen.Ptr:
			fs.nextShim(&el.Value, id, be)
		}
		popstate()
	}
	// we'll need this at the top level as well
	fs.Identities[id] = be
}

func (fs *FileSet) nextShim(ref *gen.Elem, id string, be *gen.BaseElem) {
	if (*ref).TypeName() == id {
		vn := (*ref).Varname()
		*ref = be.Copy()
		(*ref).SetVarname(vn)
	} else {
		switch el := (*ref).(type) {
		case *gen.Struct:
			for i := range el.Fields {
				fs.nextShim(&el.Fields[i].FieldElem, id, be)
			}
		case *gen.Array:
			fs.nextShim(&el.Els, id, be)
		case *gen.Slice:
			fs.nextShim(&el.Els, id, be)
		case *gen.Map:
			fs.nextShim(&el.Value, id, be)
		case *gen.Ptr:
			fs.nextShim(&el.Value, id, be)
		}
	}
}

// propInline identifies and inlines candidates
func (fs *FileSet) propInline() {
	for name, el := range fs.Identities {
		pushstate(name)
		switch el := el.(type) {
		case *gen.Struct:
			for i := range el.Fields {
				fs.nextInline(&el.Fields[i].FieldElem, name)
			}
		case *gen.Array:
			fs.nextInline(&el.Els, name)
		case *gen.Slice:
			fs.nextInline(&el.Els, name)
		case *gen.Map:
			fs.nextInline(&el.Value, name)
		case *gen.Ptr:
			fs.nextInline(&el.Value, name)
		}
		popstate()
	}
}

const fatalloop = `detected infinite recursion in inlining loop!
Please file a bug at github.com/tinylib/msgp/issues!
Thanks!
`

func (fs *FileSet) nextInline(ref *gen.Elem, root string) {
	switch el := (*ref).(type) {
	case *gen.BaseElem:
		// ensure that we're not inlining
		// a type into itself
		typ := el.TypeName()
		if el.Value == gen.IDENT && typ != root {
			if node, ok := fs.Identities[typ]; ok && node.Complexity() < maxComplex {
				infof("inlining %s\n", typ)

				// This should never happen; it will cause
				// infinite recursion.
				if node == *ref {
					panic(fatalloop)
				}

				*ref = node.Copy()
				fs.nextInline(ref, node.TypeName())
			} else if !ok && !el.Resolved() {
				// this is the point at which we're sure that
				// we've got a type that isn't a primitive,
				// a library builtin, or a processed type
				warnf("unresolved identifier: %s\n", typ)
			}
		}
	case *gen.Struct:
		for i := range el.Fields {
			fs.nextInline(&el.Fields[i].FieldElem, root)
		}
	case *gen.Array:
		fs.nextInline(&el.Els, root)
	case *gen.Slice:
		fs.nextInline(&el.Els, root)
	case *gen.Map:
		fs.nextInline(&el.Value, root)
	case *gen.Ptr:
		fs.nextInline(&el.Value, root)
	default:
		panic("bad elem type")
	}
}
