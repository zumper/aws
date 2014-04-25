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

	file := build.TypesV2("ec2", srv)
	printer.Fprint(os.Stdout, token.NewFileSet(), file)
	//ast.Print(token.NewFileSet(), file)

	//	req, resp := build.Resolve(srv)

	/*
		reqMap := make(map[string]string)
		for _, r := range req {
			fmt.Printf("req: '%v'\n", r)
			reqMap[r] = ""
		}
		inter := make(map[string]string)
		for _, r := range resp {
			fmt.Printf("resp: '%v'\n", r)
			if _, ok := reqMap[r]; ok {
				inter[r] = ""
			}
		}
		for r, _ := range inter {
			fmt.Printf("BOTH: '%v'\n", r)
		}

	*/

	/*
		req, resp := build.Resolve(srv)
		for r, dep := range req {
			fmt.Printf("Req: %v\n", r)
			for _, d := range dep {
				fmt.Printf("\t%s\n", d)
			}
		}
		for r, dep := range resp {
			fmt.Printf("Resp: %v\n", r)
			for _, d := range dep {
				fmt.Printf("\t%s\n", d)
			}
		}
	*/
}
