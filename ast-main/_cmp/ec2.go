package ec2

import (
	"strconv"
)

type DescribeInstances struct {
	InstanceId []string
	MaxResults *int
	NextToken  *string
	Filter     []FilterElem
}

type FilterElem struct {
	Item []ValueType
}

type ValueType struct {
	Name  string
	Value string
}

func (t DescribeInstances) Params() map[string]string {
	params := make(map[string]string)
	params["Action"] = "DescribeInstances"
	for i, val := range t.InstanceId {
		params["InstanceId"+"."+strconv.Itoa(i)] = val
	}
	if t.MaxResults != nil {
		params["MaxResults"] = strconv.Itoa(*t.MaxResults)
	}
	if t.NextToken != nil {
		params["NextToken"] = *t.NextToken
	}
	return params
}
