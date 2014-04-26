func main() {
	srcFd, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	srcBytes, err := ioutil.ReadAll(srcFd)
	src := string(srcBytes)
	fsetSrc := token.NewFileSet()
	fileSrc, err := parser.ParseFile(fsetSrc, "", src, 0)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	printer.Fprint(os.Stdout, fsetSrc, fileSrc)
	fmt.Println("------------")
	ast.Print(fsetSrc, fileSrc)
	fmt.Println("------------")
}
