# AWS Client for Go #
## Versions ##
### v0.0.2 ###
* Generate EC2 client code, handles Action and any non-nested []strings
* Basic query command line for DescribeInstances with InstanceIds

### v0.0.1 ###
* Basic AWS Query V2 signing from github.com/stu-art/awsclient
* Parse EC2 WSDL
* Generate valid Go structs from EC2 WSDL

## Instructions for v0.0.2 ##

Use `build` to generate Go code representing EC2 data structures from the [EC2 WSDL].

    mkdir -p gen/20140201/ec2
    go run build/run/main.go ./2014-02-01-ec2.wsdl > gen/20140201/ec2/ec2.go

Perform a DescribeInstances call

    go run gen/run/main.go $AWS_ACCESS $AWS_SECRET $REGION $INSTANCEID0 $INSTANCEID1

## Instructions for v0.0.1 ##

Use `query` to sign Query Requests using the [Query V2 signing protocol][QV2].

    go run query/run/main.go $ACCESS $SECRET [$TOKEN]

In v0.0.1, the client does not use the generated code.

Use `build` to generate Go code representing EC2 data structures from the [EC2 WSDL].

    mkdir -p gen/20140201/ec2
    go run build/run/main.go ./2014-02-01-ec2.wsdl > gen/20140201/ec2/ec2.go

Make sure the data structure is valid Go code.

    go run gen/run/main.go

[EC2 WSDL]: https://s3.amazonaws.com/ec2-downloads/ec2.wsdl "EC2 2014-02-01 WSDL"
[QV2]: https://docs.aws.amazon.com/general/latest/gr/signature-version-2.html "Signature Version 2 Signing Process"