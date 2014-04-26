	if top {
		m := newMapAssign(RET_NAME, ACTION_KEY, typ)
		f.Body.List = append(f.Body.List, m)
		for _, field := range fields {
			if field.slice && field.typ == "string" {
				s := addStrListToMap(RECV_NAME, field.name, RET_NAME)
				f.Body.List = append(f.Body.List, s)
			}
		}
	} else {
		// For types which aren't actions
		f.Body.List = append(f.Body.List, newDotStrConcat(ARG_NAME, typ))
		f.Type.Params = &ast.FieldList{}
		nf := newField(ARG_NAME, ARG_TYPE)
		f.Type.Params.List = append(f.Type.Params.List, nf)

	}
