func newFunc(name, typ string, fields []fieldSpec, top bool) *ast.FuncDecl {
	var f ast.FuncDecl
	f.Recv = &ast.FieldList{}
	f.Recv.List = append(f.Recv.List, newField(RECV_NAME, typ))
	f.Name = ast.NewIdent(name)
	f.Type = &ast.FuncType{}
	f.Body = &ast.BlockStmt{}
	f.Body.List = append(f.Body.List, newMakeMapStrStr(RET_NAME))
	// Next Slide
	f.Body.List = append(f.Body.List, newReturn(RET_NAME))
	f.Type.Results = &ast.FieldList{}
	f.Type.Results.List = append(f.Type.Results.List, newStrStrMap())
	return &f
}
