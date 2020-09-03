package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"os"
	"strings"
	"fmt"
)


func main() {
	awsSession := session.Must(session.NewSessionWithOptions(session.Options{
			SharedConfigState: session.SharedConfigEnable,
		}))
	ec2Client := ec2.New(awsSession)
	clusterString := os.Getenv("KOPS_CLUSTER")
	if clusterString == "" {
		panic(fmt.Errorf("KOPS_CLUSTER environment variable is required: got %s", clusterString))
	}
	clusters := strings.Split(clusterString, ",")
	for i, c := range clusters {
		clusters[i] = fmt.Sprintf("kubernetes.io/cluster/%s", c)
	}
	filters := []*ec2.Filter{}

	filters = append(filters, &ec2.Filter{
		Name: aws.String("tag-key"),
		Values: aws.StringSlice(clusters),
		})
	desc := &ec2.DescribeNatGatewaysInput{}
	desc = desc.SetFilter(filters)

	var nats []*ec2.NatGateway
	err := ec2Client.DescribeNatGatewaysPages(
		desc,
		func(page *ec2.DescribeNatGatewaysOutput, lastPage bool) bool {
			for _, v := range page.NatGateways {
				nats = append(nats, v)
			}
			if lastPage {
				return false
			}
			return true
		})

	if err != nil {
		panic(err)
	}
	for _, nat := range nats {
		fmt.Printf("Nat Ip: %s\n", aws.StringValue(nat.NatGatewayAddresses[0].PublicIp))
	}
}
