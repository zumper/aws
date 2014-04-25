// Derived from github.com/stu-art/awsclient/awsclient.go
// Copyright 2012 Stuart Tettemer and 2014 the aws Authors.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package sign

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"

	"strings"
	"time"

	"github.com/zumper/aws"
	"github.com/zumper/aws/query"
)

// Time format constants, see time.Format
const (
	ISO8601 = "2006-01-02T15:04:05Z"
)

func V2(creds aws.Creds, service aws.Service, req map[string]string,
	date time.Time, expires bool) map[string]string {

	params := map[string]string{
		"AWSAccessKeyId":   creds.Access,
		"SignatureMethod":  "HmacSHA256",
		"SignatureVersion": "2",
		"Version":          service.Version,
	}
	if len(creds.SecurityToken) > 0 {
		params["SecurityToken"] = creds.SecurityToken
	}
	if expires {
		params["Expires"] = date.UTC().Format(ISO8601)
	} else {
		params["Timestamp"] = date.UTC().Format(ISO8601)
	}
	for k, v := range req {
		params[k] = v
	}
	toSign := strings.Join([]string{
		"GET", // Method
		strings.ToLower(service.Endpoint),
		"/", // Path
		query.String(params)}, "\n")

	signature := hmac.New(sha256.New, []byte(creds.Secret))
	signature.Write([]byte(toSign))

	params["Signature"] = base64.StdEncoding.EncodeToString(signature.Sum(nil))
	return params
}
