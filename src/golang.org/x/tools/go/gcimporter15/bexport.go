// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build go1.6

// Binary package export.
// This file was derived from $GOROOT/src/cmd/compile/internal/gc/bexport.go;
// see that file for specification of the format.

package gcimporter

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"go/ast"
	"go/constant"
	"go/token"
	"go/types"
	"log"
	"math"
	"math/big"
	"sort"
	"strings"
)

// If debugFormat is set, each integer and string value is preceded by a marker
// and position information in the encoding. This mechanism permits an importer
// to recognize immediately when it is out of sync. The importer recognizes this
// mode automatically (i.e., it can import export data produced with debugging
// support even if debugFormat is not set at the time of import). This mode will
// lead to massively larger export data (by a factor of 2 to 3) and should only
// be enabled during development and debugging.
//
// NOTE: This flag is the first flag to enable if importing dies because of
// (suspected) format errors, and whenever a change is made to the format.
const debugFormat = false // default: false

// If posInfoFormat is set, position information (file, lineno) is written
// for each exported object, including methods and struct fields. Currently
// disabled because it may lead to different object files depending on which
// directory they are built under, which causes tests checking for hermetic
// builds to fail (e.g. TestCgoConsistentResults for cmd/go).
// TODO(gri) determine what to do here.
const posInfoFormat = false

// If trace is set, debugging output is printed to std out.
const trace = false // default: false

const exportVersion = "v0"

type exporter struct {
	fset *token.FileSet
	out  bytes.Buffer

	// object -> index maps, indexed in order of serialization
	strIndex map[string]int
	pkgIndex map[*types.Package]int
	typIndex map[types.Type]int

	// position encoding
	prevFile string
	prevLine int

	// debugging support
	written int // bytes written
	indent  int // for trace
}

// BExportData returns binary export data for pkg.
// If no file set is provided, position info will be missing.
func BExportData(fset *token.FileSet, pkg *types.Package) []byte {
	p := exporter{
		fset:     fset,
		strIndex: map[string]int{"": 0}, // empty string is mapped to 0
		pkgIndex: make(map[*types.Package]int),
		typIndex: make(map[types.Type]int),
	}

	// first byte indicates low-level encoding format
	var format byte = 'c' // compact
	if debugFormat {
		format = 'd'
	}
	p.rawByte(format)

	// posInfo exported or not?
	p.bool(posInfoFormat)

	// --- generic export data ---

	if trace {
		p.tracef("\n--- generic export data ---\n")
		if p.indent != 0 {
			log.Fatalf("gcimporter: incorrect indentation %d", p.indent)
		}
	}

	if trace {
		p.tracef("version = ")
	}
	p.string(exportVersion)
	if trace {
		p.tracef("\n")
	}

	// populate type map with predeclared "known" types
	for index, typ := range predeclared {
		p.typIndex[typ] = index
	}
	if len(p.typIndex) != len(predeclared) {
		log.Fatalf("gcimporter: duplicate entries in type map?")
	}

	// write package data
	p.pkg(pkg, true)
	if trace {
		p.tracef("\n")
	}

	// write objects
	objcount := 0
	scope := pkg.Scope()
	for _, name := range scope.Names() {
		if !ast.IsExported(name) {
			continue
		}
		if trace {
			p.tracef("\n")
		}
		p.obj(scope.Lookup(name))
		objcount++
	}

	// indicate end of list
	if trace {
		p.tracef("\n")
	}
	p.tag(endTag)

	// for self-verification only (redundant)
	p.int(objcount)

	if trace {
		p.tracef("\n")
	}

	// --- end of export data ---

	return p.out.Bytes()
}

func (p *exporter) pkg(pkg *types.Package, emptypath bool) {
	if pkg == nil {
		log.Fatalf("gcimporter: unexpected nil pkg")
	}

	// if we saw the package before, write its index (>= 0)
	if i, ok := p.pkgIndex[pkg]; ok {
		p.index('P', i)
		return
	}

	// otherwise, remember the package, write the package tag (< 0) and package data
	if trace {
		p.tracef("P%d = { ", len(p.pkgIndex))
		defer p.tracef("} ")
	}
	p.pkgIndex[pkg] = len(p.pkgIndex)

	p.tag(packageTag)
	p.string(pkg.Name())
	if emptypath {
		p.string("")
	} else {
		p.string(pkg.Path())
	}
}

func (p *exporter) obj(obj types.Object) {
	switch obj := obj.(type) {
	case *types.Const:
		p.tag(constTag)
		p.pos(obj)
		p.qualifiedName(obj)
		p.typ(obj.Type())
		p.value(obj.Val())

	case *types.TypeName:
		p.tag(typeTag)
		p.typ(obj.Type())

	case *types.Var:
		p.tag(varTag)
		p.pos(obj)
		p.qualifiedName(obj)
		p.typ(obj.Type())

	case *types.Func:
		p.tag(funcTag)
		p.pos(obj)
		p.qualifiedName(obj)
		sig := obj.Type().(*types.Signature)
		p.paramList(sig.Params(), sig.Variadic())
		p.paramList(sig.Results(), false)

	default:
		log.Fatalf("gcimporter: unexpected object %v (%T)", obj, obj)
	}
}

func (p *exporter) pos(obj types.Object) {
	if !posInfoFormat {
		return
	}

	var file string
	var line int
	if p.fset != nil {
		pos := p.fset.Position(obj.Pos())
		file = pos.Filename
		line = pos.Line
	}

	if file == p.prevFile && line != p.prevLine {
		// common case: write delta-encoded line number
		p.int(line - p.prevLine) // != 0
	} else {
		// uncommon case: filename changed, or line didn't change
		p.int(0)
		p.string(file)
		p.int(line)
		p.prevFile = file
	}
	p.prevLine = line
}

func (p *exporter) qualifiedName(obj types.Object) {
	p.string(obj.Name())
	p.pkg(obj.Pkg(), false)
}

func (p *exporter) typ(t types.Type) {
	if t == nil {
		log.Fatalf("gcimporter: nil type")
	}

	// Possible optimization: Anonymous pointer types *T where
	// T is a named type are common. We could canonicalize all
	// such types *T to a single type PT = *T. This would lead
	// to at most one *T entry in typIndex, and all future *T's
	// would be encoded as the respective index directly. Would
	// save 1 byte (pointerTag) per *T and reduce the typIndex
	// size (at the cost of a canonicalization map). We can do
	// this later, without encoding format change.

	// if we saw the type before, write its index (>= 0)
	if i, ok := p.typIndex[t]; ok {
		p.index('T', i)
		return
	}

	// otherwise, remember the type, write the type tag (< 0) and type data
	index := len(p.typIndex)
	if trace {
		p.tracef("T%d = {>\n", index)
		defer p.tracef("<\n} ")
	}
	p.typIndex[t] = index

	switch t := t.(type) {
	case *types.Named:
		p.tag(namedTag)
		p.pos(t.Obj())
		p.qualifiedName(t.Obj())
		p.typ(t.Underlying())
		if !types.IsInterface(t) {
			p.assocMethods(t)
		}

	case *types.Array:
		p.tag(arrayTag)
		p.int64(t.Len())
		p.typ(t.Elem())

	case *types.Slice:
		p.tag(sliceTag)
		p.typ(t.Elem())

	case *dddSlice:
		p.tag(dddTag)
		p.typ(t.elem)

	case *types.Struct:
		p.tag(structTag)
		p.fieldList(t)

	case *types.Pointer:
		p.tag(pointerTag)
		p.typ(t.Elem())

	case *types.Signature:
		p.tag(signatureTag)
		p.paramList(t.Params(), t.Variadic())
		p.paramList(t.Results(), false)

	case *types.Interface:
		p.tag(interfaceTag)
		p.iface(t)

	case *types.Map:
		p.tag(mapTag)
		p.typ(t.Key())
		p.typ(t.Elem())

	case *types.Chan:
		p.tag(chanTag)
		p.int(int(3 - t.Dir())) // hack
		p.typ(t.Elem())

	default:
		log.Fatalf("gcimporter: unexpected type %T: %s", t, t)
	}
}

func (p *exporter) assocMethods(named *types.Named) {
	// Sort methods (for determinism).
	var methods []*types.Func
	for i := 0; i < named.NumMethods(); i++ {
		methods = append(methods, named.Method(i))
	}
	sort.Sort(methodsByName(methods))

	p.int(len(methods))

	if trace && methods != nil {
		p.tracef("associated methods {>\n")
	}

	for i, m := range methods {
		if trace && i > 0 {
			p.tracef("\n")
		}

		p.pos(m)
		name := m.Name()
		p.string(name)
		if !exported(name) {
			p.pkg(m.Pkg(), false)
		}

		sig := m.Type().(*types.Signature)
		p.paramList(types.NewTuple(sig.Recv()), false)
		p.paramList(sig.Params(), sig.Variadic())
		p.paramList(sig.Results(), false)
	}

	if trace && methods != nil {
		p.tracef("<\n} ")
	}
}

type methodsByName []*types.Func

func (x methodsByName) Len() int           { return len(x) }
func (x methodsByName) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }
func (x methodsByName) Less(i, j int) bool { return x[i].Name() < x[j].Name() }

func (p *exporter) fieldList(t *types.Struct) {
	if trace && t.NumFields() > 0 {
		p.tracef("fields {>\n")
		defer p.tracef("<\n} ")
	}

	p.int(t.NumFields())
	for i := 0; i < t.NumFields(); i++ {
		if trace && i > 0 {
			p.tracef("\n")
		}
		p.field(t.Field(i))
		p.string(t.Tag(i))
	}
}

func (p *exporter) field(f *types.Var) {
	if !f.IsField() {
		log.Fatalf("gcimporter: field expected")
	}

	p.pos(f)
	p.fieldName(f)
	p.typ(f.Type())
}

func (p *exporter) iface(t *types.Interface) {
	// TODO(gri): enable importer to load embedded interfaces,
	// then emit Embeddeds and ExplicitMethods separately here.
	p.int(0)

	n := t.NumMethods()
	if trace && n > 0 {
		p.tracef("methods {>\n")
		defer p.tracef("<\n} ")
	}
	p.int(n)
	for i := 0; i < n; i++ {
		if trace && i > 0 {
			p.tracef("\n")
		}
		p.method(t.Method(i))
	}
}

func (p *exporter) method(m *types.Func) {
	sig := m.Type().(*types.Signature)
	if sig.Recv() == nil {
		log.Fatalf("gcimporter: method expected")
	}

	p.pos(m)
	p.string(m.Name())
	if m.Name() != "_" && !ast.IsExported(m.Name()) {
		p.pkg(m.Pkg(), false)
	}

	// interface method; no need to encode receiver.
	p.paramList(sig.Params(), sig.Variadic())
	p.paramList(sig.Results(), false)
}

// fieldName is like qualifiedName but it doesn't record the package
// for blank (_) or exported names.
func (p *exporter) fieldName(f *types.Var) {
	name := f.Name()

	// anonymous field with unexported base type name: use "?" as field name
	// (bname != "" per spec, but we are conservative in case of errors)
	if f.Anonymous() {
		base := f.Type()
		if ptr, ok := base.(*types.Pointer); ok {
			base = ptr.Elem()
		}
		if named, ok := base.(*types.Named); ok && !named.Obj().Exported() {
			name = "?"
		}
	}

	p.string(name)
	if name == "?" || name != "_" && !f.Exported() {
		p.pkg(f.Pkg(), false)
	}
}

func (p *exporter) paramList(params *types.Tuple, variadic bool) {
	// use negative length to indicate unnamed parameters
	// (look at the first parameter only since either all
	// names are present or all are absent)
	n := params.Len()
	if n > 0 && params.At(0).Name() == "" {
		n = -n
	}
	p.int(n)
	for i := 0; i < params.Len(); i++ {
		q := params.At(i)
		t := q.Type()
		if variadic && i == params.Len()-1 {
			t = &dddSlice{t.(*types.Slice).Elem()}
		}
		p.typ(t)
		if n > 0 {
			p.string(q.Name())
			p.pkg(q.Pkg(), false)
		}
		p.string("") // no compiler-specific info
	}
}

func (p *exporter) value(x constant.Value) {
	if trace {
		p.tracef("= ")
	}

	switch x.Kind() {
	case constant.Bool:
		tag := falseTag
		if constant.BoolVal(x) {
			tag = trueTag
		}
		p.tag(tag)

	case constant.Int:
		if v, exact := constant.Int64Val(x); exact {
			// common case: x fits into an int64 - use compact encoding
			p.tag(int64Tag)
			p.int64(v)
			return
		}
		// uncommon case: large x - use float encoding
		// (powers of 2 will be encoded efficiently with exponent)
		p.tag(floatTag)
		p.float(constant.ToFloat(x))

	case constant.Float:
		p.tag(floatTag)
		p.float(x)

	case constant.Complex:
		p.tag(complexTag)
		p.float(constant.Real(x))
		p.float(constant.Imag(x))

	case constant.String:
		p.tag(stringTag)
		p.string(constant.StringVal(x))

	case constant.Unknown:
		// package contains type errors
		p.tag(unknownTag)

	default:
		log.Fatalf("gcimporter: unexpected value %v (%T)", x, x)
	}
}

func (p *exporter) float(x constant.Value) {
	if x.Kind() != constant.Float {
		log.Fatalf("gcimporter: unexpected constant %v, want float", x)
	}
	// extract sign (there is no -0)
	sign := constant.Sign(x)
	if sign == 0 {
		// x == 0
		p.int(0)
		return
	}
	// x != 0

	var f big.Float
	if v, exact := constant.Float64Val(x); exact {
		// float64
		f.SetFloat64(v)
	} else if num, denom := constant.Num(x), constant.Denom(x); num.Kind() == constant.Int {
		// TODO(gri): add big.Rat accessor to constant.Value.
		r := valueToRat(num)
		f.SetRat(r.Quo(r, valueToRat(denom)))
	} else {
		// Value too large to represent as a fraction => inaccessible.
		// TODO(gri): add big.Float accessor to constant.Value.
		f.SetFloat64(math.MaxFloat64) // FIXME
	}

	// extract exponent such that 0.5 <= m < 1.0
	var m big.Float
	exp := f.MantExp(&m)

	// extract mantissa as *big.Int
	// - set exponent large enough so mant satisfies mant.IsInt()
	// - get *big.Int from mant
	m.SetMantExp(&m, int(m.MinPrec()))
	mant, acc := m.Int(nil)
	if acc != big.Exact {
		log.Fatalf("gcimporter: internal error")
	}

	p.int(sign)
	p.int(exp)
	p.string(string(mant.Bytes()))
}

func valueToRat(x constant.Value) *big.Rat {
	// Convert little-endian to big-endian.
	// I can't believe this is necessary.
	bytes := constant.Bytes(x)
	for i := 0; i < len(bytes)/2; i++ {
		bytes[i], bytes[len(bytes)-1-i] = bytes[len(bytes)-1-i], bytes[i]
	}
	return new(big.Rat).SetInt(new(big.Int).SetBytes(bytes))
}

func (p *exporter) bool(b bool) bool {
	if trace {
		p.tracef("[")
		defer p.tracef("= %v] ", b)
	}

	x := 0
	if b {
		x = 1
	}
	p.int(x)
	return b
}

// ----------------------------------------------------------------------------
// Low-level encoders

func (p *exporter) index(marker byte, index int) {
	if index < 0 {
		log.Fatalf("gcimporter: invalid index < 0")
	}
	if debugFormat {
		p.marker('t')
	}
	if trace {
		p.tracef("%c%d ", marker, index)
	}
	p.rawInt64(int64(index))
}

func (p *exporter) tag(tag int) {
	if tag >= 0 {
		log.Fatalf("gcimporter: invalid tag >= 0")
	}
	if debugFormat {
		p.marker('t')
	}
	if trace {
		p.tracef("%s ", tagString[-tag])
	}
	p.rawInt64(int64(tag))
}

func (p *exporter) int(x int) {
	p.int64(int64(x))
}

func (p *exporter) int64(x int64) {
	if debugFormat {
		p.marker('i')
	}
	if trace {
		p.tracef("%d ", x)
	}
	p.rawInt64(x)
}

func (p *exporter) string(s string) {
	if debugFormat {
		p.marker('s')
	}
	if trace {
		p.tracef("%q ", s)
	}
	// if we saw the string before, write its index (>= 0)
	// (the empty string is mapped to 0)
	if i, ok := p.strIndex[s]; ok {
		p.rawInt64(int64(i))
		return
	}
	// otherwise, remember string and write its negative length and bytes
	p.strIndex[s] = len(p.strIndex)
	p.rawInt64(-int64(len(s)))
	for i := 0; i < len(s); i++ {
		p.rawByte(s[i])
	}
}

// marker emits a marker byte and position information which makes
// it easy for a reader to detect if it is "out of sync". Used for
// debugFormat format only.
func (p *exporter) marker(m byte) {
	p.rawByte(m)
	// Enable this for help tracking down the location
	// of an incorrect marker when running in debugFormat.
	if false && trace {
		p.tracef("#%d ", p.written)
	}
	p.rawInt64(int64(p.written))
}

// rawInt64 should only be used by low-level encoders
func (p *exporter) rawInt64(x int64) {
	var tmp [binary.MaxVarintLen64]byte
	n := binary.PutVarint(tmp[:], x)
	for i := 0; i < n; i++ {
		p.rawByte(tmp[i])
	}
}

// rawByte is the bottleneck interface to write to p.out.
// rawByte escapes b as follows (any encoding does that
// hides '$'):
//
//	'$'  => '|' 'S'
//	'|'  => '|' '|'
//
// Necessary so other tools can find the end of the
// export data by searching for "$$".
// rawByte should only be used by low-level encoders.
func (p *exporter) rawByte(b byte) {
	switch b {
	case '$':
		// write '$' as '|' 'S'
		b = 'S'
		fallthrough
	case '|':
		// write '|' as '|' '|'
		p.out.WriteByte('|')
		p.written++
	}
	p.out.WriteByte(b)
	p.written++
}

// tracef is like fmt.Printf but it rewrites the format string
// to take care of indentation.
func (p *exporter) tracef(format string, args ...interface{}) {
	if strings.IndexAny(format, "<>\n") >= 0 {
		var buf bytes.Buffer
		for i := 0; i < len(format); i++ {
			// no need to deal with runes
			ch := format[i]
			switch ch {
			case '>':
				p.indent++
				continue
			case '<':
				p.indent--
				continue
			}
			buf.WriteByte(ch)
			if ch == '\n' {
				for j := p.indent; j > 0; j-- {
					buf.WriteString(".  ")
				}
			}
		}
		format = buf.String()
	}
	fmt.Printf(format, args...)
}

// Debugging support.
// (tagString is only used when tracing is enabled)
var tagString = [...]string{
	// Packages:
	-packageTag: "package",

	// Types:
	-namedTag:     "named type",
	-arrayTag:     "array",
	-sliceTag:     "slice",
	-dddTag:       "ddd",
	-structTag:    "struct",
	-pointerTag:   "pointer",
	-signatureTag: "signature",
	-interfaceTag: "interface",
	-mapTag:       "map",
	-chanTag:      "chan",

	// Values:
	-falseTag:    "false",
	-trueTag:     "true",
	-int64Tag:    "int64",
	-floatTag:    "float",
	-fractionTag: "fraction",
	-complexTag:  "complex",
	-stringTag:   "string",
	-unknownTag:  "unknown",
}
