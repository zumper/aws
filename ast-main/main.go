package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"io/ioutil"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("USAGE: %s <path-to-go-file>\n", os.Args[0])
		return
	}
	fmt.Printf("Opening '%s'\n", os.Args[1])
	srcFd, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Printf("stu:%v\n", err)
		return
	}
	srcBytes, err := ioutil.ReadAll(srcFd)
	src := string(srcBytes)
	fsetSrc := token.NewFileSet()
	fileSrc, err := parser.ParseFile(fsetSrc, "", src, 0)
	if err != nil {
		fmt.Printf("err:%v\n", err)
		return
	}
	printer.Fprint(os.Stdout, fsetSrc, fileSrc)
	fmt.Printf("------------\n")
	ast.Print(fsetSrc, fileSrc)
	fmt.Printf("------------\n")

}
