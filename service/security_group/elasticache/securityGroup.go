package elasticache

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/elasticache"
)

type RedisSecurityGroup struct {
	Arn              string
	SecurityGroupIds []string
}

func GetSecurityGroup(resourceName string) (error, []RedisSecurityGroup) {
	// Load AWS credentials from the environment
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("ap-northeast-2"),
	})
	if err != nil {
		fmt.Println("Error creating AWS session:", err)
		return err, nil
	}

	svc := elasticache.New(sess)
	input := &elasticache.DescribeCacheClustersInput{}

	result, err := svc.DescribeCacheClusters(input)
	if err != nil {
		fmt.Println("Error calling DescribeCacheClusters", err)
		return err, nil
	}

	elasticaches := []RedisSecurityGroup{}
	for _, cluster := range result.CacheClusters {
		if strings.Contains(*cluster.ARN, resourceName) {
			fmt.Println(cluster)
			elasticache := RedisSecurityGroup{
				Arn:              *cluster.ARN,
				SecurityGroupIds: []string{},
			}
			for _, securityGroup := range cluster.SecurityGroups {
				elasticache.SecurityGroupIds = append(elasticache.SecurityGroupIds, *securityGroup.SecurityGroupId)
			}
			elasticaches = append(elasticaches, elasticache)
		}
	}

	return nil, elasticaches
}
