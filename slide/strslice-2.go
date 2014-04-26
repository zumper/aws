	call := newPkgCallExpr("strconv", "Itoa", IDX_NAME)
	assign := &ast.AssignStmt{
		Lhs: []ast.Expr{
			&ast.IndexExpr{
				X: ast.NewIdent(mapName),
				Index: &ast.BinaryExpr{
					X:  newStrConcat(key, "."),
					Op: token.ADD,
					Y:  call,
				},
			},
		},
		Tok: token.ASSIGN,
		Rhs: []ast.Expr{
			ast.NewIdent(VAL_NAME),
		},
	} // Continued Next Slide
