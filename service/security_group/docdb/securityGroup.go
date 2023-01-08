package docdb

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/docdb"
)

type DocdbSecurityGroup struct {
	Arn              string
	SecurityGroupIds []string
}

func GetSecurityGroup(resourceName string) (error, []DocdbSecurityGroup) {
	// Load AWS credentials from the environment
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("ap-northeast-2"),
	})
	if err != nil {
		fmt.Println("Error creating AWS session:", err)
		return err, nil
	}

	svc := docdb.New(sess)
	input := &docdb.DescribeDBClustersInput{}

	result, err := svc.DescribeDBClusters(input)
	if err != nil {
		fmt.Println("Error calling DescribeClusters", err)
		return err, nil
	}

	// fmt.Println(result)

	// docdb 나 rds 의 경우 클러스터의 보안그룹을 인스턴스들이 따라가는 것인지?
	// 혹은 인스턴스 별로 개별 보안그룹을 줄 수 있는 것인지 궁금하다.

	docdbs := []DocdbSecurityGroup{}
	for _, cluster := range result.DBClusters {
		if strings.Contains(*cluster.DBClusterArn, resourceName) && *cluster.Engine == "docdb" {
			docdb := DocdbSecurityGroup{
				Arn:              *cluster.DBClusterArn,
				SecurityGroupIds: []string{},
			}
			for _, securityGroupId := range cluster.VpcSecurityGroups {
				docdb.SecurityGroupIds = append(docdb.SecurityGroupIds, *securityGroupId.VpcSecurityGroupId)
			}
			docdbs = append(docdbs, docdb)
		}
	}

	if len(docdbs) == 0 {
		fmt.Println("DocumentDB Cluster Not Found")
	}

	return nil, docdbs
}
