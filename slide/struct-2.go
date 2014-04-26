		var f ast.Field
		f.Names = append(f.Names, ast.NewIdent(spec.name))
		if spec.optional && !spec.slice {
			var expr ast.StarExpr
			expr.X = ast.NewIdent(spec.typ)
			f.Type = &expr
		} else if spec.slice {
			var at ast.ArrayType
			at.Elt = ast.NewIdent(spec.typ)
			f.Type = &at
		} else {
			f.Type = ast.NewIdent(spec.typ)
		}
		flist.List = append(flist.List, &f)
