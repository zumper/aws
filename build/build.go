// Copyright 2014 The aws Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package build

import (
	"go/ast"
	"go/token"
	"strings"

	"github.com/zumper/aws/parse"
)

const (
	DATETIME = "xs:dateTime"
)

var BUILTINS map[string]string

func init() {
	BUILTINS = map[string]string{
		"xs:boolean": "bool",
		DATETIME:     "time.Time",
		"xs:double":  "float64",
		"xs:long":    "int64",
		"xs:int":     "int32",
		"xs:integer": "int64",
	}
}

type fieldSpec struct {
	name, typ                string
	slice, optional, builtin bool
}

const METH_NAME = "Params"

func Types(name string, srv parse.Service) *ast.File {
	file := &ast.File{
		Name: ast.NewIdent(name),
	}
	file.Decls = append(file.Decls, newImport("time"))
	var timeUsed bool
	for _, elem := range srv.Element {
		ct := srv.ComplexType[Unqualify(elem.Type)]
		fields, importTime := buildFields(ct, srv.ComplexType)
		timeUsed = timeUsed || importTime
		elemName := strings.Title(elem.Name)
		file.Decls = append(file.Decls, newStruct(elemName, fields))
		_, op := srv.Operation[elemName]
		file.Decls = append(file.Decls, newFunc(METH_NAME, elemName, fields, op))
	}

	for _, ct := range srv.ComplexType {
		fields, importTime := buildFields(ct, srv.ComplexType)
		timeUsed = timeUsed || importTime
		file.Decls = append(file.Decls, newStruct(strings.Title(ct.Name), fields))
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

func isContainer(ct parse.ComplexType, fname string, ctmap map[string]parse.ComplexType) (fieldSpec, bool) {
	var fs fieldSpec
	var container bool
	if strings.HasSuffix(fname, "Set") &&
		(strings.HasSuffix(ct.Name, "SetType") || strings.HasSuffix(ct.Name, "InfoType")) &&
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

func newFunc(name, typ string, fields []fieldSpec, op bool) *ast.FuncDecl {
	var f ast.FuncDecl
	f.Recv = &ast.FieldList{}
	f.Recv.List = append(f.Recv.List, newField(RECV_NAME, typ))

	f.Name = ast.NewIdent(name)
	f.Type = &ast.FuncType{}
	f.Body = &ast.BlockStmt{}

	f.Body.List = append(f.Body.List, newMakeMapStrStr(RET_NAME))

	if op {
		f.Body.List = append(f.Body.List, newMapAssign(RET_NAME, ACTION_KEY, typ))
	} else {
		f.Body.List = append(f.Body.List, newNonZeroStrConcat(RET_NAME, typ))
	}

	f.Body.List = append(f.Body.List, newReturn(RET_NAME))

	f.Type.Params = &ast.FieldList{}
	f.Type.Params.List = append(f.Type.Params.List, newField(ARG_NAME, ARG_TYPE))
	f.Type.Results = &ast.FieldList{}
	f.Type.Results.List = append(f.Type.Results.List, newStrStrMap())

	return &f
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

func newNonZeroStrConcat(name, val string) *ast.IfStmt {
	cond := &ast.BinaryExpr{
		X: &ast.CallExpr{
			Fun:  ast.NewIdent("len"),
			Args: []ast.Expr{ast.NewIdent(name)},
		},
		Op: token.GTR,
		Y: &ast.BasicLit{
			Kind:  token.INT,
			Value: "0",
		},
	}
	body := &ast.BlockStmt{
		List: []ast.Stmt{
			&ast.AssignStmt{
				Lhs: []ast.Expr{ast.NewIdent(name)},
				Tok: token.ADD_ASSIGN,
				Rhs: []ast.Expr{
					&ast.BinaryExpr{
						X: &ast.BinaryExpr{
							X: &ast.BasicLit{
								Kind:  token.STRING,
								Value: "\".\"",
							},
							Op: token.ADD,
							Y: &ast.BasicLit{
								Kind:  token.STRING,
								Value: "\"" + val + "\"",
							},
						},
						Op: token.ADD,
						Y: &ast.BasicLit{
							Kind:  token.STRING,
							Value: "\".\"",
						},
					},
				},
			},
		},
	}
	return &ast.IfStmt{
		Cond: cond,
		Body: body,
	}
}

func newMapAddOrStrConcat(name, key, val string) *ast.IfStmt {
	cond := &ast.BinaryExpr{
		X: &ast.CallExpr{
			Fun:  ast.NewIdent("len"),
			Args: []ast.Expr{ast.NewIdent(name)},
		},
		Op: token.GTR,
		Y: &ast.BasicLit{
			Kind:  token.INT,
			Value: "0",
		},
	}
	body := &ast.BlockStmt{
		List: []ast.Stmt{
			&ast.AssignStmt{
				Lhs: []ast.Expr{ast.NewIdent(name)},
				Tok: token.ADD_ASSIGN,
				Rhs: []ast.Expr{
					&ast.BinaryExpr{
						X: &ast.BinaryExpr{
							X: &ast.BasicLit{
								Kind:  token.STRING,
								Value: "\".\"",
							},
							Op: token.ADD,
							Y: &ast.BasicLit{
								Kind:  token.STRING,
								Value: "\"" + val + "\"",
							},
						},
						Op: token.ADD,
						Y: &ast.BasicLit{
							Kind:  token.STRING,
							Value: "\".\"",
						},
					},
				},
			},
		},
	}
	els := &ast.BlockStmt{
		List: []ast.Stmt{newMapAssign(name, key, val)},
	}
	return &ast.IfStmt{
		Cond: cond,
		Body: body,
		Else: els,
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
