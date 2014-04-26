func Generate(name string, wsdl WSDL) *ast.File {
	file := &ast.File{
		Name: ast.NewIdent(name),
	}
	file.Decls = append(file.Decls, newImport("time"))
	file.Decls = append(file.Decls, newImport("strconv"))
	return file
}
