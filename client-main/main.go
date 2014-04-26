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
	"github.com/zumper/aws/gen/20140201/ec2"
)

func main() {
	if len(os.Args) < 5 {
		fmt.Printf("USAGE: %s ACCESS SECRET REGION InstanceId...\n",
			os.Args[0])
		return
	}
	creds := aws.Creds{os.Args[1], os.Args[2], ""}
	service := aws.Service{
		"ec2",
		os.Args[3], "ec2." + os.Args[3] + ".amazonaws.com",
		"2014-02-01",
	}
	describe := ec2.DescribeInstances{
		InstanceId: os.Args[4:],
	}
	signed := aws.SignV2(creds, service, describe.Params(), time.Now(), false)
	url := "https://" + service.Endpoint + "/?" + aws.QueryString(signed)
	fmt.Println(url)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	fmt.Println(string(body))
}
