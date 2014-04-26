rm gen/20140201/ec2/ec2.go
tree gen
./gen-main/gen-main ./2014-02-01-ec2.wsdl > gen/20140201/ec2/ec2.go
tree gen
head gen/20140201/ec2/ec2.go
go build -o ec2-query client-main/main.go
sleep 15
./ec2-query $(head -1 ~/.aws-gophercon) $(tail -1 ~/.aws-gophercon) us-west-1 i-ddbebd82 i-dcbebd83 | less
