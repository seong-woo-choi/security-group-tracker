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
	// Load AWS credentials from the environment
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("ap-northeast-2"),
	})
	if err != nil {
		fmt.Println("Error creating AWS session:", err)
		return err, nil
	}

	svc := elbv2.New(sess)

	// Set up the DescribeLoadBalancers input
	input := &elbv2.DescribeLoadBalancersInput{}

	// Call DescribeLoadBalancers
	result, err := svc.DescribeLoadBalancers(input)
	if err != nil {
		fmt.Println("Error calling DescribeLoadBalancers:", err)
		return err, nil
	}

	fmt.Println(resourceName)
	// Print the security group ID of the ALB
	albs := []AlbSecurityGroupId{}
	for _, lb := range result.LoadBalancers {
		if strings.Contains(*lb.LoadBalancerName, resourceName) {
			alb := AlbSecurityGroupId{
				AlbName:          *lb.LoadBalancerName,
				SecurityGroupIds: []string{},
			}
			for _, sg := range lb.SecurityGroups {
				// fmt.Printf("%v: %v\n", *lb.LoadBalancerName, *sg)
				alb.SecurityGroupIds = append(alb.SecurityGroupIds, *sg)
			}
			albs = append(albs, alb)
		}
	}
	return nil, albs
}

/*
{
	[
		{
			"albname": "1"
			"security_group_id": ["a", "b", "c"]
		},
		{
			"albname": "2"
			"security_group_id": ["d", "e", "f"]
		},
	]
}
*/
