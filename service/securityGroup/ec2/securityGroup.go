package elasticache

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

type Ec2SecurityGroup struct {
	Ec2GroupName     string
	SecurityGroupIds []string
}

func GetSecurityGroup(resourceName string) (error, []Ec2SecurityGroup) {
	resourceName = "*" + resourceName + "*"
	ec2s := []Ec2SecurityGroup{}

	// Load AWS credentials from the environment
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("ap-northeast-2"),
	})
	if err != nil {
		fmt.Println("Error creating AWS session:", err)
		return err, ec2s
	}

	svc := ec2.New(sess)
	input := &ec2.DescribeSecurityGroupsInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("group-name"),
				Values: []*string{aws.String(resourceName)},
			},
		},
	}

	result, err := svc.DescribeSecurityGroups(input)
	if err != nil {
		fmt.Println("Error calling DescribeSecurityGroups", err)
		return err, ec2s
	}

	for _, ec2Value := range result.SecurityGroups {
		ec2 := Ec2SecurityGroup{
			Ec2GroupName:     *ec2Value.GroupName,
			SecurityGroupIds: []string{},
		}
		ec2.SecurityGroupIds = append(ec2.SecurityGroupIds, *ec2Value.GroupId)
		ec2s = append(ec2s, ec2)
	}

	return nil, ec2s
}
