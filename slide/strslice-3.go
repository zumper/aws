	rang := &ast.RangeStmt{
		Key:   k,
		Value: v,
		Tok:   token.DEFINE,
		X:     sel,
		Body: &ast.BlockStmt{
			List: []ast.Stmt{assign},
		},
	}
	return rang
