package docdb

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/docdb"
)

type DocdbSecurityGroup struct {
	DocdbArnName     string
	SecurityGroupIds []string
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

	fmt.Println(result.DBInstances)

	for _, instance := range result.DBInstances {
		if strings.Contains(*instance.DBInstanceArn, resourceName) {
			docdb := DocdbSecurityGroup{
				DocdbArnName:     *instance.DBInstanceArn,
				SecurityGroupIds: []string{},
			}
			for _, securityGroupId := range instance.VpcSecurityGroups {
				docdb.SecurityGroupIds = append(docdb.SecurityGroupIds, *securityGroupId.VpcSecurityGroupId)
			}
			docdbs = append(docdbs, docdb)
		}
	}

	return nil, docdbs
}
