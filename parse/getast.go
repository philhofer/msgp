package parse

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"reflect"
	"strings"

	"github.com/tinylib/msgp/gen"
	"github.com/ttacon/chalk"
)

// A FileSet is the in-memory representation of a
// parsed file.
type FileSet struct {
	Package              string              // package name
	Specs                map[string]ast.Expr // type specs in file
	Identities           map[string]gen.Elem // processed from specs
	Directives           []string            // raw preprocessor directives
	EncodeValueReceivers bool                // whether generated encode function reveivers are values instead of pointers
}

// File parses a file at the relative path
// provided and produces a new *FileSet.
// (No exported structs is considered an error.)
// If you pass in a path to a directory, the entire
// directory will be parsed.
func File(name string) (*FileSet, error) {
	fs := &FileSet{
		Specs:      make(map[string]ast.Expr),
		Identities: make(map[string]gen.Elem),
	}

	fset := token.NewFileSet()
	finfo, err := os.Stat(name)
	if err != nil {
		return nil, err
	}
	if finfo.IsDir() {
		pkgs, err := parser.ParseDir(fset, name, nil, parser.ParseComments)
		if err != nil {
			return nil, err
		}
		if len(pkgs) != 1 {
			return nil, fmt.Errorf("multiple packages in directory: %s", name)
		}
		var one *ast.Package
		for _, nm := range pkgs {
			one = nm
			break
		}
		fs.Package = one.Name
		for _, fl := range one.Files {
			fs.Directives = append(fs.Directives, yieldComments(fl.Comments)...)
			ast.FileExports(fl)
			fs.getTypeSpecs(fl)
		}
	} else {
		f, err := parser.ParseFile(fset, name, nil, parser.ParseComments)
		if err != nil {
			return nil, err
		}
		fs.Package = f.Name.Name
		fs.Directives = yieldComments(f.Comments)
		ast.FileExports(f)
		fs.getTypeSpecs(f)
	}

	if len(fs.Specs) == 0 {
		return nil, fmt.Errorf("no exported definitions in %s", name)
	}

	fs.process()
	fs.applyDirectives()
	fs.propInline()

	return fs, nil
}

// ApplyDirectives applies all of the preprocessor
// directives to the file set in the order that they
// appear in the source file. It is assumed that
// f.Identities is already populated.
func (f *FileSet) applyDirectives() {
	for _, d := range f.Directives {
		chunks := strings.Split(d, " ")
		if len(chunks) > 0 {
			if fn, ok := directives[chunks[0]]; ok {
				err := fn(chunks, f)
				if err != nil {
					warnf("error applying directive: %s\n", err)
				}
			}
		}
	}
}

// process takes the contents of f.Specs and
// uses them to populate f.Identities
func (f *FileSet) process() {
	// generate elements
	for name, def := range f.Specs {
		infof("parsing %s...\n", name)
		el := f.parseExpr(def)
		if el != nil {
			el.Alias(name)
			f.Identities[name] = el
		} else {
			warnf(" \u26a0 unable to parse %s\n", name)
		}
	}
}

// GetElems creates a FileSet from a filename
// or directory and returns copies of the processed
// elements with initialized variable names.
func Elems(filename string) ([]gen.Elem, string, error) {
	fs, err := File(filename)
	if err != nil {
		return nil, "", err
	}

	g := make([]gen.Elem, 0, len(fs.Identities))

	for _, el := range fs.Identities {
		el.SetVarname("z")
		el.SetEncodeValueReceivers(fs.EncodeValueReceivers)
		g = append(g, el)
	}

	return g, fs.Package, nil
}

// getTypeSpecs extracts all of the *ast.TypeSpecs in the file
// into fs.Identities, but does not set the actual element
func (fs *FileSet) getTypeSpecs(f *ast.File) {

	// check all declarations...
	for i := range f.Decls {

		// for GenDecls...
		if g, ok := f.Decls[i].(*ast.GenDecl); ok {

			// and check the specs...
			for _, s := range g.Specs {

				// for ast.TypeSpecs....
				if ts, ok := s.(*ast.TypeSpec); ok {
					switch ts.Type.(type) {

					// this is the list of parse-able
					// type specs
					case *ast.StructType,
						*ast.ArrayType,
						*ast.StarExpr,
						*ast.MapType,
						*ast.Ident:
						fs.Specs[ts.Name.Name] = ts.Type

					}
				}
			}
		}
	}
}

func (fs *FileSet) parseFieldList(fl *ast.FieldList) []gen.StructField {
	if fl == nil || fl.NumFields() == 0 {
		return nil
	}
	out := make([]gen.StructField, 0, fl.NumFields())
	for i, field := range fl.List {
		fds := fs.getField(field)
		if len(fds) > 0 {
			out = append(out, fds...)
		} else {
			warnf(" \u26a0 ignored struct field %d\n", i)
		}
	}
	return out
}

// translate *ast.Field into []gen.StructField
func (fs *FileSet) getField(f *ast.Field) []gen.StructField {
	sf := make([]gen.StructField, 1)
	var extension bool
	// parse tag; otherwise field name is field tag
	if f.Tag != nil {
		body := reflect.StructTag(strings.Trim(f.Tag.Value, "`")).Get("msg")
		tags := strings.Split(body, ",")
		if len(tags) == 2 && tags[1] == "extension" {
			extension = true
		}
		// ignore "-" fields
		if tags[0] == "-" {
			return nil
		}
		sf[0].FieldTag = tags[0]
	}

	ex := fs.parseExpr(f.Type)
	if ex == nil {
		return nil
	}

	// parse field name
	switch len(f.Names) {
	case 0:
		sf[0].FieldName = embedded(f.Type)
	case 1:
		sf[0].FieldName = f.Names[0].Name
	default:
		// this is for a multiple in-line declaration,
		// e.g. type A struct { One, Two int }
		sf = sf[0:0]
		for _, nm := range f.Names {
			sf = append(sf, gen.StructField{
				FieldTag:  nm.Name,
				FieldName: nm.Name,
				FieldElem: ex.Copy(),
			})
		}
		return sf
	}
	sf[0].FieldElem = ex
	if sf[0].FieldTag == "" {
		sf[0].FieldTag = sf[0].FieldName
	}

	// validate extension
	if extension {
		switch ex := ex.(type) {
		case *gen.Ptr:
			if b, ok := ex.Value.(*gen.BaseElem); ok {
				b.Value = gen.Ext
			} else {
				warnf(" \u26a0 field %q couldn't be cast as an extension\n", sf[0].FieldName)
				return nil
			}
		case *gen.BaseElem:
			ex.Value = gen.Ext
		default:
			warnf(" \u26a0 field %q couldn't be cast as an extension\n", sf[0].FieldName)
			return nil
		}
	}
	return sf
}

// extract embedded field name
//
// so, for a struct like
//
//	type A struct {
//		io.Writer
//  }
//
// we want "Writer"
func embedded(f ast.Expr) string {
	switch f := f.(type) {
	case *ast.Ident:
		return f.Name
	case *ast.StarExpr:
		return embedded(f.X)
	case *ast.SelectorExpr:
		return f.Sel.Name
	default:
		// other possibilities are disallowed
		return ""
	}
}

// stringify a field type name
func stringify(e ast.Expr) string {
	switch e := e.(type) {
	case *ast.Ident:
		return e.Name
	case *ast.StarExpr:
		return "*" + stringify(e.X)
	case *ast.SelectorExpr:
		return stringify(e.X) + "." + e.Sel.Name
	case *ast.ArrayType:
		if e.Len == nil {
			return "[]" + stringify(e.Elt)
		}
		return fmt.Sprintf("[%s]%s", stringify(e.Len), stringify(e.Elt))
	case *ast.InterfaceType:
		if e.Methods == nil || e.Methods.NumFields() == 0 {
			return "interface{}"
		}
	}
	return ""
}

// recursively translate ast.Expr to gen.Elem; nil means type not supported
// expected input types:
// - *ast.MapType (map[T]J)
// - *ast.Ident (name)
// - *ast.ArrayType ([(sz)]T)
// - *ast.StarExpr (*T)
// - *ast.StructType (struct {})
// - *ast.SelectorExpr (a.B)
// - *ast.InterfaceType (interface {})
func (fs *FileSet) parseExpr(e ast.Expr) gen.Elem {
	switch e := e.(type) {

	case *ast.MapType:
		if k, ok := e.Key.(*ast.Ident); ok && k.Name == "string" {
			if in := fs.parseExpr(e.Value); in != nil {
				return &gen.Map{Value: in}
			}
		}
		return nil

	case *ast.Ident:
		b := gen.Ident(e.Name)

		// if the name isn't one of the type
		// specs, warn
		if b.Value == gen.IDENT {
			if _, ok := fs.Specs[e.Name]; !ok {
				warnf(" \u26a0 non-local identifier: %s\n", e.Name)
			}
		}
		return b

	case *ast.ArrayType:

		// special case for []byte
		if e.Len == nil {
			if i, ok := e.Elt.(*ast.Ident); ok && i.Name == "byte" {
				return &gen.BaseElem{Value: gen.Bytes}
			}
		}

		// return early if we don't know
		// what the slice element type is
		els := fs.parseExpr(e.Elt)
		if els == nil {
			return nil
		}

		// array and not a slice
		if e.Len != nil {
			switch s := e.Len.(type) {
			case *ast.BasicLit:
				return &gen.Array{
					Size: s.Value,
					Els:  els,
				}

			case *ast.Ident:
				return &gen.Array{
					Size: s.String(),
					Els:  els,
				}

			case *ast.SelectorExpr:
				return &gen.Array{
					Size: stringify(s),
					Els:  els,
				}

			default:
				return nil
			}
		}
		return &gen.Slice{Els: els}

	case *ast.StarExpr:
		if v := fs.parseExpr(e.X); v != nil {
			return &gen.Ptr{Value: v}
		}
		return nil

	case *ast.StructType:
		if fields := fs.parseFieldList(e.Fields); len(fields) > 0 {
			return &gen.Struct{Fields: fields}
		}
		return nil

	case *ast.SelectorExpr:
		return gen.Ident(stringify(e))

	case *ast.InterfaceType:
		// support `interface{}`
		if len(e.Methods.List) == 0 {
			return &gen.BaseElem{Value: gen.Intf}
		}
		return nil

	default: // other types not supported
		return nil
	}
}

func infof(s string, v ...interface{})  { fmt.Printf(chalk.Green.Color(s), v...) }
func warnf(s string, v ...interface{})  { fmt.Printf(chalk.Yellow.Color(s), v...) }
func warnln(s string)                   { fmt.Println(chalk.Yellow.Color(s)) }
func fatalf(s string, v ...interface{}) { fmt.Printf(chalk.Red.Color(s), v...) }
