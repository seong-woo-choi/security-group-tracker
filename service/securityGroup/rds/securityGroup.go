package msk

import (
	"fmt"
	"go-sdk/service/securityGroup/securityGroupAvailable"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rds"
)

type RdsSecurityGroup struct {
	RdsArnName       string
	SecurityGroupIds []map[string]int
}

func GetSecurityGroup(resourceName string) (error, []RdsSecurityGroup) {
	rdss := []RdsSecurityGroup{}
	// Load AWS credentials from the environment
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("ap-northeast-2"),
	})
	if err != nil {
		fmt.Println("Error creating AWS session:", err)
		return err, rdss
	}

	svc := rds.New(sess)

	input := &rds.DescribeDBInstancesInput{
		Filters: []*rds.Filter{
			{
				Name:   aws.String("engine"),
				Values: []*string{aws.String("aurora-mysql")},
			},
		},
	}

	result, err := svc.DescribeDBInstances(input)
	if err != nil {
		fmt.Println(err)
		return err, rdss
	}

	// Iterate through the list of RDS instances
	for _, instance := range result.DBInstances {
		// Check if the RDS instance name matches the name you specified
		if strings.Contains(*instance.DBInstanceArn, resourceName) {
			rds := RdsSecurityGroup{
				RdsArnName:       *instance.DBInstanceArn,
				SecurityGroupIds: []map[string]int{},
			}
			for _, securityGroupId := range instance.VpcSecurityGroups {
				err, countInboundRules := securityGroupAvailable.CountInboundRules(*securityGroupId.VpcSecurityGroupId)
				if err != nil {
					return err, rdss
				}
				rds.SecurityGroupIds = append(rds.SecurityGroupIds, map[string]int{*securityGroupId.VpcSecurityGroupId: countInboundRules})
			}
			rdss = append(rdss, rds)
		}
	}

	return nil, rdss
}
