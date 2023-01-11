package elasticache

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

type Ec2SecurityGroup struct {
	Ec2Name          string
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

	// DescribeSecurityGroups 가 아니라 DescribeEc2 로 해야한다..
	svc := ec2.New(sess)
	input := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			{
				Name: aws.String("tag:Name"),
				Values: []*string{
					aws.String(resourceName),
				},
			},
		},
	}

	result, err := svc.DescribeInstances(input)
	if err != nil {
		fmt.Println("Error calling DescribeSecurityGroups", err)
		return err, ec2s
	}

	for _, a := range result.Reservations {
		for _, b := range a.Instances {
			for _, c := range b.Tags {
				if *c.Key == "Name" {
					ec2 := Ec2SecurityGroup{
						Ec2Name:          *c.Value,
						SecurityGroupIds: []string{},
					}
					for _, sg := range b.SecurityGroups {
						ec2.SecurityGroupIds = append(ec2.SecurityGroupIds, *sg.GroupId)
					}
					ec2s = append(ec2s, ec2)
				}
			}
		}
	}

	return nil, ec2s
}
