func newStruct(name string, fields []fieldSpec) *ast.GenDecl {
	var gd ast.GenDecl
	gd.Tok = token.TYPE
	typeSpec := ast.TypeSpec{Name: ast.NewIdent(name)}
	typeSpec.Name.Obj = ast.NewObj(ast.Typ, name)
	var flist ast.FieldList
	for _, spec := range fields {
		// Next slide
	}
	var st ast.StructType
	st.Fields = &flist
	typeSpec.Type = &st
	gd.Specs = []ast.Spec{&typeSpec}
	return &gd
}
