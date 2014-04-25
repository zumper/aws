// Copyright 2014 The aws Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package build

import (
	"go/ast"
	"go/token"
	"sort"
	"strings"

	"github.com/zumper/aws/parse"
)

const (
	DATETIME = "xs:dateTime"
)

var BUILTINS map[string]string
var BUILTINS_UNQUAL map[string]string

func init() {
	BUILTINS = map[string]string{
		"xs:boolean": "bool",
		DATETIME:     "time.Time",
		"xs:double":  "float64",
		"xs:long":    "int64",
		"xs:int":     "int32",
		"xs:integer": "int64",
		"xs:string":  "string",
	}

	BUILTINS_UNQUAL = map[string]string{
		"boolean":  "bool",
		"dateTime": "time.Time",
		"double":   "float64",
		"long":     "int64",
		"int":      "int32",
		"integer":  "int64",
		"string":   "string",
	}
}

type fieldSpec struct {
	name, typ                string
	slice, optional, builtin bool
}

const METH_NAME = "Params"

func isUnqualBuiltin(typ string) bool {
	_, builtin := BUILTINS_UNQUAL[typ]
	return builtin
}

func Resolve(srv parse.Service) (req, resp map[string][]string) {
	req, resp = make(map[string][]string), make(map[string][]string)
	for _, op := range srv.Operation {
		inName := Unqualify(op.Input.Message)
		msgIn := srv.Message[inName]
		inName = ResolveMsg(inName, srv)
		req[inName] = []string{} // Some types have no deps
		for _, r := range resolve(Unqualify(msgIn.Part.Element), srv) {
			req[inName] = append(req[inName], r)
		}

		outName := Unqualify(op.Output.Message)
		msgOut := srv.Message[outName]
		outName = ResolveMsg(outName, srv)
		resp[outName] = []string{} // Some types have not deps
		for _, r := range resolve(Unqualify(msgOut.Part.Element), srv) {
			resp[outName] = append(resp[outName], r)
		}
	}
	return
}

func resolve(typ string, srv parse.Service) []string {
	var dep []string

	if isUnqualBuiltin(typ) {
		return dep
	}
	if e, ok := srv.Element[Unqualify(typ)]; ok {
		etype := Unqualify(e.Type)
		if !isUnqualBuiltin(etype) {
			for _, d := range resolve(Unqualify(e.Type), srv) {
				d = Unqualify(d)
				if !isUnqualBuiltin(d) {
					dep = append(dep, d)
				}
			}
		}
	} else if c, ok := srv.ComplexType[Unqualify(typ)]; ok {
		// ignore groups and choices
		for _, e := range c.Element {
			etype := Unqualify(e.Type)
			if !isUnqualBuiltin(etype) {
				dep = append(dep, etype)
			}
			for _, d := range resolve(etype, srv) {
				d = Unqualify(d)
				if !isUnqualBuiltin(d) {
					dep = append(dep, d)
				}
			}
		}
	}
	return dep
}

func ResolveMsg(name string, srv parse.Service) string {
	return Unqualify(srv.Message[name].Part.Element)
}

func TypesV2(name string, srv parse.Service) *ast.File {
	file := &ast.File{
		Name: ast.NewIdent(name),
	}
	file.Decls = append(file.Decls, newImport("time"))
	file.Decls = append(file.Decls, newImport("strconv")) // TODO remove if unnecessary
	var timeUsed bool

	req, resp := Resolve(srv)
	var reqOps []string
	for r := range req {
		reqOps = append(reqOps, r)
	}
	sort.Strings(reqOps)

	done := make(map[string]interface{})

	for _, op := range reqOps {
		opElem := srv.Element[op]
		opType := Unqualify(opElem.Type)
		ct := srv.ComplexType[opType]
		fields, importTime := buildFields(ct, srv.ComplexType)
		file.Decls = append(file.Decls, newStruct(op, fields))
		file.Decls = append(file.Decls, newFunc(METH_NAME, op, fields, true))
		timeUsed = timeUsed || importTime
		for _, dep := range req[op] {
			if _, ok := done[dep]; ok {
				continue
			} else {
				done[dep] = nil
			}
			ct := srv.ComplexType[dep]
			fields, importTime := buildFields(ct, srv.ComplexType)
			timeUsed = timeUsed || importTime
			elemName := strings.Title(dep)
			file.Decls = append(file.Decls, newStruct(elemName, fields))
			file.Decls = append(file.Decls, newFunc(METH_NAME, elemName, fields, false))
		}
	}
	var respRoot []string
	for r := range resp {
		respRoot = append(respRoot, r)
	}
	sort.Strings(respRoot)

	for _, root := range respRoot {
		rootElem := srv.Element[root]
		rootType := Unqualify(rootElem.Type)
		done[rootType] = nil
		ct := srv.ComplexType[rootType]
		fields, importTime := buildFields(ct, srv.ComplexType)
		timeUsed = timeUsed || importTime
		elemName := strings.Title(rootElem.Name)
		file.Decls = append(file.Decls, newStruct(elemName, fields))
		for _, dep := range resp[root] {
			if _, ok := done[dep]; ok {
				continue
			} else {
				done[dep] = nil
			}
			ct := srv.ComplexType[dep]
			fields, importTime := buildFields(ct, srv.ComplexType)
			timeUsed = timeUsed || importTime
			elemName := strings.Title(dep)
			file.Decls = append(file.Decls, newStruct(elemName, fields))
		}
	}
	if !timeUsed {
		file.Decls = file.Decls[1:]
	}
	return file
}

func buildFields(ct parse.ComplexType, ctmap map[string]parse.ComplexType) ([]fieldSpec, bool) {
	var fields []fieldSpec
	var importTime bool
	cont, isCont := isContainer(ct, ct.Name, ctmap)
	if isCont {
		fields = append(fields, cont)
	} else {
		for _, elem := range ct.Element { // TODO group
			var fs fieldSpec
			fs.name = strings.Title(elem.Name)
			if fs.typ, fs.builtin = BUILTINS[elem.Type]; fs.builtin {
				if elem.Type == DATETIME {
					importTime = true
				}
			} else {
				fs.typ = Unqualify(elem.Type)
			}

			if cfield, ok := ctmap[fs.typ]; ok {
				cont, isCont := isContainer(cfield, elem.Name, ctmap)
				if isCont {
					fields = append(fields, cont)
					continue
				}
			}
			if elem.MinOccurs == "0" {
				fs.optional = true
			}
			if elem.MaxOccurs == "unbounded" {
				fs.slice = true
			}
			fields = append(fields, fs)
		}
	}
	return fields, importTime
}

func hasSetSuffix(name string) bool {
	return (strings.HasSuffix(name, "SetType") ||
		strings.HasSuffix(name, "InfoType") ||
		strings.HasSuffix(name, "SetRequestType") ||
		strings.HasSuffix(name, "SetResponseType"))
}

func isContainer(ct parse.ComplexType, fname string, ctmap map[string]parse.ComplexType) (fieldSpec, bool) {
	var fs fieldSpec
	var container bool
	if strings.HasSuffix(fname, "Set") &&
		hasSetSuffix(ct.Name) &&
		len(ct.Element) == 1 &&
		len(ct.Group.Ref) == 0 &&
		len(ct.Choice) == 0 &&
		ct.Element[0].Name == "item" &&
		ct.Element[0].MaxOccurs == "unbounded" {

		container = true
		fs.name = strings.Title(strings.TrimSuffix(fname, "Set"))
		fs.typ = Unqualify(ct.Element[0].Type)

		if cttyp, ok := ctmap[fs.typ]; ok &&
			len(cttyp.Element) == 1 &&
			len(cttyp.Group.Ref) == 0 &&
			len(cttyp.Choice) == 0 {
			fs.typ = Unqualify(cttyp.Element[0].Type)
			fs.name = strings.Title(cttyp.Element[0].Name)
		}
		fs.slice = true
		fs.optional = ct.Element[0].MinOccurs == "0"
	}
	return fs, container
}

const RECV_NAME = "t"

const ARG_NAME = "prefix"
const ARG_TYPE = "string"

const RET_NAME = "params"

const ACTION_KEY = "Action"

func newFunc(name, typ string, fields []fieldSpec, top bool) *ast.FuncDecl {
	var f ast.FuncDecl
	f.Recv = &ast.FieldList{}
	f.Recv.List = append(f.Recv.List, newField(RECV_NAME, typ))

	f.Name = ast.NewIdent(name)
	f.Type = &ast.FuncType{}
	f.Body = &ast.BlockStmt{}

	f.Body.List = append(f.Body.List, newMakeMapStrStr(RET_NAME))

	if top {
		f.Body.List = append(f.Body.List, newMapAssign(RET_NAME, ACTION_KEY, typ))
		for _, field := range fields {
			if field.slice && field.typ == "string" {
				f.Body.List = append(f.Body.List,
					addStrListToMap(RECV_NAME, field.name, RET_NAME))
			}
		}
	} else {
		f.Body.List = append(f.Body.List, newDotStrConcat(ARG_NAME, typ))
		f.Type.Params = &ast.FieldList{}
		f.Type.Params.List = append(f.Type.Params.List, newField(ARG_NAME, ARG_TYPE))

	}
	f.Body.List = append(f.Body.List, newReturn(RET_NAME))

	f.Type.Results = &ast.FieldList{}
	f.Type.Results.List = append(f.Type.Results.List, newStrStrMap())

	return &f
}

// TODO import strconv
const IDX_NAME = "i"
const VAL_NAME = "val"

func addStrListToMap(typ, key, mapName string) *ast.RangeStmt {
	v := ast.NewIdent(VAL_NAME)
	v.Obj = &ast.Object{
		Kind: ast.Var,
		Name: VAL_NAME,
	}
	sel := &ast.SelectorExpr{
		X:   ast.NewIdent(typ),
		Sel: ast.NewIdent(key),
	}
	k := ast.NewIdent(IDX_NAME)
	k.Obj = &ast.Object{
		Kind: ast.Var,
		Name: IDX_NAME,
		Decl: &ast.AssignStmt{
			Lhs: []ast.Expr{k, v},
			Tok: token.DEFINE,
			Rhs: []ast.Expr{
				&ast.UnaryExpr{
					Op: token.RANGE,
					X:  sel,
				},
			},
		},
	}
	rang := &ast.RangeStmt{
		Key:   k,
		Value: v,
		Tok:   token.DEFINE,
		X:     sel,
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						&ast.IndexExpr{
							X: ast.NewIdent(mapName),
							Index: &ast.BinaryExpr{
								X:  newStrConcat(key, "."),
								Op: token.ADD,
								Y:  newPkgCallExpr("strconv", "Itoa", IDX_NAME),
							},
						},
					},
					Tok: token.ASSIGN,
					Rhs: []ast.Expr{
						ast.NewIdent(VAL_NAME),
					},
				},
			},
		},
	}
	return rang
}

func newPkgCallExpr(pkg, fun, arg string) *ast.CallExpr {
	return &ast.CallExpr{
		Fun: &ast.SelectorExpr{
			X:   ast.NewIdent(pkg),
			Sel: ast.NewIdent(fun),
		},
		Args: []ast.Expr{
			ast.NewIdent(arg),
		},
	}
}

func newVar(name string) *ast.Ident {
	var obj ast.Object
	obj.Name = name
	obj.Kind = ast.Var

	ident := ast.NewIdent(name)
	ident.Obj = &obj
	return ident
}

func newField(name, typ string) *ast.Field {
	var field ast.Field
	field.Names = append(field.Names, newVar(name))
	field.Type = ast.NewIdent(typ)
	return &field
}

func newMakeMapStrStr(name string) *ast.AssignStmt {
	var Map ast.MapType
	Map.Key = ast.NewIdent("string")
	Map.Value = ast.NewIdent("string")

	var mke ast.CallExpr
	mke.Fun = ast.NewIdent("make")
	mke.Args = append(mke.Args, &Map)

	var assign ast.AssignStmt
	assign.Lhs = append(assign.Lhs, newVar(name))
	assign.Tok = token.DEFINE
	assign.Rhs = append(assign.Rhs, &mke)
	return &assign
}

func newDotStrConcat(name, val string) *ast.AssignStmt {
	return &ast.AssignStmt{
		Lhs: []ast.Expr{ast.NewIdent(name)},
		Tok: token.ADD_ASSIGN,
		Rhs: []ast.Expr{&ast.BasicLit{
			Kind:  token.STRING,
			Value: "\".\"",
		},
		},
	}
}

func newStrConcat(a, b string) *ast.BinaryExpr {
	return &ast.BinaryExpr{
		X: &ast.BasicLit{
			Kind:  token.STRING,
			Value: "\"" + a + "\"",
		},
		Op: token.ADD,
		Y: &ast.BasicLit{
			Kind:  token.STRING,
			Value: "\"" + b + "\"",
		},
	}
}

func newMapAssign(name, key, val string) *ast.AssignStmt {
	idx := ast.IndexExpr{
		X: ast.NewIdent(name),
		Index: &ast.BasicLit{
			Kind:  token.STRING,
			Value: "\"" + key + "\"",
		},
	}
	var assign ast.AssignStmt
	assign.Lhs = append(assign.Lhs, &idx)
	assign.Tok = token.ASSIGN
	assign.Rhs = append(assign.Rhs, &ast.BasicLit{
		Kind:  token.STRING,
		Value: "\"" + val + "\"",
	})
	return &assign
}

func newReturn(name string) *ast.ReturnStmt {
	var ret ast.ReturnStmt
	ret.Results = append(ret.Results, ast.NewIdent(name))
	return &ret
}

func newStrStrMap() *ast.Field {
	var Map ast.MapType
	Map.Key = ast.NewIdent("string")
	Map.Value = ast.NewIdent("string")

	var field ast.Field
	field.Type = &Map
	return &field
}

func newImport(name string) *ast.GenDecl {
	var gd ast.GenDecl
	gd.Tok = token.IMPORT
	var is ast.ImportSpec
	var bl ast.BasicLit
	bl.Kind = token.STRING
	bl.Value = "\"" + name + "\""
	is.Path = &bl
	gd.Specs = append(gd.Specs, &is)
	return &gd
}

func newStruct(name string, fields []fieldSpec) *ast.GenDecl {
	var gd ast.GenDecl
	gd.Tok = token.TYPE
	var ts ast.TypeSpec
	ts.Name = ast.NewIdent(name)
	ts.Name.Obj = ast.NewObj(ast.Typ, name)
	var st ast.StructType
	var flist ast.FieldList

	for _, spec := range fields {
		var f ast.Field
		f.Names = append(f.Names, ast.NewIdent(spec.name))

		if spec.optional && !spec.slice {
			var t ast.StarExpr
			t.X = ast.NewIdent(spec.typ)
			f.Type = &t
		} else if spec.slice {
			var at ast.ArrayType
			at.Elt = ast.NewIdent(spec.typ)
			f.Type = &at
		} else {
			f.Type = ast.NewIdent(spec.typ)
		}

		flist.List = append(flist.List, &f)
	}
	st.Fields = &flist
	ts.Type = &st
	gd.Specs = []ast.Spec{&ts}
	return &gd
}

func Unqualify(msg string) string {
	idx := strings.Index(msg, ":")
	if idx < 0 {
		idx = 0
	} else {
		idx++
	}
	return msg[idx:]
}
