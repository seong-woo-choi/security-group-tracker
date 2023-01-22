package docdb

import (
	"fmt"
	"go-sdk/service/securityGroup/securityGroupAvailable"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/docdb"
)

type DocdbSecurityGroup struct {
	DocdbArnName     string
	SecurityGroupIds []map[string]int
}

func GetSecurityGroup(resourceName string) (error, []DocdbSecurityGroup) {
	docdbs := []DocdbSecurityGroup{}

	// Load AWS credentials from the environment
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("ap-northeast-2"),
	})
	if err != nil {
		fmt.Println("Error creating AWS session:", err)
		return err, docdbs
	}

	svc := docdb.New(sess)

	input := &docdb.DescribeDBInstancesInput{
		Filters: []*docdb.Filter{
			{
				Name:   aws.String("engine"),
				Values: []*string{(aws.String("docdb"))},
			},
		},
	}

	result, err := svc.DescribeDBInstances(input)
	if err != nil {
		fmt.Println("Error calling DescribeClusters", err)
		return err, docdbs
	}

	for _, instance := range result.DBInstances {
		if strings.Contains(*instance.DBInstanceArn, resourceName) {
			docdb := DocdbSecurityGroup{
				DocdbArnName:     *instance.DBInstanceArn,
				SecurityGroupIds: []map[string]int{},
			}
			for _, securityGroupId := range instance.VpcSecurityGroups {
				err, countInboundRules := securityGroupAvailable.CountInboundRules(*securityGroupId.VpcSecurityGroupId)
				if err != nil {
					return err, docdbs
				}
				docdb.SecurityGroupIds = append(docdb.SecurityGroupIds, map[string]int{*securityGroupId.VpcSecurityGroupId: countInboundRules})
			}
			docdbs = append(docdbs, docdb)
		}
	}

	return nil, docdbs
}
