// Copyright 2014 The aws Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/zumper/aws"
	"github.com/zumper/aws/query"
	"github.com/zumper/aws/sign"
)

func main() {
	if len(os.Args) != 3 && len(os.Args) != 4 {
		fmt.Printf("USAGE: %s ACCESS SECRET [\n", os.Args[0])
		return
	}
	access, secret := os.Args[1], os.Args[2]
	var token string
	if len(os.Args) == 4 {
		token = os.Args[3]
	}
	ec2 := aws.Service{"ec2", "us-west-1", "ec2.us-west-1.amazonaws.com", "2014-02-01"}
	creds := aws.Creds{access, secret, token}
	req := aws.QueryRequest{"DescribeInstances", nil}
	params := sign.V2(creds, ec2, req, time.Now(), false)
	url := "https://" + ec2.Endpoint + "/?" + query.String(params)
	fmt.Println(url)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("%r\n", err)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("%r\n", err)
		return
	}
	fmt.Println(string(body))
}
