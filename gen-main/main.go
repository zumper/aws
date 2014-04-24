// Copyright 2014 The aws Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"go/printer"
	"go/token"
	"os"

	"github.com/zumper/aws"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("USAGE: %s <path-to-wsdl>\n", os.Args[0])
		return
	}
	wsdlXml, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Printf("err:%v\n", err)
		return
	}

	wsdl, err := aws.NewWSDL(wsdlXml)
	if err != nil {
		fmt.Printf("err:%v\n", err)
		return
	}

	file := aws.TypesV2("ec2", wsdl)
	printer.Fprint(os.Stdout, token.NewFileSet(), file)
}
