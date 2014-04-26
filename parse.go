// Copyright 2014 The aws Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package aws

import (
	"encoding/xml"
	"io"
)

const (
	BINDING     = "binding"
	SERVICE     = "service"
	MESSAGE     = "message"
	OPERATION   = "operation"
	TYPES       = "types"
	SCHEMA      = "schema"
	ELEMENT     = "element"
	COMPLEXTYPE = "complexType"
	GROUP       = "group"
)

type WSDL struct {
	Operation   map[string]Operation
	Message     map[string]Message
	Element     map[string]Element
	ComplexType map[string]ComplexType
	Group       map[string]Group
}

func makeWSDL() WSDL {
	return WSDL{
		make(map[string]Operation),
		make(map[string]Message),
		make(map[string]Element),
		make(map[string]ComplexType),
		make(map[string]Group),
	}
}

func NewWSDL(in io.Reader) (WSDL, error) {
	wsdl := makeWSDL()
	dec := xml.NewDecoder(in)
	types, schema := 0, 0
	var token xml.Token
	var err error
	for token, err = dec.Token(); err == nil; token, err = dec.Token() {
		switch elem := token.(type) {
		case xml.StartElement:
			switch elem.Name.Local {
			case BINDING, SERVICE:
				err = dec.Skip()
			case MESSAGE:
				var msg Message
				err = dec.DecodeElement(&msg, &elem)
				if err == nil {
					wsdl.Message[msg.Name] = msg
				}
				wsdl.Message[msg.Name] = msg
			case OPERATION:
				var op Operation
				err = dec.DecodeElement(&op, &elem)
				if err == nil {
					wsdl.Operation[op.Name] = op
				}
				wsdl.Operation[op.Name] = op
			case ELEMENT:
				if types > 0 && schema > 0 {
					var el Element
					err = dec.DecodeElement(&el, &elem)
					if err == nil {
						wsdl.Element[el.Name] = el
					}
				}
			case GROUP:
				var group Group
				err = dec.DecodeElement(&group, &elem)
				if err == nil {
					wsdl.Group[group.Name] = group
				}
			case COMPLEXTYPE:
				var ct ComplexType
				err = dec.DecodeElement(&ct, &elem)
				if err == nil {
					wsdl.ComplexType[ct.Name] = ct
				}
			case TYPES:
				types++
			case SCHEMA:
				schema++
			}
		case xml.EndElement:
			switch elem.Name.Local {
			case TYPES:
				types--
			case SCHEMA:
				schema--
			}
		}

		if err != nil {
			return wsdl, err
		}
	}
	if err == io.EOF {
		err = nil
	}
	return wsdl, err
}

type Message struct {
	Name string   `xml:"name,attr"`
	Part PartElem `xml:"part"`
}

type PartElem struct {
	Name    string `xml:"name,attr"`
	Element string `xml:"element,attr"`
}

type Operation struct {
	Name   string     `xml:"name,attr"`
	Input  InputElem  `xml:"input"`
	Output OutputElem `xml:"output"`
}

type InputElem struct {
	Message string `xml:"message,attr"`
}

type OutputElem struct {
	Message string `xml:"message,attr"`
}

type Element struct {
	Name      string `xml:"name,attr"`
	Type      string `xml:"type,attr"`
	MinOccurs string `xml:"minOccurs,attr"`
	MaxOccurs string `xml:"maxOccurs,attr"`
}

type Group struct {
	Name   string   `xml:"name,attr"`
	Choice []Choice `xml:"choice>element"`
}

type Choice struct {
	Name string `xml:"name,attr"`
	Type string `xml:"type,attr"`
}

type ComplexType struct {
	Name    string    `xml:"name,attr"`
	Element []Element `xml:"sequence>element"`
	Group   GroupRef  `xml:"sequence>group"`
	Choice  []Choice  `xml:"sequence>choice>element"`
}

type GroupRef struct {
	Ref string `xml:"ref,attr"`
}
