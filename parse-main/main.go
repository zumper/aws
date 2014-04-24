// Copyright 2014 The aws Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/zumper/aws"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("USAGE: %s <path-to-wsdl>\n", os.Args[0])
		return
	}
	file, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Printf("err:%v\n", err)
		return
	}

	wsdl, err := aws.NewWSDL(file)
	if err != nil {
		fmt.Printf("err:%v\n", err)
		return
	}
	for _, op := range wsdl.Operation {
		inMsg := wsdl.Message[Unqualify(op.Input.Message)]
		outMsg := wsdl.Message[Unqualify(op.Output.Message)]
		fmt.Printf("op: %s\nin: %s\nout: %s\n", op.Name, inMsg, outMsg)

		inElem := wsdl.Element[Unqualify(inMsg.Part.Element)]
		outElem := wsdl.Element[Unqualify(outMsg.Part.Element)]
		fmt.Printf("in-e: %s\nout-e: %s\n", inElem, outElem)

		inc := wsdl.ComplexType[Unqualify(inElem.Type)]
		outc := wsdl.ComplexType[Unqualify(outElem.Type)]
		fmt.Printf("in-c: %s\nout-c: %s\n", inc, outc)

		fields := make(map[string]string)
		for _, elem := range inc.Element {
			etype := Unqualify(elem.Type)
			fields[elem.Name] = etype
			fmt.Printf("\t%s->%s\n", elem.Name, etype)
			switch elem.Type {
			case STR, BOOL, INT:
			default:
				if IsStringMap(wsdl.ComplexType[etype], wsdl.ComplexType) {
					fmt.Printf("STRING MAP!!!\n")
				}
			}
		}

		fmt.Printf("++++++++++++++++++\n")
	}

	fmt.Printf("-----------------\n")
	for name, val := range wsdl.Message {
		fmt.Printf("%s\t%s\n", name, val)
	}
}

const (
	STR  = "xs:string"
	BOOL = "xs:bool"
	INT  = "xs:int"
)

func Unqualify(msg string) string {
	idx := strings.Index(msg, ":")
	if idx < 0 {
		idx = 0
	} else {
		idx++
	}
	return msg[idx:]
}

func IsStringMap(c aws.ComplexType, types map[string]aws.ComplexType) bool {
	if len(c.Group.Ref) > 0 || len(c.Choice) > 0 || len(c.Element) != 1 || c.Element[0].Name != "item" {
		return false
	}
	item := types[Unqualify(c.Element[0].Type)]
	var key, value bool
	for _, e := range item.Element {
		switch e.Name {
		case "key":
			key = true
		case "value":
			value = true
		}
	}
	return key && value
}
