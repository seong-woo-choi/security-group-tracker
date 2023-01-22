package securityGroupAvailable

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func CountInboundRules(securityGroupId string) (error, int) {
	// Load AWS credentials from the environment
	count := 0
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("ap-northeast-2"),
	})
	if err != nil {
		fmt.Println("Error creating AWS session:", err)
		return err, 0
	}

	svc := ec2.New(sess)

	input := &ec2.DescribeSecurityGroupsInput{
		GroupIds: []*string{
			aws.String(securityGroupId),
		},
	}

	result, err := svc.DescribeSecurityGroups(input)
	if err != nil {
		return err, 0
	}
	for _, group := range result.SecurityGroups {
		for _, permissions := range group.IpPermissions {
			if len(permissions.IpRanges) != 0 {
				count += len(permissions.IpRanges)
			}
			if len(permissions.UserIdGroupPairs) != 0 {
				count += len(permissions.UserIdGroupPairs)
			}
		}
	}
	return nil, count
}
