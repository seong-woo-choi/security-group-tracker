package msk

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rds"
)

type RdsSecurityGroup struct {
	Arn              string
	SecurityGroupIds []string
}

func GetSecurityGroup(resourceName string) (error, []RdsSecurityGroup) {
	// Load AWS credentials from the environment
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("ap-northeast-2"),
	})
	if err != nil {
		fmt.Println("Error creating AWS session:", err)
		return err, nil
	}

	svc := rds.New(sess)

	result, err := svc.DescribeDBInstances(nil)
	if err != nil {
		fmt.Println(err)
		return err, nil
	}

	// fmt.Println(result.DBInstances)

	rdss := []RdsSecurityGroup{}

	// Iterate through the list of RDS instances
	for _, instance := range result.DBInstances {
		// Check if the RDS instance name matches the name you specified
		if strings.Contains(*instance.DBInstanceArn, resourceName) && *instance.Engine == "aurora-mysql" {
			rds := RdsSecurityGroup{
				Arn:              *instance.DBInstanceArn,
				SecurityGroupIds: []string{},
			}

			for _, securityGroupId := range instance.VpcSecurityGroups {
				rds.SecurityGroupIds = append(rds.SecurityGroupIds, *securityGroupId.VpcSecurityGroupId)
			}

			rdss = append(rdss, rds)
		}
	}

	return nil, rdss
}
