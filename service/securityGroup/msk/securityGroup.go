package msk

import (
	"fmt"
	"go-sdk/service/securityGroup/securityGroupAvailable"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kafka"
)

type MskSecurityGroupId struct {
	MskClusterName   string
	SecurityGroupIds []map[string]int
}

func GetSecurityGroup(resourceName string) (error, []MskSecurityGroupId) {
	msks := []MskSecurityGroupId{}

	// 1. Load MSK Cluster List
	// 2. Find specific Cluster Name in MSK Cluster List

	// Load AWS credentials from the environment
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("ap-northeast-2"),
	})
	if err != nil {
		fmt.Println("Error creating AWS session:", err)
		return err, msks
	}

	// Create an MSK client
	svc := kafka.New(sess)

	// Set up the ListClusters input
	input := &kafka.ListClustersInput{}

	// Call ListClusters
	listResult, err := svc.ListClusters(input)
	if err != nil {
		fmt.Println("Error calling ListClusters:", err)
		return err, msks
	}

	// Find the cluster with the desired name
	var cluster []*kafka.ClusterInfo
	for _, c := range listResult.ClusterInfoList {
		if strings.Contains(*c.ClusterName, resourceName) {
			cluster = append(cluster, c)
			continue
		}
	}

	if cluster == nil {
		return nil, msks
	}

	for _, val := range cluster {
		msk := MskSecurityGroupId{
			MskClusterName:   *val.ClusterName,
			SecurityGroupIds: []map[string]int{},
		}
		for _, securityGroupId := range val.BrokerNodeGroupInfo.SecurityGroups {
			err, countInboundRules := securityGroupAvailable.CountInboundRules(*securityGroupId)
			if err != nil {
				return err, msks
			}
			msk.SecurityGroupIds = append(msk.SecurityGroupIds, map[string]int{*securityGroupId: countInboundRules})
		}
		msks = append(msks, msk)
	}

	return nil, msks
}
