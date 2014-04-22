// Copyright 2014 The aws Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"github.com/zumper/aws/gen/20140201/ec2"
)

func main() {
	dt := ec2.DeleteTagsResponse{}
	dt.Return = true
	dt.RequestId = "rid"
	fmt.Printf("%v\n", dt)
}
