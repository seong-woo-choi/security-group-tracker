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
	input := &docdb.DescribeDBClustersInput{}

	result, err := svc.DescribeDBClusters(input)
	if err != nil {
		fmt.Println("Error calling DescribeClusters", err)
		return err, docdbs
	}

	for _, cluster := range result.DBClusters {
		if strings.Contains(*cluster.DBClusterArn, resourceName) && *cluster.Engine == "docdb" {
			docdb := DocdbSecurityGroup{
				DocdbArnName:     *cluster.DBClusterArn,
				SecurityGroupIds: []string{},
			}
			for _, securityGroupId := range cluster.VpcSecurityGroups {
				docdb.SecurityGroupIds = append(docdb.SecurityGroupIds, *securityGroupId.VpcSecurityGroupId)
			}
			docdbs = append(docdbs, docdb)
		}
	}

	return nil, docdbs
}
