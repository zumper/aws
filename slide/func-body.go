	params["Action"] = "DescribeInstances"
	for i, val := range t.InstanceId {
		params["InstanceId"+"."+strconv.Itoa(i)] = val
	}

