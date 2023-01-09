package msk

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kafka"
)

type MskSecurityGroupId struct {
	ClusterName      string
	SecurityGroupIds []string
}

func GetSecurityGroup(resourceName string) (error, []MskSecurityGroupId) {
	// 1. Load MSK Cluster List
	// 2. Find Cluster to Search Name in MSK Cluster List
	// 3. Text to Stdout from specific kafka cluster

	// Load AWS credentials from the environment
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("ap-northeast-2"),
	})
	if err != nil {
		fmt.Println("Error creating AWS session:", err)
		return err, nil
	}

	// Create an MSK client
	svc := kafka.New(sess)

	// Set up the ListClusters input
	input := &kafka.ListClustersInput{}

	// Call ListClusters
	listResult, err := svc.ListClusters(input)
	if err != nil {
		fmt.Println("Error calling ListClusters:", err)
		return err, nil
	}

	//var cluster []string

	// Find the cluster with the desired name
	var cluster []*kafka.ClusterInfo
	for _, c := range listResult.ClusterInfoList {
		if strings.Contains(*c.ClusterName, resourceName) {
			cluster = append(cluster, c)
			continue
		}
	}

	msks := []MskSecurityGroupId{}

	if cluster == nil {
		return err, msks
	}

	for _, val := range cluster {
		msk := MskSecurityGroupId{
			ClusterName:      *val.ClusterName,
			SecurityGroupIds: []string{},
		}
		for _, securityGroupId := range val.BrokerNodeGroupInfo.SecurityGroups {
			msk.SecurityGroupIds = append(msk.SecurityGroupIds, *securityGroupId)
		}
		msks = append(msks, msk)
	}

	return nil, msks
}
