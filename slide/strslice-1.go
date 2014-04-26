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
	} // Continued Next Slide
