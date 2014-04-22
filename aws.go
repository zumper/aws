// Derived from github.com/stu-art/awsclient/awsclient.go
// Copyright 2012 Stuart Tettemer and 2014 the aws authors.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package aws

type Service struct {
	Name, Region, Endpoint, Version string
}

type Creds struct {
	Access, Secret, SecurityToken string
}

type QueryRequest struct {
	Action string
	Params map[string]string
}
