// Copyright 2014 The aws Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	//	"go/ast"
	"go/printer"
	"go/token"
	"os"

	"github.com/zumper/aws/build"
	"github.com/zumper/aws/parse"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("USAGE: %s <path-to-wsdl>\n", os.Args[0])
		return
	}
	wsdl, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Printf("err:%v\n", err)
		return
	}

	srv, err := parse.WSDL(wsdl)
	if err != nil {
		fmt.Printf("err:%v\n", err)
		return
	}

	file := build.Types("ec2", srv)
	if err != nil {
		fmt.Printf("err:%v\n", err)
		return
	}

	printer.Fprint(os.Stdout, token.NewFileSet(), file)
	//	ast.Print(token.NewFileSet(), file)
}
