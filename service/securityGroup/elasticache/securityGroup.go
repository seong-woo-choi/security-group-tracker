package elasticache

import (
	"fmt"
	"go-sdk/service/securityGroup/securityGroupAvailable"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/elasticache"
)

type RedisSecurityGroup struct {
	RedisArnName     string
	SecurityGroupIds []map[string]int
}

func GetSecurityGroup(resourceName string) (error, []RedisSecurityGroup) {
	elasticaches := []RedisSecurityGroup{}

	// Load AWS credentials from the environment
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("ap-northeast-2"),
	})
	if err != nil {
		fmt.Println("Error creating AWS session:", err)
		return err, elasticaches
	}

	svc := elasticache.New(sess)
	input := &elasticache.DescribeCacheClustersInput{}

	result, err := svc.DescribeCacheClusters(input)
	if err != nil {
		fmt.Println("Error calling DescribeCacheClusters", err)
		return err, elasticaches
	}

	for _, cluster := range result.CacheClusters {
		if strings.Contains(*cluster.ARN, resourceName) {
			elasticache := RedisSecurityGroup{
				RedisArnName:     *cluster.ARN,
				SecurityGroupIds: []map[string]int{},
			}
			for _, securityGroup := range cluster.SecurityGroups {
				err, countInboundRules := securityGroupAvailable.CountInboundRules(*securityGroup.SecurityGroupId)
				if err != nil {
					return err, elasticaches
				}
				elasticache.SecurityGroupIds = append(elasticache.SecurityGroupIds, map[string]int{*securityGroup.SecurityGroupId: countInboundRules})
			}
			elasticaches = append(elasticaches, elasticache)
		}
	}

	return nil, elasticaches
}
