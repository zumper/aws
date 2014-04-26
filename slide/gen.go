package ec2

import "time"
import "strconv"

type DescribeInstances struct {
	InstanceId []string
	Filter     []FilterType
	NextToken  *string
	MaxResults *int32
}

func (t DescribeInstances) Params() map[string]string {
	params := make(map[string]string)
	params["Action"] = "DescribeInstances"
	for i, val := range t.InstanceId {
		params["InstanceId"+"."+strconv.Itoa(i)] = val
	}
	return params
}
