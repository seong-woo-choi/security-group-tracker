package msk

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/elbv2"
)

type AlbSecurityGroupId struct {
	AlbName          string
	SecurityGroupIds []string
}

func GetSecurityGroup(resourceName string) (error, []AlbSecurityGroupId) {
	albs := []AlbSecurityGroupId{}

	// Load AWS credentials from the environment
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("ap-northeast-2"),
	})
	if err != nil {
		fmt.Println("Error creating AWS session:", err)
		return err, albs
	}

	svc := elbv2.New(sess)

	// Set up the DescribeLoadBalancers input
	input := &elbv2.DescribeLoadBalancersInput{}

	// Call DescribeLoadBalancers
	result, err := svc.DescribeLoadBalancers(input)
	if err != nil {
		fmt.Println("Error calling DescribeLoadBalancers:", err)
		return err, albs
	}

	// Print the security group ID of the ALB
	for _, lb := range result.LoadBalancers {
		if strings.Contains(*lb.LoadBalancerName, resourceName) {
			alb := AlbSecurityGroupId{
				AlbName:          *lb.LoadBalancerName,
				SecurityGroupIds: []string{},
			}
			for _, sg := range lb.SecurityGroups {
				alb.SecurityGroupIds = append(alb.SecurityGroupIds, *sg)
			}
			albs = append(albs, alb)
		}
	}
	return nil, albs
}
