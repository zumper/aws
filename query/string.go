// Copyright 2014 The aws Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package query

import (
	"sort"
	"strings"
)

func String(params ...map[string]string) string {
	var pslice []string
	for _, p := range params {
		for k, v := range p {
			pslice = append(pslice, k+"="+Escape(v))
		}
	}
	sort.Strings(pslice)
	return strings.Join(pslice, "&")
}
