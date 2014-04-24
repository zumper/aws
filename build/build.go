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
	STR      = "xs:string"
	BOOL     = "xs:boolean"
	INT      = "xs:int"
	INTEGER  = "xs:integer"
	DATETIME = "xs:dateTime"
	DOUBLE   = "xs:double"
	LONG     = "xs:long"
)

type fieldSpec struct {
	name, typ       string
	slice, optional bool
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
		file.Decls = append(file.Decls, newStruct(strings.Title(elem.Name), fields))
		file.Decls = append(file.Decls, newFunc(METH_NAME, strings.Title(elem.Name)))
	}
	for _, ct := range srv.ComplexType {
		fields, importTime := buildFields(ct, srv.ComplexType)
		timeUsed = timeUsed || importTime
		file.Decls = append(file.Decls, newStruct(strings.Title(ct.Name), fields))
		file.Decls = append(file.Decls, newFunc(METH_NAME, strings.Title(ct.Name)))
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
			switch elem.Type {
			case BOOL:
				fs.typ = "bool"
			case DATETIME:
				importTime = true
				fs.typ = "time.Time"
			case DOUBLE:
				fs.typ = "float64"
			case LONG:
				fs.typ = "int64"
			case INT:
				fs.typ = "int32"
			case INTEGER:
				fs.typ = "int64"
			default:
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

func newFunc(name, typ string) *ast.FuncDecl {
	var f ast.FuncDecl
	f.Recv = &ast.FieldList{}
	f.Recv.List = append(f.Recv.List, newRecvr(typ))

	f.Name = ast.NewIdent(name)
	f.Type = &ast.FuncType{}
	f.Body = &ast.BlockStmt{}
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

func newRecvr(typ string) *ast.Field {
	var obj ast.Object
	obj.Name = RECV_NAME
	obj.Kind = ast.Var

	recvIdent := ast.NewIdent(RECV_NAME)
	recvIdent.Obj = &obj

	var field ast.Field
	field.Names = append(field.Names, recvIdent)
	field.Type = ast.NewIdent(typ)
	return &field
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
